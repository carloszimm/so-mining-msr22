package util

import (
	"log"
	"reflect"

	"github.com/carloszimm/stack-mining/types"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetFieldString(v *types.Post, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}
