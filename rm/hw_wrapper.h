#ifndef GO_RM_HW_WRAPPER_H
#define GO_RM_HW_WRAPPER_H

#include "./redismodule.h"

// Hand write wrapper because these function not defined in API.md

//void REDISMODULE_API_FUNC(RedisModule_LogIOError)(RedisModuleIO *io, const char *levelstr, const char *fmt, ...);
void LogIOError(RedisModuleIO *io, const char *levelstr, const char *fmt){
    RedisModule_LogIOError(io,levelstr,fmt,0);
}

//int REDISMODULE_API_FUNC(RedisModule_StringAppendBuffer)(RedisModuleCtx *ctx, RedisModuleString *str, const char *buf, size_t len);
int StringAppendBuffer(RedisModuleCtx *ctx, RedisModuleString *str, const char *buf, size_t len){
    return RedisModule_StringAppendBuffer(ctx,str,buf,len);
}

//void REDISMODULE_API_FUNC(RedisModule_RetainString)(RedisModuleCtx *ctx, RedisModuleString *str);
void RetainString(RedisModuleCtx *ctx, RedisModuleString *str){
    return RedisModule_RetainString(ctx,str);
}

// int REDISMODULE_API_FUNC(RedisModule_StringCompare)(RedisModuleString *a, RedisModuleString *b);
int StringCompare(RedisModuleString *a, RedisModuleString *b){
    return RedisModule_StringCompare(a,b);
}

//RedisModuleCtx *REDISMODULE_API_FUNC(RedisModule_GetContextFromIO)(RedisModuleIO *io);
RedisModuleCtx * GetContextFromIO(RedisModuleIO *io){
    return RedisModule_GetContextFromIO(io);
}

#endif
