package gocontainer

import (
	"github.com/facebookgo/inject"
	"log"
)

type Service interface {
	StartUp()
	Shutdown()
}

type ContainerAware interface {
	SetContainer(c ServiceContainer)
}

type ServiceContainer interface {
	Ready() error
	GetService(id string) (interface{}, bool)
	RegisterService(id string, svc interface{})
	RegisterServices(services map[string]interface{})
}

type ServiceRegistry struct {
	graph    inject.Graph
	services map[string]interface{}
	order    map[int]string
}

func (reg *ServiceRegistry) GetService(id string) (interface{}, bool) {
	svc, ok := reg.services[id]
	return svc, ok
}

func (reg *ServiceRegistry) RegisterService(id string, svc interface{}) {
	err := reg.graph.Provide(&inject.Object{Name: id, Value: svc, Complete: false})
	if err != nil {
		panic(err.Error())
	}
	reg.order[len(reg.order)] = id
	reg.services[id] = svc
}

func (reg *ServiceRegistry) RegisterServices(services map[string]interface{}) {
	for id, svc := range reg.services {
		reg.RegisterService(id, svc)
	}
}

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

func NewContainer() *ServiceRegistry {
	log.Printf("Initializing service container...")
	return &ServiceRegistry{services: make(map[string]interface{}), order: make(map[int]string)}
}
