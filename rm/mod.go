package rm

type Module struct {
    Name        string
    Version     int
    Commands    []Command
    ModuleTypes []ModuleType
    //
    BeforeInit  func(Ctx) error
    AfterInit   func(Ctx) error
}

type Command struct {
    Name     string
    Action   CmdFunc
    Flags    string
    FirstKey int
    LastKey  int
    KeyStep  int
}

func NewMod() *Module {
    return &Module{}
}

// This module will be loaded
var Mod *Module
