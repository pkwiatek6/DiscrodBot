package data

//RollHistory cotains the data fo player rolls
type RollHistory struct {
	//Rolls holds the last roll made by the character
	Rolls []int `bson:"rolls" json:"rolls"`
	//DC holds the last DC for the last roll
	DC int `bson:"dc" json:"dc"`
	//Reason holds the last reason for the last roll
	Reason string `bson:"reason" json:"reason"`
}
