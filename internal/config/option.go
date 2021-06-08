package config

import (
	"roc/internal/etcd"
	"roc/internal/namespace"
	"roc/internal/x"
)

type Option struct {
	e             *etcd.Etcd
	enableDynamic bool
	schema        string
	prefix        string
	version       string
	backupPath    string
}

type Options func(option *Option)

func EnableDynamic() Options {
	return func(option *Option) {
		option.enableDynamic = true
	}
}

func Schema(schema string) Options {
	return func(option *Option) {
		option.schema = schema
	}
}

func Version(version string) Options {
	return func(option *Option) {
		option.version = version
	}
}

func Backup(path string) Options {
	return func(option *Option) {
		option.backupPath = path
	}
}

func Prefix(prefix string) Options {
	return func(option *Option) {
		option.prefix = prefix
	}
}

func newOpts(opts ...Options) Option {
	opt := Option{}

	for i := range opts {
		opts[i](&opt)
	}

	if opt.schema == "" {
		opt.schema = namespace.DefaultConfigSchema
	}

	if opt.version == "" {
		opt.version = namespace.DefaultVersion
	}

	opt.schema += "." + opt.version

	if opt.backupPath == "" {
		opt.backupPath = "./"
	}

	if opt.prefix == "" {
		opt.prefix = x.GetProjectName()
	}

	return opt
}
