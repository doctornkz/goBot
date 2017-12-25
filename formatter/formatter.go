package formatter

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/kljensen/snowball"
)

// Message struct
type Message struct {
	Header       string
	Words        []string
	Wordsauthors []string
	Newusers     []string
	Leftusers    []string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var messageDecoded map[string]interface{}

// Formatter function
func Formatter(messageEncoded []byte) string {

	if err := json.Unmarshal(messageEncoded, &messageDecoded); err != nil {
		panic(err)
	}

	message := Message{}
	json.Unmarshal(messageEncoded, &message)

	header := message.Header
	//log.Println("---")
	words := (message.Words)
	wordsauthors := strings.Join(message.Wordsauthors, ", ")
	wordsauthors = fmt.Printf("%-20s", wordsauthors)
	newusers := strings.Join(message.Newusers, ", ")
	leftusers := strings.Join(message.Leftusers, ", ")
	simplified := strings.Join(wordsSimplifier(words), ", ")
	return header + "\r\n" + "BUZZWORDS: " + simplified + "\r\n" + "AUTHORS: " + wordsauthors + "\r\n" + "NEW USERS: " + newusers + "\r\n" + "LEFT USERS: " + leftusers
}

func wordsSimplifier(words []string) []string {

	dict := make(map[string]map[string]int)
	output := make([]string, 0, 50)
	for _, word := range words { // Dict generator
		word = strings.ToLower(word)
		word = regexp.MustCompile(`[a-z]|[@$%&*~#=/_"!?.\ ,:;\-\\+1234567890(){}\[\]]`).ReplaceAllString(word, "") // Ugly but works
		stemmed, err := snowball.Stem(word, "russian", true)
		check(err)
		if dict[stemmed] == nil {
			wordinmap := make(map[string]int)
			wordinmap[word] = 1
			dict[stemmed] = wordinmap
		} else {
			wordinmap, ok := dict[stemmed]
			if !ok {
				wordinmap := make(map[string]int)
				dict[stemmed] = wordinmap
			} else {
				wordinmap[word]++
			}

		}

	}

	for stemmed, forms := range dict {
		max := 0
		maxWord := stemmed
		for word, count := range forms {
			if max < count {
				max = count
				maxWord = word
			}
		}
		output = append(output, maxWord)

	}
	sort.Strings(output)
	return output
}
