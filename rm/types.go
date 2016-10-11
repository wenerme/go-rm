package rm

// #include <stdlib.h>
// inline intptr_t PtrToInt(void* ptr){return (intptr_t)ptr;}
import (
    "fmt"
    "os"
    "unsafe"
    "github.com/wenerme/letsgo/cutil"
)

type Ctx uintptr
type CallReply uintptr
type String uintptr
type Key uintptr
type IO uintptr
type Digest uintptr

type CmdFunc func(args CmdContext) int

type ZsetKey Key
type HashKey Key
type ListKey Key
type StringKey Key

func CreateString(ptr unsafe.Pointer) String {
    return String(cutil.PtrToUintptr(ptr))
}
func CreateCallReply(ptr unsafe.Pointer) CallReply {
    return CallReply(cutil.PtrToUintptr(ptr))
}
func NullString() String {
    return CreateString(NullPointer())
}
func NullPointer() unsafe.Pointer {
    return unsafe.Pointer(uintptr(0))
}

type DataType struct {
    // Must match [-_0-9A-Za-z]{9} suggest <typename>-<Vendor> not A{9}
    Name       string
    EncVer     int
    Desc       string
    RdbLoad    func(rdb IO, encver int) unsafe.Pointer `json:"-"`
    RdbSave    func(rdb IO, value unsafe.Pointer) `json:"-"`
    AofRewrite func(aof IO, key String, value unsafe.Pointer) `json:"-"`
    Digest     func(digest Digest, value unsafe.Pointer) `json:"-"`
    Free       func(value unsafe.Pointer) `json:"-"`
}
type LogLevel int

const (
    LOG_DEBUG LogLevel = iota
    LOG_VERBOSE
    LOG_NOTICE
    LOG_WARNING
)

type CmdContext struct {
    Ctx  Ctx
    Args [] String
}

func init() {
    //LogDebug("Init Go Redis module")
}

var LogDebug = func(format string, args... interface{}) {
    fmt.Fprintf(os.Stdout, format + "\n", args...)
}
var LogError = func(format string, args... interface{}) {
    fmt.Fprintf(os.Stderr, format + "\n", args...)
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
func (v String)IsNull() bool {
    return uintptr(v) == 0
}
func (v Ctx)IsNull() bool {
    return uintptr(v) == 0
}
func (v CallReply)IsNull() bool {
    return uintptr(v) == 0
}
func (v IO)IsNull() bool {
    return uintptr(v) == 0
}
func (v Key)IsNull() bool {
    return uintptr(v) == 0
}
