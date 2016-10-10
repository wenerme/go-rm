package rm

// #include "./rm.h"
import "C"
import "unsafe"

// Close a key handle.
// void RM_CloseKey(RedisModuleKey *key);
func (key Key)CloseKey()(){C.CloseKey(key)}

// Return the type of the key. If the key pointer is NULL then
// `REDISMODULE_KEYTYPE_EMPTY` is returned.
// int RM_KeyType(RedisModuleKey *key);
func (key Key)KeyType()(int){return int(C.KeyType(key))}

// Return the length of the value associated with the key.
// For strings this is the length of the string. For all the other types
// is the number of elements (just counting keys for hashes).
//
// If the key pointer is NULL or the key is empty, zero is returned.
// size_t RM_ValueLength(RedisModuleKey *key);
func (key Key)ValueLength()(int){return int(C.ValueLength(key))}

// If the key is open for writing, remove it, and setup the key to
// accept new writes as an empty key (that will be created on demand).
// On success `REDISMODULE_OK` is returned. If the key is not open for
// writing `REDISMODULE_ERR` is returned.
// int RM_DeleteKey(RedisModuleKey *key);
func (key Key)DeleteKey()(int){return int(C.DeleteKey(key))}

// Return the key expire value, as milliseconds of remaining TTL.
// If no TTL is associated with the key or if the key is empty,
// `REDISMODULE_NO_EXPIRE` is returned.
// mstime_t RM_GetExpire(RedisModuleKey *key);
func (key Key)GetExpire()(uint64){return uint64(C.GetExpire(key))}

// Set a new expire for the key. If the special expire
// `REDISMODULE_NO_EXPIRE` is set, the expire is cancelled if there was
// one (the same as the PERSIST command).
//
// Note that the expire must be provided as a positive integer representing
// the number of milliseconds of TTL the key should have.
//
// The function returns `REDISMODULE_OK` on success or `REDISMODULE_ERR` if
// the key was not open for writing or is an empty key.
// int RM_SetExpire(RedisModuleKey *key, mstime_t expire);
func (key Key)SetExpire(expire uint64)(int){return int(C.SetExpire(key,expire))}

// If the key is open for writing, set the specified string 'str' as the
// value of the key, deleting the old value if any.
// On success `REDISMODULE_OK` is returned. If the key is not open for
// writing or there is an active iterator, `REDISMODULE_ERR` is returned.
// int RM_StringSet(RedisModuleKey *key, RedisModuleString *str);
func (key Key)StringSet(str String)(int){return int(C.StringSet(key,str))}

// Prepare the key associated string value for DMA access, and returns
// a pointer and size (by reference), that the user can use to read or
// modify the string in-place accessing it directly via pointer.
//
// The 'mode' is composed by bitwise OR-ing the following flags:
//
// `REDISMODULE_READ` -- Read access
// `REDISMODULE_WRITE` -- Write access
//
// If the DMA is not requested for writing, the pointer returned should
// only be accessed in a read-only fashion.
//
// On error (wrong type) NULL is returned.
//
// DMA access rules:
//
// 1. No other key writing function should be called since the moment
// the pointer is obtained, for all the time we want to use DMA access
// to read or modify the string.
//
// 2. Each time `RM_StringTruncate()` is called, to continue with the DMA
// access, `RM_StringDMA()` should be called again to re-obtain
// a new pointer and length.
//
// 3. If the returned pointer is not NULL, but the length is zero, no
// byte can be touched (the string is empty, or the key itself is empty)
// so a `RM_StringTruncate()` call should be used if there is to enlarge
// the string, and later call StringDMA() again to get the pointer.
// char *RM_StringDMA(RedisModuleKey *key, size_t *len, int mode);
func (key Key)StringDMA(len *int,mode int)(string){return string(C.StringDMA(key,len,mode))}

// If the string is open for writing and is of string type, resize it, padding
// with zero bytes if the new length is greater than the old one.
//
// After this call, `RM_StringDMA()` must be called again to continue
// DMA access with the new pointer.
//
// The function returns `REDISMODULE_OK` on success, and `REDISMODULE_ERR` on
// error, that is, the key is not open for writing, is not a string
// or resizing for more than 512 MB is requested.
//
// If the key is empty, a string key is created with the new string value
// unless the new length value requested is zero.
// int RM_StringTruncate(RedisModuleKey *key, size_t newlen);
func (key Key)StringTruncate(newlen int)(int){return int(C.StringTruncate(key,newlen))}

