package transfer

import "encoding/json"

// JSON is a convenient way to represent an encoding and data tuple in JSON.
type JSON struct {
	EncodingType EncodingType `json:"$encoding"`
	Data         string       `json:"data"`

	// Factories allows this struct to be configured to use a different set of
	// encoder/decoders than the default.
	Factories EncodeDecoderFactories `json:"-"`
}

// Decode finds the given encoder for this JSON data and decodes the data using
// it.
func (t JSON) Decode() ([]byte, error) {
	factories := t.Factories
	if factories == nil {
		factories = Encoders
	}

	encoder, found := factories[t.EncodingType]
	if !found {
		return nil, ErrUnknownEncodingType
	}

	return encoder().DecodeFromTransfer(t.Data)
}

// JSONOrStr is like the JSON type, but also allows NoEncodingType to be
// represented as a raw JSON string.
type JSONOrStr struct{ JSON }

func (tos JSONOrStr) MarshalJSON() ([]byte, error) {
	if tos.EncodingType == NoEncodingType {
		return json.Marshal(tos.Data)
	}

	return json.Marshal(tos.JSON)
}

func (tos *JSONOrStr) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		tos.JSON.Data = s
		return nil
	}

	return json.Unmarshal(data, &tos.JSON)
}
