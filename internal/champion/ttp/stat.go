package ttp

type StatIndex int
type StatFormulaIndex int

const (
	AbilityPower                 StatIndex = 0
	Armor                        StatIndex = 1
	Attack                       StatIndex = 2
	AttackSpeed                  StatIndex = 3
	AttackWindupTime             StatIndex = 4
	MagicResist                  StatIndex = 5
	MoveSpeed                    StatIndex = 6
	CritChance                   StatIndex = 7
	CritDamage                   StatIndex = 8
	CooldownReduction            StatIndex = 9
	AbilityHaste                 StatIndex = 10
	MaxHealth                    StatIndex = 11
	CurrentHealth                StatIndex = 12
	PercentMissingHealth         StatIndex = 13
	Unknown14                    StatIndex = 14
	LifeSteal                    StatIndex = 15
	OmniVamp                     StatIndex = 17
	PhysicalVamp                 StatIndex = 18
	MagicPenetrationFlat         StatIndex = 19
	MagicPenetrationPercent      StatIndex = 20
	BonusMagicPenetrationPercent StatIndex = 21
	MagicLethality               StatIndex = 22
	ArmorPenetrationFlat         StatIndex = 23
	ArmorPenetrationPercent      StatIndex = 24
	BonusArmorPenetrationPercent StatIndex = 25
	PhysicalLethality            StatIndex = 26
	Tenacity                     StatIndex = 27
	AttackRange                  StatIndex = 28
	HealthRegenRate              StatIndex = 29
	ResourceRegenRate            StatIndex = 30
	Unknown31                    StatIndex = 31
	Unknown32                    StatIndex = 32
	DodgeChance                  StatIndex = 33
)

const (
	Base     StatFormulaIndex = 0
	Bonus    StatFormulaIndex = 1
	Total    StatFormulaIndex = 2
	Unknown3 StatFormulaIndex = 3
	Unknown4 StatFormulaIndex = 4
)

func (f StatFormulaIndex) toString() string {
	switch f {
	case Base:
		return "Base"
	case Bonus:
		return "Bonus"
	case Total:
		return "Total"
	case Unknown3:
		return "Unknown3"
	case Unknown4:
		return "Unknown4"
	default:
		return ""
	}
}

func (s StatIndex) toString() string {
	switch s {
	case AbilityPower:
		return "AP"
	case Armor:
		return "Armor"
	case Attack:
		return "AD"
	case AttackSpeed:
		return "Attack Speed"
	case AttackWindupTime:
		return "AttackWindupTime"
	case MagicResist:
		return "Magic Resist"
	case MoveSpeed:
		return "Move Speed"
	case CritChance:
		return "Critical Chance"
	case CritDamage:
		return "Critical Damage"
	case CooldownReduction:
		return "Cooldown Reduction"
	case AbilityHaste:
		return "Ability Haste"
	case MaxHealth:
		return "Max Health"
	case CurrentHealth:
		return "Current Health"
	case PercentMissingHealth:
		return "Percent Missing Health"
	case Unknown14:
		return "Unknown14"
	case LifeSteal:
		return "Life Steal"
	case OmniVamp:
		return "OmniVamp"
	case PhysicalVamp:
		return "Physical Vamp"
	case MagicPenetrationFlat:
		return "Magic Penetration Flat"
	case MagicPenetrationPercent:
		return "Magic Penetration Percent"
	case BonusMagicPenetrationPercent:
		return "Bonus Magic Penetration Percent"
	case MagicLethality:
		return "Magic Lethality"
	case ArmorPenetrationFlat:
		return "Armor Penetration Flat"
	case ArmorPenetrationPercent:
		return "Armor Penetration Percent"
	case BonusArmorPenetrationPercent:
		return "Bonus Armor Penetration Percent"
	case PhysicalLethality:
		return "Physical Lethality"
	case Tenacity:
		return "Tenacity"
	case AttackRange:
		return "AttackRange"
	case HealthRegenRate:
		return "HealthRegenRate"
	case ResourceRegenRate:
		return "ResourceRegenRate"
	case Unknown31:
		return "Unknown31"
	case Unknown32:
		return "Unknown32"
	case DodgeChance:
		return "Dodge Chance"
	default:
		return ""
	}
}
