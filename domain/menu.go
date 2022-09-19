package domain

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/darkcat013/pr-dining-hall/utils"
)

var Menu []Food

func InitializeMenu(jsonPath string) {

	file, err := os.Open(jsonPath)
	if err != nil {
		utils.Log.Fatal("Error opening " + jsonPath)
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)
	json.Unmarshal(bytes, &Menu)

	if Menu == nil {
		utils.Log.Fatal("Failed to decode menu from " + jsonPath)
	}
	utils.Log.Info("Menu decoded and set")
}
