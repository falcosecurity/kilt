package hocon

import (
	"encoding/json"
	"github.com/falcosecurity/kilt/pkg/kilt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestSimpleRuntime(t *testing.T) {
	targetInfoString, _ := ioutil.ReadFile("./fixtures/input.json")
	definitionString, _ := ioutil.ReadFile("./fixtures/kilt.cfg")
	k := NewKiltHocon(string(definitionString))
	info := new(kilt.TargetInfo)
	_ = json.Unmarshal(targetInfoString, info)
	r, _ := k.Runtime(info)

	assert.Equal(t, 1, len(r.Uploads), "expected 1 executable")
	assert.Equal(t, "https://storage.googleapis.com/kubernetes-release/release/v1.19.0/bin/linux/amd64/",
		r.Uploads[0].Payload.Contents)
	assert.Equal(t, kilt.URL, r.Uploads[0].Payload.Type)

	assert.Equal(t, 1, len(r.Executables))
	assert.Equal(t, "/bin/kubectl", r.Executables[0].Run[0])
}

func TestSimpleBuild(t *testing.T) {
	targetInfoString, _ := ioutil.ReadFile("./fixtures/input.json")
	definitionString, _ := ioutil.ReadFile("./fixtures/kilt.cfg")
	k := NewKiltHocon(string(definitionString))
	info := new(kilt.TargetInfo)
	_ = json.Unmarshal(targetInfoString, info)
	b, _ := k.Build(info)

	assert.Equal(t, "busybox:latest", b.Image)
	assert.Equal(t, "/falco/pdig", b.EntryPoint[0])
	assert.Equal(t, "true", b.EnvironmentVariables["TEST"])
	assert.Equal(t, 1, len(b.Resources))
}