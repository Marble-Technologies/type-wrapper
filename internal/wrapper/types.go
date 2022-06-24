package wrapper

import (
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Package struct {
	*packages.Package
	Dir     string
	Structs []*Struct
}

type Struct struct {
	Name   string
	Fields []*Field
}

type Field struct {
	Type types.Type
	Tag  *Tag
	Name string
}

type Tag struct {
	Getter *string
	Setter *string
}
