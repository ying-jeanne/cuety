package main

import (
	"bytes"
	gojson "encoding/json"
	"export"
	"fmt"
	"io/ioutil"
	"os"

	"cuelang.org/go/cue"
)

func getJSONContent(inputJSONFileName string) map[string]interface{} {
	var result map[string]interface{}
	// Open jsonFile
	jsonFile, err := os.Open(inputJSONFileName)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	gojson.Unmarshal([]byte(byteValue), &result)
	return result
}

func getCueRootValue(inputCUEFileName string) cue.Value {
	spec, _ := ioutil.ReadFile(inputCUEFileName)
	config := string(spec)
	var r cue.Runtime
	rinstance, _ := r.Compile("test", config)
	return rinstance.Value()
}

func writeResultToFile(result *map[string]interface{}, resultFileName string) {
	jsonString, _ := gojson.MarshalIndent(result, "", "    ")
	_ = ioutil.WriteFile(resultFileName, jsonString, 0644)
}

func compareResults(res string, ref string) bool {
	one, err := ioutil.ReadFile(res)
	if err != nil {
		panic(err)
	}
	two, err2 := ioutil.ReadFile(ref)
	if err2 != nil {
		panic(err2)
	}
	return bytes.Equal(one, two)
}

func addDefaultValue(res *map[string]interface{}, rvalue cue.Value) {

}

var (
	refereceFileName  string = "reference.json"
	resultFileName    string = "result.cue"
	inputJSONFileName string = "test.json"
	inputCUEFileName  string = "test.cue"
	inputPackageName  string = "cuety"
)

func main() {
	export.RemoveDefaultValue(inputPackageName, resultFileName)

	//import.FillDefaultValue this one should be equivalent to cue eval file1 file2
	// result := getJSONContent(inputJSONFileName)
	// rvalue := getCueRootValue(inputCUEFileName)
	// fmt.Println(result)
	// writeResultToFile(&result, resultFileName)
	// isEqual := compareResults(resultFileName, refereceFileName)
	// if isEqual {
	// fmt.Println("Test pass")
	// } else {
	// fmt.Println("Test fail")
	// }
	// addDefaultValue(&result, rvalue)
	// writeResultToFile(&result, "defaultback.json")
}
