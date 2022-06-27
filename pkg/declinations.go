package pkg

import "log"

type Declinations struct {
	data map[string][]string
}

func NewDeclinations() *Declinations {
	return &Declinations{
		data: map[string][]string{
			"мiй":      {"мiй", "мого", "моєму", "моїм"},
			"моя":      {"моя", "моєї", "моїй", "мою"},
			"моє":      {"моє", "мого", "моєму", "моїм"},
			"мої":      {"мої", "моїх", "моїм", "моїми"},
			"твiй":     {"твiй", "твого", "твоєму", "твоїм"},
			"твоя":     {"твоя", "твоєї", "твоїй", "твою"},
			"твоє":     {"твоє", "твого", "твоєму", "твоїм"},
			"твої":     {"твої", "твоїх", "твоїм", "твоїми"},
			"свiй":     {"свiй", "свого", "своєму", "своїм"},
			"своя":     {"своя", "своєї", "своїй", "свою"},
			"своє":     {"своє", "свого", "своєму", "своїм"},
			"свої":     {"свої", "своїх", "своїм", "своїми"},
			"їxнiй":    {"їxнiй", "їxнього", "їxньому", "їxнiм"},
			"їxня":     {"їxня", "їxньої", "їxнiй", "їxню", "їxньою"},
			"їxнє":     {"їxнє", "їxнього", "їxньому", "їxнє", "їxнiм"},
			"їxнi":     {"їxнi", "їxнiх", "їxнiм", "їxнiми"},
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
