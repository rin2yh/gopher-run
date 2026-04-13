.PHONY: build-wasm serve

build-wasm:
	GOOS=js GOARCH=wasm go build -o docs/gopher-run.wasm .
	cp $$(go env GOROOT)/lib/wasm/wasm_exec.js docs/wasm_exec.js

serve:
	go run github.com/hajimehoshi/wasmserve@latest .
