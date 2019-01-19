package rm

//#include "./rm.h"
import "C"
import (
	"fmt"
	"github.com/wenerme/letsgo/cutil"
	"regexp"
	"syscall"
	"unsafe"
)

/* ---------------- Defines common between core and modules --------------- */

/* Error status return values. */
const OK = C.REDISMODULE_OK
const ERR = C.REDISMODULE_ERR

/* API versions. */
const APIVER_1 = C.REDISMODULE_APIVER_1

/* API flags and constants */
const READ = C.REDISMODULE_READ
const WRITE = C.REDISMODULE_WRITE

const LIST_HEAD = C.REDISMODULE_LIST_HEAD
const LIST_TAIL = C.REDISMODULE_LIST_TAIL

/* Key types. */
const (
	// Return the type of the key. If the key pointer is NULL then `REDISMODULE_KEYTYPE_EMPTY` is returned.
	KEYTYPE_EMPTY = iota
	KEYTYPE_STRING
	KEYTYPE_LIST
	KEYTYPE_HASH
	KEYTYPE_SET
	KEYTYPE_ZSET
	KEYTYPE_MODULE
)

/* Reply types. */
const (
	REPLY_UNKNOWN = iota - 1
	REPLY_STRING
	REPLY_ERROR
	REPLY_INTEGER
	REPLY_ARRAY
	REPLY_NULL
)

/* Postponed array length. */
const POSTPONED_ARRAY_LEN = C.REDISMODULE_POSTPONED_ARRAY_LEN

/* Expire */
const NO_EXPIRE = C.REDISMODULE_NO_EXPIRE

/* Sorted set API flags. */
const (
	ZADD_XX = 1 << iota
	ZADD_NX
	ZADD_ADDED
	ZADD_UPDATED
	ZADD_NOP
)

/* Hash API flags. */
const (
	HASH_NONE = 0
	// Set if non-exists
	HASH_NX = 1 << iota
	// Set if exists
	HASH_XX
	// Use *char as args, ! do not use this flag
	HASH_CFIELDS
	// Check field exists
	HASH_EXISTS
)

/* Error messages. */
//const ERRORMSG_WRONGTYPE = C.REDISMODULE_ERRORMSG_WRONGTYPE
const ERRORMSG_WRONGTYPE = "WRONGTYPE Operation against a key holding the wrong kind of value"

//const POSITIVE_INFINITE = C.REDISMODULE_POSITIVE_INFINITE
//const NEGATIVE_INFINITE = C.REDISMODULE_NEGATIVE_INFINITE

func getErrno() syscall.Errno {
	return syscall.Errno(C.get_errno())
}

/* ------------------------- End of common defines ------------------------ */

// Use like malloc(). Memory allocated with this function is reported in
// Redis INFO memory, used for keys eviction according to maxmemory settings
// and in general is taken into account as memory allocated by Redis.
// You should avoid to use malloc().
// void *RM_Alloc(size_t bytes);
func Alloc(bytes int) unsafe.Pointer {
	return unsafe.Pointer(C.Alloc(C.size_t(bytes)))
}

// Use like realloc() for memory obtained with `RedisModule_Alloc()`.
// void* RM_Realloc(void *ptr, size_t bytes);
func Realloc(ptr unsafe.Pointer, bytes int) unsafe.Pointer {
	return unsafe.Pointer(C.Realloc(ptr, C.size_t(bytes)))
}

// Use like free() for memory obtained by `RedisModule_Alloc()` and
// `RedisModule_Realloc()`. However you should never try to free with
// `RedisModule_Free()` memory allocated with malloc() inside your module.
// void RM_Free(void *ptr);
func Free(ptr unsafe.Pointer) {
	C.Free(ptr)
}

// Like strdup() but returns memory allocated with `RedisModule_Alloc()`.
// char *RM_Strdup(const char *str);
func Strdup(str unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.Strdup((*C.char)(str)))
}

// Lookup the requested module API and store the function pointer into the
// target pointer. The function returns `REDISMODULE_ERR` if there is no such
// named API, otherwise `REDISMODULE_OK`.
//
// This function is not meant to be used by modules developer, it is only
// used implicitly by including redismodule.h.
// int RM_GetApi(const char *funcname, void **targetPtrPtr);
//func GetApi(funcname string,targetPtrPtr unsafe.Pointer)(int){return int(C.GetApi(funcname,targetPtrPtr))}

// =============================================================================
// ========================== Context functions
// =============================================================================

// int `RedisModule_OnLoad(RedisModuleCtx` *ctx) {
//          // some code here ...
//          BalancedTreeType = `RM_CreateDataType(`...);
//      }
// moduleType *RM_CreateDataType(RedisModuleCtx *ctx, const char *name, int encver, moduleTypeLoadFunc rdb_load, moduleTypeSaveFunc rdb_save, moduleTypeRewriteFunc aof_rewrite, moduleTypeDigestFunc digest, moduleTypeFreeFunc free);
// NOTE
//func (ctx Ctx)CreateDataType(name string,encver int,rdb_load RedisModuleTypeLoadFunc,rdb_save RedisModuleTypeSaveFunc,aof_rewrite RedisModuleTypeRewriteFunc,digest RedisModuleTypeDigestFunc,free RedisModuleTypeFreeFunc)(/* TODO RedisModuleType* */unsafe.Pointer){return /* TODO RedisModuleType* */unsafe.Pointer(C.CreateDataType((*C.struct_RedisModuleCtx)(ctx.ptr()),name,encver,rdb_load,rdb_save,aof_rewrite,digest,free))}

// Return heap allocated memory that will be freed automatically when the
// module callback function returns. Mostly suitable for small allocations
// that are short living and must be released when the callback returns
// anyway. The returned memory is aligned to the architecture word size
// if at least word size bytes are requested, otherwise it is just
// aligned to the next power of two, so for example a 3 bytes request is
// 4 bytes aligned while a 2 bytes request is 2 bytes aligned.
//
// There is no realloc style function since when this is needed to use the
// pool allocator is not a good idea.
//
// The function returns NULL if `bytes` is 0.
// void *RM_PoolAlloc(RedisModuleCtx *ctx, size_t bytes);
func (ctx Ctx) PoolAlloc(bytes int) unsafe.Pointer {
	return unsafe.Pointer(C.PoolAlloc((*C.struct_RedisModuleCtx)((*C.struct_RedisModuleCtx)(ctx.ptr())), C.size_t(bytes)))
}

// Return non-zero if a module command, that was declared with the
// flag "getkeys-api", is called in a special way to get the keys positions
// and not to get executed. Otherwise zero is returned.
// int RM_IsKeysPositionRequest(RedisModuleCtx *ctx);
func (ctx Ctx) IsKeysPositionRequest() int {
	return int(C.IsKeysPositionRequest((*C.struct_RedisModuleCtx)(ctx.ptr())))
}

// When a module command is called in order to obtain the position of
// keys, since it was flagged as "getkeys-api" during the registration,
// the command implementation checks for this special call using the
// `RedisModule_IsKeysPositionRequest()` API and uses this function in
// order to report keys, like in the following example:
//
//  if (`RedisModule_IsKeysPositionRequest(ctx))` {
//      `RedisModule_KeyAtPos(ctx`,1);
//      `RedisModule_KeyAtPos(ctx`,2);
//  }
//
//  Note: in the example below the get keys API would not be needed since
//  keys are at fixed positions. This interface is only used for commands
//  with a more complex structure.
// void RM_KeyAtPos(RedisModuleCtx *ctx, int pos);
func (ctx Ctx) KeyAtPos(pos int) {
	C.KeyAtPos((*C.struct_RedisModuleCtx)(ctx.ptr()), C.int(pos))
}

