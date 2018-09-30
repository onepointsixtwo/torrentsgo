package model

import (
	"testing"
)

func TestOrderedMapOverrides(t *testing.T) {
	m := NewOrderedMap()
	m.Add("key", "value")

	value, ok := m.GetExists("key")
	if !ok {
		t.Errorf("Expected to be able to get value for key")
	}
	if value != "value" {
		t.Errorf("Incorrect value found %v", value)
	}

	m.Add("key2", 123)
	intValue := m.Get("key2")
	if intValue != 123 {
		t.Errorf("Error geting int value - incorrect value %v", intValue)
	}
}

func TestOrderedMapIterator(t *testing.T) {
	m := NewOrderedMap()
	m.Add("firstKey", "value1")
	m.Add("secondKey", "value2")
	m.Add("thirdKey", "value3")

	index := 0
	m.Iterate(func(key string, value interface{}) {
		if index == 0 {
			if key != "firstKey" || value != "value1" {
				t.Errorf("Unexpected first iteration: %v, %v", key, value)
			}
		} else if index == 1 {
			if key != "secondKey" || value != "value2" {
				t.Errorf("Unexpected first iteration: %v, %v", key, value)
			}
		} else if index == 2 {
			if key != "thirdKey" || value != "value3" {
				t.Errorf("Unexpected first iteration: %v, %v", key, value)
			}
		}
		index = index + 1
	})
}
