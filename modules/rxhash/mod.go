package rxhash

import (
	"github.com/wenerme/go-rm/rm"
)

var commands []rm.Command
var dataTypes []rm.DataType

func CreateModule() *rm.Module {
	mod := rm.NewMod()
	mod.Name = "rxhash"
	mod.Version = 1
	mod.SemVer = "1.0.1-BETA"
	mod.Author = "wenerme"
	mod.Website = "http://github.com/wenerme/go-rm"
	mod.Desc = `This module will extends redis hash function`
	mod.Commands = commands
	mod.DataTypes = dataTypes
	return mod
}
