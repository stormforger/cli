package cmd

import (
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/internal/pflagutil"
)

const bundlingHelpInfo = `Bundling
--------
If you use the .mjs file extension, this command will automatically bundle your
JavaScript file using ECMAScript modules. See 'forge test-case build' for more details.
`

var (
	testCaseBuildCmd = &cobra.Command{
		Use:     "build FILE",
		Short:   "Build a test case",
		Example: "forge test-case build --define ENV=\"prod\" index.mjs",
		Long: `Builds a test case bundle from a javascript module file.

If the reference file has the .mjs file extension, you can import other
JavaScript files and predefine variables using ECMAScript modules.
'forge test-case build' will compile a single JavaScript out of it, resolving the
imports transparently and adding defined variables, if used.

This is also done automatically for you when using other 'forge test-case'
commands, so you don't need to call this command directly.

Imports (ECMAScript modules)
----------------------------
Using 'forge test-case build' allows importing other JavaScript files via the 'import'
statement, if the reference file ends in '.mjs':

    import helloWorldScenario from "./modules/scenarios.js"
    definition.session("helloworld", helloWorldScenario);

In 'scenarios.js' we have to export the function 'helloWorldScenario':

    function helloWorldScenario(session) {
      session.get("/hello");
    }
    export default helloWorldScenario;

Defines
-------
We use https://esbuild.github.io/ for compiling the various JavaScript files
into a single representation. Esbuild allows defining variables so your test
cases becomes more dynamic.

    const config = {
      env: ENV || "staging",
    }

In this example, configure config.env to the value "staging", if ENV is not
defined. If you pass a define (e.g. '--define ENV=\"prod\"') this will now
configure config.env to "prod".

To use multiple defines, pass multiple '--define' flags.

A few caveats:

* the compiled output no longer contains the fallback to "staging"; Esbuild
  removed this dead code.
* To use strings as defines, you may need to quote your values twice or escape
  them, otherwise your shell eats them.
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := MainTestcaseBuild(os.Stdout, args[0], buildOpts.Defines); err != nil {
				log.Fatalf("ERROR: %v\n", err)
			}
		},
	}

	buildOpts struct {
		Defines map[string]string
	}
)

func init() {
	TestCaseCmd.AddCommand(testCaseBuildCmd)

	testCaseBuildCmd.PersistentFlags().Var(&pflagutil.KeyValueFlag{Map: &buildOpts.Defines}, "define", "Defines a list of K=V while parsing: debug=false")
}

func MainTestcaseBuild(w io.Writer, file string, defines map[string]string) error {
	bundler := testCaseFileBundler{Defines: defines}
	bundle, err := bundler.Bundle(file, "test_case.js")
	if err != nil {
		return err
	}

	if _, err := io.Copy(w, bundle.Content); err != nil {
		return err
	}
	return nil
}
