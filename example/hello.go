package main

import (
    "github.com/wenerme/go-rm/rm"
    "fmt"
    "github.com/urfave/cli"
)

func main() {
    fmt.Println("Enter main")
    app := cli.NewApp()
    _ = app
}

func init() {
    fmt.Println("Init main package")
    rm.Mod = CreateMyMod()
}
func CreateMyMod() *rm.Module {
    mod := rm.NewMod()
    mod.Name = "HelloWorld"
    mod.Version = 1

    mod.Commands = []rm.Command{
        {
            Name:   "hsetget",
            Flags:  "write fast deny-oom",
            FirstKey:1, LastKey:1, KeyStep:1,
            Action: func(ctx rm.CmdContext) int {
                fmt.Println("Call hgetset")
                ctx.Ctx.ReplyWithNull()
                return rm.OK
            },
        },
    }
    return mod
}
