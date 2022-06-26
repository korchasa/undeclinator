package pkg

import "log"

type Declinations struct {
	data map[string][]string
}

func NewDeclinations() *Declinations {
	return &Declinations{
		data: map[string][]string{
			"наш/наше": {"наш", "наше", "нашого", "нашому", "нашим", "на нашому"},
			"наша":     {"наша", "нашої", "нашiй", "нашу", "нашою", "на нашiй"},
			"нашi":     {"нашi", "наших", "нашим", "нашими", "на наших"},
		},
	}
}

func (d *Declinations) GetAllWords() []string {
	var words []string
	for _, w := range d.data {
		words = append(words, w...)
	}
	return words
}

func (d *Declinations) GetVariants(word string) []string {
	for _, w := range d.data {
		for _, ww := range w {
			if ww == word {
				return w
			}
		}
	}
	log.Fatalf("word `%s` not found", word)
	return nil
}
