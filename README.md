# 
	* Install tinygo
	* 
go mod init github.com/dbgjerez/istio-wasm-filter-golang
go mod tidy
tinygo build -o filter.wasm -scheduler=none -target=wasi main.go
docker build -f Dockerfile.local -t b0rr3g0/wasm-go .
docker run -p 18000:18000 -p 38140:38140 b0rr3g0/wasm-go
