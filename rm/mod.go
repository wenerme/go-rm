package rm

import (
    "time"
    "strconv"
)

type Module struct {
    Name       string
    Version    int
    SemVer     string
    Commands   []Command `json:",omitempty"`
    DataTypes  []DataType `json:",omitempty"`

    // TODO add args
    BeforeInit func(Ctx) error `json:"-"`
    AfterInit  func(Ctx) error `json:"-"`
    OnUnload   func()

    // Compilation date
    Compiled   time.Time
    // List of all authors who contributed
    //Authors []Author
    // Copyright of the binary if any
    Copyright  string
    // Name of Author (Note: Use App.Authors, this is deprecated)
    Author     string
    Website    string
    // Email of Author (Note: Use App.Authors, this is deprecated)
    Email      string
    // Long description for this module
    Desc       string
}

type Command struct {
    Desc     string
    Name     string
    Action   CmdFunc `json:"-"`
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

func init() {
    if Mod != nil {
        if Mod.SemVer == "" {
            Mod.SemVer = strconv.Itoa(Mod.Version)
        }
    }
}
