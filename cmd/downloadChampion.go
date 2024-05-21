/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

const (
	champURL       = "https://raw.communitydragon.org/%s/game/data/characters/%s/%s.bin.json"
	pbeFileName    = "pbe.json"
	latestFileName = "live.json"
)

var clean bool

var downloadChampionCmd = &cobra.Command{
	Use:   "champion [name]",
	Short: "Download data for a specific champion",
	Long: `Download the latest data for a specific champion from the PBE and Live server.

This command fetches the latest champion data and optionally removes noise from the data.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		champName := strings.ToLower(args[0])

		pbe, err := download(champName, PBE)
		if err != nil {
			log.Fatalf("Failed to Download PBE data: %v", err)
		}

		live, err := download(champName, Latest)
		if err != nil {
			log.Fatalf("Failed to Download Live data: %v", err)
		}

		if clean {
			pbe = removeNoise(pbe)
			live = removeNoise(live)
		}

		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current directory: %v", err)
		}

		dir := fmt.Sprintf("%s/data", wd)

		if err := saveToFile(dir, champName, latestFileName, live); err != nil {
			log.Fatalf("Failed to save file %s: %v", latestFileName, err)
		}

		if err := saveToFile(dir, champName, pbeFileName, pbe); err != nil {
			log.Fatalf("Failed to save file %s: %v", pbeFileName, err)
		}

		fmt.Printf("Success!\n")
	},
	Example: `  # Download data for a specific champion
  pots download champion Sion
  
  # Download data for a specific champion and remove noise
  pots download champion Ahri --clean
  
  # Download data for a specific champion using shorthand flag for clean
  pots download champion Jinx -c`,
}

func init() {
	downloadCmd.AddCommand(downloadChampionCmd)
	downloadChampionCmd.Flags().BoolVarP(&clean, "clean", "c", false, "automatically cleans noise from data")
}

func download(champName string, patch Patch) ([]byte, error) {
	if !validate(patch) {
		return nil, fmt.Errorf("invalid patch: %s", patch)
	}

	url := fmt.Sprintf(champURL, patch, champName, champName)
	res, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from %s: %v", url, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}

func validate(patch Patch) bool {
	switch patch {
	case PBE, Latest:
		return true
	}
	return false
}

func removeNoise(championData []byte) []byte {
	replacements := map[string]string{
		`"mFormat":"\{.*?\}",`:     "",
		`"mAllStartingItemIds":.*`: "}}",
	}

	text := championData

	for pat, rep := range replacements {
		re := regexp.MustCompile(pat)
		text = re.ReplaceAll(text, []byte(rep))
	}

	return text
}

func saveToFile(dir, championName, fileName string, data []byte) error {
	finalDir := fmt.Sprintf("%s/%s", dir, championName)

	if err := os.MkdirAll(finalDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", finalDir, err)
	}

	filePath := fmt.Sprintf("%s/%s", finalDir, fileName)

	if err := os.WriteFile(filePath, data, 0666); err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}

	return nil
}
