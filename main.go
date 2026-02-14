// package declaration
package main

// imports
import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"
)

// structure definitions
type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Job struct {
	Title   string   `json:"title"`
	Company string   `json:"company"`
	Start   string   `json:"start"`
	End     string   `json:"end"`
	Points  []string `json:"points"`
}

type Reference struct {
	Name     string `json:"name"`
	Position string `json:"position"`
	Email    string `json:"email"`
}

type Resume struct {
	Name       string      `json:"name"`
	Telephone  string      `json:"telephone"`
	Email      string      `json:"email"`
	LinkedIn   string      `json:"linkedin"`
	GitHub     string      `json:"github"`
	Projects   []Project   `json:"projects"`
	Jobs       []Job       `json:"jobs"`
	References []Reference `json:"references"`
	Skills     []string    `json:"skills"`
}

// data
var t *template.Template = nil
var resume Resume = Resume{}

// function definitions
func init() {

	// periodically reload template
	go func() {
		for {
			// initialized data
			var err error = nil
			var data []byte = nil

			// load the resume json
			data, err = os.ReadFile("resume.json")
			if err != nil {

				// logs
				fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())

				// error
				os.Exit(1)
			}

			// parse the json
			err = json.Unmarshal(data, &resume)
			if err != nil {

				// logs
				fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())

				// error
				os.Exit(1)
			}

			// log
			fmt.Printf("Refreshed templates\n")

			// parse the templates
			t = template.Must(template.ParseGlob("*.html"))

			// wait
			time.Sleep(1 * time.Second)
		}
	}()
}

func index(w http.ResponseWriter, r *http.Request) {

	// respond
	t.ExecuteTemplate(w, "template.html", resume)
}

// entry point
func main() {

	// serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// routes
	http.HandleFunc("/", index)

	// listen and serve
	http.ListenAndServe(":8080", nil)
}