// And is supposed to always return `REDISMODULE_OK`.
//
// The set of flags 'strflags' specify the behavior of the command, and should
// be passed as a C string compoesd of space separated words, like for
// example "write deny-oom". The set of flags are:
//
// * **"write"**:     The command may modify the data set (it may also read
//                    from it).
// * **"readonly"**:  The command returns data from keys but never writes.
// * **"admin"**:     The command is an administrative command (may change
//                    replication or perform similar tasks).
// * **"deny-oom"**:  The command may use additional memory and should be
//                    denied during out of memory conditions.
// * **"deny-script"**:   Don't allow this command in Lua scripts.
// * **"allow-loading"**: Allow this command while the server is loading data.
//                        Only commands not interacting with the data set
//                        should be allowed to run in this mode. If not sure
//                        don't use this flag.
// * **"pubsub"**:    The command publishes things on Pub/Sub channels.
// * **"random"**:    The command may have different outputs even starting
//                    from the same input arguments and key values.
// * **"allow-stale"**: The command is allowed to run on slaves that don't
//                      serve stale data. Don't use if you don't know what
//                      this means.
// * **"no-monitor"**: Don't propoagate the command on monitor. Use this if
//                     the command has sensible data among the arguments.
// * **"fast"**:      The command time complexity is not greater
//                    than O(log(N)) where N is the size of the collection or
//                    anything else representing the normal scalability
//                    issue with the command.
// * **"getkeys-api"**: The command implements the interface to return
//                      the arguments that are keys. Used when start/stop/step
//                      is not enough because of the command syntax.
// * **"no-cluster"**: The command should not register in Redis Cluster
//                     since is not designed to work with it because, for
//                     example, is unable to report the position of the
//                     keys, programmatically creates key names, or any
//                     other reason.
// int RM_CreateCommand(RedisModuleCtx *ctx, const char *name, RedisModuleCmdFunc cmdfunc, const char *strflags, int firstkey, int lastkey, int keystep);
//func (ctx Ctx)CreateCommand(name string, cmdfunc CmdFunc, strflags string, firstkey int, lastkey int, keystep int) (int) {
//    return int(C.CreateCommand((*C.struct_RedisModuleCtx)(ctx.ptr()), name, cmdfunc, strflags, firstkey, lastkey, keystep))
//}

// Called by `RM_Init()` to setup the `ctx->module` structure.
//
// This is an internal function, Redis modules developers don't need
// to use it.
// void RM_SetModuleAttribs(RedisModuleCtx *ctx, const char *name, int ver, int apiver);
func (ctx Ctx) SetModuleAttribs(name string, ver int, apiver int) {
	C.SetModuleAttribs((*C.struct_RedisModuleCtx)(ctx.ptr()), C.CString(name), C.int(ver), C.int(apiver))
}

// Enable automatic memory management. See API.md for more information.
//
// The function must be called as the first function of a command implementation
// that wants to use automatic memory.
// void RM_AutoMemory(RedisModuleCtx *ctx);
func (ctx Ctx) AutoMemory() {
	C.AutoMemory((*C.struct_RedisModuleCtx)(ctx.ptr()))
}

// Create a new module string object. The returned string must be freed
// with `RedisModule_FreeString()`, unless automatic memory is enabled.
//
// The string is created by copying the `len` bytes starting
// at `ptr`. No reference is retained to the passed buffer.
// RedisModuleString *RM_CreateString(RedisModuleCtx *ctx, const char *ptr, size_t len);
func (ctx Ctx) CreateString(ptr string, len int) String {
	c := C.CString(ptr)
	defer C.free(unsafe.Pointer(c))
	return CreateString(unsafe.Pointer(C.CreateString((*C.struct_RedisModuleCtx)(ctx.ptr()), c, C.size_t(len))))
}

// Like `RedisModule_CreatString()`, but creates a string starting from a long long
// integer instead of taking a buffer and its length.
//
// The returned string must be released with `RedisModule_FreeString()` or by
// enabling automatic memory management.
// RedisModuleString *RM_CreateStringFromLongLong(RedisModuleCtx *ctx, long long ll);
func (ctx Ctx) CreateStringFromLongLong(ll int64) String {
	return CreateString(unsafe.Pointer(C.CreateStringFromLongLong((*C.struct_RedisModuleCtx)(ctx.ptr()), C.longlong(ll))))
}

// Like `RedisModule_CreatString()`, but creates a string starting from an existing
// RedisModuleString.
//
// The returned string must be released with `RedisModule_FreeString()` or by
// enabling automatic memory management.
// RedisModuleString *RM_CreateStringFromString(RedisModuleCtx *ctx, const RedisModuleString *str);
func (ctx Ctx) CreateStringFromString(str String) String {
	return CreateString(unsafe.Pointer(C.CreateStringFromString((*C.struct_RedisModuleCtx)(ctx.ptr()), (*C.struct_RedisModuleString)(str.ptr()))))
}

// Free a module string object obtained with one of the Redis modules API calls
// that return new string objects.
//
// It is possible to call this function even when automatic memory management
// is enabled. In that case the string will be released ASAP and removed
// from the pool of string to release at the end.
// void RM_FreeString(RedisModuleCtx *ctx, RedisModuleString *str);
func (ctx Ctx) FreeString(str String) {
	C.FreeString((*C.struct_RedisModuleCtx)(ctx.ptr()), (*C.struct_RedisModuleString)(str.ptr()))
}

//
// int RM_WrongArity(RedisModuleCtx *ctx);
func (ctx Ctx) WrongArity() int {
	return int(C.WrongArity((*C.struct_RedisModuleCtx)(ctx.ptr())))
}

// Send an integer reply to the client, with the specified long long value.
// The function always returns `REDISMODULE_OK`.
// int RM_ReplyWithLongLong(RedisModuleCtx *ctx, long long ll);
func (ctx Ctx) ReplyWithLongLong(ll int64) int {
	return int(C.ReplyWithLongLong((*C.struct_RedisModuleCtx)(ctx.ptr()), C.longlong(ll)))
}

// The function always returns `REDISMODULE_OK`.
// int RM_ReplyWithError(RedisModuleCtx *ctx, const char *err);
func (ctx Ctx) ReplyWithError(err string) int {
	//if Mod.Debug {
	// TODO Check has Error code like ERR
	//}
	c := C.CString(err)
	defer C.free(unsafe.Pointer(c))
	return int(C.ReplyWithError((*C.struct_RedisModuleCtx)(ctx.ptr()), c))
}

// Reply with a simple string (+... \r\n in RESP protocol). This replies
// are suitable only when sending a small non-binary string with small
// overhead, like "OK" or similar replies.
//
// The function always returns `REDISMODULE_OK`.
// int RM_ReplyWithSimpleString(RedisModuleCtx *ctx, const char *msg);
func (ctx Ctx) ReplyWithSimpleString(msg string) int {
	c := C.CString(msg)
	defer C.free(unsafe.Pointer(c))
	return int(C.ReplyWithSimpleString((*C.struct_RedisModuleCtx)(ctx.ptr()), c))
}

