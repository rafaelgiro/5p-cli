package ttp

type Tooltip string

type SpellObject struct {
	ScriptName string            `mapstructure:"mScriptName"`
	Spell      SpellDataResource `mapstructure:"mSpell"`
}

type SpellDataResource struct {
	CastRange         []float64         `mapstructure:"castRange"`
	CooldownTime      Cooldowns         `mapstructure:"cooldownTime"`
	CastTime          float64           `mapstructure:"mCastTime"`
	DataValues        SpellValues       `mapstructure:"mDataValues"`
	EffectAmount      SpellEffectAmount `mapstructure:"mEffectAmount"`
	SpellCalculations SpellCalc         `mapstructure:"mSpellCalculations"`
	Mana              Costs             `mapstructure:"mana"`
}

func (ttp *Tooltip) Calc(spl SpellDataResource) error {
	initialCleanup(ttp)

	spl.DataValues.toTooltip(ttp)
	spl.EffectAmount.toTooltip(ttp)
	spl.Mana.toTooltip(ttp)
	spl.CooldownTime.toTooltip(ttp)
	spl.SpellCalculations.toTooltip(ttp, spl)

	finalCleanup(ttp)
	return nil
}

func (ttp Tooltip) ToString() string {
	tp := string(ttp)
	return tp
}

// func spellCalculations(ttp *string, spl SpellDataResource) {
// 	for key, val := range spl.SpellCalculations {
// 		if val.FormulaParts != nil {
// 			formulaParts(ttp, key, val.FormulaParts, spl.DataValues)
// 		}
// 		if val.ModifiedGameCalculation != "" {
// 			if val.Multiplier.Number != 0 {
// 				// TODO
// 				// multiplierNumber(ttp, key, val.ModifiedGameCalculation, val)
// 			}
// 			if val.Multiplier.DataValue != "" {
// 				// TODO
// 				// multiplierData(ttp, key, val.ModifiedGameCalculation, spl, val.Multiplier)
// 			}
// 		}

// 	}
// }

// func formulaParts(ttp *string, key string, fps []FormulaPart, spl []SpellDataValue) {
// 	for _, fp := range fps {
// 		if len(fp.DataValue) != 0 {
// 			dataValCalc(ttp, key, fp.DataValue, spl)
// 		} else if fp.Breakpoints != nil {
// 			breakpoints(ttp, key, fp)
// 		}
// 	}
// }

// func dataValCalc(ttp *string, ttpKey, dvKey string, spl []SpellDataValue) {
// 	for _, val := range spl {
// 		if dvKey == val.Name {
// 			strValues := make([]string, len(val.Values))
// 			for i, v := range val.Values {
// 				strValues[i] = fmt.Sprint(v)
// 			}

// 			str := strings.Join(strValues, "/")
// 			n := strings.Replace(*ttp, fmt.Sprintf("@%s@", ttpKey), str, -1)
// 			*ttp = n
// 		}
// 	}
// }

// func breakpoints(ttp *string, ttpKey string, fp FormulaPart) {
// 	base := fp.Level1Value
// 	vals := []float64{}
// 	lvs := []float64{}

// 	for i, bp := range fp.Breakpoints {
// 		lvs = append(lvs, float64(bp.Level))

// 		if i == 0 {
// 			vals = append(vals, base+bp.AdditionalBonusAtThisLevel)
// 		} else {
// 			vals = append(vals, vals[i-1]+bp.AdditionalBonusAtThisLevel)
// 		}
// 	}

// 	str := fmt.Sprintf("%s (at level %s)", arrayToString(vals, "/"), arrayToString(lvs, "/"))
// 	n := strings.Replace(*ttp, fmt.Sprintf("@%s@", ttpKey), str, -1)
// 	*ttp = n
// }

// func arrayToString(a []float64, delim string) string {
// 	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
// }
