kilt.zip: cmd/handler/handler
	cd cmd/handler/ && zip ../../kilt.zip ./handler

cmd/handler/handler:
	cd cmd/handler && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

clean:
	rm kilt.zip


.PHONY: clean