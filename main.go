package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func Error(w http.ResponseWriter, e error, code int) {
	if e == nil {
		return
	}
	log.Printf("Error: %q", e)
	w.WriteHeader(code)
	_, err := fmt.Fprintf(w, "%v", e)
	if err != nil {
		panic(err)
	}
}

func main() {

	http.HandleFunc("/unit", func(w http.ResponseWriter, r *http.Request) {

		pdfNames := []string{}

		// Copie de la requête dans le fichier HTML
		d := json.NewDecoder(r.Body)
		defer r.Body.Close()

		d.Decode(&pdfNames)

		// Unit les PDFs entre eux
		cmd := exec.Command("pdfunite", pdfNames...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		e := cmd.Run()
		Error(w, e, http.StatusInternalServerError)

		log.Printf("PDF done with %q\n", pdfNames)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/pdf", func(w http.ResponseWriter, r *http.Request) {
		hfname := fmt.Sprintf("%x.html", &w)
		pfname := fmt.Sprintf("%x.pdf", &w)
		defer r.Body.Close()

		// Création du fichier HTML
		f, e := os.Create(hfname)
		Error(w, e, http.StatusInternalServerError)
		defer os.Remove(f.Name())

		// Copie de la requête dans le fichier HTML
		_, e = io.Copy(f, r.Body)
		Error(w, e, http.StatusInternalServerError)
		f.Close()

		// Transformation du HTML en PDF
		cmd := exec.Command("wkhtmltopdf", f.Name(), pfname)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		e = cmd.Run()
		Error(w, e, http.StatusInternalServerError)

		// Ouverture du nouveau fichier PDF
		pf, e := os.Open(pfname)
		Error(w, e, http.StatusInternalServerError)
		defer os.Remove(pf.Name())

		// Renvoie du fichier PDF dans la requête
		_, e = io.Copy(w, pf)
		Error(w, e, http.StatusInternalServerError)
		pf.Close()
	})

	log.Printf("listening on port :8000")
	http.ListenAndServe(":8000", nil)
}