// Push an element into a list, on head or tail depending on 'where' argumnet.
// If the key pointer is about an empty key opened for writing, the key
// is created. On error (key opened for read-only operations or of the wrong
// type) `REDISMODULE_ERR` is returned, otherwise `REDISMODULE_OK` is returned.
// int RM_ListPush(RedisModuleKey *key, int where, RedisModuleString *ele);
func (key Key)ListPush(where int,ele String)(int){return int(C.ListPush(key,where,ele))}

// Pop an element from the list, and returns it as a module string object
// that the user should be free with `RM_FreeString()` or by enabling
// automatic memory. 'where' specifies if the element should be popped from
// head or tail. The command returns NULL if:
// 1) The list is empty.
// 2) The key was not open for writing.
// 3) The key is not a list.
// RedisModuleString *RM_ListPop(RedisModuleKey *key, int where);
func (key Key)ListPop(where int)(String){return String(C.ListPop(key,where))}





// Add a new element into a sorted set, with the specified 'score'.
// If the element already exists, the score is updated.
//
// A new sorted set is created at value if the key is an empty open key
// setup for writing.
//
// Additional flags can be passed to the function via a pointer, the flags
// are both used to receive input and to communicate state when the function
// returns. 'flagsptr' can be NULL if no special flags are used.
//
// The input flags are:
//
// `REDISMODULE_ZADD_XX`: Element must already exist. Do nothing otherwise.
// `REDISMODULE_ZADD_NX`: Element must not exist. Do nothing otherwise.
//
// The output flags are:
//
// `REDISMODULE_ZADD_ADDED`: The new element was added to the sorted set.
// `REDISMODULE_ZADD_UPDATED`: The score of the element was updated.
// `REDISMODULE_ZADD_NOP`: No operation was performed because XX or NX flags.
//
// On success the function returns `REDISMODULE_OK`. On the following errors
// `REDISMODULE_ERR` is returned:
//
// * The key was not opened for writing.
// * The key is of the wrong type.
// * 'score' double value is not a number (NaN).
// int RM_ZsetAdd(RedisModuleKey *key, double score, RedisModuleString *ele, int *flagsptr);
func (key Key)ZsetAdd(score float64,ele String,flagsptr *int)(int){return int(C.ZsetAdd(key,score,ele,flagsptr))}

// This function works exactly like `RM_ZsetAdd()`, but instead of setting
// a new score, the score of the existing element is incremented, or if the
// element does not already exist, it is added assuming the old score was
// zero.
//
// The input and output flags, and the return value, have the same exact
// meaning, with the only difference that this function will return
// `REDISMODULE_ERR` even when 'score' is a valid double number, but adding it
// to the existing score resuts into a NaN (not a number) condition.
//
// This function has an additional field 'newscore', if not NULL is filled
// with the new score of the element after the increment, if no error
// is returned.
// int RM_ZsetIncrby(RedisModuleKey *key, double score, RedisModuleString *ele, int *flagsptr, double *newscore);
func (key Key)ZsetIncrby(score float64,ele String,flagsptr *int,newscore *float64)(int){return int(C.ZsetIncrby(key,score,ele,flagsptr,newscore))}

// Remove the specified element from the sorted set.
// The function returns `REDISMODULE_OK` on success, and `REDISMODULE_ERR`
// on one of the following conditions:
//
// * The key was not opened for writing.
// * The key is of the wrong type.
//
// The return value does NOT indicate the fact the element was really
// removed (since it existed) or not, just if the function was executed
// with success.
//
// In order to know if the element was removed, the additional argument
// 'deleted' must be passed, that populates the integer by reference
// setting it to 1 or 0 depending on the outcome of the operation.
// The 'deleted' argument can be NULL if the caller is not interested
// to know if the element was really removed.
//
// Empty keys will be handled correctly by doing nothing.
// int RM_ZsetRem(RedisModuleKey *key, RedisModuleString *ele, int *deleted);
func (key Key)ZsetRem(ele String,deleted *int)(int){return int(C.ZsetRem(key,ele,deleted))}

