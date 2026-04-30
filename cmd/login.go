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
	"time"

	"github.com/Taterbro/insighta/internal/api"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		login()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type PollResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	UserDetails  api.Account `json:"user_details"`
}

func login() {
	url := fmt.Sprintf("%s/auth/github?platform=cli", api.BackendUrl)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Could not contact server")
		return
	}
	defer resp.Body.Close()

	var result struct {
		Status string `json:"status"`
		Data   struct {
			URL   string `json:"url"`
			State string `json:"state"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Bad server response")
		return
	}

	fmt.Println("Opening browser...")
	err = browser.OpenURL(result.Data.URL)
	if err != nil {
		fmt.Println("Could not open browser")
		return
	}

	fmt.Println("Loading; waiting for login...")

	poll(result.Data.State)
}

func poll(state string) {
	pollURL := fmt.Sprintf(
		"%s/auth/callback/poll?state=%s",
		api.BackendUrl,
		state,
	)

	for {
		resp, err := http.Get(pollURL)
		if err != nil {
			fmt.Println("Polling failed...")
			time.Sleep(2 * time.Second)
			continue
		}

		var result struct {
			Status string `json:"status"`
			Data   PollResponse
		}

		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			fmt.Println("Login failed")
			fmt.Println("errror: ", err)
		}
		resp.Body.Close()

		if result.Status == "pending" {
			time.Sleep(3 * time.Second)
			continue
		}

		if result.Status == "success" {
			fmt.Printf("Logged in as %s", result.Data.UserDetails.Username)
			SaveCredentials(result.Data)
			return
		}

		fmt.Println("Login failed")
		return
	}
}
func SaveCredentials(creds PollResponse) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir := filepath.Join(home, ".insighta")
	file := filepath.Join(dir, "credentials.json")

	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0600)
}
