/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Taterbro/insighta_cli/internal/api"
	"github.com/Taterbro/insighta_cli/internal/helpers"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// getprofileCmd represents the getprofile command
var getCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a single profile by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runProfilesGet,
}

func init() {
	profilesCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getprofileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getprofileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runProfilesGet(cmd *cobra.Command, args []string) error {
	id := args[0]

	baseURL := api.BackendUrl
	if baseURL == "" {
		return fmt.Errorf("BACKEND_URL not set")
	}

	req, err := http.NewRequest("GET", baseURL+"/profiles/"+id, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-API-Version", "2")

	// loader
	stop := make(chan bool)
	go helpers.StartSpinner(stop)

	resp, err := helpers.WithAuthRetry(req)
	stop <- true

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			ID                 string    `json:"id"`
			Name               string    `json:"name"`
			Gender             string    `json:"gender"`
			GenderProbability  float64   `json:"gender_probability"`
			Age                int       `json:"age"`
			AgeGroup           string    `json:"age_group"`
			CountryID          string    `json:"country_id"`
			CountryName        string    `json:"country_name"`
			CountryProbability float64   `json:"country_probability"`
			CreatedAt          time.Time `json:"created_at"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)

	table.Header([]string{"Field", "Value"})

	table.Append([]string{"ID", result.Data.ID})
	table.Append([]string{"Name", result.Data.Name})
	table.Append([]string{"Gender", result.Data.Gender})
	table.Append([]string{"Gender Probability", fmt.Sprintf("%.2f", result.Data.GenderProbability)})
	table.Append([]string{"Age", fmt.Sprintf("%d", result.Data.Age)})
	table.Append([]string{"Age Group", result.Data.AgeGroup})
	table.Append([]string{"Country ID", result.Data.CountryID})
	table.Append([]string{"Country Name", result.Data.CountryName})
	table.Append([]string{"Country Probability", fmt.Sprintf("%.2f", result.Data.CountryProbability)})
	table.Append([]string{"Created At", result.Data.CreatedAt.Format("2006-01-02 15:04:05")})

	table.Render()

	return nil
}
