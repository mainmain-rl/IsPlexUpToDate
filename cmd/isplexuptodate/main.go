package main

import (
	"fmt"
	"isplexuptodate/internal/discord"
	"isplexuptodate/internal/plex"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	plexURL := os.Getenv("PLEX_URL")                      // ex: https://plex.media.svc.cluster.local:32400
	discordWebhookURL := os.Getenv("DISCORD_WEBHOOK_URL") // ex: https://discord.com/api/webhooks/...

	if plexURL == "" || discordWebhookURL == "" {
		log.Fatal("Missing Env variables.")
	}

	client := &http.Client{Timeout: 10 * time.Second}

	localVersion, machineIdentifier, err := plex.GetLocalPlexVersion(client, plexURL)
	log.Print("Is Plex Up To Date ?")
	log.Printf("Plex URL: %v", plexURL)
	log.Printf("Plex machineIdentifier: %v", machineIdentifier)
	log.Printf("Plex local version: %v", localVersion)
	if err != nil {
		log.Fatalf("Error getting local version: %v", err)
	}

	latestVersion, err := plex.GetLatestPlexVersion(client)
	log.Printf("Plex latest version: %v", latestVersion)
	if err != nil {
		log.Fatalf("Error when getting the latest version : %v", err)
	}

	if localVersion != latestVersion {
		errMsg := fmt.Sprintf("Update available ! Current: %s | Latest: %s", localVersion, latestVersion)
		log.Println(errMsg)
		discord.SendDiscordNotification(client, discordWebhookURL, "Plex Update Available", errMsg)
	} else {
		log.Printf("Plex is up to date ! (version %s)", localVersion)
	}
}
