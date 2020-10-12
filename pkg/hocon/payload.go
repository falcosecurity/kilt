package hocon

import (
	"fmt"
	"github.com/falcosecurity/kilt/pkg/kilt"
	"github.com/go-akka/configuration/hocon"
)

func retrievePayload(object *hocon.HoconObject) (*kilt.Payload, error) {
	payload := new(kilt.Payload)
	payload.Type = kilt.Unknown
	if object.GetKey("url").IsString() {
		payload.Contents = object.GetKey("url").GetString()
		payload.Type = kilt.URL
	} else if object.GetKey("file").IsString() {
		payload.Contents = object.GetKey("file").GetString()
		payload.Type = kilt.LocalPath
	} else if object.GetKey("payload").IsString() {
		payload.Contents = object.GetKey("payload").GetString()
		payload.Type = kilt.Base64
	} else if object.GetKey("text").IsString() {
		payload.Contents = object.GetKey("text").GetString()
		payload.Type = kilt.Text
	}
	if object.GetKey("gzipped") != nil && object.GetKey("gzipped").IsString() && object.GetKey("gzipped").GetBoolean() {
		payload.Gzipped = true
	}

	if payload.Type == kilt.Unknown {
		return nil, fmt.Errorf("could not identify payload type for %s", object.ToString(1))
	}


	return payload, nil
}