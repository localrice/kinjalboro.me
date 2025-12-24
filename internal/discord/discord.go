package discord

import (
	"encoding/json"
	"io"
	"net/http"
)

const userID = "749570497242595410"

// response represents the minimal Lanyard API response
// which is required to extract the Discord online status
// see lanyard_example.json for full API response
type response struct {
	Data struct {
		DiscordStatus string `json:"discord_status"`
	} `json:"data"`
}

// GetOnlineStatus returns "online" or "offline"
// as reported by the Lanyard API.
// status: online, idle, dnd default to returning "online"
// If API is unreachable or any type of parsing error,
// it defaults to returning "offline"
func GetOnlineStatus() string {
	resp, err := http.Get("https://lanyard.rest/v1/users/" + userID)

	if err != nil {
		return "offline"
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "offline"
	}

	var responseJSON response
	if err := json.Unmarshal(body, &responseJSON); err != nil {
		return "offline"
	}

	discordStatus := responseJSON.Data.DiscordStatus

	if discordStatus == "offline" {
		return "offline"
	}

	return "online"

}
