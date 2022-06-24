package test

type Tester struct {
	field1 string `wrapper:"getter,setter"`
	field2 int32  `wrapper:"getter:GetSecondField,setter:SetSecondField"`
	field3 *bool
}
