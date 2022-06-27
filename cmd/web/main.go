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
	Q                 *Question
	Chosen            string
	QuestionsAnswered int
	CorrectAnswered   int
}

func NewModel() *Model {
	m := &Model{}
	fmt.Printf("!!! Generate new model: %v\n", m)
	return m
}

type Question struct {
	SentencePrefix string
	SentenceSuffix string
	Variants       []string
	Correct        string
}

func NewQuestion() *Question {
	q := questions[rand.Intn(len(questions))]
	parts := strings.Split(q.Sentence, q.Word)
	m := &Question{
		SentencePrefix: strings.Trim(parts[0], " "),
		SentenceSuffix: strings.Trim(parts[1], " "),
		Variants:       declinations.GetVariants(q.Word),
		Correct:        q.Word,
	}
	fmt.Printf("!!! Generate new question: %v\n", m)
	return m
}

// SetupModel Helper function to get the model from the socket data.
func SetupModel(s live.Socket) *Model {
	m, ok := s.Assigns().(*Model)
	if !ok {
		m = NewModel()
	}
	return m
}

// mount initialises the thermostat state. Data returned to the mount function will
// automatically be assigned to the socket.
func mount(_ context.Context, s live.Socket) (interface{}, error) {
	return SetupModel(s), nil
}

func handleStart(_ context.Context, s live.Socket, _ live.Params) (interface{}, error) {
	model := SetupModel(s)
	model.CorrectAnswered = 0
	model.QuestionsAnswered = 0
	model.Q = NewQuestion()
	fmt.Printf("handleStart return: %#v\n", model)
	return model, nil
}

func handleSelectVariant(_ context.Context, s live.Socket, p live.Params) (interface{}, error) {
	model := SetupModel(s)
	log.Printf("New state: %#v", model)
	chosen := p["chosen"].(string)
	log.Printf("params: %#v", p["chosen"])
	model.Chosen = chosen
	model.QuestionsAnswered++
	if model.Q.Correct == chosen {
		model.CorrectAnswered++
	}
	return model, nil
}

func handleNextQuestion(_ context.Context, s live.Socket, _ live.Params) (interface{}, error) {
	model := SetupModel(s)
	model.Q = NewQuestion()
	model.Chosen = ""
	fmt.Printf("handleNextQuestion return: %#v\n", model)
	return model, nil
}

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
		{{ if not .Assigns.Q }}
			<section class="section">
				<div class="container">
					<button 
						class="button is-rounded is-large" 
						live-click="start">Let's start</button>
				</div>
			</section>
		{{ else }}
			<section class="section">
				<div class="container">
					<h2 class="title is-2">	
						<span>{{.Assigns.Q.SentencePrefix}}&nbsp;</span>
						{{ if eq .Assigns.Chosen "" }}	
							<span>...</span>
						{{ else }}
	                        {{ if eq .Assigns.Chosen .Assigns.Q.Correct }}
								<span class="has-text-success">{{.Assigns.Q.Correct}}</span>
							{{ else }}
								<span class="has-text-danger">{{.Assigns.Q.Correct}}</span>
							{{ end }}	
						{{ end }}
						<span>&nbsp;{{.Assigns.Q.SentenceSuffix}}</span>
					</h2>
					<div class="buttons">
						{{ if eq $.Assigns.Chosen "" }}	
							{{ range .Assigns.Q.Variants }}
								<button 
									class="button is-info is-rounded is-medium" 
									live-click="selectVariant"
									live-value-chosen="{{.}}">{{.}}</button>	
							{{ end }}		
						{{ else }}
							<button 
								class="button is-rounded is-large" 
								live-click="nextQuestion">далі &raquo;</button>	
						{{ end }}	
					</div>
				</div>
			</section>	
		{{ end }}
		<section class="section">
			<div class="container">
				<h3 class="is-size-2">Правильні відповіді: {{.Assigns.CorrectAnswered}}&nbsp;iз&nbsp;{{.Assigns.QuestionsAnswered}}</h3>
			</div>
		</section>
		<!--<section class="section">
			<pre>{{ toPrettyJson .Assigns }}</pre>
		</section>-->	
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

	h.HandleEvent("selectVariant", handleSelectVariant)
	h.HandleEvent("nextQuestion", handleNextQuestion)
	h.HandleEvent("start", handleStart)
	h.HandleError(func(ctx context.Context, err error) {
		log.Printf("Error: %v\n", err)
	})

	http.Handle("/", live.NewHttpHandler(live.NewCookieStore("session-name", []byte("weak-secret")), h))

	// This serves the JS needed to make live work.
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})

	addr := "localhost:28081"
	fmt.Println("http://" + addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
