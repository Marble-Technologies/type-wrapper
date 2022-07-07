package wrapper

import (
	"bytes"
	"fmt"
	"go/types"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/afero"
	"golang.org/x/tools/go/packages"
)

type generator struct {
	writer        *writer
	typ           string
	wrapperType   string
	interfaceName string
	output        string
	receiver      string
	lock          string
	reader        bool
}

type genParameters struct {
	Receiver      string
	Struct        string
	WrapperStruct string
	Interface     string
	Field         string
	GetterMethod  string
	SetterMethod  string
	Type          string
	ZeroValue     string // used only when generating getter
	Lock          string
	Reader        bool
}

func newGenerator(fs afero.Fs, pkg *Package, options ...Option) *generator {
	g := new(generator)
	for _, opt := range options {
		opt(g)
	}

	path := g.outputFilePath(pkg.Dir)
	g.writer = newWriter(fs, path)

	return g
}

// Generate generates a file and wrapper methods.
func Generate(fs afero.Fs, pkg *Package, options ...Option) error {
	g := newGenerator(fs, pkg, options...)

	importMap := make(map[string]*packages.Package, len(pkg.Imports))
	for _, imp := range pkg.Imports {
		// temporary assign nil
		importMap[imp.Name] = imp
	}

	imports := make([]*packages.Package, 0, len(importMap))

	wrappers := make([]string, 0)
	ifaces := make([]string, 0)

	for _, st := range pkg.Structs {
		if st.Name != g.typ {
			continue
		}
		typeParams := g.setupTypeParameters(pkg, st)
		structType, err := g.generateStruct(typeParams)
		if err != nil {
			return err
		}
		wrappers = append(wrappers, structType)

		for _, field := range st.Fields {
			if field.Tag == nil {
				continue
			}

			params := g.setupParameters(pkg, st, field)
			if field.Tag.Getter != nil {
				getter, err := g.generateGetter(params)
				if err != nil {
					return err
				}
				wrappers = append(wrappers, getter)

				iface, err := g.generateGetterInterface(params)
				if err != nil {
					return err
				}
				ifaces = append(ifaces, iface)
			}
			if field.Tag.Setter != nil {
				setter, err := g.generateSetter(params)
				if err != nil {
					return err
				}
				wrappers = append(wrappers, setter)

				iface, err := g.generateSetterInterface(params)
				if err != nil {
					return err
				}
				ifaces = append(ifaces, iface)
			}

			if splitted := strings.Split(strings.TrimPrefix(params.Type, "*"), "."); len(splitted) > 1 {
				otherPackage := splitted[0]
				imports = append(imports, importMap[otherPackage])
			}
		}

		if g.reader {
			readerFunc, err := g.generateReader(typeParams)
			if err != nil {
				return err
			}
			wrappers = append(wrappers, readerFunc)
		}

		iface, err := g.generateInterface(typeParams, ifaces)
		if err != nil {
			return err
		}
		wrappers = append([]string{iface}, wrappers...)
	}

	return g.writer.write(pkg.Name, g.generateImportStrings(imports), wrappers)
}

func (g *generator) outputFilePath(dir string) string {
	output := g.output
	if output == "" {
		// Use snake_case name of type as output file if output file is not specified.
		// type TestStruct will be test_struct_wrapper.go
		var firstCapMatcher = regexp.MustCompile("(.)([A-Z][a-z]+)")
		var articleCapMatcher = regexp.MustCompile("([a-z0-9])([A-Z])")

		name := firstCapMatcher.ReplaceAllString(g.typ, "${1}_${2}")
		name = articleCapMatcher.ReplaceAllString(name, "${1}_${2}")
		output = strings.ToLower(fmt.Sprintf("%s_wrapper.go", name))
	}

	return filepath.Join(dir, output)
}

