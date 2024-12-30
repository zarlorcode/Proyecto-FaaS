package config

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

const (
	adminAPIURL = "http://localhost:9180/apisix/admin/routes"
	adminAPIKey = "your_admin_key"
)

func CreateAPISIXRoute(routeID, uri, upstreamURL string) error {
	payload := map[string]interface{}{
		"uri": uri,
		"upstream": map[string]interface{}{
			"type": "roundrobin",
			"nodes": map[string]int{
				upstreamURL: 1,
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", adminAPIURL+"/"+routeID, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", adminAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to create route: %s", resp.Status)
	}
	return nil
}
