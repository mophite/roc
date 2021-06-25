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

package config

import (
    "github.com/go-roc/roc/etcd"
    "github.com/go-roc/roc/internal/namespace"
    "github.com/go-roc/roc/internal/x"
)

type Option struct {

    //etcd
    e *etcd.Etcd

    //disableDynamic switch
    disableDynamic bool

    //config schema on etcd
    //schema usually is global config dir
    schema string

    //private config on etcd
    private string

    //global config on etcd
    global string

    //global prefix
    prefix string

    //config version
    version string

    //backup file path
    backupPath string

    //backup file name
    backupName string

    //localFile switch
    //if true,will load local config.json to
    localFile bool
}

type Options func(option *Option)

func DisableDynamic() Options {
    return func(option *Option) {
        option.disableDynamic = true
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
        option.private = prefix
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

    opt.schema += "/" + opt.version

    if opt.backupPath == "" {
        opt.backupPath = "./"
    }

    if opt.backupName == "" {
        opt.backupName = "config.json"
    }

    opt.private = x.GetProjectName()

    opt.private = opt.schema + "/" + opt.private + "/"

    if opt.global == "" {
        opt.global = "global"
    }

    opt.global = opt.schema + "/" + opt.global + "/"

    opt.e = etcd.DefaultEtcd

    if opt.prefix == "" {
        opt.prefix = "roc"
    }

    return opt
}
