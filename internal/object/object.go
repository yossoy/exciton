package object

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type ObjectKey = string

type ObjectMap struct {
	l       *sync.RWMutex
	objects map[ObjectKey]interface{}
}

func (m *ObjectMap) NewKey() ObjectKey {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return id.String()
}

func (m *ObjectMap) Get(key ObjectKey) interface{} {
	m.l.RLock()
	defer m.l.RUnlock()
	i, ok := m.objects[key]
	if ok {
		return i
	}
	return nil
}

func (m *ObjectMap) Put(key ObjectKey, item interface{}) error {
	m.l.Lock()
	defer m.l.Unlock()
	if _, ok := m.objects[key]; ok {
		return errors.New("key is already exists in ObjectMap")
	}
	m.objects[key] = item
	return nil
}

func (m *ObjectMap) Delete(key ObjectKey) (interface{}, bool, error) {
	m.l.Lock()
	m.l.Unlock()
	obj, ok := m.objects[key]
	if !ok {
		return nil, len(m.objects) == 0, errors.New("key is not exist in ObjectMap")
	}
	delete(m.objects, key)
	return obj, len(m.objects) == 0, nil
}

func NewObjectMap() *ObjectMap {
	return &ObjectMap{
		l:       new(sync.RWMutex),
		objects: make(map[ObjectKey]interface{}),
	}
}

var (
	Apps    = NewObjectMap()
	Windows = NewObjectMap()
	Menus   = NewObjectMap()
)

const SingletonName = "*singleton*"
