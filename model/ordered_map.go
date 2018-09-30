package model

// Types

type OrderedMap struct {
	internalMap map[string]interface{}
	keyOrder    []string
}

type MapIterator func(key string, value interface{})

// Initialiser

func NewOrderedMap() *OrderedMap {
	internalMap := make(map[string]interface{})
	keyOrder := make([]string, 0)
	return &OrderedMap{internalMap, keyOrder}
}

// Public Methods

func (m *OrderedMap) Add(key string, value interface{}) {
	m.keyOrder = append(m.keyOrder, key)
	m.internalMap[key] = value
}

func (m *OrderedMap) Get(key string) interface{} {
	return m.internalMap[key]
}

func (m *OrderedMap) GetExists(key string) (interface{}, bool) {
	val, ok := m.internalMap[key]
	return val, ok
}

func (m *OrderedMap) Iterate(it MapIterator) {
	length := len(m.keyOrder)
	for i := 0; i < length; i++ {
		key := m.keyOrder[i]
		value := m.internalMap[key]
		it(key, value)
	}
}