func (g *generator) generateSetter(
	params *genParameters,
) (string, error) {
	var lockingCode string

	if params.Lock != "" {
		lockingCode = ` {{.Receiver}}.{{.Lock}}.Lock()
		defer {{.Receiver}}.{{.Lock}}.Unlock()
		`
	}

	var tpl = `
	func ({{.Receiver}} {{.WrapperStruct}}) {{.SetterMethod}}(val {{.Type}}) {		
	` +
		lockingCode + // inject locing code
		`{{.Receiver}}.{{.Struct}}.{{.Field}} = val
	}`

	t := template.Must(template.New("setter").Parse(tpl))
	buf := new(bytes.Buffer)

	if err := t.Execute(buf, params); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (g *generator) generateGetter(
	params *genParameters,
) (string, error) {
	var lockingCode string
	if params.Lock != "" {
		lockingCode = `{{.Receiver}}.{{.Lock}}.Lock()
		defer {{.Receiver}}.{{.Lock}}.Unlock()
		`
	}

	var tpl = `
	func ({{.Receiver}} {{.WrapperStruct}}) {{.GetterMethod}}() {{.Type}} {		
		` +
		lockingCode + // inject locing code
		`return {{.Receiver}}.{{.Struct}}.{{.Field}}
	}`

	t := template.Must(template.New("getter").Parse(tpl))
	buf := new(bytes.Buffer)

	if err := t.Execute(buf, params); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (g *generator) generateGetterInterface(
	params *genParameters,
) (string, error) {
	if params.Interface == "" {
		return "", nil
	}
	var tpl = `{{.GetterMethod}}() {{.Type}}
		`

	t := template.Must(template.New("getter-interface").Parse(tpl))
	buf := new(bytes.Buffer)

	if err := t.Execute(buf, params); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (g *generator) generateSetterInterface(
	params *genParameters,
) (string, error) {
	if params.Interface == "" {
		return "", nil
	}
	var tpl = `{{.SetterMethod}}(val {{.Type}})
		`

	t := template.Must(template.New("setter-interface").Parse(tpl))
	buf := new(bytes.Buffer)

	if err := t.Execute(buf, params); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (g *generator) generateReader(
	params *genParameters,
) (string, error) {
	var lockingCode string
	if params.Lock != "" {
		lockingCode = ` {{.Receiver}}.{{.Lock}}.Lock()
		defer {{.Receiver}}.{{.Lock}}.Unlock()
		`
	}

	var tpl = `
	func ({{.Receiver}} {{.WrapperStruct}}) Read(p []byte) (int, error) {		
	` +
		lockingCode + // inject locking code
		`{{.Receiver}}.DataType = "{{.Struct}}"
		data, err := json.Marshal({{.Receiver}}) 
		if err != nil {
			return 0, err
		}
		n := copy(p, data)
		return n, nil
	}`

	t := template.Must(template.New("reader").Parse(tpl))
	buf := new(bytes.Buffer)

	if err := t.Execute(buf, params); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (g *generator) generateStruct(
	params *genParameters,
) (string, error) {
	var datatype string
	if g.reader {
		datatype = `// The name of the original type, it gets initalized when calling Read() function, DO NOT USE IT
		DataType string ` + "`json:\"_data_type,omitempty\"`"
	}
	var tpl = `
	// {{.WrapperStruct}} encapulates the type {{.Struct}} 
	type {{.WrapperStruct}} struct {
		` + datatype + `
		{{.Struct}}		
	}
	`

	t := template.Must(template.New("struct").Parse(tpl))
	buf := new(bytes.Buffer)

	if err := t.Execute(buf, params); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (g *generator) generateInterface(params *genParameters, methods []string) (string, error) {
	if params.Interface == "" {
		return "", nil
	}
	methodsList := strings.Join(methods, "")
	var reader string
	if g.reader {
		reader = `Read(p []byte) (int, error)
		`
	}
	var tpl = `
	type {{.Interface}} interface {` + methodsList + reader + `
	}
	`
	t := template.Must(template.New("struct").Parse(tpl))
	buf := new(bytes.Buffer)

	if err := t.Execute(buf, params); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (g *generator) setupParameters(
	pkg *Package,
	st *Struct,
	field *Field,
) *genParameters {
	typeName := g.typeName(pkg.Types, field.Type)
	getter, setter := g.methodNames(field)
	return &genParameters{
		Receiver:      g.receiverName(st.Name),
		Struct:        st.Name,
		WrapperStruct: g.wrapperType,
		Field:         field.Name,
		GetterMethod:  getter,
		SetterMethod:  setter,
		Type:          typeName,
		ZeroValue:     g.zeroValue(field.Type, typeName),
		Lock:          g.lock,
		Reader:        g.reader,
		Interface:     g.interfaceName,
	}
}

func (g *generator) setupTypeParameters(
	pkg *Package,
	st *Struct,
) *genParameters {
	return &genParameters{
		Receiver:      g.receiverName(st.Name),
		Struct:        st.Name,
		WrapperStruct: g.wrapperType,
		Lock:          g.lock,
		Reader:        g.reader,
		Interface:     g.interfaceName,
	}
}

func (g *generator) receiverName(structName string) string {
	if g.receiver != "" {
		// Do nothing if receiver name specified in args.
		return g.receiver
	}

	// Use the first letter of struct as receiver if receiver name is not specified.
	return strings.ToLower(string(structName[0]))
}

func makeExportable(name string) string {
	return strings.ToUpper(name[0:1]) + name[1:]
}

func (g *generator) methodNames(field *Field) (getter, setter string) {
	if getterName := field.Tag.Getter; getterName != nil && *getterName != "" {
		getter = *getterName
	} else {
		getter = makeExportable(field.Name)
	}

	if setterName := field.Tag.Setter; setterName != nil && *setterName != "" {
		setter = *setterName
	} else {
		setter = "Set" + makeExportable(field.Name)
	}

	return getter, setter
}

func (g *generator) typeName(pkg *types.Package, t types.Type) string {
	return types.TypeString(t, func(p *types.Package) string {
		// type is defined in same package
		if pkg == p {
			return ""
		}
		// path string(like example.com/user/project/package) into slice
		return p.Name()
	})
}

func (g *generator) zeroValue(t types.Type, typeString string) string {
	switch t := t.(type) {
	case *types.Pointer:
		return "nil"
	case *types.Array:
		return "nil"
	case *types.Slice:
		return "nil"
	case *types.Chan:
		return "nil"
	case *types.Interface:
		return "nil"
	case *types.Map:
		return "nil"
	case *types.Signature:
		return "nil"
	case *types.Struct:
		return typeString + "{}"
	case *types.Basic:
		info := types.Typ[t.Kind()].Info()
		switch {
		case types.IsNumeric&info != 0:
			return "0"
		case types.IsBoolean&info != 0:
			return "false"
		case types.IsString&info != 0:
			return `""`
		}
	case *types.Named:
		if types.Identical(t, types.Universe.Lookup("error").Type()) {
			return "nil"
		}

		return g.zeroValue(t.Underlying(), typeString)
	}

	return "nil"
}

func (g *generator) generateImportStrings(pkgs []*packages.Package) []string {
	// Ensure imports are same order as previous if there are no declaration changes.
	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].Name < pkgs[j].Name
	})

	if g.reader {
		var added bool
		for _, pkg := range pkgs {
			if pkg.PkgPath == "encoding/json" {
				added = true
				break
			}
		}
		if !added {
			pkgs = append([]*packages.Package{
				{
					ID:      "encoding/json",
					Name:    "json",
					PkgPath: "encoding/json",
				},
			}, pkgs...)
		}
	}

	imports := make([]string, len(pkgs))
	for i, pkg := range pkgs {
		if pkg.Name == filepath.Base(pkg.PkgPath) {
			imports[i] = pkg.PkgPath
		} else {
			imports[i] = fmt.Sprintf("%s \"%s\"", pkg.Name, pkg.PkgPath)
		}
	}

	return imports
}
