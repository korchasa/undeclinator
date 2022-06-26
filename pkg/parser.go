package pkg

import (
	"strings"
)

type Parser struct {
	corpus string
	words  []string
}

func NewParser(corpus string, words []string) *Parser {
	return &Parser{corpus: corpus, words: words}
}

func (p *Parser) Parse() []Match {
	t := strings.ReplaceAll(p.corpus, "\n", " ")
	t = strings.ReplaceAll(t, "  ", " ")
	t = strings.ReplaceAll(t, "  ", " ")
	t = strings.ReplaceAll(t, "  ", " ")
	tt := splitAny(t, ".!?")
	var ms []Match
	for _, m := range tt {
		for _, w := range p.words {
			if strings.Contains(m, " "+w+" ") {
				ms = append(ms, Match{Sentence: strings.Trim(m, " ") + ".", Word: w})
			}
		}
	}
	return ms
}

func splitAny(s string, seps string) []string {
	splitter := func(r rune) bool {
		return strings.ContainsRune(seps, r)
	}
	return strings.FieldsFunc(s, splitter)
}

type Match struct {
	Sentence string
	Word     string
}
