package rm

import "C"
import "fmt"

//export RedisModule_OnLoad
func RedisModule_OnLoad(ctx uintptr) C.int {
    return C.int(Ctx(ctx).Load(Mod))
}

//export cmd_func_call
func cmd_func_call(id C.int, ctx uintptr, argv uintptr, argc int) C.int {
    fmt.Println("Recv command function callback")
    return C.int(getCommand(int(id)).Action(CmdContext{Ctx:Ctx(ctx)}))
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
