package models

import (
	"encoding/json"
	"os"
)

// CreateGroup создает группу, считывая
// файл в формате JSON
func CreateGroup(source string) (Group, error){
	data, err := os.ReadFile(source)
	if err != nil{
		return Group{}, err
	}
	res := Group{}
	if err := json.Unmarshal(data, &res); err != nil{
		return Group{}, err
	}
	return res, nil
}
