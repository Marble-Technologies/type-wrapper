# Type Wrapper

[![test](https://github.com/Marble-Technologies/type-wrapper/actions/workflows/test.yml/badge.svg)](https://github.com/Marble-Technologies/type-wrapper/actions/workflows/test.yml)
[![release](https://github.com/Marble-Technologies/type-wrapper/actions/workflows/release.yml/badge.svg)](https://github.com/Marble-Technologies/type-wrapper/actions/workflows/release.yml)

type-wrapper is type wrapper generator for [Go programming language](https://golang.org/).

## What is type-wrapper?

type-wrap is a tool that generates:
- Wrapper struct from any structs.
- Getter and Setter methods
- `Read` function to read the JSON `[]byte` of the struct
- Interface for the wrapper struct

Sometimes you might make struct fields unexported in order for values of fields not to be accessed
or modified from anywhere in your codebases, and define getters or setters for values to be handled in a desired way.

But writing wrappers for structs with many fields is time-consuming, but not exciting or creative.

type-wrapper frees you from tedious, monotonous tasks.

## Installation

To get the latest released version

### Go 1.18+

```bash
go install github.com/Marble-Technologies/type-wrapper@latest
```

## Usage

### Declare Struct with `wrapper` Tag

`type-wrapper` generates wrapper methods from defined structs, so you need to declare a struct and fields with `wrapper` tag.

Values for `wrapper` tag is `getter` and `setter`, `getter` is for generating getter method and `setter` is for setter methods.

Here is an example:

```go
type MyStruct struct {
    Field1 string    `wrapper:"getter"`
    field2 *int      `wrapper:"setter"`
    Field3 time.Time `wrapper:"getter,setter"`
}
```

Generated wrapper will be
```go

// MyStructWrapper encapulates the type MyStruct
type MyStructWrapper struct {
	MyStruct
}

func (m MyStructWrapper) Field1() string {
	return m.MyStruct.Field1
}

func (m MyStructWrapper) SetField2(val *int) {
	m.MyStruct.field2 = val
}

func (m MyStructWrapper) Field3() time.Time {
	return m.MyStruct.Field3
}

func (m MyStructWrapper) SetField3(val time.Time) {
	m.MyStruct.Field3 = val
}
```

getter and setter methods won't be generated if wrapper tag isn't specified. But you can explicitly skip generation by using `-` for tag value.

```go
type MyStruct struct {
    ignoredField `accessor:"-"`
}
```

Following to [convention](https://golang.org/doc/effective_go#Getters), setter's name is `Set<FieldName>()` and getter's name is `<FieldName>()` by default, in other words, Set will be put into setter's name and Get will not be put into getter's name.

You can customize names for setter and getter if you want.
```go
type MyStruct struct {
    Field1 string `wrapper:"getter:GetFirstField"`
    Field2 int    `wrapper:"setter:ChangeSecondField"`
}
```

Generated methods will be

```go
type MyStructWrapper struct {
	MyStruct
}

func (m MyStructWrapper) GetFirstField() string {
	return m.MyStruct.Field1
}

func (m MyStructWrapper) ChangeSecondField(val int) {
	m.MyStruct.Field2 = val
}
```

getter and setter methods won't be generated if `wrapper` tag isn't specified.
But you can explicitly skip generation by using `-` for tag value.

```go
type MyStruct struct {
    ignoredField `wrapper:"-"`
}
```

### Generate the `interface` of the Wrapper type
If an interface name provided to wrapper tool, it will generate the interface which the wrapper type implements,

Here is an example ()
```go
type MyStruct struct {
	Field1 string    `wrapper:"getter"`
	Field2 *int      `wrapper:"setter"`
	Field3 time.Time `wrapper:"getter,setter:SetTime"`
}
```

```go
type IStruct interface {
	Field1() string
	SetField2(val *int)
	Field3() time.Time
	SetTime(val time.Time)
}

// MyStructWrapper encapulates the type MyStruct
type MyStructWrapper struct {
	MyStruct
}

func (m MyStructWrapper) Field1() string {
	return m.MyStruct.Field1
}

func (m MyStructWrapper) SetField2(val *int) {
	m.MyStruct.Field2 = val
}

func (m MyStructWrapper) Field3() time.Time {
	return m.MyStruct.Field3
}

func (m MyStructWrapper) SetTime(val time.Time) {
	m.MyStruct.Field3 = val
}
```


### Generate the `Read` function
`type-wrapper` can generate `Read` method which uses `encoding/json` package to marshal the original type
Here is an example ()
```go
type MyStruct struct {
	Field1 string    `wrapper:"getter" json:"name,omitempty"`
	Field2 *int      `wrapper:"setter" json:"value,omitempty"`
	Field3 time.Time `wrapper:"getter,setter:SetTime" json:"time,omitempty"`
}
```

```go
import (
	"encoding/json"
	"time"
)

type IStruct interface {
	Field1() string
	SetField2(val *int)
	Field3() time.Time
	SetTime(val time.Time)
	Read(p []byte) (int, error)
}

// MyStructWrapper encapulates the type MyStruct
type MyStructWrapper struct {
	MyStruct
	// The name of the original type, it gets initalized when calling Json() function, DO NOT USE IT
	DataType string `json:"_data_type,omitempty"`
}

func (m MyStructWrapper) Field1() string {
	return m.MyStruct.Field1
}

func (m MyStructWrapper) SetField2(val *int) {
	m.MyStruct.Field2 = val
}

func (m MyStructWrapper) Field3() time.Time {
	return m.MyStruct.Field3
}

func (m MyStructWrapper) SetTime(val time.Time) {
	m.MyStruct.Field3 = val
}

func (m MyStructWrapper) Read(p []byte) (int, error) {
	m.DataType = "MyStruct"
	data, err := json.Marshal(m)
	if err != nil {
		return 0, err
	}
	n := copy(p, data)
	return n, nil
}
```

### Run `type-wrapper` command

```
$ type-wrapper [flags] source-dir

source-dir
  source-dir is the directory where the definition of the target struct is located.
  If source-dir is not specified, current directory is set as source-dir.

Flags:
  -interface string
        wrapper interface name to be generated
  -reader
        implement io.Reader interface
  -lock string
        lock name
  -output string
        output file name; default <type_name>_wrapper.go
  -receiver string
        receiver name; default first letter of type name
  -type string
        type name; must be set
  -version
        show the version of wrap
  -wrapper string
        wrapper type name; default <type_name>Wrapper
```

Example:

```shell
$ type-wrapper -type MyStruct -wrapper WStruct -interface IStruct -reader -receiver myStruct -output my_struct_wrapper.go path/to/target
```

#### go generate

You can also generate wrappers by using `go generate`.

```go
package mypackage

//go:generate type-wrapper -type MyStruct -wrapper WStruct -interface IStruct -reader -receiver myStruct -output my_struct_wrapper.go 

type MyStruct struct {
    field1 string `wrapper:"getter"`
    field2 *int   `wrapper:"setter"`
}
```

Then run go generate for your package.

## Credits
This project has been inspired by [accessory](https://github.com/masaushi/accessory) project and it uses most of its source code.

We acknowledge and are grateful to [masaushi](https://github.com/masaushi) and [yumm007](https://github.com/yumm007) for their great work

## License
The `type-wrapper` project (and all code) is licensed under the [MIT License](LICENSE).


