package rm

// #include <stdlib.h>
// inline intptr_t PtrToInt(void* ptr){return (intptr_t)ptr;}
import (
	"fmt"
	"github.com/wenerme/letsgo/cutil"
	"os"
	"unsafe"
)

type Ctx uintptr
type CallReply uintptr
type String uintptr
type Key uintptr
type IO uintptr
type Digest uintptr
type ModuleType uintptr

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
	Desc string
	// A 9 characters data type name that MUST be unique in the Redis
	// Modules ecosystem. Be creative... and there will be no collisions. Use
	// the charset A-Z a-z 9-0, plus the two "-_" characters. A good
	// idea is to use, for example `<typename>-<vendor>`. For example
	// "tree-AntZ" may mean "Tree data structure by @antirez". To use both
	// lower case and upper case letters helps in order to prevent collisions.
	//
	// Note: the module name "AAAAAAAAA" is reserved and produces an error, it
	// happens to be pretty lame as well.
	Name string
	// Encoding version, which is, the version of the serialization
	// that a module used in order to persist data. As long as the "name"
	// matches, the RDB loading will be dispatched to the type callbacks
	// whatever 'encver' is used, however the module can understand if
	// the encoding it must load are of an older version of the module.
	// For example the module "tree-AntZ" initially used encver=0. Later
	// after an upgrade, it started to serialize data in a different format
	// and to register the type with encver=1. However this module may
	// still load old data produced by an older version if the rdb_load
	// callback is able to check the encver value and act accordingly.
	// The encver must be a positive value between 0 and 1023.
	EncVer int
	// A callback function pointer that loads data from RDB files.
	RdbLoad func(rdb IO, encver int) unsafe.Pointer `json:"-"`
	// A callback function pointer that saves data to RDB files.
	RdbSave func(rdb IO, value unsafe.Pointer) `json:"-"`
	// A callback function pointer that rewrites data as commands.
	AofRewrite func(aof IO, key String, value unsafe.Pointer) `json:"-"`
	// A callback function pointer that is used for `DEBUG DIGEST`.
	Digest func(digest Digest, value unsafe.Pointer) `json:"-"`
	// A callback function pointer that can free a type value.
	Free func(value unsafe.Pointer) `json:"-"`
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
	Args []String
}

func init() {
	//LogDebug("Init Go Redis module")
}

var LogDebug = func(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}
var LogError = func(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func (v String) ptr() unsafe.Pointer {
	return unsafe.Pointer(v)
}
func (v Ctx) ptr() unsafe.Pointer {
	return unsafe.Pointer(v)
}
func (v CallReply) ptr() unsafe.Pointer {
	return unsafe.Pointer(v)
}
func (v IO) ptr() unsafe.Pointer {
	return unsafe.Pointer(v)
}
func (v Key) ptr() unsafe.Pointer {
	return unsafe.Pointer(v)
}
func (v ModuleType) ptr() unsafe.Pointer {
	return unsafe.Pointer(v)
}
func (v String) IsNull() bool {
	return uintptr(v) == 0
}
func (v Ctx) IsNull() bool {
	return uintptr(v) == 0
}
func (v CallReply) IsNull() bool {
	return uintptr(v) == 0
}
func (v IO) IsNull() bool {
	return uintptr(v) == 0
}
func (v Key) IsNull() bool {
	return uintptr(v) == 0
}
func (v ModuleType) IsNull() bool {
	return uintptr(v) == 0
}
