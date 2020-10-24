package data

//Character contains all the data a charcter needs
type Character struct {
	Name          string
	Natrue        string
	Clan          string
	Demeanor      string
	Attributes    attributes
	Abilities     abilities
	Advantages    advantages
	Merits        []string
	Flaw          []string
	Path          uint8
	PermWillpower uint8
	Willpower     uint8
	MaxBloodpool  uint8
	Bloodpool     uint8
	Health        uint8
}

type attributes struct {
	//physical attributes
	Stength   uint8
	Dexterity uint8
	Stamina   uint8
	//social attributes
	Charisma     uint8
	Manipulation uint8
	Appearance   uint8
	//mental atributes
	Perception   uint8
	Intelligence uint8
	Wits         uint8
}

type abilities struct {
	Talents   map[string]uint8
	Skills    map[string]uint8
	Knowledge map[string]uint8
}

type advantages struct {
	Disciplines map[string]uint8
	Backgrounds map[string]uint8
	Virtues     map[string]uint8
}
