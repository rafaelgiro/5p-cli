package champion

import (
	"regexp"
)

type Champion struct {
	Name string
}

type Strings struct {
	Entries map[string]string `json:"entries"`
}

const (
	champURL   = "https://raw.communitydragon.org/%s/game/data/characters/%s/%s.bin.json"
	stringsURL = "https://raw.communitydragon.org/%s/game/en_us/data/menu/en_us/main.stringtable.json"
)

var FilteredCharacters = []string{"tft", "tutorial", "trinket", "bw_", "cherry_", "durian_", "ha_", "hexgate", "item_", "kingporo", "nexus", "pet", "slime_", "sru_", "ultbook", "sruap_", "srx", "test", "practicetool_", "preseason_", "spellbook", "sr_infernal", "summonerbeacon"}

func RemoveNoise(data []byte) []byte {
	replacements := map[string]string{
		`"mFormat":"\{.*?}`:              `"mFormat":"{loveusion}`,
		`,"mAllStartingItemIds":.*`:      "}}",
		`,"mAllRecommendableItemIds":.*`: "}}",
		`"EventToTrack":.*?,`:            `"EventToTrack": 0,`,
		`"searchTagsSecondary":".*?"`:    `"searchTagsSecondary":""`,
	}

	for pat, rep := range replacements {
		re := regexp.MustCompile(pat)
		data = re.ReplaceAll(data, []byte(rep))
	}

	return data
}
