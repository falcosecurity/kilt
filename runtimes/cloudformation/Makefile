kilt.zip: cmd/handler/handler
	cd cmd/handler/ && zip ../../kilt.zip ./handler

cmd/handler/handler:
	cd cmd/handler && go build .

clean:
	rm kilt.zip


.PHONY: clean