package context

import (
	"fmt"
	"roc/internal/trace"
	"roc/internal/trace/simple"
	"roc/parcel/metadata"
	"roc/rlog/log"
)

type Context struct {
	*metadata.Metadata

	Trace trace.Trace `json:"trace"`
}

func Background() *Context {
	return new(Context)
}

func (c *Context) WithMetadata(service, method, tracing string, meta map[string]string) error {
	m, err := metadata.EncodeMetadata(service, method, tracing, meta)
	if err != nil {
		return err
	}
	c.Metadata = m
	c.Trace = simple.NewSimple(tracing)

	return nil
}

func NewConext(service, method, tracing string, meta map[string]string) (*Context, error) {
	m, err := metadata.EncodeMetadata(service, method, tracing, meta)
	if err != nil {
		return nil, err
	}
	return &Context{
		Metadata: m,
		Trace:    simple.NewSimple(tracing),
	}, nil
}

func FromMetadata(b []byte) *Context {
	m := metadata.DecodeMetadata(b)
	return &Context{
		Trace:    simple.NewSimple(m.Tracing()),
		Metadata: m,
	}
}

func (c *Context) Debug(msg ...interface{}) {
	c.Trace.Carrier()
	log.Debug(c.Trace.TraceId() + " |" + fmt.Sprintln(msg...))
}

func (c *Context) Info(msg ...interface{}) {
	c.Trace.Carrier()
	log.Info(c.Trace.TraceId() + " |" + fmt.Sprintln(msg...))
}

func (c *Context) Error(msg ...interface{}) {
	c.Trace.Carrier()
	log.Error(c.Trace.TraceId() + " |" + fmt.Sprintln(msg...))
}

func (c *Context) Debugf(f string, msg ...interface{}) {
	c.Trace.Carrier()
	log.Debug(c.Trace.TraceId() + " |" + fmt.Sprintf(f+"\n", msg...))
}

func (c *Context) Infof(f string, msg ...interface{}) {
	c.Trace.Carrier()
	log.Info(c.Trace.TraceId() + " |" + fmt.Sprintf(f+"\n", msg...))
}

func (c *Context) Errorf(f string, msg ...interface{}) {
	c.Trace.Carrier()
	log.Error(c.Trace.TraceId() + " |" + fmt.Sprintf(f+"\n", msg...))
}
