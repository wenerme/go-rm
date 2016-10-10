#ifndef GO_RM_RM_H
#define GO_RM_RM_H

#include "stdlib.h"
#include "./callbacks.h"
#include "./redismodule.h"
#include "./wrapper.h"


int CreateCommandCallID(RedisModuleCtx *ctx, const char *name, int id, const char *strflags, int firstkey, int lastkey, int keystep) {
  return RedisModule_CreateCommand(ctx, name, cb_cmd_func[id], strflags, firstkey, lastkey, keystep);
}

#define LOG_DEBUG   "debug"
#define LOG_VERBOSE "verbose"
#define LOG_NOTICE  "notice"
#define LOG_WARNING "warning"
void CtxLog(RedisModuleCtx* ctx,int level,const char* fmt){
    char* l;
    switch(level){
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
    RedisModule_Log(ctx,l,fmt);
}


#endif
