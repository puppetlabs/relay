package secrets

import "github.com/puppetlabs/horsehead/v2/encoding/transfer"

// Secret is the envelope type for a secret.
type Secret struct {
	Key   string             `json:"key"`
	Value transfer.JSONOrStr `json:"value"`
}
