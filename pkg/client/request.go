package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/puppetlabs/errawr-go/v2/pkg/encoding"
	"github.com/puppetlabs/relay/pkg/debug"
	"github.com/puppetlabs/relay/pkg/errors"
)

type RequestOptions struct {
	method           string
	path             string
	headers          map[string]string
	BodyEncodingType BodyEncodingType
	body             interface{}
	responseBody     interface{}
}

type RequestOptionSetter func(*RequestOptions)

func WithMethod(method string) RequestOptionSetter {
	return func(opts *RequestOptions) {
		opts.method = method
	}
}

func WithPath(path string) RequestOptionSetter {
	return func(opts *RequestOptions) {
		opts.path = path
	}
}

func WithHeaders(headers map[string]string) RequestOptionSetter {
	return func(opts *RequestOptions) {
		opts.headers = headers
	}
}

func WithBody(body interface{}) RequestOptionSetter {
	return func(opts *RequestOptions) {
		opts.body = body
	}
}

func WithBodyEncodingType(bodyEncodingType BodyEncodingType) RequestOptionSetter {
	return func(opts *RequestOptions) {
		opts.BodyEncodingType = bodyEncodingType
	}
}

func WithResponseInto(responseBody interface{}) RequestOptionSetter {
	return func(opts *RequestOptions) {
		opts.responseBody = responseBody
	}
}

type BodyEncoding interface {
	ContentType() string
	Encode(interface{}) (io.ReadWriter, errors.Error)
}

type BodyEncodingType string

const (
	BodyEncodingTypeJSON BodyEncodingType = "json"
	BodyEncodingTypeYAML BodyEncodingType = "yaml"
)

var mapEncodingTypeToEncoding = map[BodyEncodingType]BodyEncoding{
	BodyEncodingTypeJSON: &JSONBodyEncoding{},
	BodyEncodingTypeYAML: &YAMLBodyEncoding{},
}

type JSONBodyEncoding struct{}

func (j *JSONBodyEncoding) ContentType() string {
	return fmt.Sprintf("application/vnd.puppet.relay.%s+json", APIVersion)
}

func (j *JSONBodyEncoding) Encode(body interface{}) (io.ReadWriter, errors.Error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, errors.NewClientInternalError().WithCause(err)
		}
	}

	return buf, nil
}

type YAMLBodyEncoding struct{}

func (y *YAMLBodyEncoding) ContentType() string {
	return fmt.Sprintf("application/vnd.puppet.relay.%s+yaml", APIVersion)
}

func (y *YAMLBodyEncoding) Encode(body interface{}) (io.ReadWriter, errors.Error) {
	var buf io.ReadWriter

	bodyString, ok := body.(string)

	if !ok {
		return nil, errors.NewClientInternalError()
	}

	if body != nil {
		buf = bytes.NewBufferString(bodyString)
	}

	return buf, nil
}

func (c *Client) Request(setters ...RequestOptionSetter) errors.Error {
	const (
		defaultMethod           = http.MethodGet
		defaultBodyEncodingType = BodyEncodingTypeJSON
	)

	opts := &RequestOptions{
		method:           defaultMethod,
		BodyEncodingType: defaultBodyEncodingType,
	}

	for _, setter := range setters {
		setter(opts)
	}

	contextConfig, ok := c.config.ContextConfig[c.config.CurrentContext]
	if !ok {
		return errors.NewClientInternalError()
	}

	rel := &url.URL{Path: opts.path}
	u := contextConfig.Domains.APIDomain.ResolveReference(rel)

	encoding, ok := mapEncodingTypeToEncoding[opts.BodyEncodingType]

	if !ok {
		encodingTypeError := errors.NewClientInvalidEncodingType(string(opts.BodyEncodingType))

		return errors.NewClientInternalError().WithCause(encodingTypeError)
	}

	buf, buferr := encoding.Encode(opts.body)

	if buferr != nil {
		return buferr
	}

	req, reqerr := http.NewRequest(opts.method, u.String(), buf)

	if reqerr != nil {
		return errors.NewClientInternalError().WithCause(reqerr)
	}

	// defaults
	req.Header.Set("Accept", fmt.Sprintf("application/vnd.puppet.relay.%s+json", APIVersion))

	if opts.body != nil {
		req.Header.Set("Content-Type", encoding.ContentType())
	}

	// authorization
	token, terr := c.getToken()

	if terr != nil {
		return errors.NewClientInternalError().WithCause(terr)
	}

	if token != nil {
		req.Header.Set("Authorization", token.Bearer())
	}

	// overrides
	for name, value := range opts.headers {
		req.Header.Set(name, value)
	}

	// temporary but very useful debugging solution until we get real logging in place
	debug.LogDump(httputil.DumpRequestOut(req, true))

	resp, resperr := c.httpClient.Do(req)

	if resperr != nil {
		return errors.NewClientRequestError().WithCause(resperr)
	}

	// temporary but very useful debugging solution until we get real logging in place
	debug.LogDump(httputil.DumpResponse(resp, true))

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseError(resp)
	}

	if resp.Body != nil && opts.responseBody != nil {
		jerr := json.NewDecoder(resp.Body).Decode(opts.responseBody)

		if jerr != nil {
			return errors.NewClientInternalError().WithCause(jerr)
		}
	}

	return nil
}

func (c *Client) SetAuthorization() errors.Error {
	token, terr := c.getToken()

	if terr != nil {
		return errors.NewClientInternalError().WithCause(terr)
	}

	c.Api.GetConfig().AddDefaultHeader("Authorization", token.Bearer())

	return nil
}

type errorEnvelope struct {
	Error *encoding.ErrorDisplayEnvelope `json:"error"`
}

func parseError(resp *http.Response) errors.Error {
	// read body to buffer
	bytes, berr := ioutil.ReadAll(resp.Body)

	if berr != nil {
		return errors.NewClientRequestError()
	}

	// Attempt to parse relay api error envelope containing an errawr
	env := &errorEnvelope{}
	if err := json.Unmarshal(bytes, env); err == nil {
		return env.Error.AsError()
	}

	cause := errors.NewClientBadRequestBody(string(bytes))

	// otherwise return generic errors based on response code
	switch resp.StatusCode {
	case http.StatusNotFound:
		return errors.NewClientResponseNotFound().WithCause(cause)
	case http.StatusUnauthorized:
		return errors.NewClientUserNotAuthenticated().WithCause(cause)
	case http.StatusForbidden:
		return errors.NewClientUserNotAuthorized().WithCause(cause)
	}

	return errors.NewClientRequestError().WithCause(cause)
}
