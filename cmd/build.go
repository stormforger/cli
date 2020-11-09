package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/internal/esbundle"
)

var (
	buildCmd = &cobra.Command{
		Use:   "build FILE",
		Short: "Build a test case",
		Run:   runBuildCmd,
		Long: `Build a test case

If the references file has the .mjs file extension, you can import other javascript
files and predefine variables. 'forge build' will compile a single javascript out of
it, resolving the imports transparently.

This will also be automatically done for you when using the 'forge testcase' commands.

Imports
-------

Using 'forge build' allows importing other javascript files via the 'import' statement:

    import helloWorldScenario from "./modules/scenarios.js"
    definition.session("helloworld", helloWorldScenario);

In 'scenarios.js' we have to export the function 'helloWorldScenario':

    function helloWorldScenario(session) {
      session.get("/hello");
    }
    export default helloWorldScenario;

Defines
-------

We use https://esbuild.github.io/ for compiling the various javascript files into a single representation.
Esbuild allows defining variables so your test cases becomes more dynamic.

    const config = {
      env: ENV || "staging",
    }

In this example, configure config.env to the value "staging", if ENV is not defined.
If you pass a define (e.g. 'forge build --define ENV=prod input.mjs') this will now
configure config.env to "prod".
Note that the compiled output no longer contains the fallback to "staging"; Esbuild removed this dead code.

`,
	}

	buildOpts struct {
		Replacements []string
	}
)

func init() {
	RootCmd.AddCommand(buildCmd)

	buildCmd.PersistentFlags().StringArrayVar(&buildOpts.Replacements, "define", []string{}, "Substitute a list of K=V while parsing: debug=false")
}

func runBuildCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("Missing argument: Entry file")
	}

	defines := make(map[string]string)
	for _, kv := range buildOpts.Replacements {
		equals := strings.IndexByte(kv, '=')
		if equals == -1 {
			log.Fatalf("Missing \"=\": %q", kv)
		}

		defines[kv[:equals]] = kv[equals+1:]
	}

	res, err := esbundle.Bundle(args[0], defines)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(res.CompiledContent)
}
