package encoding

type ErrorDomain struct {
	Key   string `json:"key"`
	Title string `json:"title"`
}

type ErrorSection struct {
	Key   string `json:"key"`
	Title string `json:"title"`
}

type ErrorDescription struct {
	Friendly  string `json:"friendly"`
	Technical string `json:"technical"`
}

type HTTPErrorMetadata struct {
	Status  int                 `json:"status"`
	Headers map[string][]string `json:"headers"`
}

type ErrorMetadata struct {
	HTTPErrorMetadata *HTTPErrorMetadata `json:"http,omitempty"`
}

type ErrorArgument struct {
	Value       interface{} `json:"value"`
	Description string      `json:"description"`
}

type ErrorArguments map[string]*ErrorArgument

type ErrorItems map[string]*ErrorTransitEnvelope
