# Go-container

Simple dependency injection service container for golang

## Usage

```go
package main

import (
	"fmt"
	"github.com/ncrypthic/gocontainer"
)

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

func main() {
	container := gocontainer.NewContainer()
	container.RegisterService("cool", &CoolStruct{Name: "a cool struct"})
	container.RegisterService("seed", &Seed{"the seed"})
	container.RegisterService("base", &BaseStruct{"base deps"})
	container.Ready() // Populate dependency tree
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
```
