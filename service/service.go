package service

import (
	"net/http"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/gorilla/mux"

	"github.com/go-roc/roc/config"
	"github.com/go-roc/roc/internal/endpoint"
	"github.com/go-roc/roc/internal/etcd"
	"github.com/go-roc/roc/internal/namespace"
	"github.com/go-roc/roc/internal/registry"
	"github.com/go-roc/roc/parcel/codec"
	"github.com/go-roc/roc/rlog/log"
	"github.com/go-roc/roc/service/client"
	"github.com/go-roc/roc/service/server"
)

const SupportPackageIsVersion1 = 1

//DefaultApiPrefix it must be unique in all of your handler path
var DefaultApiPrefix = "/roc/"

type Option struct {

	//etcd config
	etcdConfig *clientv3.Config

	//config options
	configOpt []config.Options

	//roc service client,for rpc to server
	client *client.Client

	//roc service server,listen and wait call
	server *server.Server

	//http api router
	router *mux.Router

	//server tcp address
	tcpAddress string

	//server http address
	httpAddress string
}

type Options func(option *Option)

// EtcdConfig setting global etcd config first
func EtcdConfig(e *clientv3.Config) Options {
	return func(option *Option) {
		option.etcdConfig = e
	}
}

func TCPAddress(address string) Options {
	return func(option *Option) {
		option.tcpAddress = address
	}
}

func HttpAddress(address string) Options {
	return func(option *Option) {
		option.httpAddress = address
	}
}

func Codec(cc codec.Codec) Options {
	return func(option *Option) {
		//encode or decode method net packet
		codec.DefaultCodec = nil
		codec.DefaultCodec = cc
	}
}

func ConfigOption(opts ...config.Options) Options {
	return func(option *Option) {
		option.configOpt = opts
	}
}

func Version(version string) Options {
	return func(option *Option) {
		namespace.DefaultVersion = version
	}
}

func Registry(r registry.Registry) Options {
	return func(option *Option) {
		registry.DefaultRegistry = nil
		registry.DefaultRegistry = r
	}
}

func ApiPrefix(apiPrefix string) Options {
	return func(option *Option) {
		if !strings.HasPrefix(apiPrefix, "/") {
			apiPrefix = "/" + apiPrefix
		}
		if !strings.HasSuffix(apiPrefix, "/") {
			apiPrefix += "/"
		}
		DefaultApiPrefix = apiPrefix
	}
}

func Server(s *server.Server) Options {
	return func(option *Option) {
		option.server = s
	}
}

func Client(c *client.Client) Options {
	return func(option *Option) {
		option.client = c
	}
}

// Deprecated: newServer will created endpoint
func Endpoint(e *endpoint.Endpoint) Options {
	return func(option *Option) {
		endpoint.DefaultLocalEndpoint = e
	}
}

func newOpts(opts ...Options) Option {
	opt := Option{}

	for i := range opts {
		opts[i](&opt)
	}

	if opt.etcdConfig == nil {
		opt.etcdConfig = &clientv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: time.Second * 5,
		}
	}

	// init etcd.DefaultEtcd
	err := etcd.NewEtcd(time.Second*5, 5, opt.etcdConfig)
	if err != nil {
		panic("etcdConfig occur error: " + err.Error())
	}

	registry.DefaultRegistry = registry.NewRegistry()

	err = config.NewConfig(opt.configOpt...)
	if err != nil {
		panic("config NewConfig occur error: " + err.Error())
	}

	if opt.server == nil {
		opt.server = server.NewServer(
			server.TCPAddress(opt.tcpAddress),
			server.HttpAddress(opt.httpAddress),
		)
	}

	opt.router = mux.NewRouter()

	return opt
}

type Service struct {
	opts Option
}

func New(opts ...Options) *Service {
	return &Service{opts: newOpts(opts...)}
}

func (s *Service) Client() *client.Client {
	if s.opts.client == nil {
		s.opts.client = client.NewClient()
	}

	return s.opts.client
}

func (s *Service) Server() *server.Server {
	return s.opts.server
}

