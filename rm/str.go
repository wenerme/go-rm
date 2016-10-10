package rm

//#include "./rm.h"
import "C"
import (
    "unsafe"
)


// Given a string module object, this function returns the string pointer
// and length of the string. The returned pointer and length should only
// be used for read only accesses and never modified.
// const char *RM_StringPtrLen(RedisModuleString *str, size_t *len);
func (str String)StringPtrLen(len *int)(string){return string(C.StringPtrLen(str,len))}

// Convert the string into a long long integer, storing it at `*ll`.
// Returns `REDISMODULE_OK` on success. If the string can't be parsed
// as a valid, strict long long (no spaces before/after), `REDISMODULE_ERR`
// is returned.
// int RM_StringToLongLong(RedisModuleString *str, long long *ll);
func (str String)StringToLongLong(ll *int64)(int){return int(C.StringToLongLong(str,ll))}

// Convert the string into a double, storing it at `*d`.
// Returns `REDISMODULE_OK` on success or `REDISMODULE_ERR` if the string is
// not a valid string representation of a double value.
// int RM_StringToDouble(RedisModuleString *str, double *d);
func (str String)StringToDouble(d *float64)(int){return int(C.StringToDouble(str,d))}

