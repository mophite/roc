package mock

import (
    "time"

    "context"
    "fmt"
    "net"
    "sync"

    "github.com/coreos/etcd/clientv3"
    pb "github.com/coreos/etcd/etcdserver/etcdserverpb"
    "github.com/coreos/etcd/mvcc/mvccpb"
    "github.com/go-roc/roc/internal/etcd"

    "google.golang.org/grpc"
)

// _ "github.com/go-roc/roc/internal/etcd/mock"
//to setup mock etcd server
var gMockServer *server

func init() {
    //start server
    var err error
    gMockServer, err = StartMockServer()
    if err != nil {
        panic(err)
    }

    //malloc etcd.DefaultEtcd
    err = etcd.NewEtcd(
        time.Second*5, 30, &clientv3.Config{
            Endpoints: []string{"localhost:2379"},
        },
    )

    if err != nil {
        panic(err)
    }
}

type server struct {
    mu         sync.RWMutex
    wg         sync.WaitGroup
    ln         net.Listener
    Address    string
    GrpcServer *grpc.Server
}

func StartMockServer() (ms *server, err error) {
    return startMockServer("localhost:2379")
}

func startMockServer(addr string) (ms *server, err error) {
    ln, err := net.Listen("tcp", addr)
    if err != nil {
        return nil, fmt.Errorf("failed to listen %v", err)
    }

    ms = &server{
        wg: sync.WaitGroup{},
        ln: ln, Address: ln.Addr().String(),
    }

    ms.start()
    return ms, nil
}

func (ms *server) start() {
    ms.mu.Lock()
    defer ms.mu.Unlock()

    svr := grpc.NewServer()
    pb.RegisterKVServer(svr, &mockKVServer{data: make(map[string]string)})
    ms.GrpcServer = svr

    ms.wg.Add(1)
    go func(svr *grpc.Server, l net.Listener) {
        svr.Serve(l)
    }(ms.GrpcServer, ms.ln)
}

func Stop() {
    gMockServer.mu.Lock()
    defer gMockServer.mu.Unlock()

    if gMockServer.ln == nil {
        return
    }

    gMockServer.GrpcServer.Stop()
    gMockServer.GrpcServer = nil
    gMockServer.ln = nil
    gMockServer.wg.Done()
}

type mockKVServer struct {
    sync.RWMutex
    data map[string]string
}

// Range range one to return
func (m *mockKVServer) Range(c context.Context, req *pb.RangeRequest) (*pb.RangeResponse, error) {
    m.RLock()
    defer m.RUnlock()

    var kv []*mvccpb.KeyValue
    v, ok := m.data[string(req.GetKey())]
    if ok {
        tmp := new(mvccpb.KeyValue)
        tmp.Key = req.Key
        tmp.Value = []byte(v)
        kv = append(kv, tmp)
    }
    return &pb.RangeResponse{
        Kvs:   kv,
        Count: 1,
    }, nil
}

// Put put one to cache data
func (m *mockKVServer) Put(c context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
    m.Lock()
    defer m.Unlock()

    pre := new(mvccpb.KeyValue)
    v, ok := m.data[string(req.Key)]
    if ok {
        pre.Key = req.Key
        pre.Value = []byte(v)
        pre.Lease = req.Lease
    }

    m.data[string(req.Key)] = string(req.Value)
    return &pb.PutResponse{PrevKv: pre}, nil
}

// DeleteRange No need to implement for mock
func (m *mockKVServer) DeleteRange(c context.Context, req *pb.DeleteRangeRequest) (*pb.DeleteRangeResponse, error) {
    return &pb.DeleteRangeResponse{}, nil
}

// Txn No need to implement
func (m *mockKVServer) Txn(c context.Context, req *pb.TxnRequest) (*pb.TxnResponse, error) {
    return &pb.TxnResponse{}, nil
}

// Compact No need to implement
func (m *mockKVServer) Compact(c context.Context, req *pb.CompactionRequest) (*pb.CompactionResponse, error) {
    return &pb.CompactionResponse{}, nil
}
