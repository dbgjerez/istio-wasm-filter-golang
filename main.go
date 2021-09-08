package main

import (
	"fmt"
	"time"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

const (
	headerContentTypeJSON = "application/json"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	// Embed the default plugin context here,
	types.DefaultPluginContext
	shouldEchoBody bool
}

// Decide if execute an echo example or json wrapper
func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	if ctx.shouldEchoBody {
		return &echoBodyContext{}
	}
	return &setBodyContext{}
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	data, err := proxywasm.GetPluginConfiguration()
	if err != nil {
		proxywasm.LogCriticalf("error reading plugin configuration: %v", err)
	}
	ctx.shouldEchoBody = string(data) == "echo"
	return types.OnPluginStartStatusOK
}

type setBodyContext struct {
	// Embed the default root http context here,
	types.DefaultHttpContext
	totalRequestBodySize int
	contentType          string
}

// Override types.DefaultHttpContext.
func (ctx *setBodyContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	if _, err := proxywasm.GetHttpRequestHeader("content-length"); err != nil {
		if err := proxywasm.SendHttpResponse(400, nil, []byte("content must be provided")); err != nil {
			panic(err)
		}
		return types.ActionPause
	}

	// Remove Content-Length in order to prevent severs from crashing if we set different body from downstream.
	if err := proxywasm.RemoveHttpRequestHeader("content-length"); err != nil {
		panic(err)
	}

	contentType, err := proxywasm.GetHttpRequestHeader("Content-Type")
	if err != nil {
		proxywasm.LogCriticalf("error reading contentType: %v", err)
	}
	ctx.contentType = contentType

	return types.ActionContinue
}

const (
	timeFormat = "%d-%02d-%02dT%02d:%02d:%02d"
)

func now(t time.Time) string {
	return fmt.Sprintf(timeFormat,
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

// Override types.DefaultHttpContext.
func (ctx *setBodyContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	if !endOfStream {
		return types.ActionPause
	}

	originalBody, err := proxywasm.GetHttpRequestBody(0, bodySize)
	if err != nil {
		proxywasm.LogErrorf("failed to get request body: %v", err)
		return types.ActionContinue
	}

	res := []byte(`{"result": `)
	res = append(res, originalBody...)
	res = append(res, `, "time": "`...)
	res = append(res, now(time.Now())...)
	res = append(res, `"}`...)

	err = proxywasm.ReplaceHttpRequestBody(res)

	if err != nil {
		proxywasm.LogErrorf("failed to replace body: %v", err)
		return types.ActionContinue
	}
	return types.ActionContinue
}

type echoBodyContext struct {
	// mbed the default plugin context
	// so that you don't need to reimplement all the methods by yourself.
	types.DefaultHttpContext
	totalRequestBodySize int
}

// Override types.DefaultHttpContext.
func (ctx *echoBodyContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	ctx.totalRequestBodySize += bodySize
	if !endOfStream {
		// Wait until we see the entire body to replace.
		return types.ActionPause
	}

	// Send the request body as the response body.
	body, _ := proxywasm.GetHttpRequestBody(0, ctx.totalRequestBodySize)
	if err := proxywasm.SendHttpResponse(200, nil, body); err != nil {
		panic(err)
	}
	return types.ActionPause
}
