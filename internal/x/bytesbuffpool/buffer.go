package bytesbuffpool

import (
	"bytes"
	"sync"
)

var bytesBufferPool = sync.Pool{New: func() interface{} {
	return new(bytes.Buffer)
}}

func Put(b *bytes.Buffer) {
	b.Reset()
	bytesBufferPool.Put(b)
}

func Get() *bytes.Buffer {
	return bytesBufferPool.Get().(*bytes.Buffer)
}