// Reply with an array type of 'len' elements. However 'len' other calls
// to `ReplyWith*` style functions must follow in order to emit the elements
// of the array.
//
// When producing arrays with a number of element that is not known beforehand
// the function can be called with the special count
// `REDISMODULE_POSTPONED_ARRAY_LEN`, and the actual number of elements can be
// later set with `RedisModule_ReplySetArrayLength()` (which will set the
// latest "open" count if there are multiple ones).
//
// The function always returns `REDISMODULE_OK`.
// int RM_ReplyWithArray(RedisModuleCtx *ctx, long len);
func (ctx Ctx) ReplyWithArray(len int64) int {
	return int(C.ReplyWithArray((*C.struct_RedisModuleCtx)(ctx.ptr()), C.long(len)))
}

// When `RedisModule_ReplyWithArray()` is used with the argument
// `REDISMODULE_POSTPONED_ARRAY_LEN`, because we don't know beforehand the number
// of items we are going to output as elements of the array, this function
// will take care to set the array length.
//
// Since it is possible to have multiple array replies pending with unknown
// length, this function guarantees to always set the latest array length
// that was created in a postponed way.
//
// For example in order to output an array like [1,[10,20,30]] we
// could write:
//
//  `RedisModule_ReplyWithArray(ctx`,`REDISMODULE_POSTPONED_ARRAY_LEN`);
//  `RedisModule_ReplyWithLongLong(ctx`,1);
//  `RedisModule_ReplyWithArray(ctx`,`REDISMODULE_POSTPONED_ARRAY_LEN`);
//  `RedisModule_ReplyWithLongLong(ctx`,10);
//  `RedisModule_ReplyWithLongLong(ctx`,20);
//  `RedisModule_ReplyWithLongLong(ctx`,30);
//  `RedisModule_ReplySetArrayLength(ctx`,3); // Set len of 10,20,30 array.
//  `RedisModule_ReplySetArrayLength(ctx`,2); // Set len of top array
//
// Note that in the above example there is no reason to postpone the array
// length, since we produce a fixed number of elements, but in the practice
// the code may use an interator or other ways of creating the output so
// that is not easy to calculate in advance the number of elements.
// void RM_ReplySetArrayLength(RedisModuleCtx *ctx, long len);
func (ctx Ctx) ReplySetArrayLength(len int64) {
	C.ReplySetArrayLength((*C.struct_RedisModuleCtx)(ctx.ptr()), C.long(len))
}

// Reply with a bulk string, taking in input a C buffer pointer and length.
//
// The function always returns `REDISMODULE_OK`.
// int RM_ReplyWithStringBuffer(RedisModuleCtx *ctx, const char *buf, size_t len);
func (ctx Ctx) ReplyWithStringBuffer(buf []byte, len int) int {
	// TODO free
	return int(C.ReplyWithStringBuffer((*C.struct_RedisModuleCtx)(ctx.ptr()), (*C.char)(C.CBytes(buf)), C.size_t(len)))
}

// Reply with a bulk string, taking in input a RedisModuleString object.
//
// The function always returns `REDISMODULE_OK`.
// int RM_ReplyWithString(RedisModuleCtx *ctx, RedisModuleString *str);
func (ctx Ctx) ReplyWithString(str String) int {
	return int(C.ReplyWithString((*C.struct_RedisModuleCtx)(ctx.ptr()), (*C.struct_RedisModuleString)(str.ptr())))
}

// Reply to the client with a NULL. In the RESP protocol a NULL is encoded
// as the string "$-1\r\n".
//
// The function always returns `REDISMODULE_OK`.
// int RM_ReplyWithNull(RedisModuleCtx *ctx);
func (ctx Ctx) ReplyWithNull() int {
	return int(C.ReplyWithNull((*C.struct_RedisModuleCtx)(ctx.ptr())))
}

// Reply exactly what a Redis command returned us with `RedisModule_Call()`.
// This function is useful when we use `RedisModule_Call()` in order to
// execute some command, as we want to reply to the client exactly the
// same reply we obtained by the command.
//
// The function always returns `REDISMODULE_OK`.
// int RM_ReplyWithCallReply(RedisModuleCtx *ctx, RedisModuleCallReply *reply);
func (ctx Ctx) ReplyWithCallReply(reply CallReply) int {
	return int(C.ReplyWithCallReply((*C.struct_RedisModuleCtx)(ctx.ptr()), (*C.struct_RedisModuleCallReply)(reply.ptr())))
}

// Send a string reply obtained converting the double 'd' into a bulk string.
// This function is basically equivalent to converting a double into
// a string into a C buffer, and then calling the function
// `RedisModule_ReplyWithStringBuffer()` with the buffer and length.
//
// The function always returns `REDISMODULE_OK`.
// int RM_ReplyWithDouble(RedisModuleCtx *ctx, double d);
func (ctx Ctx) ReplyWithDouble(d float64) int {
	return int(C.ReplyWithDouble((*C.struct_RedisModuleCtx)(ctx.ptr()), C.double(d)))
}

// Replicate the specified command and arguments to slaves and AOF, as effect
// of execution of the calling command implementation.
//
// The replicated commands are always wrapped into the MULTI/EXEC that
// contains all the commands replicated in a given module command
// execution. However the commands replicated with `RedisModule_Call()`
// are the first items, the ones replicated with `RedisModule_Replicate()`
// will all follow before the EXEC.
//
// Modules should try to use one interface or the other.
//
// This command follows exactly the same interface of `RedisModule_Call()`,
// so a set of format specifiers must be passed, followed by arguments
// matching the provided format specifiers.
//
// Please refer to `RedisModule_Call()` for more information.
//
// The command returns `REDISMODULE_ERR` if the format specifiers are invalid
// or the command name does not belong to a known command.
// int RM_Replicate(RedisModuleCtx *ctx, const char *cmdname, const char *fmt, ...);
func (ctx Ctx) Replicate(cmdname string, format string, args ...interface{}) int {
	c := C.CString(cmdname)
	defer C.free(unsafe.Pointer(c))
	msg := fmt.Sprintf(format, args...)
	s := C.CString(msg)
	defer C.free(unsafe.Pointer(s))
	return int(C.Replicate((*C.struct_RedisModuleCtx)(ctx.ptr()), c, s))
}

// This function will replicate the command exactly as it was invoked
// by the client. Note that this function will not wrap the command into
// a MULTI/EXEC stanza, so it should not be mixed with other replication
// commands.
//
// Basically this form of replication is useful when you want to propagate
// the command to the slaves and AOF file exactly as it was called, since
// the command can just be re-executed to deterministically re-create the
// new state starting from the old one.
//
// The function always returns `REDISMODULE_OK`.
// int RM_ReplicateVerbatim(RedisModuleCtx *ctx);
func (ctx Ctx) ReplicateVerbatim() int {
	return int(C.ReplicateVerbatim((*C.struct_RedisModuleCtx)(ctx.ptr())))
}

// Return the ID of the current client calling the currently active module
// command. The returned ID has a few guarantees:
//
// 1. The ID is different for each different client, so if the same client
//    executes a module command multiple times, it can be recognized as
//    having the same ID, otherwise the ID will be different.
// 2. The ID increases monotonically. Clients connecting to the server later
//    are guaranteed to get IDs greater than any past ID previously seen.
//
// Valid IDs are from 1 to 2^64-1. If 0 is returned it means there is no way
// to fetch the ID in the context the function was currently called.
// unsigned long long RM_GetClientId(RedisModuleCtx *ctx);
func (ctx Ctx) GetClientId() uint64 {
	return uint64(C.GetClientId((*C.struct_RedisModuleCtx)(ctx.ptr())))
}

// Return the currently selected DB.
// int RM_GetSelectedDb(RedisModuleCtx *ctx);
func (ctx Ctx) GetSelectedDb() int {
	return int(C.GetSelectedDb((*C.struct_RedisModuleCtx)(ctx.ptr())))
}

