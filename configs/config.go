package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/carloszimm/stack-mining/internal/util"
	"github.com/go-playground/validator/v10"
)

const DataExplorerPath string = "data explorer"
const ResultPath string = "results"

type Config struct {
	Dir         string `json:"dir" validate:"required"`
	Field       string `json:"field" validate:"required"`
	MinTopics   int    `json:"minTopics" validate:"min=0"`
	MaxTopics   int    `json:"maxTopics" validate:"min=0,gtefield=MinTopics"`
	SampleWords int    `json:"sampleWords" validate:"required"`
	FileName    string `json:"fileName" validate:"required"`
}

func ReadConfig() []Config {
	dat, err := os.ReadFile(filepath.Join("configs", "config.json"))
	util.CheckError(err)

	var configs []Config
	json.Unmarshal(dat, &configs)

	validateFields(configs)

	return configs
}

func validateFields(configs []Config) {
	v := validator.New()

	err := v.Var(configs, "required,dive")
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Fatal(e)
		}
	}
}
