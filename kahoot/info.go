package kahoot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func AccessToken(email, password string) string {
	//setup request
	client := &http.Client{}
	rawauth := map[string]string{"username": email, "password": password, "grant_type": "password"}
	authentication, err := json.Marshal(rawauth)
	if err != nil {
		panic(err)
	}
	//make POST request
	request, err := http.NewRequest("POST", "https://create.kahoot.it/rest/authenticate", bytes.NewReader(authentication))
	request.Header.Add("content-type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	//decode response
	var values map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&values)
	if err != nil {
		panic(err)
	}
	return values["access_token"].(string)
}
func ReturnData(token, quizid string) map[string]interface{} {
	//setup request
	client := &http.Client{}
	//make request
	request, err := http.NewRequest("GET", fmt.Sprintf("https://create.kahoot.it/rest/kahoots/%s", quizid), nil)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("authorization", token)
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	//decode response
	defer response.Body.Close()
	var data map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		panic(err)
	}
	return data
}
func ParseData(data map[string]interface{}) [][]string {
	//since all data is of interface type, lots of type assertions
	var results [][]string
	colormap := map[int]string{0: "red", 1: "blue", 2: "yellow", 3: "blue"}
	for _, rawvalue := range data["questions"].([]interface{}) {
		var questiondata []string
		value := rawvalue.(map[string]interface{})
		for i, choice := range value["choices"].([]interface{}) {
			choiceinfo := choice.(map[string]interface{})
			correct := choiceinfo["correct"].(bool)
			if correct == true {
				questiondata = append(questiondata, value["question"].(string), choiceinfo["answer"].(string), strconv.Itoa(i), colormap[i])
				break
			}
		}
		results = append(results, questiondata)
	}
	return results
}
