package cmd

import (
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/internal/pflagutil"
)

var (
	buildCmd = &cobra.Command{
		Use:     "build FILE",
		Short:   "Build a test case",
		Run:     runBuildCmd,
		Example: "forge build --define ENV=\"prod\" index.mjs",
		Long: `Build a test case bundle

If the reference file has the .mjs file extension, you can import other
JavaScript files and predefine variables using ECMAScript modules.
'forge build' will compile a single JavaScript out of it, resolving the
imports transparently and adding defined variables, if used.

This is also done automatically for you when using the 'forge test-case'
commands.

Imports (ECMAScript modules)
----------------------------
Using 'forge build' allows importing other JavaScript files via the 'import'
statement, if your first files ends in '.mjs':

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
	}

	buildOpts struct {
		Replacements map[string]string
	}
)

func init() {
	RootCmd.AddCommand(buildCmd)

	buildCmd.PersistentFlags().Var(&pflagutil.KeyValueFlag{Map: &buildOpts.Replacements}, "define", "Substitute a list of K=V while parsing: debug=false")
}

func runBuildCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("Missing argument: Entry file")
	}

	bundler := testCaseFileBundler{Replacements: buildOpts.Replacements}
	bundle, err := bundler.Bundle(args[0], "test_case.js")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(os.Stdout, bundle.Content); err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
}
