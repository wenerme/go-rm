#ifndef GO_RM_RM_H
#define GO_RM_RM_H

#include <errno.h>
#include <stdlib.h>
#include <stdint.h>
#include "./redismodule.h"
#include "./callbacks.h"
#include "./wrapper.h"
#include "./varargs.h"

int ReplyWithOK(RedisModuleCtx* ctx){return RedisModule_ReplyWithSimpleString(ctx,"OK");}
int ModuleTypeSetValuePtr(RedisModuleKey* key,RedisModuleType* mt,uintptr_t value){return RedisModule_ModuleTypeSetValue(key,mt,(void*)(value));}

int CreateCommandCallID(RedisModuleCtx *ctx,int id, const char *name,  const char *strflags, int firstkey, int lastkey, int keystep) {
  return RedisModule_CreateCommand(ctx, name, cb_cmd_func[id], strflags, firstkey, lastkey, keystep);
}

uintptr_t CreateDataTypeCallID(RedisModuleCtx* ctx,int id,const char* name,int encver){
    RedisModuleTypeLoadFunc rdb_load =cb_mt_rdb_load[id];
    RedisModuleTypeSaveFunc rdb_save =cb_mt_rdb_save[id];
    RedisModuleTypeRewriteFunc aof_rewrite =cb_mt_aof_rewrite[id];
    RedisModuleTypeDigestFunc digest =cb_mt_digest[id];
    RedisModuleTypeFreeFunc free =cb_mt_free[id];
    return (uintptr_t)CreateDataType(ctx,name,encver,rdb_load,rdb_save,aof_rewrite,digest,free);
}

#define LOG_DEBUG   "debug"
#define LOG_VERBOSE "verbose"
#define LOG_NOTICE  "notice"
#define LOG_WARNING "warning"

void CtxLog(RedisModuleCtx *ctx, int level, const char *fmt) {
  char *l;
  switch (level) {
    default:
    case 0:
      l = LOG_DEBUG;
      break;
    case 1:
      l = LOG_VERBOSE;
      break;
    case 2:
      l = LOG_NOTICE;
      break;
    case 3:
      l = LOG_WARNING;
      break;
  }
  RedisModule_Log(ctx, l, fmt);
}

int get_errno(){
    return errno;
}

#endif