// On success retrieve the double score associated at the sorted set element
// 'ele' and returns `REDISMODULE_OK`. Otherwise `REDISMODULE_ERR` is returned
// to signal one of the following conditions:
//
// * There is no such element 'ele' in the sorted set.
// * The key is not a sorted set.
// * The key is an open empty key.
// int RM_ZsetScore(RedisModuleKey *key, RedisModuleString *ele, double *score);
func (key Key)ZsetScore(ele String,score *float64)(int){return int(C.ZsetScore(key,ele,score))}

// Stop a sorted set iteration.
// void RM_ZsetRangeStop(RedisModuleKey *key);
func (key Key)ZsetRangeStop()(){C.ZsetRangeStop(key)}

// Return the "End of range" flag value to signal the end of the iteration.
// int RM_ZsetRangeEndReached(RedisModuleKey *key);
func (key Key)ZsetRangeEndReached()(int){return int(C.ZsetRangeEndReached(key))}

// Setup a sorted set iterator seeking the first element in the specified
// range. Returns `REDISMODULE_OK` if the iterator was correctly initialized
// otherwise `REDISMODULE_ERR` is returned in the following conditions:
//
// 1. The value stored at key is not a sorted set or the key is empty.
//
// The range is specified according to the two double values 'min' and 'max'.
// Both can be infinite using the following two macros:
//
// `REDISMODULE_POSITIVE_INFINITE` for positive infinite value
// `REDISMODULE_NEGATIVE_INFINITE` for negative infinite value
//
// 'minex' and 'maxex' parameters, if true, respectively setup a range
// where the min and max value are exclusive (not included) instead of
// inclusive.
// int RM_ZsetFirstInScoreRange(RedisModuleKey *key, double min, double max, int minex, int maxex);
func (key Key)ZsetFirstInScoreRange(min float64,max float64,minex int,maxex int)(int){return int(C.ZsetFirstInScoreRange(key,min,max,minex,maxex))}

// Exactly like `RedisModule_ZsetFirstInScoreRange()` but the last element of
// the range is selected for the start of the iteration instead.
// int RM_ZsetLastInScoreRange(RedisModuleKey *key, double min, double max, int minex, int maxex);
func (key Key)ZsetLastInScoreRange(min float64,max float64,minex int,maxex int)(int){return int(C.ZsetLastInScoreRange(key,min,max,minex,maxex))}

// Setup a sorted set iterator seeking the first element in the specified
// lexicographical range. Returns `REDISMODULE_OK` if the iterator was correctly
// initialized otherwise `REDISMODULE_ERR` is returned in the
// following conditions:
//
// 1. The value stored at key is not a sorted set or the key is empty.
// 2. The lexicographical range 'min' and 'max' format is invalid.
//
// 'min' and 'max' should be provided as two RedisModuleString objects
// in the same format as the parameters passed to the ZRANGEBYLEX command.
// The function does not take ownership of the objects, so they can be released
// ASAP after the iterator is setup.
// int RM_ZsetFirstInLexRange(RedisModuleKey *key, RedisModuleString *min, RedisModuleString *max);
func (key Key)ZsetFirstInLexRange(min String,max String)(int){return int(C.ZsetFirstInLexRange(key,min,max))}

// Exactly like `RedisModule_ZsetFirstInLexRange()` but the last element of
// the range is selected for the start of the iteration instead.
// int RM_ZsetLastInLexRange(RedisModuleKey *key, RedisModuleString *min, RedisModuleString *max);
func (key Key)ZsetLastInLexRange(min String,max String)(int){return int(C.ZsetLastInLexRange(key,min,max))}

// Return the current sorted set element of an active sorted set iterator
// or NULL if the range specified in the iterator does not include any
// element.
// RedisModuleString *RM_ZsetRangeCurrentElement(RedisModuleKey *key, double *score);
func (key Key)ZsetRangeCurrentElement(score *float64)(String){return String(C.ZsetRangeCurrentElement(key,score))}

// Go to the next element of the sorted set iterator. Returns 1 if there was
// a next element, 0 if we are already at the latest element or the range
// does not include any item at all.
// int RM_ZsetRangeNext(RedisModuleKey *key);
func (key Key)ZsetRangeNext()(int){return int(C.ZsetRangeNext(key))}

