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

type Patch string

const (
	PBE    Patch = "pbe"
	Latest Patch = "latest"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Please provide a champion name.")
		} else if len(args) > 1 {
			log.Fatal("Please provide only a single champion name.")
		}

		champName := strings.ToLower(args[0])

		pbe := removeNoise(download(champName, PBE))
		latest := removeNoise(download(champName, Latest))

		os.WriteFile("temp/pbe.json", pbe, 0666)
		os.WriteFile("temp/latest.json", latest, 0666)

	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func download(champName string, patch Patch) []byte {
	if !validate(patch) {
		log.Fatalf("Invalid patch: %s", patch)
	}

	res, err := http.Get(fmt.Sprintf("https://raw.communitydragon.org/%s/game/data/characters/%s/%s.bin.json", patch, champName, champName))

	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	return body
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
