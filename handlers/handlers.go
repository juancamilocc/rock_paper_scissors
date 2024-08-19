package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"strconv"

	"github.com/juancamilocc/rock_paper_scissors/rps"
)

type Player struct {
	Name string
}

var player Player

func Home(w http.ResponseWriter, r *http.Request) {
	restartValue()
	renderTemplate(w, "home.html", nil)
}

func NewGame(w http.ResponseWriter, r *http.Request) {
	restartValue()
	renderTemplate(w, "new-game.html", nil)
}

func Game(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		// Leer los datos del formulario
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		player.Name = r.Form.Get("name")
	}

	// Redirecíon a otra ruta
	if player.Name == "" {
		http.Redirect(w, r, "/new", http.StatusFound)
	}

	renderTemplate(w, "game.html", player)
}

func Play(w http.ResponseWriter, r *http.Request) {
	playerChoice, _ := strconv.Atoi(r.URL.Query().Get("c"))
	result := rps.PlayRound(playerChoice)

	out, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func About(w http.ResponseWriter, r *http.Request) {
	restartValue()
	renderTemplate(w, "about.html", nil)
}

// Manage error pages
var errorTemplates = template.Must(template.ParseGlob("templates/**/*.html"))

func handlerError(w http.ResponseWriter, name string, status int) {
	w.WriteHeader(status)
	errorTemplates.ExecuteTemplate(w, name, nil)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	handlerError(w, "404", http.StatusNotFound)
}

const baseDir = "templates/"

func renderTemplate(w http.ResponseWriter, name string, data any) {
	templates := template.Must(template.ParseFiles(baseDir+"base.html", baseDir+name))
	// Header
	w.Header().Set("Content-Type", "text/html")

	// Render template in the answer
	err := templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		handlerError(w, "500", http.StatusInternalServerError)
	}
}

// Restart values
func restartValue() {
	player.Name = ""
	rps.ComputerScore = 0
	rps.PlayerScore = 0
}
