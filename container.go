// Package gocontainer is a dependency injection and service container library for go.
// This library allow you to register service objects using string id, inject them
// into other service objects using an `inject:"some_service_id"` struct tags and
// get the service from the container by its id
package gocontainer

import (
	"github.com/facebookgo/inject"
)

// Service interface
type Service interface {
	StartUp()
	Shutdown()
}

// ContainerAware is an interface to make a service to have reference to the container
type ContainerAware interface {
	SetContainer(c ServiceContainer)
}

// ServiceContainer is an interface for custom service container
type ServiceContainer interface {
	// Lock-up the container for registration
	Ready() error
	// Find a service by specified id in the container and return them
	// Should return false if the service is not found
	GetService(id string) (interface{}, bool)
	// Register a service to the container
	RegisterService(id string, svc interface{})
	// Register multiple services into the container
	// services is a map of service id and service object
	RegisterServices(services map[string]interface{})
}

// ServiceRegistry is default implementation of ServiceContainer interface
type ServiceRegistry struct {
	graph    inject.Graph
	services map[string]interface{}
	order    map[int]string
}

// GetService find a service by its id in ServiceRegistry
func (reg *ServiceRegistry) GetService(id string) (interface{}, bool) {
	svc, ok := reg.services[id]
	return svc, ok
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

// Ready locks up the container to prevent service registration
func (reg *ServiceRegistry) Ready() error {
	err := reg.graph.Populate()
	if err != nil {
		return err
	}
	for i := 0; i < len(reg.order); i++ {
		k := reg.order[i]
		obj := reg.services[k]
		switch s := obj.(type) {
		case Service:
			containerAware, ok := s.(ContainerAware)
			if ok {
				containerAware.SetContainer(reg)
			}
			s.StartUp()
		}
	}
	return nil
}

// NewContainer creates a new empty ServiceRegistry
func NewContainer() *ServiceRegistry {
	return &ServiceRegistry{services: make(map[string]interface{}), order: make(map[int]string)}
}
