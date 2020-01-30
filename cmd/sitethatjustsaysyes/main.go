package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"text/template"

	log "github.com/sirupsen/logrus"
)

type answer struct {
	Title string
}

var (
	answerTmpl = template.Must(template.New("answer").Parse(`
<!DOCTYPE html><html>
<head><style>
.imtw div {
  margin-top: 25%;
}
.imtw div p {
  text-align: center;
  font-family: serif;
  font-size: xx-large;
}
.imtw div p.suble {
	font-size: small;
	color: grey;
}
</style><title>{{ .Title }}</title></head>
<body class="imtw">
	<div><p>Yes.</p></div>
</body></html>
<!-- Hello, Min! -->
`))
)

type answerer struct {
	payload string
}

func (a *answerer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	fmt.Fprint(w, a.payload)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set.")
	}

	title := os.Getenv("SITE_TITLE")

	if title == "" {
		title = "Is this website improperly configured?"
	}

	var tpl bytes.Buffer
	if err := answerTmpl.Execute(&tpl, answer{Title: title}); err != nil {
		log.WithError(err).Fatal("Failed parsing the hardcoded template")
	}

	http.Handle("/", &answerer{
		payload: tpl.String(),
	})

	hp := fmt.Sprintf(":%v", port)
	log.Printf("Listening on %v", hp)
	http.ListenAndServe(hp, nil)
}
