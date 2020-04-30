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

// JSONInterface allows arbitrary embedding of encoded data within any JSON
// type. The accepted interface values for the data correspond to those listed
// in the Go documentation for encoding/json.Unmarshal.
type JSONInterface struct {
	Data interface{}
}

func (ji JSONInterface) MarshalJSON() ([]byte, error) {
	switch dt := ji.Data.(type) {
	case map[string]interface{}:
		m := make(map[string]interface{}, len(dt))
		for k, v := range dt {
			m[k] = JSONInterface{Data: v}
		}

		return json.Marshal(m)
	case []interface{}:
		s := make([]interface{}, len(dt))
		for i, v := range dt {
			s[i] = JSONInterface{Data: v}
		}

		return json.Marshal(s)
	case string:
		jos, err := EncodeJSON([]byte(dt))
		if err != nil {
			return nil, err
		}

		return json.Marshal(jos)
	default:
		return json.Marshal(dt)
	}
}

func (ji *JSONInterface) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	var decode func(v interface{}) (interface{}, error)
	decode = func(v interface{}) (interface{}, error) {
		switch vt := v.(type) {
		case map[string]interface{}:
			ty, ok := vt["$encoding"].(string)
			if ok {
				d, _ := vt["data"].(string)

				b, err := JSON{EncodingType: EncodingType(ty), Data: d}.Decode()
				if err != nil {
					return nil, err
				}

				// Plop this back into one of the interface types supported by
				// json.Marshal and json.Unmarshal (string).
				v = string(b)
			}

			for k, v := range vt {
				v, err := decode(v)
				if err != nil {
					return nil, err
				}

				vt[k] = v
			}
		case []interface{}:
			for i, v := range vt {
				v, err := decode(v)
				if err != nil {
					return nil, err
				}

				vt[i] = v
			}
		}

		return v, nil
	}

	v, err := decode(v)
	if err != nil {
		return err
	}

	ji.Data = v
	return nil
}
