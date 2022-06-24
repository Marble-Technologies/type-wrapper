package test

import (
	"time"

	"github.com/Marble-Technologies/type-wrapper/cmd/testdata/import_packages/sub1"
	sub "github.com/Marble-Technologies/type-wrapper/cmd/testdata/import_packages/sub2"
	. "github.com/Marble-Technologies/type-wrapper/cmd/testdata/import_packages/sub3"
	_ "github.com/Marble-Technologies/type-wrapper/cmd/testdata/import_packages/sub4"
)

type Tester struct {
	field1 time.Time       `wrapper:"getter,setter"`
	field2 *sub1.SubTester `wrapper:"getter,setter"`
	field3 *sub.SubTester  `wrapper:"getter,setter"`
	field4 *SubTester      `wrapper:"getter,setter"`
	field5 *bool
}
