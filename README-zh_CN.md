# Go Redis module
go-rm 旨在通过 Golang 来实现 redis 模块.

## 演示

```bash
# 确保你安装了最新的 redis, 不是最新版,而是 github 上的最新编译版本
# 如果你能够使用 brew 命令,那么你可以通过以下方式来安装
# brew reinstall redis --HEAD

# 构建 redis 模块
go get -u -v -buildmode=c-shared github.com/wenerme/go-rm/modules/hashex

# 启动 redis-server 并加载刚刚编译的模块,使用 debug 日志级别
redis-server --loadmodule ~/go/pkg/*/github.com/wenerme/go-rm/modules/hashex* --loglevel debug
```

__客户端__

```
$ redis-cli hsetget a a 5
(nil)
$ redis-cli hsetget a a 4
"5"
$ redis-cli hsetget a a 3
"3"
```

## 如何实现一个 Redis 模块

实现一个 Redis 模块非常简单,就像是写一个 cli 程序一样,以下代码实现了上面演示的功能,源代码在[这里](https://github.com/wenerme/go-rm/blob/master/modules/hashex/hashex.go).

```go
package main

import "github.com/wenerme/go-rm/rm"

func main() {
    // 避免改代码被运行
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

## 幻想

* 实现一个用于管理模块的命令,提供下述命令
    * mod.search
        * 从仓库(github?)搜索模块
        * 仓库的结构类似于这样
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
        * 下载模块到 ~/.redismodule
        * 因为模块是用 Go 写的,因此大多数平台都能使用
        * 可以使用 tag 或者是提交 id 来标识版本
    * mod.install
        * 调用 redis 的命令来安装模块
    * ...
* 集群管理模块
    * 用于简化 redis 3 的集群 创建/管理/监控
* 实现一个 json 数据类型,用于演示如果添加新的 redis 类型,支持以下命令
    * json.fmt key template
    * json.path key path \[pretty]
    * json.get key \[pretty]
    * json.set key value
        * 该操作会验证 value 是否为 json
