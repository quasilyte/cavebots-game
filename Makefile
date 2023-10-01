wasm:
	GOARCH=wasm GOOS=js go build -ldflags="-s -w" -trimpath -o _web/main.wasm ./cmd/game

itchio-wasm: wasm
	cd _web && \
		mkdir -p ../bin && \
		rm -f ../bin/cavebots.zip && \
		zip ../bin/cavebots.zip -r main.wasm index.html wasm_exec.js
