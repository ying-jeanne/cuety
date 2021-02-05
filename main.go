package main

import (
	"bytes"
	gojson "encoding/json"
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

func removeDefaultValue(result *map[string]interface{}, rvalue cue.Value) {
	for key, value := range *result {
		lvField, _ := rvalue.FieldByName(key, false)
		lv := lvField.Value
		if lv.Exists() {
			if defaultv, isdefault := lv.Default(); isdefault {
				// make sure the datatype to cast to
				interValue := value
				_ = defaultv.Decode(&interValue)
				if interValue == value {
					delete(*result, key)
				}
			}
		}
	}
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
	// value := &Sum{B: 4, C: 8}

	// r := &cue.Runtime{}
	// codec := gocodec.New(r, nil)
	// v, err := codec.ExtractType(value)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// _ = codec.Complete(v, value)

	// fmt.Println(res, rvalue)

	// fmt.Printf("%+v", v)
	// out = &bytes.Buffer{}
	// d := json.NewDecoder(nil, tcname, strings.NewReader(tcin))
	// for {
	// 	east, err := d.Extract()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	toString(out, east, err)
	// 	if err != nil {
	// 		break
	// 	}
	// }
	// fmt.Println(out.String())
}

var (
	refereceFileName  string = "reference.json"
	resultFileName    string = "result.json"
	inputJSONFileName string = "test.json"
	inputCUEFileName  string = "test.cue"
)

func main() {
	result := getJSONContent(inputJSONFileName)
	rvalue := getCueRootValue(inputCUEFileName)
	removeDefaultValue(&result, rvalue)
	writeResultToFile(&result, resultFileName)
	isEqual := compareResults(resultFileName, refereceFileName)
	if isEqual {
		fmt.Println("Test pass")
	} else {
		fmt.Println("Test fail")
	}
	addDefaultValue(&result, rvalue)
	writeResultToFile(&result, "defaultback.json")
}