// Change the currently selected DB. Returns an error if the id
// is out of range.
//
// Note that the client will retain the currently selected DB even after
// the Redis command implemented by the module calling this function
// returns.
//
// If the module command wishes to change something in a different DB and
// returns back to the original one, it should call `RedisModule_GetSelectedDb()`
// before in order to restore the old DB number before returning.
// int RM_SelectDb(RedisModuleCtx *ctx, int newid);
func (ctx Ctx) SelectDb(newid int) int {
	return int(C.SelectDb((*C.struct_RedisModuleCtx)(ctx.ptr()), C.int(newid)))
}

// Return an handle representing a Redis key, so that it is possible
// to call other APIs with the key handle as argument to perform
// operations on the key.
//
// The return value is the handle repesenting the key, that must be
// closed with `RM_CloseKey()`.
//
// If the key does not exist and WRITE mode is requested, the handle
// is still returned, since it is possible to perform operations on
// a yet not existing key (that will be created, for example, after
// a list push operation). If the mode is just READ instead, and the
// key does not exist, NULL is returned. However it is still safe to
// call `RedisModule_CloseKey()` and `RedisModule_KeyType()` on a NULL
// value.
// void *RM_OpenKey(RedisModuleCtx *ctx, robj *keyname, int mode);
func (ctx Ctx) OpenKey(keyname String, mode int) Key {
	return Key(C.OpenKey((*C.struct_RedisModuleCtx)(ctx.ptr()), (*C.struct_RedisModuleString)(keyname.ptr()), C.int(mode)))
}

// Exported API to call any Redis command from modules.
// On success a RedisModuleCallReply object is returned, otherwise
// NULL is returned and errno is set to the following values:
//
// EINVAL: command non existing, wrong arity, wrong format specifier.
// EPERM:  operation in Cluster instance with key in non local slot.
// RedisModuleCallReply *RM_Call(RedisModuleCtx *ctx, const char *cmdname, const char *fmt, ...);
func (ctx Ctx) Call(cmdname string, format string, args ...interface{}) CallReply {
	c := C.CString(cmdname)
	defer C.free(unsafe.Pointer(c))

	f := C.CString(format)
	defer C.free(unsafe.Pointer(f))

	args = append(args, uintptr(0))
	p, err := cutil.VarArgsPtr(args...)
	defer C.free(p)
	if err != nil {
		LogError("Call failed: %v", err)
		return CreateCallReplyError(syscall.EINVAL)
	}
	return CreateCallReply(unsafe.Pointer(C.CallVar(
		(*C.struct_RedisModuleCtx)(ctx.ptr()),
		c,
		f,
		C.int(len(args)),
		(*C.intptr_t)(p))))
}

// Produces a log message to the standard Redis log, the format accepts
// printf-alike specifiers, while level is a string describing the log
// level to use when emitting the log, and must be one of the following:
//
// * "debug"
// * "verbose"
// * "notice"
// * "warning"
//
// If the specified log level is invalid, verbose is used by default.
// There is a fixed limit to the length of the log line this function is able
// to emit, this limti is not specified but is guaranteed to be more than
// a few lines of text.
// void RM_Log(RedisModuleCtx *ctx, const char *levelstr, const char *fmt, ...);
func (ctx Ctx) Log(l LogLevel, format string, args ...interface{}) {
	c := C.CString(fmt.Sprintf(format, args...))
	defer C.free(unsafe.Pointer(c))
	C.CtxLog((*C.struct_RedisModuleCtx)(ctx.ptr()), C.int(l), c)
}
func (ctx Ctx) LogDebug(format string, args ...interface{}) {
	ctx.Log(LOG_DEBUG, format, args...)
}
func (ctx Ctx) LogVerbose(format string, args ...interface{}) {
	ctx.Log(LOG_VERBOSE, format, args...)
}
func (ctx Ctx) LogNotice(format string, args ...interface{}) {
	ctx.Log(LOG_NOTICE, format, args...)
}
func (ctx Ctx) LogWarn(format string, args ...interface{}) {
	ctx.Log(LOG_WARNING, format, args...)
}

func (ctx Ctx) Init(name string, version int, apiVersion int) int {
	c := C.CString(name)
	defer C.free(unsafe.Pointer(c))
	return (int)(C.RedisModule_Init((*C.struct_RedisModuleCtx)(ctx.ptr()), c, (C.int)(version), (C.int)(apiVersion)))
}
func (c Ctx) Load(mod *Module, args []String) int {
	if mod == nil {
		c.LogWarn("Load Mod must not nil")
		return ERR
	}
	if c.Init(mod.Name, mod.Version, APIVER_1) == ERR {
		c.LogWarn("Init mod %s failed", mod.Name)
		return ERR
	}
	c.LogDebug("Load module %s %v", mod.Name, args)
	if mod.BeforeInit != nil {
		err := mod.BeforeInit(c, args)
		if err != nil {
			c.LogWarn("BeforeInit failed: %v", err)
			return ERR
		}
	}

	for _, cmd := range mod.Commands {
		if c.CreateCommand(cmd) == ERR {
			return ERR
		}
	}
	for _, dt := range mod.DataTypes {
		dataType := c.CreateDataType(dt)
		if dataType.IsNull() {
			return ERR
		}
		moduleDataTypes[dt.Name] = dataType
	}

	if mod.AfterInit != nil {
		err := mod.AfterInit(c, args)
		if err != nil {
			c.LogWarn("BeforeInit failed: %v", err)
			return ERR
		}
	}
	return OK
}

var moduleDataTypes = make(map[string]ModuleType)

func GetModuleDataType(name string) ModuleType {
	return moduleDataTypes[name]
}

//
func (c Ctx) CreateCommand(cmd Command) int {
	id := commandId(cmd)
	name := C.CString(cmd.Name)
	defer C.free(unsafe.Pointer(name))
	flags := C.CString(cmd.Flags)
	defer C.free(unsafe.Pointer(flags))
	c.LogVerbose("CreateCommand#%v %s", id, cmd.Usage)
	return (int)(C.CreateCommandCallID((*C.struct_RedisModuleCtx)(c.ptr()), C.int(id), name, flags, C.int(cmd.FirstKey), C.int(cmd.LastKey), C.int(cmd.KeyStep)))
}
func (c Ctx) CreateDataType(dt DataType) ModuleType {
	if m, _ := regexp.MatchString("[-_0-9A-Za-z]{9}", dt.Name); !m {
		c.LogWarn("Wrong datatype name need `[-_0-9A-Za-z]{9}` got %v", dt.Name)
		return ModuleType(0)
	}
	id := dataTypeId(dt)
	name := C.CString(dt.Name)
	defer C.free(unsafe.Pointer(name))
	c.LogVerbose("CreateDataType#%v %s", id, dt.Name)
	return ModuleType(C.CreateDataTypeCallID((*C.struct_RedisModuleCtx)(c.ptr()), C.int(id), name, C.int(dt.EncVer)))
}
func (ctx Ctx) ReplyWithOK() int {
	return int(C.ReplyWithOK((*C.struct_RedisModuleCtx)(ctx.ptr())))
}

// region CallReply

// Wrapper for the recursive free reply function. This is needed in order
// to have the first level function to return on nested replies, but only
// if called by the module API.
// void RM_FreeCallReply(RedisModuleCallReply *reply);
func (reply CallReply) FreeCallReply() {
	C.FreeCallReply((*C.struct_RedisModuleCallReply)(reply.ptr()))
}

