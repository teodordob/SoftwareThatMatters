/*
Copyright © 2022 A.J.M Brands A.J.M.Brands@student.tudelft.nl

*/
package main

import (
	"runtime/debug"

	"github.com/AJMBrands/SoftwareThatMatters/cmd"
)

func main() {
	debug.SetGCPercent(15)
	cmd.Execute()
}
