package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/spf13/afero"

	"github.com/Marble-Technologies/type-wrapper/internal/wrapper"
)

// Version is the version of `type-wrapper`, injected at build time.
var Version = ""

// newUsage returns a function to replace default usage function of FlagSet.
func newUsage(flags *flag.FlagSet) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "Usage of type-wrapper:\n")
		fmt.Fprintf(os.Stderr, "\ttype-wrapper [flags] [directory]\n")
		fmt.Fprintf(os.Stderr, "For more information, see:\n")
		fmt.Fprintf(os.Stderr, "\thttps://github.com/Marble-Technologies/type-wrapper\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flags.PrintDefaults()
	}
}

// Execute executes a whole process of generating wrapper codes.
func Execute(fs afero.Fs, args []string) {
	log.SetFlags(0 | log.Lshortfile)
	log.SetPrefix("type-wrapper: ")

	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.Usage = newUsage(flags)
	version := flags.Bool("version", false, "show the version of wrap")
	reader := flags.Bool("reader", false, "implement io.Reader interface")
	typeName := flags.String("type", "", "type name; must be set")
	wrapperTypeName := flags.String("wrapper", "", "wrapper type name; default <type_name>Wrapper")
	interfaceName := flags.String("interface", "", "wrapper interface name to be generated")
	lockName := flags.String("lock", "", "lock name")
	receiver := flags.String("receiver", "", "receiver name; default first letter of type name")
	output := flags.String("output", "", "output file name; default <type_name>_wrapper.go")

	if err := flags.Parse(args[1:]); err != nil {
		flags.Usage()
		os.Exit(1)
	}

	if *version {
		fmt.Fprintf(os.Stdout, "type-wrapper version: %s\n", getVersion())
		os.Exit(0)
	}

	if typeName == nil || len(*typeName) == 0 {
		flags.Usage()
		os.Exit(1)
	}
	if wrapperTypeName == nil || len(*wrapperTypeName) == 0 {
		*wrapperTypeName = *typeName + "Wrapper"
	}

	var dir string
	if cliArgs := flags.Args(); len(cliArgs) > 0 {
		dir = cliArgs[0]
	} else {
		// Default: process whole package in current directory.
		dir = "."
	}

	if !isDir(dir) {
		fmt.Fprintln(os.Stderr, "Specified argument is not a directory.")
		flags.Usage()
		os.Exit(1)
	}

	pkg, err := wrapper.ParsePackage(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		flags.Usage()
		os.Exit(1)
	}

	var options = []wrapper.Option{
		wrapper.Type(*typeName),
		wrapper.Output(*output),
		wrapper.Receiver(*receiver),
		wrapper.Lock(*lockName),
		wrapper.Reader(*reader),
		wrapper.Wrapper(*wrapperTypeName),
		wrapper.Interface(*interfaceName),
	}
	if err = wrapper.Generate(fs, pkg, options...); err != nil {
		log.Fatal("err", err)
	}
}

func isDir(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

func getVersion() string {
	if Version != "" {
		return Version
	}

	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" {
		return "unknown"
	}

	return info.Main.Version
}
