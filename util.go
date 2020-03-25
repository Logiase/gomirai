package gomirai

import (
	"encoding/json"
)

func Struct2Map(obj interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	j, err := json.Marshal(&obj)
	if err != nil {
		println(err)
	}
	err = json.Unmarshal(j, &m)
	if err != nil {
		println(err)
	}
	return m
}
