package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	jsonMap := make(map[string]interface{})
	raw, err := ioutil.ReadFile("test2.json")
	if err != nil {
		return
	}

	var startingString string = "type GeneratedJsonStruct struct {\n"
	var structString *string = &startingString

	json.Unmarshal(raw, &jsonMap)

	if len(jsonMap) != 0 {
		for key, value := range jsonMap {
			returnType(value, key, structString, false, " ", " ")
		}
	} else {
		var jsonInterface []interface{}
		json.Unmarshal(raw, &jsonInterface)
		if len(jsonInterface) != 0 {
			for _, jsonObject := range jsonInterface {
				returnType(jsonObject, "nil", structString, false, " ", " ")
			}
		}
	}

	indentedStruct := structIndentParser(structString)
	indentedStruct += "\n}"

	fmt.Printf("\n%s\n", indentedStruct)
}

func returnType(a interface{}, key string, str *string, fromArr bool, keySpacing string, jsonSpacing string) {
	keySpacing = strings.Replace(keySpacing, " ", "", len(key))
	switch a.(type) {
	case int:
		jsonSpacing = strings.Replace(jsonSpacing, " ", "", 3)
		*str += fmt.Sprintf("%sint%s`json:\"%s\"`\n", strings.Title(key)+keySpacing, jsonSpacing, key)
	case string:
		jsonSpacing = strings.Replace(jsonSpacing, " ", "", 6)
		*str += fmt.Sprintf("%sstring%s`json:\"%s\"`\n", strings.Title(key)+keySpacing, jsonSpacing, key)
	case bool:
		jsonSpacing = strings.Replace(jsonSpacing, " ", "", 4)
		*str += fmt.Sprintf("%sbool%s`json:\"%s\"`\n", strings.Title(key)+keySpacing, jsonSpacing, key)
	case float64:
		jsonSpacing = strings.Replace(jsonSpacing, " ", "", 7)
		*str += fmt.Sprintf("%sfloat64%s`json:\"%s\"`\n", strings.Title(key)+keySpacing, jsonSpacing, key)
	case map[string]interface{}:
		if len(a.(map[string]interface{})) > 0 {
			keySpacing := getKeySpacing(a.(map[string]interface{}))
			jsonSpacing := getJsonSpacing(a.(map[string]interface{}))
			if strings.Compare(key, "nil") != 0 && !fromArr {
				*str += fmt.Sprintf("%s struct {\n", strings.Title(key))

				for key, value := range a.(map[string]interface{}) {
					returnType(value, key, str, false, keySpacing, jsonSpacing)
				}
				*str += fmt.Sprintf("} `json:\"%s\"`\n", key)
			} else {
				for key, value := range a.(map[string]interface{}) {
					returnType(value, key, str, false, keySpacing, jsonSpacing)
				}
			}
		} else {
			jsonSpacing = strings.Replace(jsonSpacing, " ", "", 8)
			*str += fmt.Sprintf("%sstruct{}%s`json:\"%s\"`\n", strings.Title(key)+keySpacing, jsonSpacing, key)
		}
	case []interface{}:
		if len(a.([]interface{})) > 0 {
			*str += fmt.Sprintf("%s []struct {\n", strings.Title(key))
			returnType(a.([]interface{})[0], key, str, true, " ", " ")
			*str += fmt.Sprintf("} `json:\"%s\"`\n", key)
		} else {
			jsonSpacing = strings.Replace(jsonSpacing, " ", "", 13)
			*str += fmt.Sprintf("%s[]interface{}%s`json:\"%s\"`\n", strings.Title(key)+keySpacing, jsonSpacing, key)
		}
	default:
		jsonSpacing = strings.Replace(jsonSpacing, " ", "", 11)
		*str += fmt.Sprintf("%sinterface{}%s`json:\"%s\"`\n", strings.Title(key)+keySpacing, jsonSpacing, key)
	}
}

func structIndentParser(str *string) string {
	var tabs string
	var indentedStructArr []string
	var indentedStruct string
	for index, char := range *str {
		if char == '{' {
			tabs += "    "
		}

		if char == '}' {
			indentedStructArr[index-1] = strings.Replace(indentedStructArr[index-1], "    ", "", 1)
			tabs = strings.Replace(tabs, "    ", "", 1)
		}

		if char == '\n' {
			indentedStructArr = append(indentedStructArr, "\n"+tabs)
			continue
		}

		indentedStructArr = append(indentedStructArr, string(char))
	}

	indentedStructArr[len(indentedStructArr)-1] = ""

	for _, char := range indentedStructArr {
		indentedStruct += char
	}

	return indentedStruct
}

func getKeySpacing(a map[string]interface{}) string {
	var longestWord int = 0
	for key := range a {
		if len(key) > longestWord {
			longestWord = len(key)
		}
	}

	var spacing string
	for i := 0; i < longestWord+1; i++ {
		spacing += " "
	}

	return spacing
}

func getJsonSpacing(a map[string]interface{}) string {
	var longestType int = 0
	var spacing string
	for _, value := range a {
		switch value.(type) {
		case int:
			if longestType < len("int") {
				longestType = len("int")
			}
		case string:
			if longestType < len("string") {
				longestType = len("string")
			}
		case bool:
			if longestType < len("bool") {
				longestType = len("bool")
			}
		case float64:
			if longestType < len("float64") {
				longestType = len("float64")
			}
		case map[string]interface{}:
			continue
		case []interface{}:
			if len(value.([]interface{})) == 0 {
				longestType = len("[]interface{}")
				break
			}
		default:
			if longestType < len("interface{}") {
				longestType = len("interface{}")
			}
		}
	}

	for i := 0; i < longestType+1; i++ {
		spacing += " "
	}

	return spacing
}
