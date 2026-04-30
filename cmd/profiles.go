/*
Copyright © 2026 NAME HERE <Age ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/Taterbro/insighta/internal/api"
	"github.com/Taterbro/insighta/internal/helpers"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// profilesCmd represents the profiles command
var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runProfilesList(cmd)
	},
}

func init() {
	rootCmd.AddCommand(profilesCmd)
	profilesCmd.AddCommand(listCmd)
	profilesCmd.AddCommand(profilesSearchCmd)
	listCmd.Flags().String("gender", "", "")
	listCmd.Flags().String("country", "", "")
	listCmd.Flags().String("age-group", "", "")
	listCmd.Flags().Int("min-age", 0, "")
	listCmd.Flags().Int("max-age", 0, "")
	listCmd.Flags().String("sort-by", "", "")
	listCmd.Flags().String("order", "", "")
	listCmd.Flags().Int("page", 1, "")
	listCmd.Flags().Int("limit", 10, "")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// profilesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// profilesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runProfilesList(cmd *cobra.Command) error {
	baseURL := api.BackendUrl
	if baseURL == "" {
		return fmt.Errorf("BACKEND_URL not set")
	}

	// build query params
	q := url.Values{}

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
	if v, _ := cmd.Flags().GetInt("page"); v > 0 {
		q.Set("page", strconv.Itoa(v))
	}
	if v, _ := cmd.Flags().GetInt("limit"); v > 0 {
		q.Set("limit", strconv.Itoa(v))
	}

	// request setup
	req, err := http.NewRequest("GET", baseURL+"/api/profiles?"+q.Encode(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-API-Version", "2")

	// spinner
	stop := make(chan bool)
	go helpers.StartSpinner(stop)

	resp, err := helpers.WithAuthRetry(req)
	stop <- true

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Printf("request failed with status code %d", resp.StatusCode)
		return nil
	}
	// decode response
	var result struct {
		Data []struct {
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

	// render table
	// render table
	table := tablewriter.NewWriter(os.Stdout)

	table.Header([]string{
		"ID",
		"Name",
		"Gender",
		"Gender Probability",
		"Age",
		"Age Group",
		"Country ID",
		"Country Name",
		"Country Probability",
		"Created At",
	})

	for _, u := range result.Data {
		table.Append([]string{
			u.ID,
			u.Name,
			u.Gender,
			fmt.Sprintf("%.2f", u.GenderProbability),
			strconv.Itoa(u.Age),
			u.AgeGroup,
			u.CountryID,
			u.CountryName,
			fmt.Sprintf("%.2f", u.CountryProbability),
			u.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	fmt.Println()
	table.Render()

	return nil
}

var profilesSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search profiles using natural language",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		q := url.Values{}
		q.Set("q", args[0])

		baseURL := api.BackendUrl

		req, _ := http.NewRequest("GET", baseURL+"/api/profiles/search?"+q.Encode(), nil)
		req.Header.Set("X-API-Version", "2")

		resp, err := helpers.WithAuthRetry(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Printf("request failed with status code %d", resp.StatusCode)
			return nil
		}

		var result struct {
			Data []struct {
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

		table.Header([]string{
			"ID",
			"Name",
			"Gender",
			"Gender Probability",
			"Age",
			"Age Group",
			"Country ID",
			"Country Name",
			"Country Probability",
			"Created At",
		})

		for _, u := range result.Data {
			table.Append([]string{
				u.ID,
				u.Name,
				u.Gender,
				fmt.Sprintf("%.2f", u.GenderProbability),
				strconv.Itoa(u.Age),
				u.AgeGroup,
				u.CountryID,
				u.CountryName,
				fmt.Sprintf("%.2f", u.CountryProbability),
				u.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
		fmt.Println()
		table.Render()

		return nil
	},
}
