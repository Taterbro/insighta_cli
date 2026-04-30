/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Taterbro/insighta_cli/internal/api"
	"github.com/spf13/cobra"
)

// whoamiCmd represents the whoami command
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		creds, err := whoAmI()
		if err != nil {
			fmt.Println("whoami error: ", err)
		} else {
			fmt.Println(creds)
		}

	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// whoamiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// whoamiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func whoAmI() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home dir: %w", err)
	}

	path := filepath.Join(home, ".insighta", "credentials.json")

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read credentials file: %w", err)
	}

	var creds api.Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return "", fmt.Errorf("failed to parse credentials: %w", err)
	}

	// Prefer login (GitHub), fallback to username, then name
	if creds.UserDetails.Username != "" {
		return creds.UserDetails.Username, nil
	}

	return "", fmt.Errorf("no username found in credentials")
}
