#ifndef GO_RM_DECL_H
#define GO_RM_DECL_H

#include "./redismodule.h"

extern int RedisModule_OnLoad(RedisModuleCtx *ctx);

extern int cmd_func_call(int id, RedisModuleCtx *ctx, RedisModuleString **argv, int argc);

extern void *mt_rdb_load_call(int id, RedisModuleIO *rdb, int encver);
extern void mt_rdb_save_call(int id, RedisModuleIO *rdb, void *value);
extern void mt_aof_rewrite_call(int id, RedisModuleIO *aof, RedisModuleString *key, void *value);
extern void mt_digest_call(int id, RedisModuleDigest *digest, void *value);
extern void mt_free_call(int id, void *value);
#endif