// Return the reply type.
// int RM_CallReplyType(RedisModuleCallReply *reply);
func (reply CallReply) CallReplyType() int {
	return int(C.CallReplyType((*C.struct_RedisModuleCallReply)(reply.ptr())))
}

// Return the reply type length, where applicable.
// size_t RM_CallReplyLength(RedisModuleCallReply *reply);
func (reply CallReply) CallReplyLength() int {
	return int(C.CallReplyLength((*C.struct_RedisModuleCallReply)(reply.ptr())))
}

// Return the 'idx'-th nested call reply element of an array reply, or NULL
// if the reply type is wrong or the index is out of range.
// RedisModuleCallReply *RM_CallReplyArrayElement(RedisModuleCallReply *reply, size_t idx);
func (reply CallReply) CallReplyArrayElement(idx int) CallReply {
	return CreateCallReply(unsafe.Pointer(C.CallReplyArrayElement((*C.struct_RedisModuleCallReply)(reply.ptr()), C.size_t(idx))))
}

// Return the long long of an integer reply.
// long long RM_CallReplyInteger(RedisModuleCallReply *reply);
func (reply CallReply) CallReplyInteger() int64 {
	return int64(C.CallReplyInteger((*C.struct_RedisModuleCallReply)(reply.ptr())))
}

// Return the pointer and length of a string or error reply.
// const char *RM_CallReplyStringPtr(RedisModuleCallReply *reply, size_t *len);
func (reply CallReply) CallReplyStringPtr(len *int) unsafe.Pointer {
	return unsafe.Pointer(C.CallReplyStringPtr((*C.struct_RedisModuleCallReply)(reply.ptr()), (*C.size_t)(unsafe.Pointer(len))))
}

// Return a new string object from a call reply of type string, error or
// integer. Otherwise (wrong reply type) return NULL.
// RedisModuleString *RM_CreateStringFromCallReply(RedisModuleCallReply *reply);
func (reply CallReply) CreateStringFromCallReply() String {
	return CreateString(unsafe.Pointer(C.CreateStringFromCallReply((*C.struct_RedisModuleCallReply)(reply.ptr()))))
}

// Return a pointer, and a length, to the protocol returned by the command
// that returned the reply object.
// const char *RM_CallReplyProto(RedisModuleCallReply *reply, size_t *len);
func (reply CallReply) CallReplyProto(len *uint64) unsafe.Pointer {
	return unsafe.Pointer(C.CallReplyProto((*C.struct_RedisModuleCallReply)(reply.ptr()), (*C.size_t)(len)))
}

// endregion

// =============================================================================
// ========================== String functions
// =============================================================================

// Given a string module object, this function returns the string pointer
// and length of the string. The returned pointer and length should only
// be used for read only accesses and never modified.
// const char *RM_StringPtrLen(RedisModuleString *str, size_t *len);
func (str String) String() string {
	l := uint64(0)
	ptr := C.StringPtrLen((*C.struct_RedisModuleString)(str.ptr()), (*C.size_t)(&l))
	return C.GoStringN(ptr, C.int(l))
}

// Convert the string into a long long integer, storing it at `*ll`.
// Returns `REDISMODULE_OK` on success. If the string can't be parsed
// as a valid, strict long long (no spaces before/after), `REDISMODULE_ERR`
// is returned.
// int RM_StringToLongLong(RedisModuleString *str, long long *ll);
func (str String) StringToLongLong(ll *int64) int {
	return int(C.StringToLongLong((*C.struct_RedisModuleString)(str.ptr()), (*C.longlong)(ll)))
}

// Convert the string into a double, storing it at `*d`.
// Returns `REDISMODULE_OK` on success or `REDISMODULE_ERR` if the string is
// not a valid string representation of a double value.
// int RM_StringToDouble(RedisModuleString *str, double *d);
func (str String) StringToDouble(d *float64) int {
	return int(C.StringToDouble((*C.struct_RedisModuleString)(str.ptr()), (*C.double)(d)))
}

func (str String) Compare(b String) int {
	return int(C.StringCompare((*C.struct_RedisModuleString)(str.ptr()), (*C.struct_RedisModuleString)(b.ptr())))
}

// =============================================================================
// ========================== Key functions
// =============================================================================

func (key Key) IsEmpty() bool {
	return key.KeyType() == KEYTYPE_EMPTY
}

// Close a key handle.
// void RM_CloseKey(RedisModuleKey *key);
func (key Key) CloseKey() {
	C.CloseKey((*C.struct_RedisModuleKey)(key.ptr()))
}

// Return the type of the key. If the key pointer is NULL then
// `REDISMODULE_KEYTYPE_EMPTY` is returned.
// int RM_KeyType(RedisModuleKey *key);
func (key Key) KeyType() int {
	return int(C.KeyType((*C.struct_RedisModuleKey)(key.ptr())))
}

// Return the length of the value associated with the key.
// For strings this is the length of the string. For all the other types
// is the number of elements (just counting keys for hashes).
//
// If the key pointer is NULL or the key is empty, zero is returned.
// size_t RM_ValueLength(RedisModuleKey *key);
func (key Key) ValueLength() int {
	return int(C.ValueLength((*C.struct_RedisModuleKey)(key.ptr())))
}

// If the key is open for writing, remove it, and setup the key to
// accept new writes as an empty key (that will be created on demand).
// On success `REDISMODULE_OK` is returned. If the key is not open for
// writing `REDISMODULE_ERR` is returned.
// int RM_DeleteKey(RedisModuleKey *key);
func (key Key) DeleteKey() int {
	return int(C.DeleteKey((*C.struct_RedisModuleKey)(key.ptr())))
}

// Return the key expire value, as milliseconds of remaining TTL.
// If no TTL is associated with the key or if the key is empty,
// `REDISMODULE_NO_EXPIRE` is returned.
// mstime_t RM_GetExpire(RedisModuleKey *key);
func (key Key) GetExpire() uint64 {
	return uint64(C.GetExpire((*C.struct_RedisModuleKey)(key.ptr())))
}

// Set a new expire for the key. If the special expire
// `REDISMODULE_NO_EXPIRE` is set, the expire is cancelled if there was
// one (the same as the PERSIST command).
//
// Note that the expire must be provided as a positive integer representing
// the number of milliseconds of TTL the key should have.
//
// The function returns `REDISMODULE_OK` on success or `REDISMODULE_ERR` if
// the key was not open for writing or is an empty key.
// int RM_SetExpire(RedisModuleKey *key, mstime_t expire);
func (key Key) SetExpire(expire uint64) int {
	return int(C.SetExpire((*C.struct_RedisModuleKey)(key.ptr()), (C.mstime_t)(expire)))
}

// If the key is open for writing, set the specified string 'str' as the
// value of the key, deleting the old value if any.
// On success `REDISMODULE_OK` is returned. If the key is not open for
// writing or there is an active iterator, `REDISMODULE_ERR` is returned.
// int RM_StringSet(RedisModuleKey *key, RedisModuleString *str);
func (key Key) StringSet(str String) int {
	return int(C.StringSet((*C.struct_RedisModuleKey)(key.ptr()), (*C.struct_RedisModuleString)(str.ptr())))
}

