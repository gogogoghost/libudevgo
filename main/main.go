package main

import (
	"fmt"

	"github.com/gogogoghost/udev"
)

func main() {
	udev.SetLibName("libudev.so.1")
	udev.Init()
	ctx, err := udev.NewContext()
	if err != nil {
		panic(err)
	}
	enumerator, err := ctx.NewEnumerator()
	if err != nil {
		panic(err)
	}
	err = enumerator.AddFilter("block")
	if err != nil {
		panic(err)
	}
	for _, dev := range enumerator.List() {
		fmt.Println(dev.Env["DEVNAME"])
	}
}
