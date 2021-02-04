installer: deps
	go get github.com/markbates/pkger/cmd/pkger
	pkger -o cmd/kilt-installer
	cd cmd/kilt-installer && CGO_ENABLED=0 go build .
	@cp cmd/kilt-installer/kilt-installer kilt

deps: artifacts/kilt-cfn-macro.zip cmd/kilt-installer/kilt.yaml

artifacts/kilt-cfn-macro.zip: artifacts ../runtimes/cloudformation/kilt.zip
	@cp ../runtimes/cloudformation/kilt.zip $@

artifacts:
	mkdir -p artifacts

clean:
	@rm artifacts/*.zip || true
	@rm kilt || true
	@rm pkged.go || true
