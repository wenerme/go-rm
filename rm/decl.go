package rm

/*
#include <stdlib.h>
 */
import "C"
import (
    "fmt"
    "unsafe"
    "bytes"
)

//export RedisModule_OnLoad
func RedisModule_OnLoad(ctx uintptr) C.int {
    return C.int(Ctx(ctx).Load(Mod))
}

//export cmd_func_call
func cmd_func_call(id C.int, ctx uintptr, argv uintptr, argc int) C.int {
    args := make([]String, argc)
    size := int(unsafe.Sizeof(C.uintptr_t(0)))
    for i := 0; i < argc; i ++ {
        ptr := unsafe.Pointer(argv + uintptr(size * i))
        args[i] = String(uintptr(*(*C.uintptr_t)(ptr)))
    }
    c := Ctx(ctx)
    cmd := getCommand(int(id))
    buf := bytes.NewBufferString(fmt.Sprintf("CmdFuncCall(%v): %v", id, cmd.Name))
    for i := 0; i < argc; i ++ {
        buf.WriteString(" ")
        buf.WriteString(args[i].String())
    }
    c.LogDebug(buf.String())
    return C.int(cmd.Action(CmdContext{Ctx:c, Args:args}))
}
//export mt_rdb_load_call
func mt_rdb_load_call(id int, rdb uintptr, encver int) uintptr {
    return 0
}
//export mt_rdb_save_call
func mt_rdb_save_call(id int, rdb uintptr, value uintptr) {

}
//export mt_aof_rewrite_call
func mt_aof_rewrite_call(id int, aof uintptr, key uintptr, value uintptr) {

}
//export mt_digest_call
func mt_digest_call(id int, digest uintptr, value uintptr) {

}
//export mt_free
func mt_free(id int, value uintptr) {

}

func init() {
    // Preserve import "C"
    _ = C.int(0)
}
