package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	charactersURL string = "https://raw.communitydragon.org/json/%s/game/data/characters/"
)

func DownCharacters(p Patch, filter []string) ([]string, error) {
	if !Validate(p) {
		return nil, fmt.Errorf("invalid patch provided: %s", p)
	}

	url := fmt.Sprintf(charactersURL, p)
	res, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch character data from URL %s: %v", url, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body from URL %s: %v", url, err)
	}

	var characters []Character
	if err := json.Unmarshal(body, &characters); err != nil {
		return nil, fmt.Errorf("failed to unmarshal character data from response body: %v", err)
	}

	var list []string

	for _, c := range characters {
		if !contains(filter, c.Name) {
			list = append(list, c.Name)
		}
	}

	return list, nil
}

func contains(f []string, s string) bool {
	for _, a := range f {
		if strings.Contains(s, a) {
			return true
		}
	}
	return false
}
