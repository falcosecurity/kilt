all: debug-bridge process-trace

debug-bridge:
	cd debug-bridge && CGO_ENABLED=0 go build .
process-trace:
	cd process-trace && CGO_ENABLED=0 go build .