func (s *Service) Run() error {
	return s.opts.server.Run()
}

func (s *Service) Close() {
	if registry.DefaultRegistry != nil {
		_ = registry.DefaultRegistry.Deregister(endpoint.DefaultLocalEndpoint)
		registry.DefaultRegistry.Close()
		registry.DefaultRegistry = nil
	}

	if s.opts.client != nil {
		s.opts.client.Close()
	}

	if s.opts.server != nil {
		s.opts.server.Close()
	}

	etcd.DefaultEtcd.Close()
	registry.DefaultRegistry.Close()

	//todo flush rlog content
	log.Close()
}

var defaultRouter *mux.Router

// GROUP for not post http method
func (s *Service) GROUP(prefix string) *Service {
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	//cannot had post prefix
	if strings.HasPrefix(prefix, GetApiPrefix()) {
		panic("cannot contain unique prefix")
	}

	defaultRouter = s.opts.router.PathPrefix(prefix).Subrouter()
	return s
}

func (s *Service) GET(relativePath string, handler http.Handler) {

	relativePath = tidyRelativePath(relativePath)

	if strings.HasPrefix(relativePath, GetApiPrefix()) {
		panic("cannot contain unique prefix")
	}
	defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodOptions, http.MethodGet)
}

func (s *Service) POST(relativePath string, handler http.Handler) {

	relativePath = tidyRelativePath(relativePath)

	if strings.HasPrefix(relativePath, GetApiPrefix()) {
		panic("cannot contain unique prefix")
	}
	defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodPost)
}

func (s *Service) PUT(relativePath string, handler http.Handler) {

	relativePath = tidyRelativePath(relativePath)

	if strings.HasPrefix(relativePath, GetApiPrefix()) {
		panic("cannot contain unique prefix")
	}
	defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodPut)
}

func (s *Service) DELETE(relativePath string, handler http.Handler) {

	relativePath = tidyRelativePath(relativePath)

	if strings.HasPrefix(relativePath, GetApiPrefix()) {
		panic("cannot contain unique prefix")
	}
	defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodDelete)
}

func (s *Service) ANY(relativePath string, handler http.Handler) {

	relativePath = tidyRelativePath(relativePath)

	if strings.HasPrefix(relativePath, GetApiPrefix()) {
		panic("cannot contain unique prefix")
	}
	defaultRouter.PathPrefix(relativePath).Handler(handler)
}

func (s *Service) HEAD(relativePath string, handler http.Handler) {

	relativePath = tidyRelativePath(relativePath)

	if strings.HasPrefix(relativePath, GetApiPrefix()) {
		panic("cannot contain unique prefix")
	}
	defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodHead)
}

func (s *Service) PATCH(relativePath string, handler http.Handler) {

	relativePath = tidyRelativePath(relativePath)

	if strings.HasPrefix(relativePath, GetApiPrefix()) {
		panic("cannot contain unique prefix")
	}
	defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodPatch)
}

func (s *Service) CONNECT(relativePath string, handler http.Handler) {

	relativePath = tidyRelativePath(relativePath)

	if strings.HasPrefix(relativePath, GetApiPrefix()) {
		panic("cannot contain unique prefix")
	}
	defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodConnect)
}

func (s *Service) TRACE(relativePath string, handler http.Handler) {

	relativePath = tidyRelativePath(relativePath)

	if strings.HasPrefix(relativePath, GetApiPrefix()) {
		panic("cannot contain unique prefix")
	}
	defaultRouter.PathPrefix(relativePath).Handler(handler).Methods(http.MethodTrace)
}

func tidyRelativePath(relativePath string) string {
	//trim suffix "/"
	if strings.HasSuffix(relativePath, "/") {
		relativePath = strings.TrimSuffix(relativePath, "/")
	}

	//add prefix "/"
	if !strings.HasPrefix(relativePath, "/") {
		relativePath = "/" + relativePath
	}

	return relativePath
}

func GetApiPrefix() string {
	return DefaultApiPrefix
}
