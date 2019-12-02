/*
Copyright (c) 2019 Dave Hammers
*/
package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMainControlStart(t *testing.T) {
	commandLineInit()
	fmt.Println(os.Args)
	os.Args = append(os.Args, "-d")
	ctrlPathInit()
}
