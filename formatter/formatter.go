package formatter

import (
	"encoding/json"
	"regexp"
	"sort"
	"strings"

	"github.com/kljensen/snowball"
)

// DigestMessage struct
type DigestMessage struct {
	Header       string
	Words        []string
	Wordsauthors []string
	Newusers     []string
	Leftusers    []string
}

// StatusMessage structure to formatter package
type StatusMessage struct {
	Header  string
	TopList string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var messageDecoded map[string]interface{}

// DigestFormatter function
func DigestFormatter(messageEncoded []byte) string {

	if err := json.Unmarshal(messageEncoded, &messageDecoded); err != nil {
		panic(err)
	}

	message := DigestMessage{}
	json.Unmarshal(messageEncoded, &message)

	header := message.Header

	newusers, leftusers, wordsauthors, simplified := "", "", "", ""

	if len(message.Words) != 0 {
		words := message.Words
		simplified = "*buzzwords:* " + strings.Join(wordsSimplifier(words), ", ")
		wordsauthors = "*authors:* " + strings.Join(message.Wordsauthors, ", ")
	} else {
		simplified = " Nothing there"
	}

	if len(message.Newusers) != 0 {
		newusers = "*new users:* " + strings.Join(message.Newusers, ", ")
	}

	if len(message.Leftusers) != 0 {
		newusers = "*left users:* " + strings.Join(message.Leftusers, ", ")
	}

	return header + "\r\n" + simplified + "\r\n" + wordsauthors + "\r\n" + newusers + "\r\n" + leftusers
}

// StatusFormatter formatter for Status
func StatusFormatter(messageEncoded []byte) string {

	if err := json.Unmarshal(messageEncoded, &messageDecoded); err != nil {
		panic(err)
	}

	message := StatusMessage{}
	json.Unmarshal(messageEncoded, &message)

	header := message.Header
	toplist := message.TopList
	return header + "\r\n" + toplist
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
