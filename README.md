# Go-container

[![GoDoc](https://godoc.org/github.com/ncrypthic/gocontainer?status.svg)](https://godoc.org/github.com/ncrypthic/gocontainer)
[![Go Report Card](https://goreportcard.com/badge/github.com/ncrypthic/gocontainer)](https://goreportcard.com/report/github.com/ncrypthic/gocontainer)

Simple dependency injection service container for golang

## Usage

Suppose we have the following files:

- `app/config/config.go` => Configuration values collection
  ```go
  // app/config/config.go
  package config

  type Config struct {
          UserServiceUrl string
  }
  ```
- `app/user/service.go` => Service to manipulate users through external endpoint

  ```go
  // app/user/service.go
  package user

  import (
          "app/config"
  )

  type Service struct {
          // Will be injected
          Config config.Config `inject:"config"`
  }

  func (svc *Service) GetUserByID(id string) error {
          _, err := http.Get(fmt.Sprintf("%s/%s", svc.Config.SomeUrl, id))
          return err
  }
  ```
- `app/main.go` => application main file

```go
// app/main.go
package main

import (
	"fmt"
        "log"

	"github.com/ncrypthic/gocontainer"

        "app/config"
        "app/user"
)

func main() {
        config := Config{
                UserServiceUrl: "http://example.com/users",
        }
        // no need to manually pass config to user.Service struct
        userService := new(user.Service)
	container := gocontainer.NewContainer()
	container.RegisterService("config", config)
	container.RegisterService("userService", userService)
        // Populate and inject dependencies
	if err := container.Ready(); err != nil {
                log.Fatalf("Failed to populate service container! %v", err)
        }
        // do GET request to http://example.com/users/some-id
	userService.GetUserByID("some-id")
}
```

Service container allow clean intialization file by injecting dependencies to every
services in the container.

## License

MIT License
