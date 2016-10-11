# Go Redis module
go-rm will let you write redis module in golang.

## Demo

```bash
# Ensure you installed the newest redis
# for example by using brew you can 
# brew reinstall redis --with-jemalloc --HEAD

# Build redis module
go get -u -v -buildmode=c-shared github.com/wenerme/go-rm/modules/hashex

# Start redis-server and load out module with debug log
redis-server --loadmodule ~/go/pkg/*/github.com/wenerme/go-rm/modules/hashex* --loglevel debug
```

__Connect to out redis-server__

```
$ redis-cli hsetget a a 5
(nil)
$ redis-cli hsetget a a 4
"5"
$ redis-cli hsetget a a 3
"3"
```

Wow, it works, now you can distribute this redis module to you friends. :P

## How to write a module

Implement a redis module is as easy as you write a cli app in go, this is all you need to implement above command, the source code is [here](https://github.com/wenerme/go-rm/blob/master/modules/hashex/hashex.go).

```go
package main

import "github.com/wenerme/go-rm/rm"

func main() {
    // In case someone try to run this
    rm.Run()
}

func init() {
    rm.Mod = CreateMyMod()
}
func CreateMyMod() *rm.Module {
    mod := rm.NewMod()
    mod.Name = "hashex"
    mod.Version = 1
    mod.Commands = []rm.Command{
        {
            Name:   "hsetget",
            Flags:  "write fast deny-oom",
            FirstKey:1, LastKey:1, KeyStep:1,
            Action: func(cmd rm.CmdContext) int {
                ctx := cmd.Ctx
                if len(cmd.Args) != 4 {
                    return ctx.WrongArity()
                }
                ctx.AutoMemory()
                // open the key and make sure it is indeed a Hash and not empty
                key := ctx.OpenKey(cmd.Args[1], rm.READ | rm.WRITE)
                if key.KeyType() != rm.KEYTYPE_EMPTY && key.KeyType() != rm.KEYTYPE_HASH {
                    ctx.ReplyWithError(rm.ERRORMSG_WRONGTYPE)
                    return rm.ERR
                }
                // get the current value of the hash element
                var val rm.String;
                key.HashGet(rm.HASH_NONE, cmd.Args[2], (*uintptr)(&val), rm.NullString())
                // set the element to the new value
                key.HashSet(rm.HASH_NONE, cmd.Args[2], cmd.Args[3], rm.NullString())
                if val.IsNull() {
                    ctx.ReplyWithNull()
                } else {
                    ctx.ReplyWithString(val)
                }
                return rm.OK
            },
        },
    }
    return mod
}
```

## Fantasy

* A module management module, supplies
    * mod.search
        * Search module from repository(github?)
        * Repository structure like this
        ```
        /namespace
            /module-name
                /bin
                    /darwin_amd64
                        module-name.so
                        module-name.sha
                    /linux_amd64
                module-name.go     
        ```
    * mod.get
        * Download module to ~/.redismodule
        * Because module is write in go, so we can build for almost any platform
        * We can use tag/commit to version the binary, so we can download the old version too
    * mod.install
        * Install downloaded module by calling redis command
    * ...
* A cluster management module
    * Easy to create/manage/monitor redis3 cluster
* A json data type to demonstration how to add new data type in redis.
    * json.fmt key template
    * json.path key path \[pretty]
    * json.get key \[pretty]
    * json.set key value
        * this will validate the json format
