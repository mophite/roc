package roc

import (
	"github.com/go-roc/roc/internal/namespace"
)

type InvokeOptions func(*InvokeOption)

type InvokeOption struct {

	//scope is the service discovery prefix key
	scope string

	//address is witch server you want to call
	address string

	//serviceName is witch server by service serviceName
	serviceName string

	//version is witch server by version
	version string

	//buffSize effective only requestChannel
	buffSize int

	//trace is request flow trace
	//it's will be from web client,or generated on initialize
	trace string
}

// WithTracing set tracing
func WithTracing(t string) InvokeOptions {
	return func(invokeOption *InvokeOption) {
		invokeOption.trace = t
	}
}

// BuffSize set buff size for requestChannel
func InvokeBuffSize(buffSize int) InvokeOptions {
	return func(invokeOption *InvokeOption) {
		invokeOption.buffSize = buffSize
	}
}

// WithName set service discover prefix with service serviceName
func WithName(name string, version ...string) InvokeOptions {
	return func(invokeOption *InvokeOption) {
		var ver = namespace.DefaultVersion

		// if no version ,use default version number
		if len(version) == 1 {
			ver = version[0]
		}

		invokeOption.scope = namespace.SplicingScope(name, ver)
		invokeOption.serviceName = name
		invokeOption.version = ver
	}
}

// WithAddress set service discover prefix with both service serviceName and address
func WithAddress(name, address string, version ...string) InvokeOptions {
	return func(invokeOption *InvokeOption) {
		var ver = namespace.DefaultVersion

		// if no version ,use default version number
		if len(version) == 1 {
			ver = version[0]
		}

		invokeOption.scope = namespace.SplicingScope(name, ver)
		invokeOption.address = address
		invokeOption.serviceName = name
		invokeOption.version = ver
	}
}

