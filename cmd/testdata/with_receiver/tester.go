package test

type Tester struct {
	field1 string `wrapper:"getter"`
	field2 int32  `wrapper:"setter:SetSecondField"`
	field3 *bool
}
