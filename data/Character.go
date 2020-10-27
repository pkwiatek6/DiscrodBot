package data

//Character contains all the data a charcter needs
type Character struct {
	Name         string     `bson:"name" json:"name"`
	Natrue       string     `bson:"natrue" json:"natrue"`
	Clan         string     `bson:"clan" json:"clan"`
	Demeanor     string     `bson:"demeanor" json:"demeanor"`
	Attributes   attributes `bson:"attributes" json:"attributes"`
	Abilities    abilities  `bson:"abilities" json:"abilities"`
	Advantages   advantages `bson:"advantages" json:"advantages"`
	Merits       []string   `bson:"merits" json:"merits"`
	Flaw         []string   `bson:"flaw" json:"flaw"`
	Path         uint8      `bson:"path" json:"path"`
	PermWill     uint8      `bson:"permwill" json:"permwill"`
	Willpower    uint8      `bson:"willpower" json:"willpower"`
	MaxBloodpool uint8      `bson:"maxbloodpool" json:"maxbloodpool"`
	Bloodpool    uint8      `bson:"bloodpool" json:"bloodpool"`
	Health       uint8      `bson:"health" json:"health"`
}

type attributes struct {
	//physical attributes
	Stength   uint8 `bson:"stength" json:"stength"`
	Dexterity uint8 `bson:"dexterity" json:"dexterity"`
	Stamina   uint8 `bson:"stamina" json:"stamina"`
	//social attributes
	Charisma     uint8 `bson:"charisma" json:"charisma"`
	Manipulation uint8 `bson:"manipulation" json:"manipulation"`
	Appearance   uint8 `bson:"appearance" json:"appearance"`
	//mental atributes
	Perception   uint8 `bson:"perception" json:"perception"`
	Intelligence uint8 `bson:"intelligence" json:"intelligence"`
	Wits         uint8 `bson:"wits" json:"wits"`
}

type abilities struct {
	Talents   map[string]uint8 `bson:"Talents" json:"Talents"`
	Skills    map[string]uint8 `bson:"Skills" json:"Skills"`
	Knowledge map[string]uint8 `bson:"Knowledge" json:"Knowledge"`
}

type advantages struct {
	Disciplines map[string]uint8 `bson:"Disciplines" json:"Disciplines"`
	Backgrounds map[string]uint8 `bson:"Backgrounds" json:"Backgrounds"`
	Virtues     map[string]uint8 `bson:"Virtues" json:"Virtues"`
}
