package gocontainer

import (
	"fmt"
	_ "testing"
)

func ExampleUsage() {
	type Seed struct {
		Name string
	}
	type BaseStruct struct {
		Id string `inject:"seed"`
	}
	type CoolStruct struct {
		Name string
		Base BaseStruct `inject:"base"`
	}
	container := NewContainer()
	container.RegisterService("cool", &CoolStruct{Name: "a cool struct"})
	container.RegisterService("seed", &Seed{"the seed"})
	container.RegisterService("base", &BaseStruct{"base deps"})
	container.Ready()
	seed, _ := container.GetService("seed")
	base, _ := container.GetService("base")
	cool, _ := container.GetService("cool")
	fmt.Println(seed.(*Seed).Name)
	fmt.Println(base.(*BaseStruct).Id)
	fmt.Println(cool.(*CoolStruct).Name)
	// Output:
	// the seed
	// base deps
	// a cool struct
}
