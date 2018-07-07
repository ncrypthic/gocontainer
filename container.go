// Package gocontainer is a dependency injection and service container library for go.
// This library allow you to register service objects using string id, inject them
// into other service objects using an `inject:"some_service_id"` struct tags and
// get the service from the container by its id
package gocontainer

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/facebookgo/inject"
)

// Service interface
type Service interface {
	// Start up routine, will be executed on container Ready
	StartUp() error
	// Shutdown routine, will be executed on container Shutdown
	Shutdown() error
}

// ContainerAware is an interface to make a service to have reference to the container
type ContainerAware interface {
	// Inject the service container
	SetContainer(c ServiceContainer)
}

// ServiceContainer is an interface for custom service container
type ServiceContainer interface {
	// Lock-up the container for further service registration and populate current
	// dependency trees
	Ready() error
	// Find a service by specified id in the container and return them
	// Should return false if the service is not found
	GetService(id string) (interface{}, bool)
	// Register a service to the container
	RegisterService(id string, svc interface{})
	// Register multiple services into the container
	// services is a map of service id and service object
	RegisterServices(services map[string]interface{})
	// Shutdown all service objects
	Shutdown() error
	// Set shutdown duration
	HandleGracefulShutdown(time.Duration)
}

// ServiceRegistry is default implementation of ServiceContainer interface
type ServiceRegistry struct {
	graph          inject.Graph
	services       map[string]interface{}
	order          map[int]string
	exitOnShutdown bool
	gracefulPeriod time.Duration
}

// GetService find a service by its id in ServiceRegistry
func (reg *ServiceRegistry) GetService(id string) (svc interface{}, ok bool) {
	svc, ok = reg.services[id]
	return
}

// RegisterService adds a service object with specified id
func (reg *ServiceRegistry) RegisterService(id string, svc interface{}) {
	err := reg.graph.Provide(&inject.Object{Name: id, Value: svc, Complete: false})
	if err != nil {
		panic(err.Error())
	}
	reg.order[len(reg.order)] = id
	reg.services[id] = svc
}

// RegisterServices adds multiple services into the container
// services is a map of service id and service object
func (reg *ServiceRegistry) RegisterServices(services map[string]interface{}) {
	for id, svc := range reg.services {
		reg.RegisterService(id, svc)
	}
}

// Ready locks up the container to prevent futher service registration and
// populate current dependency tree
func (reg *ServiceRegistry) Ready() (err error) {
	err = reg.graph.Populate()
	if err != nil {
		return
	}
	if reg.exitOnShutdown {
		defer reg.shutdownHandler()
	}
	for i := 0; i < len(reg.order); i++ {
		k := reg.order[i]
		obj := reg.services[k]
		var startErr error
		switch s := obj.(type) {
		case Service:
			containerAware, ok := s.(ContainerAware)
			if ok {
				containerAware.SetContainer(reg)
			}
			startErr = s.StartUp()
		}
		if startErr != nil {
			fmt.Println("Error starting up service [%s] %v", k, startErr)
		}
	}
	return
}

func (reg *ServiceRegistry) Shutdown() (err error) {
	if reg.exitOnShutdown {
		fmt.Println("Gracefully shutting down the service container")
	}
	for i := 0; i < len(reg.order); i++ {
		k := reg.order[i]
		obj := reg.services[k]
		var shutdownErr error
		switch s := obj.(type) {
		case Service:
			shutdownErr = s.Shutdown()
		}
		if shutdownErr != nil {
			fmt.Printf("Error while shutting down service [%s] %v", k, shutdownErr)
		}
	}
	if reg.exitOnShutdown {
		<-time.After(reg.gracefulPeriod)
		os.Exit(0)
	}
	return
}

func (reg *ServiceRegistry) shutdownHandler() {
	go func() {
		sigchan := make(chan os.Signal, 15)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		reg.Shutdown()
	}()
}

func (reg *ServiceRegistry) HandleGracefulShutdown(d time.Duration) {
	reg.exitOnShutdown = true
	reg.gracefulPeriod = d
}

// NewContainer creates a new empty ServiceRegistry
func NewContainer() (reg ServiceContainer) {
	reg = &ServiceRegistry{services: make(map[string]interface{}), order: make(map[int]string)}
	return
}
