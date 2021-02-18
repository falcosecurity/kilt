all: debug-bridge/debug-bridge process-trace/process-trace

debug-bridge/debug-bridge:
	cd debug-bridge && CGO_ENABLED=0 go build .
process-trace/process-trace:
	cd process-trace && CGO_ENABLED=0 go build .