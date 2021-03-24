package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api"
	"github.com/stormforger/cli/api/serviceaccount"
)

var (
	serviceAccountCmd = &cobra.Command{
		Use:     "serviceaccount",
		Aliases: []string{"sa"},
		Short:   "Manage service accounts for automating access to your organisations",
	}

	serviceAccountListCmd = &cobra.Command{
		Use:     "list <organisation-ref>",
		Aliases: []string{"ls"},
		Short:   "List service accounts",
		Args:    cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			client := NewClient()
			org := lookupOrganisationUID(client, args[0])
			if org == "" {
				log.Fatal("<organisation-ref> parameter must be provided")
			}

			if list, err := MainServiceAccountsList(client, org); err != nil {
				printServiceAccountError(os.Stderr, err)
			} else {
				switch rootOpts.OutputFormat {
				case "json":
					json.NewEncoder(os.Stdout).Encode(list)
				default:
					printServiceAccountListHuman(os.Stdout, *list)
				}
			}
		},
	}

	serviceAccountsCreateCmd = &cobra.Command{
		Use:   "create <organisation-ref> <token-label>",
		Short: "Create new service accounts",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewClient()
			org := lookupOrganisationUID(client, args[0])
			token_label := args[1]
			if org == "" {
				log.Fatal("<organisation-ref> parameter must be provided")
			}
			if token_label == "" {
				log.Fatal("<token-label> parameter must be provided.")
			}

			if sa, err := MainServiceAccountsCreate(client, org, token_label); err != nil {
				printServiceAccountError(os.Stderr, err)
			} else {
				switch rootOpts.OutputFormat {
				case "json":
					json.NewEncoder(os.Stdout).Encode(sa)
				default:
					printServiceAccountHuman(os.Stdout, *sa)
				}
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(serviceAccountCmd)
	serviceAccountCmd.AddCommand(serviceAccountListCmd, serviceAccountsCreateCmd)
}

func MainServiceAccountsList(client *api.Client, org string) (*serviceaccount.List, error) {
	ok, data, err := client.ListServiceAccounts(org)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("ListServiceAccounts failed: %v", data)
	}
	list, err := serviceaccount.UnmarshalList(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal api list response: %w", err)
	}
	return &list, err
}

func MainServiceAccountsCreate(client *api.Client, org, token_label string) (*serviceaccount.ServiceAccount, error) {
	ok, data, err := client.CreateServiceAccount(org, token_label)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("CreateServiceAccount failed: %v", data)
	}

	serviceAccount, err := serviceaccount.Unmarshal(strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal api response: %w", err)
	}
	return serviceAccount, nil
}

func printServiceAccountListHuman(w io.Writer, list serviceaccount.List) {
	fmt.Fprintf(w, "Service Accounts (found: %d)\n", len(list.ServiceAccounts))
	for _, item := range list.ServiceAccounts {
		fmt.Fprintf(w, " * %s (uid: %s)\n", item.TokenLabel, item.UID)
	}
}

func printServiceAccountHuman(w io.Writer, sa serviceaccount.ServiceAccount) {
	fmt.Fprintln(w, "Token Label:\t", sa.TokenLabel)
	fmt.Fprintln(w, "UID:\t\t", sa.UID)
	fmt.Fprintln(w, "Generated At:\t", sa.GeneratedAt.Format(time.RFC3339))

	if sa.MostRecentAPIClientVersion != "" {
		fmt.Fprintln(w, "Last API Access:\t", sa.MostRecentAPIAccessAt.Format(time.RFC3339))
		fmt.Fprintln(w, "Client Version:\t", sa.MostRecentAPIClientVersion)
	}

	// AccessToken is only returned on create()
	if sa.AccessToken != "" {
		fmt.Fprintln(w, "Access Token:\t", sa.AccessToken)
	}
}

func printServiceAccountError(w io.Writer, err error) {
	switch rootOpts.OutputFormat {
	case "json":
		json.NewEncoder(os.Stderr).Encode(map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		})
	default:
		log.Fatalf("ERROR: %v", err)
	}
}
