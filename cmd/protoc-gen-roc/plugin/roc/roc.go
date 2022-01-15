// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2015 The Go Authors.  All rights reserved.
// https://github.com/golang/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Package roc outputs roc service descriptions in Go code.
// It runs as a plugin for the Go protocol buffer compiler plugin.
// It is linked in to protoc-gen-go.

package roc

import (
	"fmt"
	"strings"

	pb "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

// generatedCodeVersion indicates a version of the generated code.
// It is incremented whenever an incompatibility between the generated code and
// the roc package is introduced; the generated code references
// a constant, roc.SupportPackageIsVersionN (where N is generatedCodeVersion).
const generatedCodeVersion = 1

// Paths for packages used by code generated in this file,
// relative to the import_prefix of the generator.Generator.
const (
	contextPkgPath    = "github.com/go-roc/roc/parcel/context"
	rocServicePkgPath = "github.com/go-roc/roc/service"
	parcelPkgPath     = "github.com/go-roc/roc/parcel"
	handlerPkgPath    = "github.com/go-roc/roc/service/handler"
	invokePkgPath     = "github.com/go-roc/roc/service/invoke"
	clientPkgPath     = "github.com/go-roc/roc/service/client"
	serverPkgPath     = "github.com/go-roc/roc/service/server"
)

func init() {
	generator.RegisterPlugin(new(roc))
}

// roc is an implementation of the Go protocol buffer compiler's
// plugin architecture.  It generates bindings for roc support.
type roc struct {
	gen *generator.Generator
}

// Name returns the name of this plugin, "github.com/go-roc/roc".
func (r *roc) Name() string {
	return "roc"
}

// The names for packages imported in the generated code.
// They may vary from the final path component of the import path
// if the name is used by other packages.
var (
	contextPkg    string
	rocServicePkg string
	parcelPkg     string
	handlerPkg    string
	invokePkg     string
	clientPkg     string
	serverPkg     string
)

// Init initializes the plugin.
func (r *roc) Init(gen *generator.Generator) {
	r.gen = gen
}

// Given a type name defined in a .proto, return its object.
// Also record that we're using it, to guarantee the associated import.
func (r *roc) objectNamed(name string) generator.Object {
	r.gen.RecordTypeUse(name)
	return r.gen.ObjectNamed(name)
}

// Given a type name defined in a .proto, return its name as we will print it.
func (r *roc) typeName(str string) string {
	return r.gen.TypeName(r.objectNamed(str))
}

// P forwards to g.gen.P.
func (r *roc) P(args ...interface{}) { r.gen.P(args...) }

// Generate generates code for the services in the given file.
func (r *roc) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	invokePkg = string(r.gen.AddImport(invokePkgPath))
	handlerPkg = string(r.gen.AddImport(handlerPkgPath))
	contextPkg = string(r.gen.AddImport(contextPkgPath))
	rocServicePkg = string(r.gen.AddImport(rocServicePkgPath))
	parcelPkg = string(r.gen.AddImport(parcelPkgPath))
	clientPkg = string(r.gen.AddImport(clientPkgPath))
	serverPkg = string(r.gen.AddImport(serverPkgPath))

	r.P("// Reference imports to suppress errors if they are not otherwise used.")
	r.P("var _ ", contextPkg, ".Context")
	r.P("var _ ", invokePkg, ".Invoke")
	r.P("var _ ", handlerPkg, ".Handler")
	r.P("var _ ", rocServicePkg, ".Service")
	r.P("var _ ", parcelPkg, ".RocPacket")
	r.P("var _ ", clientPkg, ".Client")
	r.P("var _ ", serverPkg, ".Server")
	r.P()

	// Assert version compatibility.
	r.P("// This is a compile-time assertion to ensure that this generated file")
	r.P("// is compatible with the roc package it is being compiled against.")
	r.P("const _ = ", rocServicePkg, ".SupportPackageIsVersion", generatedCodeVersion)
	r.P()

	for i, service := range file.FileDescriptorProto.Service {
		r.generateService(file, service, i)
	}
}

// GenerateImports generates the import declaration for this file.
func (r *roc) GenerateImports(file *generator.FileDescriptor) {}

func unexport(s string) string { return strings.ToLower(s[:1]) + s[1:] }

// deprecationComment is the standard comment added to deprecated
// messages, fields, enums, and enum values.
var deprecationComment = "// Deprecated: Do not use."

