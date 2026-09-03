package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type PlexServerIdentity struct {
	MediaContainer struct {
		Version string `json:"version"`
	} `json:"MediaContainer"`
}

type PlexDownloads struct {
	Computer struct {
		Linux struct {
			Releases []struct {
				Version string `json:"version"`
			} `json:"releases"`
		} `json:"Linux"`
	} `json:"computer"`
}

func main() {

	plexURL := os.Getenv("PLEX_URL")         // ex: https://plex.media.svc.cluster.local:32400
	gatusBaseURL := os.Getenv("GATUS_URL")   // ex: http://gatus.monitoring.svc.cluster.local:8080/api/v1/endpoints/xxxxx/external
	gatusApiKey := os.Getenv("GATUS_APIKEY") // ex: PZUFRZ7Zozryfgou

	if plexURL == "" || gatusBaseURL == "" || gatusApiKey == "" {
		log.Fatal("Missing Env variables.")
	}

	client := &http.Client{Timeout: 10 * time.Second}

	localVersion, err := getLocalPlexVersion(client, plexURL)
	log.Printf("Local version: %v", localVersion)
	if err != nil {
		reportToGatus(client, gatusBaseURL, gatusApiKey, false, err.Error())
		log.Fatalf("Errorr getting local version: %v", err)
	}

	latestVersion, err := getLatestPlexVersion(client)
	log.Printf("Latest version: %v", latestVersion)
	if err != nil {
		reportToGatus(client, gatusBaseURL, gatusApiKey, false, err.Error())
		log.Fatalf("Error when getting the latest version : %v", err)
	}

	if localVersion != latestVersion {
		errMsg := fmt.Sprintf("Update available ! Current: %s | Latest: %s", localVersion, latestVersion)
		log.Println(errMsg)
		reportToGatus(client, gatusBaseURL, gatusApiKey, false, errMsg)
	} else {
		log.Printf("Plex is up to date (version %s)", localVersion)
		reportToGatus(client, gatusBaseURL, gatusApiKey, true, "")
	}
}

func getLocalPlexVersion(client *http.Client, plexURL string) (string, error) {
	req, err := http.NewRequest("GET", plexURL+"/identity", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data PlexServerIdentity
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	return data.MediaContainer.Version, nil
}

func getLatestPlexVersion(client *http.Client) (string, error) {
	// Endpoint officiel Plex public pour obtenir la liste des builds
	resp, err := client.Get("https://plex.tv/api/downloads/5.json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data PlexDownloads
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	releases := data.Computer.Linux.Releases
	if len(releases) == 0 {
		return "", fmt.Errorf("No version found in the Plex API")
	}

	return releases[0].Version, nil
}

func reportToGatus(client *http.Client, baseURL, token string, success bool, errMsg string) {
	bearer := fmt.Sprintf("Bearer %s", token)
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Printf("Erreur URL Gatus : %v", err)
		return
	}

	q := u.Query()
	q.Set("success", fmt.Sprintf("%t", success))
	if errMsg != "" {
		q.Set("error", errMsg)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		log.Printf("Error request Gatus : %v", err)
		return
	}
	req.Header.Set("Authorization", bearer)

	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error http call to Gatus: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		log.Printf("Gatus returned unexpected status: %s", resp.Status)
		return
	}

	log.Printf("Gatus notified (HTTP %d)", resp.StatusCode)
}
