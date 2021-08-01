package say

import (
    "github.com/go-roc/roc/_auxiliary/example/tutorials/internal/ipc"
    "github.com/go-roc/roc/_auxiliary/example/tutorials/proto/phello"
    "github.com/go-roc/roc/parcel/context"
    "github.com/go-roc/roc/rlog"
)

type Say struct{}

// SayGet if err!=nil header status code will be 500
// bench shell
/*
wrk -t100 -c1000 -d10s --latency http://127.0.0.1:9999/roc/hello/sayget
Running 10s test @ http://127.0.0.1:9999/roc/hello/sayget
  100 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    13.45ms   12.17ms 248.27ms   80.96%
    Req/Sec   456.48    316.67     3.18k    81.65%
  Latency Distribution
     50%   11.25ms
     75%   17.81ms
     90%   24.66ms
     99%   62.24ms
  412262 requests in 10.10s, 47.18MB read
  Socket errors: connect 0, read 886, write 0, timeout 0
Requests/sec:  40821.21
Transfer/sec:      4.67MB
*/
func (h *Say) SayGet(c *context.Context, req *phello.ApiReq, rsp *phello.ApiRsp) (err error) {
    rsp.Code = 200
    rsp.Msg = "success"
    return nil
}

// Say http post+rpc
/*
wrk -t100 -c1000 -d10s --latency -s say.lua http://127.0.0.1:9999/roc/hello/say
Running 10s test @ http://127.0.0.1:9999/roc/hello/say
  100 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    14.45ms    5.65ms 102.40ms   65.41%
    Req/Sec   323.82    219.29     1.76k    76.50%
  Latency Distribution
     50%   14.18ms
     75%   18.09ms
     90%   22.07ms
     99%   27.86ms
  273227 requests in 10.10s, 32.06MB read
  Socket errors: connect 0, read 811, write 0, timeout 0
Requests/sec:  27041.90
Transfer/sec:      3.17MB

*/
func (h *Say) Say(c *context.Context, req *phello.SayReq, rsp *phello.SayRsp) error {
    //send a rr rpc request
    sayRsp, err := ipc.SaySrv(c, req)
    if err != nil {
        rlog.Error(err)
        rsp.Pong = "error"
        return err
    }

    rsp.Pong = sayRsp.Pong

    return nil
}
