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
