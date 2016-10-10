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
            Action: func(ctx rm.Ctx, args rm.CmdArgs) int {
                fmt.Println("Call hgetset")
                //ctx.ReplyWithString("World")
                ctx.ReplyWithNull()
                return rm.OK
            },
        },
    }
    return mod
}
