package test

import (
	"encoding/json"
	"io"
)

type ITester interface {
	Field1() string
	SetField1(val string)
	GetSecondField() int32
	Read(p []byte) (int, error)
}

// TesterWrapper encapulates the type Tester
type TesterWrapper struct {
	// The name of the original type, it gets initalized when calling Read() function, DO NOT USE IT
	DataType string `json:"_data_type,omitempty"`
	Tester
}

func (t TesterWrapper) Field1() string {
	return t.Tester.Field1
}

func (t TesterWrapper) SetField1(val string) {
	t.Tester.Field1 = val
}

func (t TesterWrapper) GetSecondField() int32 
	return t.Tester.Field2
}

func (t TesterWrapper) Read(p []byte) (int, error) {
	t.DataType = "Tester"
	data, err := json.Marshal(t)
	if err != nil {
		return 0, err
	}
	n := copy(p, data)
	return n, nil
}
