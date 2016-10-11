package main

import (
	"github.com/wenerme/go-rm/rm"
	"github.com/wenerme/go-rm/modules/rxhash"
)

func main() {
	// In case someone try to run this
	rm.Run()
}

func init() {
	rm.Mod = rxhash.CreateModule()
}
