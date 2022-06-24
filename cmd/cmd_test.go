package cmd_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/spf13/afero"

	"github.com/Marble-Technologies/type-wrapper/cmd"
)

func TestExecute(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cmd    string
		output string
	}{
		"Getter": {
			cmd:    "type-wrapper -type Tester testdata/getter",
			output: "testdata/getter/tester_wrapper.go",
		},
		"GetterAndInterface": {
			cmd:    "type-wrapper -type Tester -interface ITester testdata/getter",
			output: "testdata/getter/tester_wrapper.go",
		},
		"GetterAndInterfaceAndJson": {
			cmd:    "type-wrapper -type Tester -interface ITester -json testdata/getter",
			output: "testdata/getter/tester_wrapper.go",
		},
		"Setter": {
			cmd:    "type-wrapper -type Tester testdata/setter",
			output: "testdata/setter/tester_wrapper.go",
		},
		"GetterAndSetter": {
			cmd:    "type-wrapper -type Tester testdata/getter_and_setter",
			output: "testdata/getter_and_setter/tester_wrapper.go",
		},
		"IgnoreFields": {
			cmd:    "type-wrapper -type Tester testdata/ignore_fields",
			output: "testdata/ignore_fields/tester_wrapper.go",
		},
		"ImportPackages": {
			cmd:    "type-wrapper -type Tester testdata/import_packages",
			output: "testdata/import_packages/tester_wrapper.go",
		},
		"WithOutput": {
			cmd:    "type-wrapper -type Tester -output my_wrapper.go testdata/with_output",
			output: "testdata/with_output/my_wrapper.go",
		},
		"WithReceiver": {
			cmd:    "type-wrapper -type Tester -receiver tester testdata/with_receiver",
			output: "testdata/with_receiver/tester_wrapper.go",
		},
		"WithReceiverAndInterfaceAndJson": {
			cmd:    "type-wrapper -type Tester -receiver tester -interface ITester -json testdata/with_receiver",
			output: "testdata/with_receiver/tester_wrapper.go",
		},
		"WithLock": {
			cmd:    "type-wrapper -type Tester -lock lock testdata/with_lock",
			output: "testdata/with_lock/tester_wrapper.go",
		},
		"WithLockAndInterfaceAndJson": {
			cmd:    "type-wrapper -type Tester -lock lock -interface ITester -json testdata/with_lock",
			output: "testdata/with_lock/tester_wrapper.go",
		},
	}

	fs := afero.NewMemMapFs()
	snapshot := cupaloy.New(
		cupaloy.SnapshotSubdirectory("testdata/.snapshots"),
		cupaloy.SnapshotFileExtension(".go"),
	)

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			args := strings.Split(tt.cmd, " ")
			cmd.Execute(fs, args)

			output, _ := filepath.Abs(tt.output)

			exists, err := afero.Exists(fs, output)
			if err != nil {
				t.Fatal(err)
			}
			if !exists {
				t.Fatalf("file %s not exists", output)
			}

			file, err := afero.ReadFile(fs, output)
			if err != nil {
				t.Fatal(err)
			}

			snapshot.SnapshotT(t, file)
		})
	}
}
