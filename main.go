package main

import "github.com/fljdin/dispatch/internal"

var (
	version string = "*unreleased*"
)

func main() {
	internal.Dispatch(version)
}
