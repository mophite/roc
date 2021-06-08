package im

import (
	"roc/_auxiliary/example/tutorials/proto/pbim"
	"sync"
)

func NewHub() *Hub {
	h := &Hub{
		lock:         new(sync.RWMutex),
		connectCount: 0,
		clients:      make(map[string]*point),
		broadCast:    make(chan *pbim.SendMessageReq),
	}
	go h.poller()
	return h
}

type Hub struct {
	lock         *sync.RWMutex
	connectCount uint32
	clients      map[string]*point
	broadCast    chan *pbim.SendMessageReq
}

type point struct {
	userName string
	message  chan *pbim.SendMessageRsp
}

func (h *Hub) count() uint32 {
	return h.connectCount
}

func (h *Hub) addClient(p *point) {
	if _, ok := h.clients[p.userName]; !ok {
		h.lock.RLock()
		h.clients[p.userName] = p
		h.connectCount += 1
		h.lock.RUnlock()
	}
}

func (h *Hub) removeClient(p *point) {
	if _, ok := h.clients[p.userName]; ok {
		h.lock.RLock()
		delete(h.clients, p.userName)
		h.connectCount -= 1
		h.lock.RUnlock()
	}
}

func (h *Hub) poller() {

	for {
		select {
		case b := <-h.broadCast:
			go func() {
				for userName, _ := range h.clients {
					h.clients[userName].message <- &pbim.SendMessageRsp{Message: b.Message}
				}
			}()

			// todo some thing
		}
	}
}
