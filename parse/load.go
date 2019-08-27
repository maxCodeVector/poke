package parse

import (
	"encoding/json"
	"io/ioutil"
)

type Game struct {
	Alice  string
	Bob    string
	Result int
}

type Record struct {
	Matches []Game
}


func LoadJsonFile(fileName string, rep int) *[] Game {

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic("path error")
	}
	var record Record
	record.Matches = make([]Game, 0, 10100*rep)
	err = json.Unmarshal(data, &record)
	if err != nil {
		panic("format error")
	}
	if rep > 1{
		tempS := record.Matches
		for x := 1; x < rep; x++ {
			record.Matches = append(record.Matches, tempS...)
		}
	}
	return &record.Matches
}
