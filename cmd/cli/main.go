package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"github.com/korchasa/undeclinator/pkg"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// нашої
var (
	questions = 0
	correct   = 0
)

func main() {
	rand.Seed(time.Now().UnixNano())

	pairs := map[string][]string{
		"наш/наше": {"наш", "наше", "нашого", "нашому", "нашим", "на нашому"},
		"наша":     {"наша", "нашої", "нашiй", "нашу", "нашою", "на нашiй"},
		"нашi":     {"нашi", "наших", "нашим", "нашими", "на наших"},
	}

	filename := os.Args[1]
	corpus, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("can't open `%s`: %v", filename, err)
	}

	var words []string
	for _, w := range pairs {
		words = append(words, w...)
	}

	p := pkg.NewParser(string(corpus), words)
	matches := p.Parse()
	fmt.Printf("Завантажено питань: %d\n", len(matches))

	for {
		questions++
		rm := matches[rand.Intn(len(matches))]
		for mask, ws := range pairs {
			for _, w := range ws {
				if rm.Word == w {
					fmt.Printf(
						"Виправте відмінок займенника: %s\n",
						strings.ReplaceAll(
							rm.Sentence,
							rm.Word,
							"<"+mask+">"))
					answer := AskUser(ws)
					if answer == w {
						fmt.Printf("Так!\n\n\n")
						correct++
					} else {
						fmt.Printf("Ви помилилися. Правильна відповідь - `%s`!\n\n\n", w)
					}
				}
			}
		}
	}
}

func AskUser(variants []string) string {

	for i, v := range variants {
		fmt.Printf("\t%d) %s\n", i+1, v)
	}
	fmt.Printf("\tq) показати результат та вийти\n")

	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		log.Fatalln(err)
	}

	if 'q' == char {
		fmt.Printf(
			"\n%.0f%% правильних відповідей (%d iз %d)\n",
			float32(correct)/float32(questions)*100,
			correct,
			questions,
		)
		os.Exit(0)
	}

	selectedIndex := int(char - '1')

	if selectedIndex >= len(variants) {
		return AskUser(variants)
	}
	return variants[selectedIndex]
}
