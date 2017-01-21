package main

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type fileInfo struct {
	Size int `json:"size"`
}

var templates = template.Must(template.ParseFiles("tmpl/form.html"))

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/form", formHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/form", http.StatusFound)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	// read file from form
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	// create temp file
	f, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name())

	io.Copy(f, file)

	// get size
	fi, err := f.Stat()
	if err != nil {
		log.Println(err)
	}

	// write JSON
	toJSON := fileInfo{int(fi.Size())}
	js, err := json.Marshal(toJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "form.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
