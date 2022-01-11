package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
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

func ReadTestGuesses(r io.Reader) []TestGuess {
	s := bufio.NewScanner(r)
	res := []TestGuess{}
	for s.Scan() {
		letters := s.Text()
		if !s.Scan() {
			break
		}
		matches := s.Text()
		for len([]rune(matches)) < 5 {
			matches += " "
		}
		rules := []Rule{}
		for i := 0; i < 5; i++ {
			r := Rule{
				pos:    i,
				letter: []rune(letters)[i],
				match:  ParseMatchSymbol([]rune(matches)[i]),
			}
			rules = append(rules, r)
		}
		res = append(res, TestGuess{letters, rules})
	}
	return res
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

type GuessRules []Rule

type TestGuess struct {
	word  string
	rules GuessRules
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

// Answer
// provides the wordle answer for guess against word
// e.g. word=drink, guess=snare, returns:
// {0, s, grey}, {1, n, yellow}, {2, a, grey}, {3, r, yellow}, {4, e, grey}
func Answer(word, guess string) GuessRules {
	freq := map[rune]int{}
	res := []Rule{}
	for i, r := range guess {
		freq[r]++
		res = append(res, Rule{
			pos:    i,
			letter: r,
			match:  AnswerMatch(word, r, i, freq[r]),
		})
	}
	return res
}

func AnswerMatch(word string, letter rune, pos, count int) Match {
	if []rune(word)[pos] == letter {
		return GREEN
	}
	// count runes for repeated letter determination
	freq := map[rune]int{}
	for _, r := range word {
		freq[r]++
	}
	if strings.ContainsRune(word, letter) && freq[letter] >= count {
		return YELLOW
	}
	return GREY
}

func EqualRules(a, b []Rule) bool {
	if len(a) != len(b) {
		return false
	}
	for i, x := range a {
		if x != b[i] {
			return false
		}
	}
	return true
}

func (r GuessRules) String() string {
	bgGrey := color.New(color.BgHiBlack, color.Bold).SprintFunc()
	bgYellow := color.New(color.FgHiWhite, color.BgHiYellow, color.Bold).SprintFunc()
	bgGreen := color.New(color.FgHiWhite, color.BgHiGreen, color.Bold).SprintFunc()
	s := ""
	for _, c := range r {
		var prFunc func(a ...interface{}) string
		switch c.match {
		case GREY:
			prFunc = bgGrey
		case YELLOW:
			prFunc = bgYellow
		case GREEN:
			prFunc = bgGreen
		}
		s += fmt.Sprintf("%s", prFunc(string(c.letter)))
	}
	return s
}

func main() {
	var (
		listMatchWords, allWords bool
		testWord                 string
	)
	flag.BoolVar(&listMatchWords, "l", false, "List all words that match the given rules")
	flag.BoolVar(&allWords, "a", false, "Consider all words for best match")
	flag.StringVar(&testWord, "t", "", "Test guesses against this answer and print wordle responses")
	flag.Parse()

	guessesFile := os.Stdin
	guessesFilename := flag.Arg(0)
	if guessesFilename != "" && guessesFilename != "-" {
		var err error
		guessesFile, err = os.Open(guessesFilename)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if testWord != "" {
		fail := 0
		tests := ReadTestGuesses(guessesFile)
		for _, t := range tests {
			actual := Answer(testWord, t.word)
			if EqualRules(actual, t.rules) {
				fmt.Printf("%s OK %v\n", t.word, t.rules)
			} else {
				fmt.Printf("%s FAIL %v != %v\n", t.word, t.rules, actual)
				fail = 1
			}
		}
		os.Exit(fail)
	}

	f, _ := os.Open("words")
	words := ReadWords(f)
	rules := ReadGuesses(guessesFile)
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
