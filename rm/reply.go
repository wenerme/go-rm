package rm
//#include "./rm.h"
import "C"
import (
    "unsafe"
)


// Wrapper for the recursive free reply function. This is needed in order
// to have the first level function to return on nested replies, but only
// if called by the module API.
// void RM_FreeCallReply(RedisModuleCallReply *reply);
func (reply CallReply)FreeCallReply()(){C.FreeCallReply(reply)}

// Return the reply type.
// int RM_CallReplyType(RedisModuleCallReply *reply);
func (reply CallReply)CallReplyType()(int){return int(C.CallReplyType(reply))}

// Return the reply type length, where applicable.
// size_t RM_CallReplyLength(RedisModuleCallReply *reply);
func (reply CallReply)CallReplyLength()(int){return int(C.CallReplyLength(reply))}

// Return the 'idx'-th nested call reply element of an array reply, or NULL
// if the reply type is wrong or the index is out of range.
// RedisModuleCallReply *RM_CallReplyArrayElement(RedisModuleCallReply *reply, size_t idx);
func (reply CallReply)CallReplyArrayElement(idx int)(CallReply){return CallReply(C.CallReplyArrayElement(reply,idx))}

// Return the long long of an integer reply.
// long long RM_CallReplyInteger(RedisModuleCallReply *reply);
func (reply CallReply)CallReplyInteger()(int64){return int64(C.CallReplyInteger(reply))}

// Return the pointer and length of a string or error reply.
// const char *RM_CallReplyStringPtr(RedisModuleCallReply *reply, size_t *len);
func (reply CallReply)CallReplyStringPtr(len *int)(string){return string(C.CallReplyStringPtr(reply,len))}

// Return a new string object from a call reply of type string, error or
// integer. Otherwise (wrong reply type) return NULL.
// RedisModuleString *RM_CreateStringFromCallReply(RedisModuleCallReply *reply);
func (reply CallReply)CreateStringFromCallReply()(String){return String(C.CreateStringFromCallReply(reply))}

// Return a pointer, and a length, to the protocol returned by the command
// that returned the reply object.
// const char *RM_CallReplyProto(RedisModuleCallReply *reply, size_t *len);
func (reply CallReply)CallReplyProto(len *int)(string){return string(C.CallReplyProto(reply,len))}
