package champion

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func handleTooltip(ttp string, spl SpellDataResource) (string, error) {
	c := removeHTMLTags(ttp)
	re := regexp.MustCompile(`@(.*?)@`)
	matches := re.FindAllStringSubmatch(c, -1)

	for _, match := range matches {
		if len(match) > 1 {
			scaling(match[1], &c, spl)
			normal(match[1], &c, spl)
		}
	}

	return c, nil
}

func removeHTMLTags(input string) string {
	// re := regexp.MustCompile(`<.*?>`)
	// result := re.ReplaceAllString(input, "")

	reSpecial := regexp.MustCompile(`@[^@]*?(?:Postfix|Prefix)@`)
	result := reSpecial.ReplaceAllString(input, "")
	return result
}

func normal(o string, ttp *string, spl SpellDataResource) {
	for _, val := range spl.DataValues {
		if val.Name == o {
			strs := make([]string, len(val.Values))
			for i, v := range val.Values {
				strs[i] = fmt.Sprint(v)
			}
			str := strings.Join(strs, "/")
			n := strings.Replace(*ttp, o, str, -1)
			*ttp = n
		}
	}
}

func scaling(o string, ttp *string, spl SpellDataResource) error {
	vre := regexp.MustCompile(`(\D+?)(\d+)`)
	m := vre.FindStringSubmatch(o)

	if len(m) == 3 {
		w := m[1]
		i, err := strconv.Atoi(m[2])
		if err != nil {
			return fmt.Errorf("failed to find index to ability variable. %s; %v", m[1], err)
		}

		if strings.ToLower(w) == "cooldown" {
			str := fmt.Sprint(spl.CooldownTime[i])
			n := strings.Replace(*ttp, o, str, -1)
			*ttp = n
		} else {
			fmt.Println(*ttp)
			for _, val := range spl.DataValues {
				if val.Name == w {
					fmt.Println(val.Values, val.Name, w, i)
					str := fmt.Sprint(val.Values[i])
					n := strings.Replace(*ttp, o, str, -1)
					*ttp = n
				}
			}
		}

	}

	return nil
}
