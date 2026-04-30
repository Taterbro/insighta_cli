package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Taterbro/insighta_cli/internal/api"
)

func loadCredentials() (*api.Credentials, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, "", err
	}

	path := filepath.Join(home, ".insighta", "credentials.json")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, path, err
	}

	var creds api.Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, path, err
	}

	return &creds, path, nil
}

func refreshAccessToken(refreshToken string) (string, string, error) {
	backendURL := api.BackendUrl

	reqBody, _ := json.Marshal(map[string]string{
		"refresh_token": refreshToken,
	})

	req, err := http.NewRequest("POST", backendURL+"/auth/refresh", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", "", fmt.Errorf("refresh expired")
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("refresh failed: %d", resp.StatusCode)
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	return result.AccessToken, result.RefreshToken, nil
}

func saveCredentials(path string, creds *api.Credentials) error {
	// 1. Load existing file (if it exists)
	existing := &api.Credentials{}

	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, existing)
	}

	// 2. Merge fields (ONLY overwrite what is provided)
	if creds.AccessToken != "" {
		existing.AccessToken = creds.AccessToken
	}

	if creds.RefreshToken != "" {
		existing.RefreshToken = creds.RefreshToken
	}

	// Only overwrite user_details if it's actually provided
	if creds.UserDetails.ID != "" {
		existing.UserDetails = creds.UserDetails
	}

	// 3. Write merged result back
	data, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func WithAuthRetry(req *http.Request) (*http.Response, error) {
	creds, path, err := loadCredentials()
	if err != nil {
		return nil, fmt.Errorf("not logged in")
	}

	// attach access token
	req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
	req.Header.Set("X-API-Version", "2")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// if not 401 → return immediately
	if resp.StatusCode != http.StatusUnauthorized {
		return resp, nil
	}

	resp.Body.Close() // important

	// try refresh
	newAccess, newRefresh, err := refreshAccessToken(creds.RefreshToken)
	if err != nil {
		// refresh failed → force login
		_ = os.Remove(path)
		return nil, fmt.Errorf("session expired, please login again")
	}

	// update credentials
	creds.AccessToken = newAccess
	creds.RefreshToken = newRefresh
	_ = saveCredentials(path, creds)

	// retry original request with new token
	req2 := req.Clone(req.Context())
	req2.Header.Set("Authorization", "Bearer "+newAccess)

	return http.DefaultClient.Do(req2)
}
