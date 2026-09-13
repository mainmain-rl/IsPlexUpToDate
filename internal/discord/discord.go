package discord

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func SendDiscordNotification(client *http.Client, webhookURL, title, message string) {
	// Create Discord webhook payload
	payload := fmt.Sprintf(`{
		"embeds": [
			{
				"title": "%s",
				"description": "%s",
				"color": %s
			}
		]
	}`, title, message, "10079754") // Blue color for Discord embed

	// Send the payload to Discord webhook
	req, err := http.NewRequest("POST", webhookURL, strings.NewReader(payload))
	if err != nil {
		log.Fatalf("Error creating Discord webhook request: %v", err)
	}

	// Set proper headers for Discord webhook
	req.Header.Set("Content-Type", "application/json")

	// If client is nil, use default client
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending Discord notification: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Fatalf("Discord webhook returned unexpected status: %d", resp.StatusCode)
	}

	log.Printf("Discord notification sent successfully")
}