// Go to the previous element of the sorted set iterator. Returns 1 if there was
// a previous element, 0 if we are already at the first element or the range
// does not include any item at all.
// int RM_ZsetRangePrev(RedisModuleKey *key);
func (key Key)ZsetRangePrev()(int){return int(C.ZsetRangePrev(key))}

// Return value:
//
// The number of fields updated (that may be less than the number of fields
// specified because of the XX or NX options).
//
// In the following case the return value is always zero:
//
// * The key was not open for writing.
// * The key was associated with a non Hash value.
// int RM_HashSet(RedisModuleKey *key, int flags, ...);
func (key Key)HashSet(flags int)(int){return int(C.HashSet(key,flags))}

// Get fields from an hash value. This function is called using a variable
// number of arguments, alternating a field name (as a StringRedisModule
// pointer) with a pointer to a StringRedisModule pointer, that is set to the
// value of the field if the field exist, or NULL if the field did not exist.
// At the end of the field/value-ptr pairs, NULL must be specified as last
// argument to signal the end of the arguments in the variadic function.
//
// This is an example usage:
//
//      RedisModuleString *first, *second;
//      `RedisModule_HashGet(mykey`,`REDISMODULE_HASH_NONE`,argv[1],&first,
//                      argv[2],&second,NULL);
//
// As with `RedisModule_HashSet()` the behavior of the command can be specified
// passing flags different than `REDISMODULE_HASH_NONE`:
//
// `REDISMODULE_HASH_CFIELD`: field names as null terminated C strings.
//
// `REDISMODULE_HASH_EXISTS`: instead of setting the value of the field
// expecting a RedisModuleString pointer to pointer, the function just
// reports if the field esists or not and expects an integer pointer
// as the second element of each pair.
//
// Example of `REDISMODULE_HASH_CFIELD`:
//
//      RedisModuleString *username, *hashedpass;
//      `RedisModule_HashGet(mykey`,"username",&username,"hp",&hashedpass, NULL);
//
// Example of `REDISMODULE_HASH_EXISTS`:
//
//      int exists;
//      `RedisModule_HashGet(mykey`,argv[1],&exists,NULL);
//
// The function returns `REDISMODULE_OK` on success and `REDISMODULE_ERR` if
// the key is not an hash value.
//
// Memory management:
//
// The returned RedisModuleString objects should be released with
// `RedisModule_FreeString()`, or by enabling automatic memory management.
// int RM_HashGet(RedisModuleKey *key, int flags, ...);
func (key Key)HashGet(flags int)(int){return int(C.HashGet(key,flags))}

// If the key is open for writing, set the specified module type object
// as the value of the key, deleting the old value if any.
// On success `REDISMODULE_OK` is returned. If the key is not open for
// writing or there is an active iterator, `REDISMODULE_ERR` is returned.
// int RM_ModuleTypeSetValue(RedisModuleKey *key, moduleType *mt, void *value);
func (key Key)ModuleTypeSetValue(mt /* TODO RedisModuleType* */unsafe.Pointer,value unsafe.Pointer)(int){return int(C.ModuleTypeSetValue(key,mt,value))}

// Assuming `RedisModule_KeyType()` returned `REDISMODULE_KEYTYPE_MODULE` on
// the key, returns the moduel type pointer of the value stored at key.
//
// If the key is NULL, is not associated with a module type, or is empty,
// then NULL is returned instead.
// moduleType *RM_ModuleTypeGetType(RedisModuleKey *key);
func (key Key)ModuleTypeGetType()(/* TODO RedisModuleType* */unsafe.Pointer){return /* TODO RedisModuleType* */unsafe.Pointer(C.ModuleTypeGetType(key))}

// Assuming `RedisModule_KeyType()` returned `REDISMODULE_KEYTYPE_MODULE` on
// the key, returns the module type low-level value stored at key, as
// it was set by the user via `RedisModule_ModuleTypeSet()`.
//
// If the key is NULL, is not associated with a module type, or is empty,
// then NULL is returned instead.
// void *RM_ModuleTypeGetValue(RedisModuleKey *key);
func (key Key)ModuleTypeGetValue()(unsafe.Pointer){return unsafe.Pointer(C.ModuleTypeGetValue(key))}
