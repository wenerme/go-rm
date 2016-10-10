package rm

// #include "./rm.h"
import "C"
import (
    "unsafe"
    "fmt"
    "os"
)
/* ---------------- Global context --------------- */
var callbacks[]CmdFunc

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
const KEYTYPE_EMPTY = C.REDISMODULE_KEYTYPE_EMPTY
const KEYTYPE_STRING = C.REDISMODULE_KEYTYPE_STRING
const KEYTYPE_LIST = C.REDISMODULE_KEYTYPE_LIST
const KEYTYPE_HASH = C.REDISMODULE_KEYTYPE_HASH
const KEYTYPE_SET = C.REDISMODULE_KEYTYPE_SET
const KEYTYPE_ZSET = C.REDISMODULE_KEYTYPE_ZSET
const KEYTYPE_MODULE = C.REDISMODULE_KEYTYPE_MODULE

/* Reply types. */
const REPLY_UNKNOWN = C.REDISMODULE_REPLY_UNKNOWN
const REPLY_STRING = C.REDISMODULE_REPLY_STRING
const REPLY_ERROR = C.REDISMODULE_REPLY_ERROR
const REPLY_INTEGER = C.REDISMODULE_REPLY_INTEGER
const REPLY_ARRAY = C.REDISMODULE_REPLY_ARRAY
const REPLY_NULL = C.REDISMODULE_REPLY_NULL

/* Postponed array length. */
const POSTPONED_ARRAY_LEN = C.REDISMODULE_POSTPONED_ARRAY_LEN

/* Expire */
const NO_EXPIRE = C.REDISMODULE_NO_EXPIRE

/* Sorted set API flags. */
const ZADD_XX = C.REDISMODULE_ZADD_XX
const ZADD_NX = C.REDISMODULE_ZADD_NX
const ZADD_ADDED = C.REDISMODULE_ZADD_ADDED
const ZADD_UPDATED = C.REDISMODULE_ZADD_UPDATED
const ZADD_NOP = C.REDISMODULE_ZADD_NOP

/* Hash API flags. */
const HASH_NONE = C.REDISMODULE_HASH_NONE
const HASH_NX = C.REDISMODULE_HASH_NX
const HASH_XX = C.REDISMODULE_HASH_XX
const HASH_CFIELDS = C.REDISMODULE_HASH_CFIELDS
const HASH_EXISTS = C.REDISMODULE_HASH_EXISTS

/* A special pointer that we can use between the core and the module to signal
 * field deletion, and that is impossible to be a valid pointer. */
//const HASH_DELETE = C.REDISMODULE_HASH_DELETE

/* Error messages. */
const ERRORMSG_WRONGTYPE = C.REDISMODULE_ERRORMSG_WRONGTYPE

//const POSITIVE_INFINITE = C.REDISMODULE_POSITIVE_INFINITE
//const NEGATIVE_INFINITE = C.REDISMODULE_NEGATIVE_INFINITE

/* ------------------------- End of common defines ------------------------ */


// Use like malloc(). Memory allocated with this function is reported in
// Redis INFO memory, used for keys eviction according to maxmemory settings
// and in general is taken into account as memory allocated by Redis.
// You should avoid to use malloc().
// void *RM_Alloc(size_t bytes);
func Alloc(bytes int)(unsafe.Pointer){return unsafe.Pointer(C.Alloc(bytes))}

// Use like realloc() for memory obtained with `RedisModule_Alloc()`.
// void* RM_Realloc(void *ptr, size_t bytes);
func Realloc(ptr unsafe.Pointer,bytes int)(unsafe.Pointer){return unsafe.Pointer(C.Realloc(ptr,bytes))}

// Use like free() for memory obtained by `RedisModule_Alloc()` and
// `RedisModule_Realloc()`. However you should never try to free with
// `RedisModule_Free()` memory allocated with malloc() inside your module.
// void RM_Free(void *ptr);
func Free(ptr unsafe.Pointer)(){C.Free(ptr)}

// Like strdup() but returns memory allocated with `RedisModule_Alloc()`.
// char *RM_Strdup(const char *str);
func Strdup(str unsafe.Pointer)(unsafe.Pointer){return unsafe.Pointer(C.Strdup(str))}

// Lookup the requested module API and store the function pointer into the
// target pointer. The function returns `REDISMODULE_ERR` if there is no such
// named API, otherwise `REDISMODULE_OK`.
//
// This function is not meant to be used by modules developer, it is only
// used implicitly by including redismodule.h.
// int RM_GetApi(const char *funcname, void **targetPtrPtr);
//func GetApi(funcname string,targetPtrPtr /* TODO void** */unsafe.Pointer)(int){return int(C.GetApi(funcname,targetPtrPtr))}


func init() {
    LogDebug("Init Go Redis module")
}

var LogErr = func(format string, args... interface{}) {
    fmt.Fprintf(os.Stderr, format + "\n", args...)
}

var LogDebug = func(format string, args... interface{}) {
    fmt.Fprintf(os.Stdout, format + "\n", args...)
}

type CmdArgs struct {
    argv unsafe.Pointer
    argc int
}
type CmdContext struct {
    Ctx Ctx
}

