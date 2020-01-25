package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type TrainDelayInfo struct {
	Name          string `json:"name"`
	Company       string `json:"company"`
	LastupdateGmt int    `json:"lastupdate_gmt"`
	Source        string `json:"source"`
}
type TrainDelayInfos []TrainDelayInfo

func getInfo() []string {
	var data []string
	url := "https://tetsudo.rti-giken.jp/free/delay.json"
	resp, err := http.Get((url))
	if err != nil {
		fmt.Println("error!")
		return data
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var trainDelayInfos TrainDelayInfos
	err = json.Unmarshal(body, &trainDelayInfos)
	if err != nil {
		log.Fatal(err)
		return data
	}
	for _, train := range trainDelayInfos {
		// if train.Name == "湘南新宿ライン" {
		data = append(data, train.Name)
		// } else if train.Name == "埼京線" {
		// 	data = append(data, train.Name)
		// }
	}
	return data
}

const sendUrl string = "https://hooks.slack.com/services/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"

func main() {
	targetTrainDelayInfo := getInfo()
	var sendText string
	if len(targetTrainDelayInfo) != 0 {
		sendText = strings.Join(targetTrainDelayInfo, ",") + "が遅延しています。"
	} else {
		sendText = "遅延はありません。"
	}
	fmt.Println(sendText)

	params := map[string]interface{}{
		"text":     sendText,
		"userName": "From golang to slack",
	}
	jsonparms, _ := json.Marshal(params)

	resp, _ := http.PostForm(
		sendUrl,
		url.Values{"payload": {string(jsonparms)}},
	)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	println(string(body))
}
