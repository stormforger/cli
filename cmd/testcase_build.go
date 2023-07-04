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

## Imports (ECMAScript modules)

You can use the CLI to create a test case that imports other JavaScript code from your local machine using ECMAScript
modules. The CLI uses https://esbuild.github.io/ for compiling the various JavaScript files into a single representation.

    // file: testcase.mjs
    definition.setTarget("http://testapp.loadtest.party");

    definition.setArrivalPhases([
      {
        duration: 5 * 60,
        rate: 1.0,
      },
    ]);

    import helloWorldScenario from "./exported_modules.js";
    definition.session("hello world", helloWorldScenario);

In 'exported_modules.js' we have to export the function 'helloWorldScenario':

    // file: exported_modules.js
    function helloWorldScenario(session) {
      session.get("/hello");
    }

    export default helloWorldScenario;

This example will be automatically processed (because of the .mjs file extension), for example when using the 'test-case create' command:

    forge test-case create my-modular-testcase testcase.mjs

## Defines

In addition to ECMAScript modules, the CLI also supports passing values via the command line (via --define)
to make test cases more dynamic. We use esbuild's define feature here:

> This feature provides a way to replace global identifiers with constant expressions.

To pass a value from the CLI, use '--define name=value' (e.g. '--define users=123' or '--define env=\"prod\"').
The name must be an expression from your test case that you want to adjust and value a valid javascript value.

To default values, we recommend wrapping them in an object that gets defaulted:

    // NOTE: --define works on global identifiers only! To make defines a global identifier,
    //       it MUST NOT be defined via var/let/const.
    defines = {}
    const config = {
        env: defines.ENV || "staging",
    }

In this example, config.env is set to "staging" by default, but by using define (e.g. --define defines.ENV=\"prod\")
config.env is set to "prod". Note the escaped quotes around the value to pass the double quotes to the javascript
environment too - this might not be necessary depending on your shell.

    forge test-case update my-modular-testcase testcase.mjs --define defines.ENV=\"prod\"

*Note*
We recommend wrapping and defaulting every defined parameter you expected as shown above, as it becomes hard to use
to reuse test cases in temporary tests if your test case has more than one or two such parameters.

To use multiple defines, pass multiple --define flags.

A few caveats:

* The compiled output no longer contains the fallback to “staging”; Esbuild removed this dead code
* You can use --define only if you also provide the test case file
* To use strings as defines, you may need to quote your values twice to escape them, otherwise your shell eats them
* If you want defaulting, if no --define is used, you have to define a global identifier first as shown above

## Documentation

You can also find these information in our docs at https://docs.stormforge.io/perftest/guides/advanced-cli-usage/.
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
