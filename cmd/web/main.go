package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/jfyne/live"
	"github.com/korchasa/undeclinator/pkg"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var questions []pkg.Match
var declinations *pkg.Declinations

// Model of our test.
type Model struct {
	SentencePrefix string
	SentenceSuffix string
	Variants       []string
	Correct        string
	Chosen         string
}

// NewModel Helper function to get the model from the socket data.
func NewModel(s live.Socket) *Model {
	fmt.Printf("state: %#v\n", s.Assigns())
	m, ok := s.Assigns().(*Model)
	// If we haven't already initialised set up.
	if !ok {
		m = initModel()
	}
	return m
}

func initModel() *Model {
	q := questions[rand.Intn(len(questions))]
	parts := strings.Split(q.Sentence, q.Word)
	m := &Model{
		SentencePrefix: strings.Trim(parts[0], " "),
		SentenceSuffix: strings.Trim(parts[1], " "),
		Variants:       declinations.GetVariants(q.Word),
		Correct:        q.Word,
	}
	fmt.Printf("!!! Generate new model: %v\n", m)
	return m
}

// mount initialises the thermostat state. Data returned to the mount function will
// automatically be assigned to the socket.
func mount(_ context.Context, s live.Socket) (interface{}, error) {
	fmt.Printf("mount\n")
	//debug.PrintStack()
	return NewModel(s), nil
}

// selectVariant on the temp down event, decrease the thermostat temperature by .1 C.
func selectVariant(_ context.Context, s live.Socket, p live.Params) (interface{}, error) {
	model := NewModel(s)
	log.Printf("New state: %#v", model)
	chosen := p["chosen"].(string)
	log.Printf("params: %#v", p["chosen"])
	model.Chosen = chosen
	return model, nil
}

// nextQuestion on the temp down event, decrease the thermostat temperature by .1 C.
func nextQuestion(_ context.Context, _ live.Socket, _ live.Params) (interface{}, error) {
	model := initModel()
	return model, nil
}

// Example shows a simple temperature control using the
// "live-click" event.
func main() {

	rand.Seed(time.Now().UnixNano())

	declinations = pkg.NewDeclinations()

	filename := os.Getenv("CORPUS")
	corpus, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("can't open `%s`: %v", filename, err)
	}

	p := pkg.NewParser(string(corpus), declinations.GetAllWords())
	questions = p.Parse()
	fmt.Printf("Завантажено питань: %d\n", len(questions))

	// Set up the handler.
	h := live.NewHandler()

	// Mount function is called on initial HTTP load and then initial web
	// socket connection. This should be used to create the initial state,
	// the socket Connected func will be true if the mount call is on a web
	// socket connection.
	h.HandleMount(mount)

	// Provide a render function. Here we are doing it manually, but there is a
	// provided WithTemplateRenderer which can be used to work with `html/template`
	h.HandleRender(func(ctx context.Context, data *live.RenderContext) (io.Reader, error) {
		tmpl, err := template.New("declinations").Funcs(sprig.FuncMap()).Parse(`
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>Її</title>
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.9.4/css/bulma.min.css">
	</head>
	<body>
		<section class="section">
			<div class="container">
				<h1 class="title is-1">Її</h1>
				<p class="subtitle">
					Виберіть правильне відмінювання займенника
				</p>
			</div>	
		</section>
		<section class="section">
			<div class="container">
				<h2 class="title is-2">
					<span>{{.Assigns.SentencePrefix}}&nbsp;</span>
					{{ if eq .Assigns.Chosen "" }}	
						<span>...</span>
					{{ else }}
                        {{ if eq .Assigns.Chosen .Assigns.Correct }}
							<span class="has-text-success">{{.Assigns.Correct}}</span>
						{{ else }}
							<span class="has-text-danger">{{.Assigns.Correct}}</span>
						{{ end }}	
					{{ end }}
					<span>&nbsp;{{.Assigns.SentenceSuffix}}</span>
				</h2>
				{{ if eq $.Assigns.Chosen "" }}
				<div class="buttons">	
					{{ range .Assigns.Variants }}
						<button 
								class="button is-info is-rounded is-medium" 
								live-click="selectVariant"
								live-value-chosen="{{.}}">{{.}}</button>	
					{{ end }}	
				</div>
				{{ else }}
				<div class="buttons">
					<button class="button is-rounded is-large" live-click="nextQuestion">далі &raquo;</button>
				</div>
				{{ end }}
			</div>	
		</section>
		<section>
			<pre>{{ toPrettyJson .Assigns }}</pre>
		</section>
	</body>	
	<script src="/live.js"></script>
</html>
        `)
		if err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return nil, err
		}
		return &buf, nil
	})

	h.HandleEvent("selectVariant", selectVariant)
	h.HandleEvent("nextQuestion", nextQuestion)

	http.Handle("/", live.NewHttpHandler(live.NewCookieStore("session-name", []byte("weak-secret")), h))

	// This serves the JS needed to make live work.
	http.Handle("/live.js", live.Javascript{})

	addr := "localhost:28081"
	fmt.Println("http://" + addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
