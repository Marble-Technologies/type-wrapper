// Code generated by type-wrapper; DO NOT EDIT.
package test

import (
	"encoding/json"
)

type ITester interface {
	GetField1() string
	SetField1(val string)
	GetField2() int32
	SetField2(val int32)
	Json() []byte
}

// TesterWrapper encapulates the type Tester
type TesterWrapper struct {
	Tester
	// The name of the original type, it gets initalized when calling Json() function, DO NOT USE IT
	DataType string `json:"_data_type,omitempty"`
}

func (t TesterWrapper) GetField1() string {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.Tester.field1
}

func (t TesterWrapper) SetField1(val string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.Tester.field1 = val
}

func (t TesterWrapper) GetField2() int32 {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.Tester.field2
}

func (t TesterWrapper) SetField2(val int32) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.Tester.field2 = val
}

func (t TesterWrapper) Json() []byte {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.DataType = "Tester"
	if data, err := json.Marshal(t); err == nil {
		return data
	}
	return []byte{}
}