// generateService generates all the code for the named service.
func (r *roc) generateService(file *generator.FileDescriptor, service *pb.ServiceDescriptorProto, index int) {
	path := fmt.Sprintf("6,%d", index) // 6 means service.

	origServerName := service.GetName()
	fullServerName := origServerName
	if pkg := file.GetPackage(); pkg != "" {
		fullServerName = pkg + "." + fullServerName
	}
	serverName := generator.CamelCase(origServerName)
	deprecated := service.GetOptions().GetDeprecated()

	r.P()
	// service interface.
	if deprecated {
		r.P("//")
		r.P(deprecationComment)
	}
	r.P("type ", serverName, "Client interface {")
	for i, method := range service.Method {
		r.gen.PrintComments(fmt.Sprintf("%s,2,%d", path, i)) // 2 means method in a service.
		clientSignature := r.generateClientSignature(serverName, method)
		if clientSignature == "" {
			continue
		}
		r.P(clientSignature)
	}
	r.P("}")
	r.P()

	// service structure.
	r.P("type ", unexport(serverName), "Client struct {")
	r.P("c *", "client.Client")
	r.P("}")
	r.P()

	// NewClient factory.
	if deprecated {
		r.P(deprecationComment)
	}
	r.P("func New", serverName, "Client (c *", "client.Client) ", serverName, "Client {")
	r.P("return &", unexport(serverName), "Client{c}")
	r.P("}")
	r.P()

	// service method implementations.
	for _, method := range service.Method {
		r.generateClientMethod(serverName, method)
	}

	// Server interface.
	serverType := serverName + "Server"
	r.P("// ", serverType, " is the server API for ", serverName, " server.")
	if deprecated {
		r.P("//")
		r.P(deprecationComment)
	}
	r.P("type ", serverType, " interface {")
	for i, method := range service.Method {
		r.gen.PrintComments(fmt.Sprintf("%s,2,%d", path, i)) // 2 means method in a service.
		r.P(r.generateServerSignature(method))
	}
	r.P("}")
	r.P()

	// Server registration.
	if deprecated {
		r.P(deprecationComment)
	}
	r.P("func Register", serverName, "Server(s *server.Server", ", h ", serverType, ") {")
	r.P("var r = &", unexport(serverName), "Handler{h:h,s:s}")

	for _, v := range service.Method {
		if !v.GetClientStreaming() && !v.GetServerStreaming() {
			r.P(`s.RegisterHandler("/"+s.Name()+"/`, strings.ToLower(serverName), "/", strings.ToLower(*v.Name), `",r.`, *v.Name, ")")
		}
		if !v.GetClientStreaming() && v.GetServerStreaming() {
			r.P(`s.RegisterStreamHandler("/"+s.Name()+"/`, strings.ToLower(serverName), "/", strings.ToLower(*v.Name), `",r.`, *v.Name, ")")
		}

		if v.GetClientStreaming() && v.GetServerStreaming() {
			r.P(`s.RegisterChannelHandler("/"+s.Name()+"/`, strings.ToLower(serverName), "/", strings.ToLower(*v.Name), `",r.`, *v.Name, ")")
		}
	}
	r.P("}")
	r.P()

	r.P("type ", unexport(serverName), "Handler struct{")
	r.P("h ", serverName, "Server")
	r.P("s *server.Server")
	r.P("}")
	r.P()

	for _, method := range service.Method {
		r.generateServerMethod(serverName, method)
	}
	r.P()
}

// generateClientSignature returns the client-side signature for a method.
func (r *roc) generateClientSignature(serverName string, method *pb.MethodDescriptorProto) string {
	var (
		origMethodName = method.GetName()
		methodName     = generator.CamelCase(origMethodName)
	)

	if !method.GetClientStreaming() && !method.GetServerStreaming() {
		var (
			reqArg   = ", req *" + r.typeName(method.GetInputType())
			respName = "*" + r.typeName(method.GetOutputType())
		)

		//if r.GetRocApiPrefix(methodName) {
		//    return ""
		//}

		return fmt.Sprintf(
			"%s(c *%s.Context%s, opts ...invoke.InvokeOptions) (%s, error)",
			methodName,
			contextPkg,
			reqArg,
			respName,
		)
	}

	if !method.GetClientStreaming() && method.GetServerStreaming() {
		var (
			reqArg   = ", req *" + r.typeName(method.GetInputType())
			respName = "chan *" + r.typeName(method.GetOutputType())
		)
		return fmt.Sprintf(
			"%s(c *%s.Context%s, opts ...invoke.InvokeOptions) %s",
			methodName,
			contextPkg,
			reqArg,
			respName,
		)
	}

	if method.GetClientStreaming() && method.GetServerStreaming() {
		var (
			reqArg   = ", req chan *" + r.typeName(method.GetInputType())
			respName = "chan *" + r.typeName(method.GetOutputType())
		)
		return fmt.Sprintf(
			"%s(c *%s.Context%s, opts ...invoke.InvokeOptions) %s",
			methodName,
			contextPkg,
			reqArg,
			respName,
		)
	}

	return ""
}

