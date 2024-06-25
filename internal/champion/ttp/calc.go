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
	Type                 string           `mapstructure:"__type"`
	DataValue            string           `mapstructure:"mDataValue,omitempty"`
	Coefficient          float64          `mapstructure:"mCoefficient,omitempty"`
	EndValue             float64          `mapstructure:"mEndValue,omitempty"`
	StartValue           float64          `mapstructure:"mStartValue,omitempty"`
	Breakpoints          []Breakpoint     `mapstructure:"mBreakpoints,omitempty"`
	InitialBonusPerLevel float64          `mapstructure:"mInitialBonusPerLevel,omitempty"`
	Level1Value          float64          `mapstructure:"mLevel1Value,omitempty"`
	Stat                 StatIndex        `mapstructure:"mStat,omitempty"`
	StatFormula          StatFormulaIndex `mapstructure:"mStatFormula,omitempty"`
	EffectIndex          int              `mapstructure:"mEffectIndex,omitempty"`
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
			r := val.FormulaParts.toString(val.DisplayAsPercent, spl, 1)
			n := strings.Replace(ttp.ToString(), fmt.Sprintf("@%s@", key), r, -1)
			*ttp = Tooltip(n)
		} else if val.Type == "GameCalculationModified" {
			r := val.Multiplier.toString(val.ModifiedGameCalculation, spl)
			n := strings.Replace(ttp.ToString(), fmt.Sprintf("@%s@", key), r, -1)
			*ttp = Tooltip(n)
		}
	}
}

func (f FormulaParts) toString(percentage bool, spl SpellDataResource, mult float64) string {
	strs := []string{}

	for _, p := range f {
		switch p.Type {
		case "NamedDataValueCalculationPart":
			strs = append(strs, nameddatavaluecalculationpart(p.DataValue, percentage, spl.DataValues, mult))
		case "StatByNamedDataValueCalculationPart":
			strs = append(strs, statbynameddatavaluecalculationpart(p, spl.DataValues, mult))
		case "StatByCoefficientCalculationPart":
			strs = append(strs, statbycoefficientcalculationpart(p, mult))
		case "ByCharLevelBreakpointsCalculationPart":
			strs = append(strs, bycharlevelbreakpointscalculationpart(p, percentage, mult))
		case "EffectValueCalculationPart":
			strs = append(strs, effectvaluecalculationpart(p, spl.EffectAmount, mult))
		default:
			strs = append(strs, fmt.Sprintf("{{NOT IMPL: %s}}", p.Type))
		}

	}

	return strings.Join(strs, " ")
}

func nameddatavaluecalculationpart(k string, percentage bool, dv SpellValues, mult float64) string {
	for _, val := range dv {
		if k == val.Name {
			return val.Values.toString(percentage, mult)
		}
	}
	return ""
}

func statbynameddatavaluecalculationpart(p FormulaPart, dv SpellValues, mult float64) string {
	for _, val := range dv {
		if p.DataValue == val.Name {
			ratio := val.Values.toString(true, mult)
			formula := p.StatFormula.toString()
			stat := p.Stat.toString()
			return fmt.Sprintf("(+ %s %s %s)", ratio, formula, stat)
		}
	}

	return ""
}

func statbycoefficientcalculationpart(p FormulaPart, mult float64) string {
	ratio := p.Coefficient * 100 * mult
	formula := p.StatFormula.toString()
	stat := p.Stat.toString()

	return fmt.Sprintf("(+ %.2f%% %s %s)", ratio, formula, stat)
}

func bycharlevelbreakpointscalculationpart(p FormulaPart, percentage bool, mult float64) string {
	vals := []float64{}
	lvs := []float64{1}
	pm := float64(1)
	if percentage {
		pm = 100
	}
	base := p.Level1Value * mult * pm

	for i, bp := range p.Breakpoints {
		lvs = append(lvs, float64(bp.Level))

		if i == 0 {
			vals = append(vals, base)
		} else {
			vals = append(vals, vals[i-1]+(bp.AdditionalBonusAtThisLevel*mult*pm))
		}

		if i == len(p.Breakpoints)-1 {
			vals = append(vals, vals[i]+(bp.AdditionalBonusAtThisLevel*mult*pm))
		}
	}

	return fmt.Sprintf("%s (at level %s)", arrayToString(vals, "/", percentage), arrayToString(lvs, "/", false))
}

func arrayToString(a []float64, delim string, percentage bool) string {
	strs := []string{}

	for _, v := range a {
		if percentage {
			strs = append(strs, fmt.Sprintf("%.2f%%", v))
		} else {
			strs = append(strs, fmt.Sprint(v))
		}
	}

	return strings.Join(strs, delim)
}

func effectvaluecalculationpart(p FormulaPart, e SpellEffectAmount, mult float64) string {
	ef := e[p.EffectIndex-1]

	return ef.toString(mult)
}

func (m Multiply) toString(mdf string, spl SpellDataResource) string {
	c := spl.SpellCalculations[mdf]

	if m.Number != 0 {
		return c.FormulaParts.toString(c.DisplayAsPercent, spl, m.Number)
	} else if m.DataValue != "" {
		return multDataVal(m.DataValue, mdf, spl)
	}

	return "{?????}"
}

func multDataVal(dvk, mdf string, spl SpellDataResource) string {
	d := spl.getDataValue(dvk)
	c := spl.SpellCalculations[mdf]

	if d.Values != nil {
		// TODO: Check if first value is sufficient
		first := d.Values[0]
		return c.FormulaParts.toString(c.DisplayAsPercent, spl, float64(first))
	}
	return "{{multDataVal}}}"
}
