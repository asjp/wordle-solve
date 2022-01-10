package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

func ReadWords(r io.ReadSeeker) []string {
	r.Seek(io.SeekStart, 0)
	var words []string
	s := bufio.NewScanner(r)
	for s.Scan() {
		words = append(words, s.Text())
	}
	return words
}

// ParseMatchSymbol
// ' ' means GREY, i.e. not in the word
// '.' means YELLOW, i.e. somewhere else in the word
// anything else means GREEN, i.e. correct letter in the right place
func ParseMatchSymbol(r rune) Match {
	if r == '.' {
		return YELLOW
	}
	if r == ' ' {
		return GREY
	}
	return GREEN
}

func ReadGuesses(r io.Reader) []Rule {
	s := bufio.NewScanner(r)
	res := []Rule{}
	for s.Scan() {
		letters := s.Text()
		if !s.Scan() {
			break
		}
		matches := s.Text()
		for len([]rune(matches)) < 5 {
			matches += " "
		}
		for i := 0; i < 5; i++ {
			r := Rule{
				pos:    i,
				letter: []rune(letters)[i],
				match:  ParseMatchSymbol([]rune(matches)[i]),
			}
			res = append(res, r)
		}
	}
	return res
}

type Match int

const (
	GREY   Match = 0
	YELLOW Match = 1
	GREEN  Match = 2
)

type Rule struct {
	pos    int
	letter rune
	match  Match
}

func Reduce(words []string, rules []Rule) []string {
	var res []string
	for _, w := range words {
		valid := true
		for _, r := range rules {
			if !MatchRule(w, r) {
				valid = false
				break
			}
		}
		if valid {
			res = append(res, w)
		}
	}
	return res
}

func MatchRule(word string, rule Rule) bool {
	inWord := strings.ContainsRune(word, rule.letter)
	exactPos := []rune(word)[rule.pos] == rule.letter

	if rule.match == GREY {
		return !inWord
	}
	if rule.match == GREEN {
		return exactPos
	}
	return inWord && !exactPos
}

// make all 243 (3^5) possible rule permutations
func MakeRulePermutations(w string) [][]Rule {
	res := [][]Rule{}
	for i := 0; i < 243; i++ {
		next := []Rule{
			{
				0, []rune(w)[0], Match(i / 81),
			},
			{
				1, []rune(w)[1], Match((i / 27) % 3),
			},
			{
				2, []rune(w)[2], Match((i / 9) % 3),
			},
			{
				3, []rune(w)[3], Match((i / 3) % 3),
			},
			{
				4, []rune(w)[4], Match(i % 3),
			},
		}
		res = append(res, next)
	}
	return res
}

func CalcExpectedWordsRemaining(available []string, rules []Rule, w string) float64 {
	perms := MakeRulePermutations(w)
	sum := 0.
	for _, p := range perms {
		reduced := Reduce(available, p)
		sum += float64(len(reduced)) * (float64(len(reduced)) / float64(len(available)))
	}
	return sum
}

func main() {
	var listMatchWords, allWords bool
	flag.BoolVar(&listMatchWords, "l", false, "List all words that match the given rules")
	flag.BoolVar(&allWords, "a", false, "Consider all words for best match")
	flag.Parse()

	f, _ := os.Open("words")
	words := ReadWords(f)
	rules := ReadGuesses(os.Stdin)
	options := Reduce(words, nil)
	if !allWords {
		options = Reduce(words, rules)
	}

	lowestExpected := math.MaxFloat64
	bestGuess := ""

	for _, w := range options {
		expectedWordsRemaining := CalcExpectedWordsRemaining(options, rules, w)
		if listMatchWords {
			fmt.Printf("%s %.3f\n", w, expectedWordsRemaining)
		}
		if expectedWordsRemaining < lowestExpected {
			lowestExpected = expectedWordsRemaining
			bestGuess = w
		}
	}

	if !listMatchWords {
		fmt.Println(bestGuess)
	}
}
