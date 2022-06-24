package test

type Tester struct {
	field1 string `wrapper:"-"`
	field2 int32  `wrapper:"getter"`
	field3 *bool
}
