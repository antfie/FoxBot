package main

import (
	"encoding/json"
	"fmt"
	"github.com/antfie/FoxBot/utils"
	"io"
	"net/http"

	"github.com/fatih/color"
)

func checkForUpdates() {
	client := &http.Client{}
	url := "https://github.com/antfie/FoxBot/releases/latest"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req) //#nosec G704 -- URL is hardcoded

	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return
	}

	var response map[string]any

	err = json.Unmarshal(body, &response)

	if err != nil {
		return
	}

	latestReleasedVersion, err := utils.StringToFloat(response["tag_name"].(string))

	if err != nil {
		return
	}

	appVersion, err := utils.StringToFloat(AppVersion)

	if err != nil {
		return
	}

	if latestReleasedVersion > appVersion {
		color.HiYellow(fmt.Sprintf("Please upgrade to the latest version (v%s) by visiting %s\n", response["tag_name"], url))
	}
}
