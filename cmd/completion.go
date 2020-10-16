package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/organisation"
	"github.com/stormforger/cli/api/testcase"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

$ source <(yourprogram completion bash)

# To load completions for each session, execute once:
Linux:
  $ yourprogram completion bash > /etc/bash_completion.d/yourprogram
MacOS:
  $ yourprogram completion bash > /usr/local/etc/bash_completion.d/yourprogram

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ yourprogram completion zsh > "${fpath[1]}/_yourprogram"

# You will need to start a new shell for this setup to take effect.

Fish:

$ yourprogram completion fish | source

# To load completions for each session, execute once:
$ yourprogram completion fish > ~/.config/fish/completions/yourprogram.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}

func init() {
	RootCmd.AddCommand(completionCmd)
}

func completeOrgaAndCase(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if strings.Contains(toComplete, "/") {
		return completionTestCases(toComplete), cobra.ShellCompDirectiveNoFileComp
	}
	return completionOrganisations(toComplete, "/"), cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
}

func completeOrga(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return completionOrganisations(toComplete, ""), cobra.ShellCompDirectiveNoFileComp
}

func completionOrganisations(toComplete, suffix string) []string {
	client := NewClient()

	success, result, err := client.ListOrganisations()
	if err != nil {
		log.Fatal(err)
	}
	if !success {
		log.Fatal(string(result))
	}

	items, err := organisation.Unmarshal(bytes.NewReader(result))
	if err != nil {
		log.Fatal(err)
	}

	out := []string{}
	for _, item := range items.Organisations {
		if toComplete == "" || strings.HasPrefix(item.Name, toComplete) {
			out = append(out, item.Name+suffix)
		}
	}

	return out
}

func completionTestCases(toComplete string) []string {
	client := NewClient()

	x := strings.Split(toComplete, "/")
	orgaName := x[0]
	tcPrefix := x[1]

	orgaUID := lookupOrganisationUID(client, orgaName)

	status, result, err := client.ListTestCases(orgaUID, "")
	if err != nil {
		log.Fatal(err)
	}

	if !status {
		fmt.Fprintln(os.Stderr, "Could not list test cases for "+orgaUID)
		fmt.Fprintln(os.Stderr, string(result))

		os.Exit(1)
	}

	items, err := testcase.Unmarshal(bytes.NewReader(result))
	if err != nil {
		log.Fatal(err)
	}

	out := []string{}
	for _, item := range items.TestCases {
		if strings.HasPrefix(item.Name, tcPrefix) {
			out = append(out, orgaName+"/"+item.Name)
		}
	}

	return out
}
