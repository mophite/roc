// Copyright (c) 2021 roc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//

package server

import (
	"errors"
	"net"
	"os"
	"github.com/go-roc/roc/parcel/codec"
	"strconv"

	"github.com/google/uuid"

	"github.com/go-roc/roc/internal/endpoint"
	"github.com/go-roc/roc/internal/namespace"
	"github.com/go-roc/roc/internal/registry"
	"github.com/go-roc/roc/internal/sig"
	"github.com/go-roc/roc/internal/transport"
	rs "github.com/go-roc/roc/internal/transport/rscoket"
	"github.com/go-roc/roc/internal/x"
	"github.com/go-roc/roc/parcel"
)

const SupportPackageIsVersion1 = 1

type option struct {

	//server id can be set,or random
	//this is the unique identifier of the service
	id string

	//server name,eg.srv.hello or api.hello
	name string

	//schema is the namespace in your service collection
	//it's the root name prefix
	schema string

	//random port
	//Random ports make you don’t need to care about specific port numbers,
	//which are more commonly used in internal services
	randPort *[2]int

	//socket tcp ip:port address
	tcpAddress string

	//websocket ip:port address
	wssAddress string

	//websocket relative path address of websocket
	wssPath string

	//server version
	version string

	//buffSize to data tunnel if it's need
	buffSize int

	//service discovery registry
	register registry.Registry

	//server transport
	transporter transport.Server

	//error packet
	//It will affect the format of the data you return
	err parcel.ErrorPackager

	//service discovery endpoint
	e *endpoint.Endpoint

	//receive system signal
	signal []os.Signal

	//wrapper some middleware
	//it's can be interrupt
	wrappers []parcel.Wrapper

	//when server exit,will do exit func
	exit []func()

	//codec packet tool
	cc codec.Codec
}

type Options func(*option)

func Codec(cc codec.Codec) Options {
	return func(option *option) {
		option.cc = cc
	}
}

func BuffSize(buffSize int) Options {
	return func(option *option) {
		option.buffSize = buffSize
	}
}

func Wrapper(wrappers ...parcel.Wrapper) Options {
	return func(option *option) {
		option.wrappers = append(option.wrappers, wrappers...)
	}
}

func Exit(exit ...func()) Options {
	return func(option *option) {
		option.exit = exit
	}
}

func Signal(signal ...os.Signal) Options {
	return func(option *option) {
		option.signal = signal
	}
}

func Transport(transport transport.Server) Options {
	return func(option *option) {
		option.transporter = transport
	}
}

func E(e *endpoint.Endpoint) Options {
	return func(option *option) {
		option.e = e
	}
}

func Port(port [2]int) Options {
	return func(option *option) {
		if port[0] > port[1] {
			panic("port index 0 must more than 1")
		}

		if port[0] < 10000 {
			panic("rand port for internal server suggest more than 10000")
		}

		option.randPort = &port
	}
}

func Error(err parcel.ErrorPackager) Options {
	return func(option *option) {
		option.err = err
	}
}

func WssAddress(address, path string) Options {
	return func(option *option) {
		option.wssAddress = address
		option.wssPath = path
	}
}

func Id(id string) Options {
	return func(option *option) {
		option.id = id
	}
}

func Namespace(name string) Options {
	return func(option *option) {
		option.name = name
	}
}

func TCPAddress(address string) Options {
	return func(option *option) {
		option.tcpAddress = address
	}
}

func Version(version string) Options {
	return func(option *option) {
		option.version = version
	}
}

func Registry(registry registry.Registry) Options {
	return func(option *option) {
		option.register = registry
	}
}

func newOpts(opts ...Options) option {
	opt := option{}

	for i := range opts {
		opts[i](&opt)
	}

	if opt.name == "" {
		opt.name = x.GetProjectName()
	}

	if opt.id == "" {
		//todo change to git commit id+timestamp
		opt.id = uuid.New().String()
	}

	ip, err := LocalIp()
	if err != nil {
		panic(err)
	}

	if opt.randPort == nil {
		opt.randPort = &[2]int{10000, 19999}
	}

	if opt.tcpAddress == "" {
		opt.tcpAddress = ip + ":" + strconv.Itoa(x.RandInt(opt.randPort[0], opt.randPort[1]))
	}

	if opt.version == "" {
		opt.version = namespace.DefaultVersion
	}

	if opt.wssAddress != "" && opt.wssPath == "" {
		opt.wssPath = "/roc/wss"
	}

	//if opt.ratelimit <= 0 {
	//	opt.ratelimit = math.MaxInt32
	//}

	if opt.schema == "" {
		opt.schema = namespace.DefaultSchema
	}

	if opt.err == nil {
		opt.err = parcel.DefaultErrorPacket
	}

	if opt.transporter == nil {
		opt.transporter = rs.NewServer(
			opt.tcpAddress,
			opt.wssAddress,
			opt.name,
			opt.buffSize,
		)
	}

	if opt.register == nil {
		opt.register = registry.NewRegistry(registry.Schema(opt.schema))
	}

	if opt.e == nil {
		opt.e = &endpoint.Endpoint{
			Id:      opt.id,
			Name:    opt.name,
			Version: opt.version,
			Address: opt.tcpAddress,
		}

		opt.e.Splicing(opt.schema)
	}

	if opt.signal == nil {
		opt.signal = sig.DefaultSignal
	}

	if opt.buffSize == 0 {
		opt.buffSize = 10
	}

	if opt.cc == nil {
		opt.cc = codec.DefaultCodec
	}

	return opt
}

func LocalIp() (string, error) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addr {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}

		}
	}

	return "", errors.New("cannot find local ip")
}
