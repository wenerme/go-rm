package rm

// #include "./rm.h"
import "C"




// Save an unsigned 64 bit value into the RDB file. This function should only
// be called in the context of the rdb_save method of modules implementing new
// data types.
// void RM_SaveUnsigned(RedisModuleIO *io, uint64_t value);
func (io IO)SaveUnsigned(value uint64)(){C.SaveUnsigned(io,value)}

// Load an unsigned 64 bit value from the RDB file. This function should only
// be called in the context of the rdb_load method of modules implementing
// new data types.
// uint64_t RM_LoadUnsigned(RedisModuleIO *io);
func (io IO)LoadUnsigned()(uint64){return uint64(C.LoadUnsigned(io))}

// Like `RedisModule_SaveUnsigned()` but for signed 64 bit values.
// void RM_SaveSigned(RedisModuleIO *io, int64_t value);
func (io IO)SaveSigned(value int64)(){C.SaveSigned(io,value)}

// Like `RedisModule_LoadUnsigned()` but for signed 64 bit values.
// int64_t RM_LoadSigned(RedisModuleIO *io);
func (io IO)LoadSigned()(int64){return int64(C.LoadSigned(io))}

// In the context of the rdb_save method of a module type, saves a
// string into the RDB file taking as input a RedisModuleString.
//
// The string can be later loaded with `RedisModule_LoadString()` or
// other Load family functions expecting a serialized string inside
// the RDB file.
// void RM_SaveString(RedisModuleIO *io, RedisModuleString *s);
func (io IO)SaveString(s String)(){C.SaveString(io,s)}

// Like `RedisModule_SaveString()` but takes a raw C pointer and length
// as input.
// void RM_SaveStringBuffer(RedisModuleIO *io, const char *str, size_t len);
func (io IO)SaveStringBuffer(str string,len int)(){C.SaveStringBuffer(io,str,len)}

// In the context of the rdb_load method of a module data type, loads a string
// from the RDB file, that was previously saved with `RedisModule_SaveString()`
// functions family.
//
// The returned string is a newly allocated RedisModuleString object, and
// the user should at some point free it with a call to `RedisModule_FreeString()`.
//
// If the data structure does not store strings as RedisModuleString objects,
// the similar function `RedisModule_LoadStringBuffer()` could be used instead.
// RedisModuleString *RM_LoadString(RedisModuleIO *io);
func (io IO)LoadString()(String){return String(C.LoadString(io))}

// Like `RedisModule_LoadString()` but returns an heap allocated string that
// was allocated with `RedisModule_Alloc()`, and can be resized or freed with
// `RedisModule_Realloc()` or `RedisModule_Free()`.
//
// The size of the string is stored at '*lenptr' if not NULL.
// The returned string is not automatically NULL termianted, it is loaded
// exactly as it was stored inisde the RDB file.
// char *RM_LoadStringBuffer(RedisModuleIO *io, size_t *lenptr);
func (io IO)LoadStringBuffer(lenptr *int)(string){return string(C.LoadStringBuffer(io,lenptr))}

// In the context of the rdb_save method of a module data type, saves a double
// value to the RDB file. The double can be a valid number, a NaN or infinity.
// It is possible to load back the value with `RedisModule_LoadDouble()`.
// void RM_SaveDouble(RedisModuleIO *io, double value);
func (io IO)SaveDouble(value float64)(){C.SaveDouble(io,value)}

// In the context of the rdb_save method of a module data type, loads back the
// double value saved by `RedisModule_SaveDouble()`.
// double RM_LoadDouble(RedisModuleIO *io);
func (io IO)LoadDouble()(float64){return float64(C.LoadDouble(io))}

// Emits a command into the AOF during the AOF rewriting process. This function
// is only called in the context of the aof_rewrite method of data types exported
// by a module. The command works exactly like `RedisModule_Call()` in the way
// the parameters are passed, but it does not return anything as the error
// handling is performed by Redis itself.
// void RM_EmitAOF(RedisModuleIO *io, const char *cmdname, const char *fmt, ...);
func (io IO)EmitAOF(cmdname string,fmt string)(){C.EmitAOF(io,cmdname,fmt)}

