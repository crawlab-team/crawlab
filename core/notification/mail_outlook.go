package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getOutlookToken(tenantID, clientID, clientSecret string) (string, error) {
	url := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)
	data := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"scope":         "https://outlook.office365.com/.default",
	}

	formData := ""
	for key, value := range data {
		if formData != "" {
			formData += "&"
		}
		formData += fmt.Sprintf("%s=%s", key, value)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(formData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if token, ok := result["access_token"].(string); ok {
		return token, nil
	}
	return "", fmt.Errorf("no access token found")
}
