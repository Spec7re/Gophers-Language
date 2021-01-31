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

	vowelSet := []byte{'a', 'e', 'i','o', 'u' }

	startVowel := bytes.Contains(vowelSet, []byte(word[:1]))
	secondVowel := bytes.Contains(vowelSet, []byte(word[1:2]))

	var gophWord strings.Builder

	if startVowel {
		gophWord.WriteString("g"+ word)
	} else if word[:2] == "xr" {
		gophWord.WriteString("ge"+ word)
	} else if startVowel == false {
		if !secondVowel && word[1:3] != "qu" {
			gophWord.WriteString(word[2:]+word[:2]+"ogo")
		} else if secondVowel && word[1:3] != "qu" {
			gophWord.WriteString(word[1:]+word[:1]+"ogo")
		} else if word[1:3] == "qu" {
			gophWord.WriteString(word[3:]+word[:3]+"ogo")
		}
	} else {
			gophWord.WriteString(word[1:]+word[:1]+"ogo")
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

	textBytes, err := json.Marshal(map[string]interface{}{"gopher-word": goph})
	if err != nil {
		return
	}

	gopherTranslated := string(textBytes)
	fmt.Println(gopherTranslated)
	fmt.Fprintf( w, "<h3>%s</h3>\n",gopherTranslated )
}

// Separate punctuation from words.
func separateSign(word string) (string, string) {
	sign := word[len(word)-1:]
	word = word[:len(word)-1]
	return word, sign
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

	var gophSen []string
	for i := range sen {
		if strings.IndexAny(sen[i], ".,!?") != -1 {
			word, sign := separateSign(sen[i])
			goph, err := translateWord(word)
			if err != nil {
				fmt.Println(err)
				return
			}
			gophSen = append(gophSen, goph + sign)
		} else {
			goph, err := translateWord(sen[i])
			if err != nil {
				fmt.Println(err)
				return
			}
			gophSen = append(gophSen, goph )
		}
	}

	gS := strings.Join(gophSen, " ")

	textBytes, err := json.Marshal(map[string]interface{}{"gopher-sentence": gS})
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
