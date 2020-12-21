installer: deps pkged.go
	cd cmd/kilt-installer && go build
	@cp cmd/kilt-installer/kilt-installer kilt

pkged.go:
	go get github.com/markbates/pkger/cmd/pkger
	pkger

deps: artifacts/kilt-cfn-macro.zip runtimes/cfn-macro/kilt.yaml

artifacts/kilt-cfn-macro.zip: artifacts ../runtimes/cloudformation/kilt.zip
	@cp ../runtimes/cloudformation/kilt.zip $@

artifacts:
	mkdir -p artifacts

clean:
	@rm artifacts/*.zip || true
	@rm kilt || true
	@rm pkged.go || true
