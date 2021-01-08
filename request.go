package main

// Use this script to send post request to the gopher API.
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	//TODO loop for test requests
	// url := "http://localhost:8080/word"
	url := "http://localhost:8080/sentence"
	fmt.Println("URL:>", url)

	// var jsonStr = []byte(`{"english-word":"mqueen"}`)
	var jsonStr = []byte(`{"english-sentence":"Apple orange xray, xfactor chair stool square mqeen!"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
