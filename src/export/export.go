package export

import (
	"fmt"
	"io/ioutil"
	"os"

	"diff"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/tools/trim"
)

var forcetrim = false
var simplify = false

/*RemoveDefaultValue is taking the destination file name as input,
and remove the default value from json input. */
func RemoveDefaultValue(pkgName string, dst string) {
	overlay := map[string]load.Source{}
	loadedInstances := load.Instances([]string{"."}, &load.Config{Package: pkgName, Overlay: overlay})
	instances := cue.Build(loadedInstances)

	for i, inst := range loadedInstances {
		root := instances[i]
		err := trim.Files(inst.Files, root, &trim.Config{
			Trace: true,
		})
		if err != nil {
			fmt.Println(err)
		}

		for _, f := range inst.Files {
			overlay[f.Filename] = load.FromFile(f)
		}

	}

	tinsts := instances
	if len(tinsts) != len(loadedInstances) {
		fmt.Println("unexpected number of new instances")
	}

	if forcetrim {
		for i, p := range instances {
			k, script := diff.Final.Diff(p.Value(), tinsts[i].Value())
			if k != diff.Identity {
				diff.Print(os.Stdout, script)
				fmt.Println("Aborting trim, output differs after trimming. This is a bug! Use -i to force trim.")
				fmt.Println("You can file a bug here: https://github.com/cuelang/cue/issues/new?assignees=&labels=NeedsInvestigation&template=bug_report.md&title=")
				os.Exit(1)
			}
		}
	}

	for _, inst := range loadedInstances {
		for _, f := range inst.Files {
			filename := f.Filename

			opts := []format.Option{}
			if simplify {
				opts = append(opts, format.Simplify())
			}

			b, err := format.Node(f, opts...)
			if err != nil {
				fmt.Errorf("error formatting file: %v", err)
			}

			if dst != "" {
				filename = dst
			}

			err = ioutil.WriteFile(filename, b, 0644)
			if err != nil {
				fmt.Errorf("writing file failed: ", err)
			}
		}
	}
}
