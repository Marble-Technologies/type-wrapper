package test

type Tester struct {
	field1 string `wrapper:"setter"`
	field2 int32  `wrapper:"setter:SetSecondField"`
	field3 *bool
}