func (r *roc) generateClientMethod(serverName string, method *pb.MethodDescriptorProto) {
	var (
		methodName = generator.CamelCase(method.GetName())
		outType    = r.typeName(method.GetOutputType())
	)

	if method.GetOptions().GetDeprecated() {
		r.P(deprecationComment)
	}

	if !method.GetServerStreaming() && !method.GetClientStreaming() {

		//if r.GetRocApiPrefix(methodName) {
		//    return
		//}

		r.P("func (cc *", unexport(serverName), "Client) ", r.generateClientSignature(serverName, method), "{")
		r.P("rsp := &", outType, "{}")
		r.P(`err := cc.c.InvokeRR(c, "/`, strings.ToLower(serverName), "/", strings.ToLower(methodName), `", req, rsp, opts...)`)
		r.P("return rsp, err")
		r.P("}")
		r.P()
		return
	}

	if !method.GetClientStreaming() && method.GetServerStreaming() {
		r.P("func (cc *", unexport(serverName), "Client) ", r.generateClientSignature(serverName, method), "{")
		r.P(`data :=cc.c.InvokeRS(c, "/`, strings.ToLower(serverName), "/", strings.ToLower(methodName), `", req, opts...)`)
		r.P("var rsp = make(chan *", outType, ",cap(data))")
		r.P("go func() {")
		r.P("for b := range data {")
		r.P("v := &", outType, "{}")
		r.P("err :=  c.Codec().Decode(b, v)")
		r.P("if err != nil {")
		r.P(" c.Errorf(\"client decode pakcet err=%v |method=%s |data=%s\", err, c.Method(), req.String())")
		r.P("continue")
		r.P("}")
		r.P("rsp <- v")
		r.P("}")
		r.P("close(rsp)")
		r.P("}()")
		r.P("return rsp")
		r.P("}")
		r.P()
	}

	if method.GetClientStreaming() && method.GetServerStreaming() {
		r.P("func (cc *", unexport(serverName), "Client) ", r.generateClientSignature(serverName, method), "{")
		r.P("var in = make(chan []byte,cap(req))")
		r.P("go func() {")
		r.P("for b := range req {")
		r.P("v, err := c.Codec().Encode(b)")
		r.P("if err != nil {")
		r.P(" c.Errorf(\"client encode pakcet err=%v |method=%s |data=%s\", err, c.Method(), b.String())")
		r.P("continue")
		r.P("}")
		r.P("in <- v")
		r.P("}")
		r.P("close(in)")
		r.P("}()")
		r.P()
		r.P(`data :=cc.c.InvokeRC(c, "/`, strings.ToLower(serverName), "/", strings.ToLower(methodName), `", in, opts...)`)
		r.P("var rsp = make(chan *", outType, ",cap(data))")
		r.P("go func() {")
		r.P("for b := range data {")
		r.P("v := &", outType, "{}")
		r.P("err := c.Codec().Decode(b, v)")
		r.P("if err != nil {")
		r.P(" c.Errorf(\"client decode pakcet err=%v |method=%s |data=%s\", err, c.Method(), string(b))")
		r.P("continue")
		r.P("}")
		r.P("rsp <- v")
		r.P("}")
		r.P("close(rsp)")
		r.P("}()")
		r.P("return rsp")
		r.P("}")
		r.P()
	}
}

