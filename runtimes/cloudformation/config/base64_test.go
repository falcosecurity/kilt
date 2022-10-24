package config

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64(t *testing.T) {
	tests := []struct {
		decompress bool
		plaintext  string
		encoded    string

	}{
		{
			decompress: false,
			plaintext: `So long, and thanks for all the fish!`,
			encoded:   `U28gbG9uZywgYW5kIHRoYW5rcyBmb3IgYWxsIHRoZSBmaXNoIQ==`,
		},
		{
			decompress: false,
			plaintext: `nanananananananaBatman!`,
			encoded:   `bmFuYW5hbmFuYW5hbmFuYUJhdG1hbiE=`,
		},
		{
			decompress: false,
			plaintext: `{[(?!weird string!?)]`,
			encoded:   `e1soPyF3ZWlyZCBzdHJpbmchPyld`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.plaintext, func(t *testing.T) {
			encoded := base64.StdEncoding.EncodeToString([]byte(tc.plaintext))
			assert.Equal(t, tc.encoded, encoded, "the encoded string does not match the expected value")

			decoded := FromBase64(encoded, tc.decompress)
			assert.Equal(t, tc.plaintext, decoded, "decoded and plaintext strings do not match")
		})
	}
}
