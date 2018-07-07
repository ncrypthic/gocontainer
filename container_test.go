package gocontainer

import (
	"fmt"
	"testing"
)

func Example() {
	type Seed struct {
		Name string
	}
	type BaseStruct struct {
		ID string `inject:"seed"`
	}
	type CoolStruct struct {
		Name string
		Base *BaseStruct `inject:"base"`
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
	fmt.Println(base.(*BaseStruct).ID)
	fmt.Println(cool.(*CoolStruct).Name)
	fmt.Println(cool.(*CoolStruct).Base.ID)
	// Output:
	// the seed
	// base deps
	// a cool struct
	// base deps
}

func TestNewContainer(t *testing.T) {
	container := NewContainer()
	switch container.(type) {
	case ServiceContainer:
	default:
		t.Error("Failed to initialize service container")
	}
}

func TestUsage(t *testing.T) {
	Example()
}
