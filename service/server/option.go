package server

import (
    "os"
    "strconv"
    "strings"

    "github.com/go-roc/roc/internal/registry"
    "github.com/go-roc/roc/internal/transport"
    rs "github.com/go-roc/roc/internal/transport/rsocket"
    "github.com/go-roc/roc/service/handler"
    "github.com/go-roc/roc/x"
    "github.com/go-roc/roc/x/fs"

    "github.com/go-roc/roc/internal/endpoint"
    "github.com/go-roc/roc/internal/sig"
    "github.com/go-roc/roc/parcel"
)

type Option struct {
    //server id can be set,or random
    //this is the unique identifier of the service
    id string

    //server name,eg.srv.hello or api.hello
    name string

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

    //http service address
    httpAddress string

    //buffSize to data tunnel if it's need
    buffSize int

    //server transport
    server transport.Server

    //error packet
    //It will affect the format of the data you return
    err parcel.ErrorPackager

    //receive system signal
    signal []os.Signal

    //wrapper some middleware
    //it's can be interrupt
    wrappers []handler.WrapperHandler

    //when server exit,will do exit func
    exit []func()

    //it must be unique in all of your handler path
    apiPrefix string

    //service discover registry
    registry registry.Registry
}

type Options func(option *Option)

func BuffSize(buffSize int) Options {
    return func(option *Option) {
        option.buffSize = buffSize
    }
}

func Wrapper(wrappers ...handler.WrapperHandler) Options {
    return func(option *Option) {
        option.wrappers = append(option.wrappers, wrappers...)
    }
}

func Exit(exit ...func()) Options {
    return func(option *Option) {
        option.exit = exit
    }
}

func Signal(signal ...os.Signal) Options {
    return func(option *Option) {
        option.signal = signal
    }
}

func Port(port [2]int) Options {
    return func(option *Option) {
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
    return func(option *Option) {
        option.err = err
    }
}

func WssAddress(address, path string) Options {
    return func(option *Option) {
        option.wssAddress = address
        option.wssPath = path
    }
}

func HttpAddress(address string) Options {
    return func(option *Option) {
        option.httpAddress = address
    }
}

func Id(id string) Options {
    return func(option *Option) {
        option.id = id
    }
}

func Namespace(name string) Options {
    return func(option *Option) {
        option.name = name
    }
}

func TCPAddress(address string) Options {
    return func(option *Option) {
        option.tcpAddress = address
    }
}

func newOpts(opts ...Options) Option {
    opt := Option{}

    for i := range opts {
        opts[i](&opt)
    }

    if opt.name == "" {
        opt.name = fs.GetProjectName()
    }

    if opt.id == "" {
        //todo change to git commit id+timestamp
        opt.id = x.NewUUID()
    }

    ip, err := x.LocalIp()
    if err != nil {
        panic(err)
    }

    if opt.randPort == nil {
        opt.randPort = &[2]int{10000, 19999}
    }

    if opt.tcpAddress == "" {
        opt.tcpAddress = ip + ":" + strconv.Itoa(x.RandInt(opt.randPort[0], opt.randPort[1]))
    }

    if opt.wssAddress != "" && opt.wssPath == "" {
        opt.wssPath = "/roc/wss"
    }

    if opt.err == nil {
        opt.err = parcel.DefaultErrorPacket
    }

    if opt.buffSize == 0 {
        opt.buffSize = 10
    }

    if opt.server == nil {
        opt.server = rs.NewServer(opt.tcpAddress, opt.wssAddress, opt.name, opt.buffSize)
    }

    endpoint.DefaultLocalEndpoint, err = endpoint.NewEndpoint(opt.id, opt.name, opt.tcpAddress)
    if err != nil {
        panic(err)
    }

    if opt.signal == nil {
        opt.signal = sig.DefaultSignal
    }

    if opt.apiPrefix == "" {
        opt.apiPrefix = "roc"
    }

    if !strings.HasPrefix(opt.apiPrefix, "/") {
        opt.apiPrefix = "/" + opt.apiPrefix
    }

    if !strings.HasSuffix(opt.apiPrefix, "/") {
        opt.apiPrefix = opt.apiPrefix + "/"
    }

    if opt.registry == nil {
        opt.registry = registry.DefaultRegistry
    }

    return opt
}
