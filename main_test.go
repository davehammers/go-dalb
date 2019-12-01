package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMainStart(t *testing.T) {
	fmt.Println(os.Args)
	// line below commented to prevent printing of debug info during compile time
	os.Args = append(os.Args, "-d")
	os.Args = append(os.Args, "-p")
	os.Args = append(os.Args, "8080")
	mainStart()
	os.Setenv("PORT", "8001")
	mainStart()
}
