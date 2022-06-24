package test

type Tester struct {
	Field1 string `wrapper:"getter,setter" json:"field_1,omitempty"`
	Field2 int32  `wrapper:"getter:GetSecondField" json:"field_2,omitempty"`
	Field3 *bool  `json:"field_3"`
}
