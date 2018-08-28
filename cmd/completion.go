package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// CompletionCmd is the command to generate bash/zsh completion
var CompletionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Create completion files for bash and zsh.",
	Long: `Create shell completion code to tab-complete forge's commands.

The generated completion files will be output in the current folder:
- forge-completion-bash.sh
- forge-completion-zsh.sh

Zsh (experimental)
------------------
1. Rename forge-completion-zsh.sh to _forge and move it into a directory
   that is part of your $fpath
2. Re-initialize your completion files:
   rm -f ~/.zcompdump; compinit

Bash
----
To enable bash completion for forge:
1. Edit your ~/.bash_profile and add a line like this:
   source /path/to/forge-completion-bash.sh
2. Start a new terminal session
`,
	Run: generateCompletionFiles,
}

const (
	completionFileNameBash string = "forge-completion-bash.sh"
	completionFileNameZsh  string = "forge-completion-zsh.sh"
)

func init() {
	RootCmd.AddCommand(CompletionCmd)
}

func generateCompletionFiles(cmd *cobra.Command, args []string) {
	log.Printf("Creating completion file for bash in %s\n", completionFileNameBash)
	err := RootCmd.GenBashCompletionFile(completionFileNameBash)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chmod(completionFileNameBash, 0777)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Creating completion file for zsh in %s\n", completionFileNameZsh)
	err = RootCmd.GenZshCompletionFile(completionFileNameZsh)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chmod(completionFileNameZsh, 0777)
	if err != nil {
		log.Fatal(err)
	}
}
