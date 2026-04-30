/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/Taterbro/insighta/internal/api"
	"github.com/Taterbro/insighta/internal/helpers"
	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: runProfilesExport,
}

func init() {
	profilesCmd.AddCommand(exportCmd)

	exportCmd.Flags().String("format", "csv", "export format")
	exportCmd.Flags().String("gender", "", "")
	exportCmd.Flags().String("country", "", "")
	exportCmd.Flags().String("age-group", "", "")
	exportCmd.Flags().Int("min-age", 0, "")
	exportCmd.Flags().Int("max-age", 0, "")
	exportCmd.Flags().String("sort-by", "", "")
	exportCmd.Flags().String("order", "", "")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runProfilesExport(cmd *cobra.Command, args []string) error {
	baseURL := api.BackendUrl
	if baseURL == "" {
		return fmt.Errorf("BACKEND_URL not set")
	}

	// query params
	q := url.Values{}

	format, _ := cmd.Flags().GetString("format")
	q.Set("format", format)

	if v, _ := cmd.Flags().GetString("gender"); v != "" {
		q.Set("gender", v)
	}
	if v, _ := cmd.Flags().GetString("country"); v != "" {
		q.Set("country_id", v)
	}
	if v, _ := cmd.Flags().GetString("age-group"); v != "" {
		q.Set("age_group", v)
	}
	if v, _ := cmd.Flags().GetInt("min-age"); v > 0 {
		q.Set("min_age", strconv.Itoa(v))
	}
	if v, _ := cmd.Flags().GetInt("max-age"); v > 0 {
		q.Set("max_age", strconv.Itoa(v))
	}
	if v, _ := cmd.Flags().GetString("sort-by"); v != "" {
		q.Set("sort_by", v)
	}
	if v, _ := cmd.Flags().GetString("order"); v != "" {
		q.Set("order", v)
	}

	// request (POST, but query params)
	req, err := http.NewRequest("POST", baseURL+"/api/profiles/export?"+q.Encode(), nil)
	if err != nil {
		return err
	}

	// spinner
	stop := make(chan bool)
	go helpers.StartSpinner(stop)

	resp, err := helpers.WithAuthRetry(req)
	stop <- true

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("export failed: status %d", resp.StatusCode)
	}

	// get filename from header OR fallback
	contentDisp := resp.Header.Get("Content-Disposition")

	filename := "profiles_export.csv"
	if strings.Contains(contentDisp, "filename=") {
		parts := strings.Split(contentDisp, "filename=")
		if len(parts) > 1 {
			filename = strings.Trim(parts[1], `"`)
		}
	}

	// create file in current directory
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// stream response into file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("\nExport saved to %s\n", filename)

	return nil
}
