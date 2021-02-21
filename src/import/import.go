package import

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

/*AddDefaultValue is taking the destination file name as input,
and remove the default value from json input. */
func AddDefaultValue(pkgName string, dst string) {
	b, err := parseArgs(cmd, args, &config{outMode: filetypes.Eval})
	exitOnErr(cmd, err, true)

	syn := []cue.Option{
		cue.Final(), // for backwards compatibility
		cue.Definitions(true),
		cue.Attributes(flagAttributes.Bool(cmd)),
		cue.Optional(flagAll.Bool(cmd) || flagOptional.Bool(cmd)),
	}

	// Keep for legacy reasons. Note that `cue eval` is to be deprecated by
	// `cue` eventually.
	opts := []format.Option{
		format.UseSpaces(4),
		format.TabIndent(false),
	}
	if flagSimplify.Bool(cmd) {
		opts = append(opts, format.Simplify())
	}
	b.encConfig.Format = opts

	e, err := encoding.NewEncoder(b.outFile, b.encConfig)
	exitOnErr(cmd, err, true)

	iter := b.instances()
	defer iter.close()
	for i := 0; iter.scan(); i++ {
		id := ""
		if len(b.insts) > 1 {
			id = iter.id()
		}
		v := iter.instance().Value()

		if flagConcrete.Bool(cmd) {
			syn = append(syn, cue.Concrete(true))
		}
		if flagHidden.Bool(cmd) || flagAll.Bool(cmd) {
			syn = append(syn, cue.Hidden(true))
		}

		errHeader := func() {
			if id != "" {
				fmt.Fprintf(cmd.OutOrStderr(), "// %s\n", id)
			}
		}

		if len(b.expressions) > 1 {
			b, _ := format.Node(b.expressions[i%len(b.expressions)])
			id = string(b)
		}
		if err := v.Err(); err != nil {
			errHeader()
			return err
		}

		// TODO(#553): this can be removed once v.Syntax() below retains line
		// information.
		if (e.IsConcrete() || flagConcrete.Bool(cmd)) && !flagIgnore.Bool(cmd) {
			if err := v.Validate(cue.Concrete(true)); err != nil {
				errHeader()
				exitOnErr(cmd, err, false)
				continue
			}
		}

		f := internal.ToFile(v.Syntax(syn...))
		f.Filename = id
		err := e.EncodeFile(f)
		if err != nil {
			errHeader()
			exitOnErr(cmd, err, false)
		}
	}
	exitOnErr(cmd, iter.err(), true)
}
