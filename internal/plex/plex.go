package plex

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func GetLocalPlexVersion(client *http.Client, plexURL string) (string, error) {
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

func GetLatestPlexVersion(client *http.Client) (string, error) {
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
