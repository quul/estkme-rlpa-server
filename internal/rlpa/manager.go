package rlpa

import (
	"errors"
	"sync"
)

type Manager interface {
	Add(id string, connection *Conn)
	Remove(id string)
	AddCallbackFunc(callbackFunc func(
		callbackType CallbackType, id string, data string,
	))
	HandleCallback(callbackType CallbackType, id string, data string)
	Get(id string) (*Conn, error)
	GetAll() []*Conn
	Len() int
}

const (
	ErrConnNotFound = "connection not found"
)

type manager struct {
	connections sync.Map
	callbacks   []func(callbackType CallbackType, id string, data string)
}

func NewManager() Manager {
	return &manager{}
}

func (m *manager) Add(id string, connection *Conn) {
	m.connections.Store(id, connection)
	m.HandleCallback(CallbackTypeConnAdd, id, "")
}

func (m *manager) Remove(id string) {
	m.connections.Delete(id)
	m.HandleCallback(CallbackTypeConnRemove, id, "")
}

// TODO: Is Callback a good idea? If there's any better way to do this?

func (m *manager) AddCallbackFunc(callbackFunc func(
	callbackType CallbackType, id string, data string,
)) {
	m.callbacks = append(m.callbacks, callbackFunc)
}

func (m *manager) HandleCallback(callbackType CallbackType, id string, data string) {
	for _, callback := range m.callbacks {
		callback(callbackType, id, data)
	}
}

func (m *manager) Get(id string) (*Conn, error) {
	conn, ok := m.connections.Load(id)
	if !ok {
		return nil, errors.New(ErrConnNotFound)
	}
	return conn.(*Conn), nil
}

func (m *manager) GetAll() []*Conn {
	var connections []*Conn
	m.connections.Range(func(_, value interface{}) bool {
		connections = append(connections, value.(*Conn))
		return true
	})
	return connections
}

func (m *manager) Len() int {
	length := 0
	m.connections.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}
