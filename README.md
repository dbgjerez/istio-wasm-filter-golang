# Prerequisites
These are the tool that you need to install to start developing it. If you only like to run the container, you will need only a docker installation. 
* golang: language used to develop the filter
* tinygo: A Golang compiler
* docker:  containers is an easy way to test envoy without installing it
* IDE: a Golang IDE

# Build
## wasm
### Dependencies
Tinygo can't update the dependencies when you are using Golang modules. In this case, the go.mod file indicates it. 

To download the dependencies, you only have to put the following command: 
```bash
go mod tidy
```
At this point, you can start to develop and upgrade the code. 

### Tinygo
Wasm is the module that we have to deploy with our envoy server. In this case, to build it, we will use tinygo.

> Tinygo is a Golang compiler intended for use in small places like WebAssembly or microcontrollers.

```bash
tinygo build -o filter.wasm -scheduler=none -target=wasi  main.go
```
As a result, we can find a filter.wasm that is the WASM filter. 

How can we test it? To test the filter, we need an Envoy instance. One way to do this is by installing an instance. Another better way is using a container platform, such as docker, to test it.

## Docker
It is mandatory before building the Dockerfile to compile the WASM filter like the previous section. 

In this case, I'm using my user, b0rr3g0, but you can use your user. 
```bash
docker build -f Dockerfile.local -t b0rr3g0/wasm-go .
```
Once the container building finishes, we can run it.

# Run
To run the container you have to use the name indicated when it was built.
```bash
docker run -p 18000:18000 b0rr3g0/wasm-go
```
