package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/puppetlabs/errawr-go/v2/pkg/encoding"
	"github.com/puppetlabs/relay/pkg/errors"
)

func (c *Client) get(path string, headers map[string]string, responseBody interface{}) errors.Error {
	return c.request(http.MethodGet, path, headers, nil, responseBody)
}

func (c *Client) post(path string, headers map[string]string, body interface{}, responseBody interface{}) errors.Error {
	return c.request(http.MethodPost, path, headers, body, responseBody)
}

func (c *Client) put(path string, headers map[string]string, body interface{}, responseBody interface{}) errors.Error {
	return c.request(http.MethodPut, path, headers, body, responseBody)
}

func (c *Client) delete(path string, headers map[string]string, responseBody interface{}) errors.Error {
	return c.request(http.MethodDelete, path, headers, nil, responseBody)
}

func (c *Client) request(method string, path string, headers map[string]string, body interface{}, responseBody interface{}) errors.Error {
	rel := &url.URL{Path: path}
	u := c.config.APIDomain.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			errors.NewClientInternalError().WithCause(err).Bug()
		}
	}

	req, reqerr := http.NewRequest(method, u.String(), buf)

	if reqerr != nil {
		return errors.NewClientInternalError().WithCause(reqerr).Bug()
	}

	// defaults
	req.Header.Set("Accept", fmt.Sprintf("application/vnd.puppet.nebula.%v+json", API_VERSION))

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// authorization
	token, terr := c.getToken()

	if terr != nil {
		return errors.NewClientInternalError().WithCause(terr).Bug()
	}

	if token != nil {
		req.Header.Set("Authorization", token.Bearer())
	}

	// overrides
	for name, value := range headers {
		req.Header.Set(name, value)
	}

	debug(httputil.DumpRequestOut(req, true))

	resp, resperr := c.httpClient.Do(req)

	if resperr != nil {
		return errors.NewClientRequestError().WithCause(resperr)
	}

	debug(httputil.DumpResponse(resp, true))

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseError(resp)
	}

	jerr := json.NewDecoder(resp.Body).Decode(responseBody)

	if jerr != nil {
		return errors.NewClientInternalError().WithCause(jerr).Bug()
	}

	return nil
}

type errorEnvelope struct {
	Error *encoding.ErrorDisplayEnvelope `json:"error"`
}

func parseError(resp *http.Response) errors.Error {
	// Attempt to parse relay api error envelope containing an errawr
	env := &errorEnvelope{}
	if err := json.NewDecoder(resp.Body).Decode(env); err == nil {
		return env.Error.AsError()
	}

	// otherwise return generic errors based on response code
	switch resp.StatusCode {
	case http.StatusNotFound:
		return errors.NewClientResponseNotFound()
	case http.StatusUnauthorized:
		return errors.NewClientUserNotAuthenticated()
	case http.StatusForbidden:
		return errors.NewClientUserNotAuthorized()
	}

	return errors.NewClientRequestError()
}

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}
