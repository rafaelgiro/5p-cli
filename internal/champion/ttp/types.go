package ttp

type SpellObject struct {
	ScriptName string            `mapstructure:"mScriptName"`
	Spell      SpellDataResource `mapstructure:"mSpell"`
}

type SpellDataResource struct {
	CastRange         []float64                  `mapstructure:"castRange"`
	CooldownTime      []float64                  `mapstructure:"cooldownTime"`
	CastTime          float64                    `mapstructure:"mCastTime"`
	DataValues        SpellValues                `mapstructure:"mDataValues"`
	EffectAmount      SpellEffectAmount          `mapstructure:"mEffectAmount"`
	SpellCalculations map[string]GameCalculation `mapstructure:"mSpellCalculations"`
	Mana              []float64                  `mapstructure:"mana"`
}

type GameCalculation struct {
	Type                    string        `mapstructure:"__type"`
	FormulaParts            []FormulaPart `mapstructure:"mFormulaParts"`
	ModifiedGameCalculation string        `mapstructure:"mModifiedGameCalculation"`
	Multiplier              Multiply      `mapstructure:"mMultiplier"`
}

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
