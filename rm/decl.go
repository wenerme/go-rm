package rm

/*
#include <stdlib.h>
 */
import "C"
import (
	"fmt"
	"unsafe"
	"bytes"
	"github.com/wenerme/letsgo/cutil"
)

//export RedisModule_OnLoad
func RedisModule_OnLoad(ctx uintptr, argv uintptr, argc int) C.int {
	return C.int(Ctx(ctx).Load(Mod, toStringSlice(argv, argc)))
}
//export redis_module_on_unload
func redis_module_on_unload() {
	if Mod != nil && Mod.OnUnload != nil {
		Mod.OnUnload()
	}
}

//export cmd_func_call
func cmd_func_call(id C.int, ctx uintptr, argv uintptr, argc int) C.int {
	args := toStringSlice(argv, argc)
	c := Ctx(ctx)
	cmd := getCommand(int(id))
	if Debug {
		buf := bytes.NewBufferString(fmt.Sprintf("CmdFuncCall(%v): %v", id, cmd.Name))
		for i := 0; i < argc; i ++ {
			buf.WriteString(" ")
			buf.WriteString(args[i].String())
		}
		c.LogDebug(buf.String())
	}
	return C.int(cmd.Action(CmdContext{Ctx:c, Args:args}))
}

func toStringSlice(argv uintptr, argc int) []String {
	args := make([]String, argc)
	size := int(unsafe.Sizeof(C.uintptr_t(0)))
	for i := 0; i < argc; i ++ {
		ptr := unsafe.Pointer(argv + uintptr(size * i))
		args[i] = String(uintptr(*(*C.uintptr_t)(ptr)))
	}
	return args
}

// typedef int (*RedisModuleCmdFunc) (RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
// typedef void *(*RedisModuleTypeLoadFunc)(RedisModuleIO *rdb, int encver);
// typedef void (*RedisModuleTypeSaveFunc)(RedisModuleIO *rdb, void *value);
// typedef void (*RedisModuleTypeRewriteFunc)(RedisModuleIO *aof, RedisModuleString *key, void *value);
// typedef void (*RedisModuleTypeDigestFunc)(RedisModuleDigest *digest, void *value);
// typedef void (*RedisModuleTypeFreeFunc)(void *value);

//export mt_rdb_load_call
func mt_rdb_load_call(id int, rdb uintptr, encver int) uintptr {
	dt := getDataType(id)
	if dt.RdbLoad != nil {
		return cutil.PtrToUintptr(dt.RdbLoad(IO(rdb), encver))
	}
	return 0
}
//export mt_rdb_save_call
func mt_rdb_save_call(id int, rdb uintptr, value uintptr) {
	dt := getDataType(id)
	if dt.RdbSave != nil {
		dt.RdbSave(IO(rdb), unsafe.Pointer(value))
	}
}
//export mt_aof_rewrite_call
func mt_aof_rewrite_call(id int, aof uintptr, key uintptr, value uintptr) {
	dt := getDataType(id)
	if dt.AofRewrite != nil {
		dt.AofRewrite(IO(aof), String(key), unsafe.Pointer(value))
	}
}
//export mt_digest_call
func mt_digest_call(id int, digest uintptr, value uintptr) {
	dt := getDataType(id)
	if dt.Digest != nil {
		dt.Digest(Digest(digest), unsafe.Pointer(value))
	}
}
//export mt_free_call
func mt_free_call(id int, value uintptr) {
	dt := getDataType(id)
	if dt.Free != nil {
		dt.Free(unsafe.Pointer(value))
	}
}

func init() {
	// Preserve import "C"
	_ = C.int(0)
}