// Prepare the key associated string value for DMA access, and returns
// a pointer and size (by reference), that the user can use to read or
// modify the string in-place accessing it directly via pointer.
//
// The 'mode' is composed by bitwise OR-ing the following flags:
//
// `REDISMODULE_READ` -- Read access
// `REDISMODULE_WRITE` -- Write access
//
// If the DMA is not requested for writing, the pointer returned should
// only be accessed in a read-only fashion.
//
// On error (wrong type) NULL is returned.
//
// DMA access rules:
//
// 1. No other key writing function should be called since the moment
// the pointer is obtained, for all the time we want to use DMA access
// to read or modify the string.
//
// 2. Each time `RM_StringTruncate()` is called, to continue with the DMA
// access, `RM_StringDMA()` should be called again to re-obtain
// a new pointer and length.
//
// 3. If the returned pointer is not NULL, but the length is zero, no
// byte can be touched (the string is empty, or the key itself is empty)
// so a `RM_StringTruncate()` call should be used if there is to enlarge
// the string, and later call StringDMA() again to get the pointer.
// char *RM_StringDMA(RedisModuleKey *key, size_t *len, int mode);
func (key Key) StringDMA(len *uint64, mode int) unsafe.Pointer {
	return unsafe.Pointer(C.StringDMA((*C.struct_RedisModuleKey)(key.ptr()), (*C.size_t)(len), C.int(mode)))
}

// If the string is open for writing and is of string type, resize it, padding
// with zero bytes if the new length is greater than the old one.
//
// After this call, `RM_StringDMA()` must be called again to continue
// DMA access with the new pointer.
//
// The function returns `REDISMODULE_OK` on success, and `REDISMODULE_ERR` on
// error, that is, the key is not open for writing, is not a string
// or resizing for more than 512 MB is requested.
//
// If the key is empty, a string key is created with the new string value
// unless the new length value requested is zero.
// int RM_StringTruncate(RedisModuleKey *key, size_t newlen);
func (key Key) StringTruncate(newlen int) int {
	return int(C.StringTruncate((*C.struct_RedisModuleKey)(key.ptr()), (C.size_t)(newlen)))
}

// Push an element into a list, on head or tail depending on 'where' argumnet.
// If the key pointer is about an empty key opened for writing, the key
// is created. On error (key opened for read-only operations or of the wrong
// type) `REDISMODULE_ERR` is returned, otherwise `REDISMODULE_OK` is returned.
// int RM_ListPush(RedisModuleKey *key, int where, RedisModuleString *ele);
func (key Key) ListPush(where int, ele String) int {
	return int(C.ListPush((*C.struct_RedisModuleKey)(key.ptr()), C.int(where), (*C.struct_RedisModuleString)(ele.ptr())))
}

// Pop an element from the list, and returns it as a module string object
// that the user should be free with `RM_FreeString()` or by enabling
// automatic memory. 'where' specifies if the element should be popped from
// head or tail. The command returns NULL if:
// 1) The list is empty.
// 2) The key was not open for writing.
// 3) The key is not a list.
// RedisModuleString *RM_ListPop(RedisModuleKey *key, int where);
func (key Key) ListPop(where int) String {
	return CreateString(unsafe.Pointer(C.ListPop((*C.struct_RedisModuleKey)(key.ptr()), C.int(where))))
}

// Add a new element into a sorted set, with the specified 'score'.
// If the element already exists, the score is updated.
//
// A new sorted set is created at value if the key is an empty open key
// setup for writing.
//
// Additional flags can be passed to the function via a pointer, the flags
// are both used to receive input and to communicate state when the function
// returns. 'flagsptr' can be NULL if no special flags are used.
//
// The input flags are:
//
// `REDISMODULE_ZADD_XX`: Element must already exist. Do nothing otherwise.
// `REDISMODULE_ZADD_NX`: Element must not exist. Do nothing otherwise.
//
// The output flags are:
//
// `REDISMODULE_ZADD_ADDED`: The new element was added to the sorted set.
// `REDISMODULE_ZADD_UPDATED`: The score of the element was updated.
// `REDISMODULE_ZADD_NOP`: No operation was performed because XX or NX flags.
//
// On success the function returns `REDISMODULE_OK`. On the following errors
// `REDISMODULE_ERR` is returned:
//
// * The key was not opened for writing.
// * The key is of the wrong type.
// * 'score' double value is not a number (NaN).
// int RM_ZsetAdd(RedisModuleKey *key, double score, RedisModuleString *ele, int *flagsptr);
func (key Key) ZsetAdd(score float64, ele String, flagsptr *int32) int {
	return int(C.ZsetAdd((*C.struct_RedisModuleKey)(key.ptr()), (C.double)(score), (*C.struct_RedisModuleString)(ele.ptr()), (*C.int)(flagsptr)))
}

// This function works exactly like `RM_ZsetAdd()`, but instead of setting
// a new score, the score of the existing element is incremented, or if the
// element does not already exist, it is added assuming the old score was
// zero.
//
// The input and output flags, and the return value, have the same exact
// meaning, with the only difference that this function will return
// `REDISMODULE_ERR` even when 'score' is a valid double number, but adding it
// to the existing score resuts into a NaN (not a number) condition.
//
// This function has an additional field 'newscore', if not NULL is filled
// with the new score of the element after the increment, if no error
// is returned.
// int RM_ZsetIncrby(RedisModuleKey *key, double score, RedisModuleString *ele, int *flagsptr, double *newscore);
func (key Key) ZsetIncrby(score float64, ele String, flagsptr *int32, newscore *float64) int {
	return int(C.ZsetIncrby((*C.struct_RedisModuleKey)(key.ptr()), (C.double)(score), (*C.struct_RedisModuleString)(ele.ptr()), (*C.int)(flagsptr), (*C.double)(newscore)))
}

// Remove the specified element from the sorted set.
// The function returns `REDISMODULE_OK` on success, and `REDISMODULE_ERR`
// on one of the following conditions:
//
// * The key was not opened for writing.
// * The key is of the wrong type.
//
// The return value does NOT indicate the fact the element was really
// removed (since it existed) or not, just if the function was executed
// with success.
//
// In order to know if the element was removed, the additional argument
// 'deleted' must be passed, that populates the integer by reference
// setting it to 1 or 0 depending on the outcome of the operation.
// The 'deleted' argument can be NULL if the caller is not interested
// to know if the element was really removed.
//
// Empty keys will be handled correctly by doing nothing.
// int RM_ZsetRem(RedisModuleKey *key, RedisModuleString *ele, int *deleted);
func (key Key) ZsetRem(ele String, deleted *int32) int {
	return int(C.ZsetRem((*C.struct_RedisModuleKey)(key.ptr()), (*C.struct_RedisModuleString)(ele.ptr()), (*C.int)(deleted)))
}

// On success retrieve the double score associated at the sorted set element
// 'ele' and returns `REDISMODULE_OK`. Otherwise `REDISMODULE_ERR` is returned
// to signal one of the following conditions:
//
// * There is no such element 'ele' in the sorted set.
// * The key is not a sorted set.
// * The key is an open empty key.
// int RM_ZsetScore(RedisModuleKey *key, RedisModuleString *ele, double *score);
func (key Key) ZsetScore(ele String, score *float64) int {
	return int(C.ZsetScore((*C.struct_RedisModuleKey)(key.ptr()), (*C.struct_RedisModuleString)(ele.ptr()), (*C.double)(score)))
}

// Stop a sorted set iteration.
// void RM_ZsetRangeStop(RedisModuleKey *key);
func (key Key) ZsetRangeStop() {
	C.ZsetRangeStop((*C.struct_RedisModuleKey)(key.ptr()))
}

// Return the "End of range" flag value to signal the end of the iteration.
// int RM_ZsetRangeEndReached(RedisModuleKey *key);
func (key Key) ZsetRangeEndReached() int {
	return int(C.ZsetRangeEndReached((*C.struct_RedisModuleKey)(key.ptr())))
}

