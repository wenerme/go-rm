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

	BeforeInit func(Ctx, []String) error `json:"-"`
	AfterInit  func(Ctx, []String) error `json:"-"`
	// When module unload
	// ! Very unstable, try to avoid this
	OnUnload   func() `json:"-"`

	// Compilation date
	Compiled   time.Time
	// List of all authors who contributed
	//Authors []Author
	// Copyright of the binary if any
	Copyright  string `json:",omitempty"`
	// Name of Author (Note: Use App.Authors, this is deprecated)
	Author     string `json:",omitempty"`
	Website    string `json:",omitempty"`
	// Email of Author (Note: Use App.Authors, this is deprecated)
	Email      string `json:",omitempty"`
	// Long description for this module
	Desc       string `json:",omitempty"`
}

type Command struct {
	Usage    string
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
