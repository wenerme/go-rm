package rm

import "unsafe"

type Ctx unsafe.Pointer
type Ctx unsafe.Pointer
type CmdFunc func(ctx Ctx, args CmdArgs) int

type CallReply unsafe.Pointer

type Key unsafe.Pointer
type ZsetKey Key
type HashKey Key
type ListKey Key
type StringKey Key

// RedisModuleString
type String unsafe.Pointer

// ModuleType pattern [-_0-9A-Za-z]{9} suggest <typename>-<Vendor> not A{9}
//type CmdFunc func(ctx Ctx, args CmdArgs) int

type IO unsafe.Pointer

type ModuleType struct {
    Name       string
    EncVer     string
    RdbLoad    func()
    RdbSave    func()
    AofRewrite func()
    Digest     func()
    Free       func()
}

type KeyType int

const (
    KeyType_Module KeyType = 0//TODO
)
