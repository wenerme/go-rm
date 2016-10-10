#ifndef GO_RM_RM_H
#define GO_RM_RM_H

#include "./callbacks.h"
#include "./redismodule.h"
#include "./wrapper.h"


int CreateCommandCallID(RedisModuleCtx *ctx, const char *name, int id, const char *strflags, int firstkey, int lastkey, int keystep) {
  return RedisModule_CreateCommand(ctx, name, cb_cmd_func[id], strflags, firstkey, lastkey, keystep);
}

#endif
