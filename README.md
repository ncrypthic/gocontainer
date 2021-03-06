# Go-container

[![GoDoc](https://godoc.org/github.com/ncrypthic/gocontainer?status.svg)](https://godoc.org/github.com/ncrypthic/gocontainer)
[![Go Report Card](https://goreportcard.com/badge/github.com/ncrypthic/gocontainer)](https://goreportcard.com/report/github.com/ncrypthic/gocontainer)

Simple dependency injection service container for golang

## Usage

```go
// app/main.go
package main

import (
        "fmt"
        "log"

        "github.com/ncrypthic/gocontainer"
)

type Config struct {
        UserServiceUrl string
}

type Service struct {
        // Will be injected by service container
        Config config.Config `inject:"config"`
}

func main() {
        config := Config{
                UserServiceUrl: "http://example.com/users",
        }
        // no need to manually pass config to user.Service struct
        userService := new(user.Service)
        container := gocontainer.NewContainer()
        container.RegisterService("config", config)
        container.RegisterService("userService", userService)
        /* To allow service container handling the application process
           exit and graceful shutdown, just uncomment the following line */
        // container.EnableGracefulShutdown(25 * time.Second )

        // Populate and inject dependencies
        if err := container.Ready(); err != nil {
                log.Fatalf("Failed to populate service container! %v", err)
        }
        // http.ListenAndServe(":8080", nil)
}
```

Service container allow clean intialization file by injecting dependencies to every
services in the container.

## License

MIT License
