/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Taterbro/insighta/internal/api"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := Logout()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Successfully logged out")
		}
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Logout() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %w", err)
	}

	path := filepath.Join(home, ".insighta", "credentials.json")

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read credentials file: %w", err)
	}

	var creds api.Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return fmt.Errorf("failed to parse credentials: %w", err)
	}

	if creds.AccessToken == "" {
		return fmt.Errorf("no access token found")
	}

	backendURL := api.BackendUrl
	if backendURL == "" {
		return fmt.Errorf("BACKEND_URL not set")
	}

	// IMPORTANT: match backend expectation -> Authorization header
	req, err := http.NewRequest("GET", backendURL+"/auth/logout", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("logout request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("logout failed with status: %d", resp.StatusCode)
	}
	if resp.StatusCode == 401 {
		return fmt.Errorf("logout failed with status: %d \nexpired or invalid token", resp.StatusCode)
	}

	// delete local credentials after successful logout
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("logout succeeded but failed to delete credentials: %w", err)
	}

	return nil
}
