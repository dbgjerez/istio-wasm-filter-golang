# Prerequisites
These are the tool that you need to install to start developing it. If you only like to run the container, you will need only a docker installation. 
* golang: language used to develop the filter
* tinygo: a go compiler for small places. No available all the golang API but enough to develop an envoy filter
* docker:  containers is an easy way to test envoy without installing it
* IDE: a Golang IDE

# Dependencies


go mod init github.com/dbgjerez/istio-wasm-filter-golang
go mod tidy
tinygo build -o filter.wasm -scheduler=none -target=wasi main.go
docker build -f Dockerfile.local -t b0rr3g0/wasm-go .
docker run -p 18000:18000 -p 38140:38140 b0rr3g0/wasm-go
