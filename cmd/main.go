package main

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
var dst = "output.cue"

func main() {
	_, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}

	// loadedInstances := load.Instances([]string{os.Args[1]}, nil)
	// Haven't figured everything out yet about how they're related to the
	// working directory, the argument here, the name of a package given,
	// imports, etc. This is more "quick n' dumb" for now - just enough to
	// test out the contained test file.
	overlay := map[string]load.Source{}
	loadedInstances := load.Instances([]string{"."}, &load.Config{Package: "cuety", Overlay: overlay})
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
