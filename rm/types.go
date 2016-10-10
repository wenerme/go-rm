package rm

import "unsafe"



// #include <stdlib.h>
// inline intptr_t PtrToInt(void* ptr){return (intptr_t)ptr;}
import (
    "fmt"
    "os"
    "github.com/wenerme/letsgo/cutil"
)

type Ctx uintptr
type CallReply uintptr
type String uintptr
type IO uintptr
type Key uintptr


type CmdFunc func(args CmdContext) int

type ZsetKey Key
type HashKey Key
type ListKey Key
type StringKey Key

func CreateString(ptr unsafe.Pointer) String {
    return String(cutil.PtrToIntptr(ptr))
}
func CreateCallReply(ptr unsafe.Pointer) CallReply {
    return CallReply(cutil.PtrToIntptr(ptr))
}

// ModuleType pattern [-_0-9A-Za-z]{9} suggest <typename>-<Vendor> not A{9}
//type CmdFunc func(ctx Ctx, args CmdArgs) int


type ModuleType struct {
    Name       string
    EncVer     string
    RdbLoad    func()
    RdbSave    func()
    AofRewrite func()
    Digest     func()
    Free       func()
}
type LogLevel int

const (
    LOG_DEBUG LogLevel = iota
    LOG_VERBOSE
    LOG_NOTICE
    LOG_WARNING
)


type CmdArgs struct {
    argv unsafe.Pointer
    argc int
}
type CmdContext struct {
    Ctx Ctx
}

func init() {
    LogDebug("Init Go Redis module")
}

var LogErr = func(format string, args... interface{}) {
    fmt.Fprintf(os.Stderr, format + "\n", args...)
}

var LogDebug = func(format string, args... interface{}) {
    fmt.Fprintf(os.Stdout, format + "\n", args...)
}

func (v String)ptr() unsafe.Pointer {
    return unsafe.Pointer(v)
}
func (v Ctx)ptr() unsafe.Pointer {
    return unsafe.Pointer(v)
}
func (v CallReply)ptr() unsafe.Pointer {
    return unsafe.Pointer(v)
}
func (v IO)ptr() unsafe.Pointer {
    return unsafe.Pointer(v)
}
func (v Key)ptr() unsafe.Pointer {
    return unsafe.Pointer(v)
}
