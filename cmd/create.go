package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Taterbro/insighta/internal/api"
	"github.com/Taterbro/insighta/internal/helpers"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user profile",
	RunE:  runProfilesCreate,
}

func init() {
	profilesCmd.AddCommand(createCmd)

	createCmd.Flags().String("name", "", "Name of the user")
	_ = createCmd.MarkFlagRequired("name")
}

func runProfilesCreate(cmd *cobra.Command, args []string) error {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	baseURL := api.BackendUrl
	if baseURL == "" {
		return fmt.Errorf("BACKEND_URL not set")
	}

	// request body
	body, _ := json.Marshal(map[string]string{
		"name": name,
	})

	req, err := http.NewRequest("POST", baseURL+"/api/profiles", bytes.NewBuffer(body))
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

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		var errRes map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&errRes)

		return fmt.Errorf("request failed: %v", errRes)
	}

	// decode response
	var result struct {
		Data struct {
			ID                 string  `json:"id"`
			Name               string  `json:"name"`
			Gender             string  `json:"gender"`
			GenderProbability  float64 `json:"gender_probability"`
			Age                int     `json:"age"`
			AgeGroup           string  `json:"age_group"`
			CountryID          string  `json:"country_id"`
			CountryName        string  `json:"country_name"`
			CountryProbability float64 `json:"country_probability"`
			CreatedAt          string  `json:"created_at"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// table output
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Field", "Value"})

	table.Append([]string{"ID", result.Data.ID})
	table.Append([]string{"Name", result.Data.Name})
	table.Append([]string{"Gender", result.Data.Gender})
	table.Append([]string{"Gender Probability", fmt.Sprintf("%.2f", result.Data.GenderProbability)})
	table.Append([]string{"Age", strconv.Itoa(result.Data.Age)})
	table.Append([]string{"Age Group", result.Data.AgeGroup})
	table.Append([]string{"Country", result.Data.CountryName})
	table.Append([]string{"Created At", result.Data.CreatedAt})

	fmt.Println()
	table.Render()

	return nil
}
