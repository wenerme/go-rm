package rm

type Module struct {
    Name     string
    Version  int
    Commands []Command
}

type Command struct {
    Name   string
    Action CmdFunc
}

func NewMod() *Module {
    return &Module{}
}

// This module will be loaded
var Mod *Module