// Setup a sorted set iterator seeking the first element in the specified
// range. Returns `REDISMODULE_OK` if the iterator was correctly initialized
// otherwise `REDISMODULE_ERR` is returned in the following conditions:
//
// 1. The value stored at key is not a sorted set or the key is empty.
//
// The range is specified according to the two double values 'min' and 'max'.
// Both can be infinite using the following two macros:
//
// `REDISMODULE_POSITIVE_INFINITE` for positive infinite value
// `REDISMODULE_NEGATIVE_INFINITE` for negative infinite value
//
// 'minex' and 'maxex' parameters, if true, respectively setup a range
// where the min and max value are exclusive (not included) instead of
// inclusive.
// int RM_ZsetFirstInScoreRange(RedisModuleKey *key, double min, double max, int minex, int maxex);
func (key Key) ZsetFirstInScoreRange(min float64, max float64, minex int, maxex int) int {
	return int(C.ZsetFirstInScoreRange((*C.struct_RedisModuleKey)(key.ptr()), (C.double)(min), (C.double)(max), (C.int)(minex), (C.int)(maxex)))
}

// Exactly like `RedisModule_ZsetFirstInScoreRange()` but the last element of
// the range is selected for the start of the iteration instead.
// int RM_ZsetLastInScoreRange(RedisModuleKey *key, double min, double max, int minex, int maxex);
func (key Key) ZsetLastInScoreRange(min float64, max float64, minex int, maxex int) int {
	return int(C.ZsetLastInScoreRange((*C.struct_RedisModuleKey)(key.ptr()), (C.double)(min), (C.double)(max), (C.int)(minex), (C.int)(maxex)))
}

// Setup a sorted set iterator seeking the first element in the specified
// lexicographical range. Returns `REDISMODULE_OK` if the iterator was correctly
// initialized otherwise `REDISMODULE_ERR` is returned in the
// following conditions:
//
// 1. The value stored at key is not a sorted set or the key is empty.
// 2. The lexicographical range 'min' and 'max' format is invalid.
//
// 'min' and 'max' should be provided as two RedisModuleString objects
// in the same format as the parameters passed to the ZRANGEBYLEX command.
// The function does not take ownership of the objects, so they can be released
// ASAP after the iterator is setup.
// int RM_ZsetFirstInLexRange(RedisModuleKey *key, RedisModuleString *min, RedisModuleString *max);
func (key Key) ZsetFirstInLexRange(min String, max String) int {
	return int(C.ZsetFirstInLexRange((*C.struct_RedisModuleKey)(key.ptr()), (*C.struct_RedisModuleString)(min.ptr()), (*C.struct_RedisModuleString)(max.ptr())))
}

// Exactly like `RedisModule_ZsetFirstInLexRange()` but the last element of
// the range is selected for the start of the iteration instead.
// int RM_ZsetLastInLexRange(RedisModuleKey *key, RedisModuleString *min, RedisModuleString *max);
func (key Key) ZsetLastInLexRange(min String, max String) int {
	return int(C.ZsetLastInLexRange((*C.struct_RedisModuleKey)(key.ptr()), (*C.struct_RedisModuleString)(min.ptr()), (*C.struct_RedisModuleString)(max.ptr())))
}

// Return the current sorted set element of an active sorted set iterator
// or NULL if the range specified in the iterator does not include any
// element.
// RedisModuleString *RM_ZsetRangeCurrentElement(RedisModuleKey *key, double *score);
func (key Key) ZsetRangeCurrentElement(score *float64) String {
	return CreateString(unsafe.Pointer(C.ZsetRangeCurrentElement((*C.struct_RedisModuleKey)(key.ptr()), (*C.double)(score))))
}

// Go to the next element of the sorted set iterator. Returns 1 if there was
// a next element, 0 if we are already at the latest element or the range
// does not include any item at all.
// int RM_ZsetRangeNext(RedisModuleKey *key);
func (key Key) ZsetRangeNext() int {
	return int(C.ZsetRangeNext((*C.struct_RedisModuleKey)(key.ptr())))
}

// Go to the previous element of the sorted set iterator. Returns 1 if there was
// a previous element, 0 if we are already at the first element or the range
// does not include any item at all.
// int RM_ZsetRangePrev(RedisModuleKey *key);
func (key Key) ZsetRangePrev() int {
	return int(C.ZsetRangePrev((*C.struct_RedisModuleKey)(key.ptr())))
}

// Return value:
//
// The number of fields updated (that may be less than the number of fields
// specified because of the XX or NX options).
//
// In the following case the return value is always zero:
//
// * The key was not open for writing.
// * The key was associated with a non Hash value.
// int RM_HashSet(RedisModuleKey *key, int flags, ...);
func (key Key) HashSet(flags int, args ...interface{}) int {
	args = append(args, uintptr(0))
	p, err := cutil.VarArgsPtr(args...)
	defer C.free(p)
	if err != nil {
		LogError("HashSet failed: %v", err)
		return ERR
	}
	return int(C.HashSetVar((*C.struct_RedisModuleKey)(key.ptr()), (C.int)(flags), C.int(len(args)), (*C.intptr_t)(p)))
}

// Get fields from an hash value. This function is called using a variable
// number of arguments, alternating a field name (as a StringRedisModule
// pointer) with a pointer to a StringRedisModule pointer, that is set to the
// value of the field if the field exist, or NULL if the field did not exist.
// At the end of the field/value-ptr pairs, NULL must be specified as last
// argument to signal the end of the arguments in the variadic function.
//
// This is an example usage:
//
//      RedisModuleString *first, *second;
//      `RedisModule_HashGet(mykey`,`REDISMODULE_HASH_NONE`,argv[1],&first,
//                      argv[2],&second,NULL);
//
// As with `RedisModule_HashSet()` the behavior of the command can be specified
// passing flags different than `REDISMODULE_HASH_NONE`:
//
// `REDISMODULE_HASH_CFIELD`: field names as null terminated C strings.
//
// `REDISMODULE_HASH_EXISTS`: instead of setting the value of the field
// expecting a RedisModuleString pointer to pointer, the function just
// reports if the field esists or not and expects an integer pointer
// as the second element of each pair.
//
// Example of `REDISMODULE_HASH_CFIELD`:
//
//      RedisModuleString *username, *hashedpass;
//      `RedisModule_HashGet(mykey`,"username",&username,"hp",&hashedpass, NULL);
//
// Example of `REDISMODULE_HASH_EXISTS`:
//
//      int exists;
//      `RedisModule_HashGet(mykey`,argv[1],&exists,NULL);
//
// The function returns `REDISMODULE_OK` on success and `REDISMODULE_ERR` if
// the key is not an hash value.
//
// Memory management:
//
// The returned RedisModuleString objects should be released with
// `RedisModule_FreeString()`, or by enabling automatic memory management.
// int RM_HashGet(RedisModuleKey *key, int flags, ...);
// args is RedisModuleString** RedisModuleString* or char* if flags include CFIELD
func (key Key) HashGet(flags int, args ...interface{}) int {
	args = append(args, uintptr(0))
	p, err := cutil.VarArgsPtr(args...)
	defer C.free(p)
	if err != nil {
		LogError("HashGet failed: %v", err)
		return ERR
	}
	return int(C.HashGetVar((*C.struct_RedisModuleKey)(key.ptr()), (C.int)(flags), C.int(len(args)), (*C.intptr_t)(p)))
}
func (key Key) HashExists(field String) bool {
	exists := 0
	key.HashGet(HASH_EXISTS, field, &exists)
	return exists == 1
}

