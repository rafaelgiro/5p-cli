package ttp

import (
	"fmt"
	"strings"
)

type SpellCalc map[string]GameCalculation

type GameCalculation struct {
	Type                    string       `mapstructure:"__type"`
	DisplayAsPercent        bool         `mapstructure:"mDisplayAsPercent"`
	FormulaParts            FormulaParts `mapstructure:"mFormulaParts"`
	ModifiedGameCalculation string       `mapstructure:"mModifiedGameCalculation"`
	Multiplier              Multiply     `mapstructure:"mMultiplier"`
}

type FormulaParts []FormulaPart

type FormulaPart struct {
	Type                 string       `mapstructure:"__type"`
	DataValue            string       `mapstructure:"mDataValue,omitempty"`
	Coefficient          float64      `mapstructure:"mCoefficient,omitempty"`
	EndValue             float64      `mapstructure:"mEndValue,omitempty"`
	StartValue           float64      `mapstructure:"mStartValue,omitempty"`
	Breakpoints          []Breakpoint `mapstructure:"mBreakpoints,omitempty"`
	InitialBonusPerLevel float64      `mapstructure:"mInitialBonusPerLevel,omitempty"`
	Level1Value          float64      `mapstructure:"mLevel1Value,omitempty"`
}

type Breakpoint struct {
	Type                       string  `mapstructure:"__type"`
	AdditionalBonusAtThisLevel float64 `mapstructure:"mAdditionalBonusAtThisLevel,omitempty"`
	BonusPerLevelAtAndAfter    float64 `mapstructure:"mBonusPerLevelAtAndAfter,omitempty"`
	Level                      int     `mapstructure:"mLevel"`
}

type Multiply struct {
	Type      string  `mapstructure:"__type"`
	Number    float64 `mapstructure:"mNumber,omitempty"`
	DataValue string  `mapstructure:"mDataValue,omitempty"`
}

func (sc SpellCalc) toTooltip(ttp *Tooltip, spl SpellDataResource) {
	for key, val := range sc {
		if val.Type == "GameCalculation" && val.FormulaParts != nil {
			r := val.FormulaParts.toString(val.DisplayAsPercent, spl)
			fmt.Println(key, r)
			n := strings.Replace(ttp.ToString(), fmt.Sprintf("@%s@", key), r, -1)
			*ttp = Tooltip(n)
		}
	}
}

func (f FormulaParts) toString(percentage bool, spl SpellDataResource) string {
	strs := []string{}

	for _, p := range f {
		switch p.Type {
		case "NamedDataValueCalculationPart":
			str := nameddatavaluecalculationpart(p.DataValue, percentage, spl.DataValues)
			strs = append(strs, str)

		default:
			strs = append(strs, fmt.Sprintf("{{NOT IMPL: %s}}", p.Type))
		}

	}

	return strings.Join(strs, " ")
}

func nameddatavaluecalculationpart(k string, percentage bool, dv SpellValues) string {
	for _, val := range dv {
		if k == val.Name {
			strValues := make([]string, len(val.Values))
			for i, v := range val.Values {
				if percentage {
					strValues[i] = fmt.Sprintf("%.1f%%", v*100)
				} else {

					strValues[i] = fmt.Sprint(v)
				}
			}

			return strings.Join(strValues, "/")
		}
	}
	return fmt.Sprintf("NOT IMPL {{%s}}", k)
}
