package test

import "sync"

type Tester struct {
	lock   sync.Cond
	field1 string `wrapper:"getter:GetField1,setter"`
	field2 int32  `wrapper:"getter:GetField2,setter"`
	field3 *bool
}