//`RedisModule_HashSet(key`,`REDISMODULE_HASH_NONE`,argv[1], `REDISMODULE_HASH_DELETE`,NULL);
func (key Key) HashDel(field String) int {
	return key.HashSet(HASH_NONE, field, String(1))
}

// If the key is open for writing, set the specified module type object
// as the value of the key, deleting the old value if any.
// On success `REDISMODULE_OK` is returned. If the key is not open for
// writing or there is an active iterator, `REDISMODULE_ERR` is returned.
// int RM_ModuleTypeSetValue(RedisModuleKey *key, moduleType *mt, void *value);
func (key Key) ModuleTypeSetValue(mt ModuleType, value unsafe.Pointer) int {
	v := cutil.PtrToUintptr(value)
	return int(C.ModuleTypeSetValuePtr((*C.struct_RedisModuleKey)(key.ptr()), (*C.struct_RedisModuleType)(mt.ptr()), C.uintptr_t(v)))
}

// Assuming `RedisModule_KeyType()` returned `REDISMODULE_KEYTYPE_MODULE` on
// the key, returns the moduel type pointer of the value stored at key.
//
// If the key is NULL, is not associated with a module type, or is empty,
// then NULL is returned instead.
// moduleType *RM_ModuleTypeGetType(RedisModuleKey *key);
func (key Key) ModuleTypeGetType() ModuleType {
	return ModuleType(cutil.PtrToUintptr(unsafe.Pointer(C.ModuleTypeGetType((*C.struct_RedisModuleKey)(key.ptr())))))
}

// Assuming `RedisModule_KeyType()` returned `REDISMODULE_KEYTYPE_MODULE` on
// the key, returns the module type low-level value stored at key, as
// it was set by the user via `RedisModule_ModuleTypeSet()`.
//
// If the key is NULL, is not associated with a module type, or is empty,
// then NULL is returned instead.
// void *RM_ModuleTypeGetValue(RedisModuleKey *key);
func (key Key) ModuleTypeGetValue() unsafe.Pointer {
	return unsafe.Pointer(C.ModuleTypeGetValue((*C.struct_RedisModuleKey)(key.ptr())))
}

// =============================================================================
// ========================== IO functions
// =============================================================================

// Save an unsigned 64 bit value into the RDB file. This function should only
// be called in the context of the rdb_save method of modules implementing new
// data types.
// void RM_SaveUnsigned(RedisModuleIO *io, uint64_t value);
func (io IO) SaveUnsigned(value uint64) {
	C.SaveUnsigned((*C.struct_RedisModuleIO)(io.ptr()), C.uint64_t(value))
}

// Load an unsigned 64 bit value from the RDB file. This function should only
// be called in the context of the rdb_load method of modules implementing
// new data types.
// uint64_t RM_LoadUnsigned(RedisModuleIO *io);
func (io IO) LoadUnsigned() uint64 {
	return uint64(C.LoadUnsigned((*C.struct_RedisModuleIO)(io.ptr())))
}

// Like `RedisModule_SaveUnsigned()` but for signed 64 bit values.
// void RM_SaveSigned(RedisModuleIO *io, int64_t value);
func (io IO) SaveSigned(value int64) {
	C.SaveSigned((*C.struct_RedisModuleIO)(io.ptr()), C.int64_t(value))
}

// Like `RedisModule_LoadUnsigned()` but for signed 64 bit values.
// int64_t RM_LoadSigned(RedisModuleIO *io);
func (io IO) LoadSigned() int64 {
	return int64(C.LoadSigned((*C.struct_RedisModuleIO)(io.ptr())))
}

// In the context of the rdb_save method of a module type, saves a
// string into the RDB file taking as input a RedisModuleString.
//
// The string can be later loaded with `RedisModule_LoadString()` or
// other Load family functions expecting a serialized string inside
// the RDB file.
// void RM_SaveString(RedisModuleIO *io, RedisModuleString *s);
func (io IO) SaveString(s String) {
	C.SaveString((*C.struct_RedisModuleIO)(io.ptr()), (*C.struct_RedisModuleString)(s.ptr()))
}

// Like `RedisModule_SaveString()` but takes a raw C pointer and length
// as input.
// void RM_SaveStringBuffer(RedisModuleIO *io, const char *str, size_t len);
func (io IO) SaveStringBuffer(str []byte, len int) {
	// TODO useless
	v := C.CBytes(str[:len])
	defer C.free(unsafe.Pointer(v))
	C.SaveStringBuffer((*C.struct_RedisModuleIO)(io.ptr()), (*C.char)(v), C.size_t(len))
}

// In the context of the rdb_load method of a module data type, loads a string
// from the RDB file, that was previously saved with `RedisModule_SaveString()`
// functions family.
//
// The returned string is a newly allocated RedisModuleString object, and
// the user should at some point free it with a call to `RedisModule_FreeString()`.
//
// If the data structure does not store strings as RedisModuleString objects,
// the similar function `RedisModule_LoadStringBuffer()` could be used instead.
// RedisModuleString *RM_LoadString(RedisModuleIO *io);
func (io IO) LoadString() String {
	return CreateString(unsafe.Pointer(C.LoadString((*C.struct_RedisModuleIO)(io.ptr()))))
}

// Like `RedisModule_LoadString()` but returns an heap allocated string that
// was allocated with `RedisModule_Alloc()`, and can be resized or freed with
// `RedisModule_Realloc()` or `RedisModule_Free()`.
//
// The size of the string is stored at '*lenptr' if not NULL.
// The returned string is not automatically NULL termianted, it is loaded
// exactly as it was stored inisde the RDB file.
// char *RM_LoadStringBuffer(RedisModuleIO *io, size_t *lenptr);
func (io IO) LoadStringBuffer(lenptr *uint64) unsafe.Pointer {
	// func C.GoBytes(unsafe.Pointer, C.int) []byte
	// TODO return byte slice
	return unsafe.Pointer(C.LoadStringBuffer((*C.struct_RedisModuleIO)(io.ptr()), (*C.size_t)(lenptr)))
}

// In the context of the rdb_save method of a module data type, saves a double
// value to the RDB file. The double can be a valid number, a NaN or infinity.
// It is possible to load back the value with `RedisModule_LoadDouble()`.
// void RM_SaveDouble(RedisModuleIO *io, double value);
func (io IO) SaveDouble(value float64) {
	C.SaveDouble((*C.struct_RedisModuleIO)(io.ptr()), C.double(value))
}

// In the context of the rdb_save method of a module data type, loads back the
// double value saved by `RedisModule_SaveDouble()`.
// double RM_LoadDouble(RedisModuleIO *io);
func (io IO) LoadDouble() float64 {
	return float64(C.LoadDouble((*C.struct_RedisModuleIO)(io.ptr())))
}

// Emits a command into the AOF during the AOF rewriting process. This function
// is only called in the context of the aof_rewrite method of data types exported
// by a module. The command works exactly like `RedisModule_Call()` in the way
// the parameters are passed, but it does not return anything as the error
// handling is performed by Redis itself.
// void RM_EmitAOF(RedisModuleIO *io, const char *cmdname, const char *fmt, ...);
func (io IO) EmitAOF(cmdname string, format string, args ...interface{}) {
	v := C.CString(fmt.Sprintf(format, args))
	defer C.free(unsafe.Pointer(v))
	n := C.CString(cmdname)
	defer C.free(unsafe.Pointer(n))
	C.EmitAOF((*C.struct_RedisModuleIO)(io.ptr()), n, v)
}

// =============================================================================
// ========================== Util functions
// =============================================================================
