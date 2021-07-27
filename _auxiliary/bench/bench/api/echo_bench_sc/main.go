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

package main

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/jjeffcaii/reactor-go/scheduler"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/core/transport"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
)

var tp transport.ClientTransporter

func init() {
	tp = rsocket.TCPClient().SetHostAndPort("127.0.0.1", 7878).Build()
}

func main() {
	const (
		parallel = 1
		round    = 100
		amount   = 100000
	)
	var clients []rsocket.Client
	for i := 0; i < parallel; i++ {
		client, err := createClient()
		if err != nil {
			return
		}
		clients = append(clients, client)
	}
	defer func() {
		for _, c := range clients {
			_ = c.Close()
		}
	}()
	if len(clients) == 0 {
		return
	}
	time.Sleep(time.Second * 1)

	var tps int64
	var next uint64

	done := make(chan struct{})

	cnt := int64(round * amount)

	// 创建订阅
	sub := rx.NewSubscriber(
		rx.OnComplete(func() {
			atomic.AddInt64(&tps, 1)
			if atomic.AddInt64(&cnt, -1) < 1 {
				close(done)
			}
		}),
		rx.OnError(func(e error) {
			if atomic.AddInt64(&cnt, -1) < 1 {
				close(done)
			}
		}),
	)

	for range [round]struct{}{} {
		go func(c rsocket.Client) {
			for range [amount]struct{}{} {
				// 这里尽量避免用Block, 直接订阅response即可, 订阅默认就是异步的, 所以不需要额外创建协程去执行。
				//
				// 原本的Block会导致一次Copy, 如果场景考虑性能又要需要阻塞获取response, 可以使用BlockUnsafe:
				//
				// go func() {
				// 	  res, closer, err := c.RequestResponse(payload.NewString("hello", "")).BlockUnsafe(context.Background())
				// 	  if err != nil {
				// 	  	  return
				// 	  }
				//    // 这里需要确保最后调用closer, 此方法会将response所占用的资源释放(response的payload为了减少GC是pooled的)。
				// 	  defer closer()
				// 	  log.Println("response:", res)
				// }()
				//

				c.RequestResponse(payload.NewString("hello", "")).SubscribeWith(context.Background(), sub)
			}
		}(clients[(int(atomic.AddUint64(&next, 1))-1)%len(clients)])
	}

	for {
		select {
		case <-done:
			return
		case <-time.After(time.Second * 1):
			fmt.Println("-----------------------tps-------------------", atomic.LoadInt64(&tps))
			atomic.SwapInt64(&tps, 0)
		}
	}
}

func createClient() (rsocket.Client, error) {
	return rsocket.Connect().
		Scheduler(scheduler.Elastic(), nil). // 这里设置发送请求时调度器为动态无界协程池
		OnClose(func(err error) {
			log.Println("*** disconnected ***")
		}).
		SetupPayload(payload.NewString("你好", "世界")).
		Transport(tp).
		Start(context.Background())
}
