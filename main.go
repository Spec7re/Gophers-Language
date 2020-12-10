package main

import (
	"encoding/json"
	"fmt"
  "strings"
	"net/http"
	"bytes"
	"io/ioutil"
)

type Word struct {
  Word string `json:"english-word"`
}

// Single word translation
func translateWord(word string) (string, error) {

	if len(word) == 0 {
     return "Error", fmt.Errorf("Must have atleast one symbol.")
  }

	if strings.IndexAny(word, "'’") != -1 {
		return "Error", fmt.Errorf("Gophers do not understand short words.")
	}

	word = strings.ToLower(word)

	gophix := "g"
	wordIndex := strings.Index(word, "xr")
	if wordIndex == 0 {
		gophix = "ge"
	} else {
		wordIndex = strings.IndexAny(word, "aeiou")
		if wordIndex >= 2 && word[wordIndex-1:wordIndex+1] == "qu" {
			wordIndex++
		}
	}
	if wordIndex == -1 {
		return "Error", fmt.Errorf( "Gophers need atlest one vowel.")
	}

	var gophWord strings.Builder

	if wordIndex == 0 {
		gophWord.WriteString(gophix)
	}
	gophWord.WriteString(word[wordIndex:len(word)])
	gophWord.WriteString(word[0:wordIndex])
	if wordIndex != 0 {
		gophWord.WriteString("ogo")
	}
	return gophWord.String(), nil
}

// Bootup and report messages
func bootup (w http.ResponseWriter, r *http.Request ) {
	fmt.Fprintf( w, "<h1>%s</h1>\n", "Starting app for translating gophers language..." )
	fmt.Fprintf( w, "<h4>%s</h4>\n", "We currently support two POST endpoints /word and /sentence" )
}

// Single word json encode/decode and translate
func handleWord(w http.ResponseWriter, r *http.Request ) {
	fmt.Fprintf( w, "<h1>%s</h1>\n", "Endpoint for gopher Word translation" )

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		fmt.Println(readErr)
		return
	}

	word := Word{}
	err := json.Unmarshal(body, &word)
	if err != nil {
		fmt.Println(err)
		return
	}

	goph, err := translateWord(word.Word)
	if err != nil {
		fmt.Println(err)
		return
	}

	textBytes, err := json.Marshal(map[string]interface{}{
       "gopher-word": goph})
	if err != nil {
		return
	}

	gopherTranslated := string(textBytes)
	fmt.Println(gopherTranslated)
	fmt.Fprintf( w, "<h3>%s</h3>\n",gopherTranslated )
}


type Sentence struct {
    Sentence string `json:"english-sentence"`
}
// Sentence json decode/encode and translation
func handleSentence(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf( w, "<h1>%s</h1>\n", "Endpoint for gopher Sentence translation" )

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		fmt.Println(readErr)
		return
	}

	sentence := Sentence{}
	err := json.Unmarshal(body, &sentence)
	if err != nil {
		fmt.Println(err)
		return
	}

  sen := strings.Split(sentence.Sentence, " ")

	var gophSen bytes.Buffer
  for i := range sen {
		goph, err := translateWord(sen[i])
		if err != nil {
			fmt.Println(err)
			return
		}
		gophSen.WriteString(goph+" ")
  }

	textBytes, err := json.Marshal(map[string]interface{}{
       "gopher-sentence": gophSen.String()})
	if err != nil {
		return
	}

	gopherTranslated := string(textBytes)
	fmt.Println(gopherTranslated)
	fmt.Fprintf( w, "<h3>%s</h3>\n",gopherTranslated )
}


func startServer(port string) {

	http.HandleFunc("/", bootup)
	http.HandleFunc("/word", handleWord )
	http.HandleFunc("/sentence", handleSentence )

	fmt.Printf("Starting server for gopher language...\n")
	port = fmt.Sprintf(":%s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println(err)
		return
	}

}

var port string = "8080"
func main() {

	startServer(port)
}
