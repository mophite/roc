package log

import (
	"roc/rlog/format"
	"roc/rlog/output"
)

type Option struct {
	call   int
	name   string
	prefix string
	format format.Formatter
	out    output.Outputor
}

type Options func(*Option)

func newOpts(opts ...Options) Option {
	opt := Option{}

	for i := range opts {
		opts[i](&opt)
	}

	if opt.format == nil {
		opt.format = format.DefaultFormat
	}

	if opt.name == "" {
		opt.name = "roc"
	}

	if opt.out == nil {
		opt.out = output.DefaultOutput
	}

	if opt.prefix == "" {
		opt.prefix = ""
	}

	if opt.call == 0 {
		opt.call = 4
	}

	return opt
}

func Call(call int) Options {
	return func(option *Option) {
		option.call = call
	}
}

func Output(out output.Outputor) Options {
	return func(option *Option) {
		option.out = out
	}
}

func Name(name string) Options {
	return func(option *Option) {
		option.name = name
	}
}

func Prefix(prefix string) Options {
	return func(option *Option) {
		option.prefix = prefix
	}
}

func Format(format format.Formatter) Options {
	return func(option *Option) {
		option.format = format
	}
}
