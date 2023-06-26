package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	typer "github.com/apooravm/goTyper/typer"
)

type JsonData struct {
	Words []string `json:"words"`
	Phrases []string `json:"phrases"`
	Passage []string `json:"passage"`
}

func main() {
	// Loading in json data
	data, err := ioutil.ReadFile("./data.json")
	if err != nil {
		fmt.Println("Error Reading in data...", err)
		return
	}

	var TextData JsonData

	err = json.Unmarshal(data, &TextData)
	if err != nil {
		fmt.Println("Error decoding json...", err)
		return
	}
	fmt.Println("Words:", TextData.Words)
	fmt.Println("Passage:", TextData.Passage)
	fmt.Println("Phrases:", TextData.Phrases)

	typer.Start(TextData.Phrases[0])
}

