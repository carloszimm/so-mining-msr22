package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/carloszimm/stack-mining/internal/util"
	"github.com/go-playground/validator/v10"
)

var (
	CONSOLIDATED_SOURCES_PATH = filepath.Join(DATA_EXPLORER_PATH, "consolidated sources")
	DATA_EXPLORER_PATH        = filepath.Join("assets", "data explorer")
	LDA_RESULT_PATH           = filepath.Join("assets", "lda-results")
	OPERATORS_PATH            = filepath.Join("assets", "operators")
	OPERATORS_RESULT_PATH     = filepath.Join("assets", "operators-search")
	OPENSORT_RESULT_PATH      = filepath.Join("assets", "opensort")
)

type Config struct {
	FileName         string `json:"fileName" validate:"required"`
	Field            string `json:"field" validate:"required"`
	CombineTitleBody bool   `json:"combineTitleBody"`
	MinTopics        int    `json:"minTopics" validate:"min=0"`
	MaxTopics        int    `json:"maxTopics" validate:"min=0,gtefield=MinTopics"`
	SampleWords      int    `json:"sampleWords" validate:"required"`
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