// generateServerSignature returns the server-side signature for a method.
func (r *roc) generateServerSignature(method *pb.MethodDescriptorProto) string {
	origMethodName := method.GetName()
	methodName := generator.CamelCase(origMethodName)

	var reqArgs []string
	if !method.GetServerStreaming() && !method.GetClientStreaming() {

		reqArgs = append(reqArgs, "c *"+contextPkg+".Context")
		reqArgs = append(reqArgs, "req *"+r.typeName(method.GetInputType()))
		reqArgs = append(reqArgs, "rsp *"+r.typeName(method.GetOutputType()))

		return methodName + "(" + strings.Join(
			reqArgs,
			", ",
		) + ")"
	}

	if !method.GetClientStreaming() && method.GetServerStreaming() {
		reqArgs = append(reqArgs, "c *"+contextPkg+".Context")
		reqArgs = append(reqArgs, "req *"+r.typeName(method.GetInputType()))
		return methodName + "(" + strings.Join(
			reqArgs,
			", ",
		) + ") " + "chan *" + r.typeName(method.GetOutputType())
	}

	if method.GetClientStreaming() && method.GetServerStreaming() {
		reqArgs = append(reqArgs, "c *"+contextPkg+".Context")
		reqArgs = append(reqArgs, "req chan *"+r.typeName(method.GetInputType()))
		return methodName + "(" + strings.Join(
			reqArgs,
			", ",
		) + ",exit chan struct{}) " + "chan *" + r.typeName(method.GetOutputType())
	}

	return ""
}

func (r *roc) generateServerMethod(serverName string, method *pb.MethodDescriptorProto) {
	var (
		methodName = generator.CamelCase(method.GetName())
		inType     = r.typeName(method.GetInputType())
		outType    = r.typeName(method.GetOutputType())
	)

	if !method.GetServerStreaming() && !method.GetClientStreaming() {

		r.P(
			"func (r *",
			unexport(serverName),
			"Handler)",
			methodName,
			"(c *",
			contextPkg,
			".Context, req *parcel.RocPacket,interrupt handler.Interceptor) (rsp proto.Message, err error) {",
		)
		r.P("var in ", inType)
		r.P("err = c.Codec().Decode(req.Bytes(), &in)")
		r.P("if err != nil {")
		r.P(" c.Errorf(\"server decode packet err=%v |method=%s |data=%s\", err, c.Method(), req.String())")
		r.P("return nil,err")
		r.P("}")
		r.P("var out = ", outType, "{}")
		r.P("if interrupt == nil {")
		r.P("r.h.", methodName, "(c, &in,&out)")
		r.P("return &out, err")
		r.P("}")
		r.P("f := func(c *context.Context, req proto.Message) proto.Message {")
		r.P("r.h.", methodName, "(c, req.(*", inType, "),&out)")
		r.P("return &out")
		r.P("}")
		r.P("return interrupt(c, &in, f)")
		r.P("}")
		r.P()
		return
	}

	if !method.GetClientStreaming() && method.GetServerStreaming() {
		r.P(
			"func (r *",
			unexport(serverName),
			"Handler)",
			methodName,
			"(c *",
			contextPkg,
			".Context, req *parcel.RocPacket) chan proto.Message {",
		)
		r.P("var in ", inType)
		r.P("err := c.Codec().Decode(req.Bytes(), &in)")
		r.P("if err != nil {")
		r.P(" c.Errorf(\"server decode packet err=%v |method=%s |data=%s\", err, c.Method(), req.String())")
		r.P("return nil")
		r.P("}")
		r.P()
		r.P("out := r.h.", methodName, "(c, &in)")
		r.P("var rsp = make(chan proto.Message,cap(out))")
		r.P()
		r.P("go func() {")
		r.P("for d := range out {")
		r.P("rsp <- d")
		r.P("}")
		r.P("close(rsp)")
		r.P("}()")
		r.P("return rsp")
		r.P("}")
		r.P()
		return
	}

	if method.GetClientStreaming() && method.GetServerStreaming() {
		r.P(
			"func (r *",
			unexport(serverName),
			"Handler)",
			methodName,
			"(c *",
			contextPkg,
			".Context, req chan *parcel.RocPacket,exit chan struct{}) chan proto.Message {",
		)
		r.P("var in = make(chan *", inType, ",cap(req))")
		r.P("go func() {")
		r.P("for b := range req {")
		r.P("var v = &", inType, "{}")
		r.P("err := c.Codec().Decode(b.Bytes(), v)")
		r.P("if err != nil {")
		r.P("c.Errorf(\"server decode packet err=%v |method=%s |data=%s\", err, c.Method(), b.String())")
		r.P("continue")
		r.P("}")
		r.P("in <- v")
		r.P("parcel.Recycle(b)")
		r.P("}")
		r.P("close(in)")
		r.P("}()")
		r.P("out := r.h.", methodName, "(c, in,exit)")
		r.P("var rsp = make(chan proto.Message,cap(out))")
		r.P()
		r.P("go func() {")
		r.P("for d := range out {")
		r.P("rsp <- d")
		r.P("}")
		r.P("close(rsp)")
		r.P("}()")
		r.P("return rsp")
		r.P("}")
		r.P()
		return
	}
}
