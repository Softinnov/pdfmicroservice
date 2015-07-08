package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func Error(w http.ResponseWriter, e error, code int) {
	log.Printf("Error: %q", e)
	w.WriteHeader(code)
	_, err := fmt.Fprintf(w, "%v", e)
	if err != nil {
		panic(err)
	}
}

func main() {

	http.HandleFunc("/pdfunite", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)

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
		if e != nil {
			Error(w, e, http.StatusInternalServerError)
			return
		}

		log.Printf("PDF done with %q\n", pdfNames)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/pdf", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)

		v := r.URL.Query()

		htmlPath := v.Get("file")
		header := v.Get("header")
		footer := v.Get("footer")

		pdfPath := strings.Replace(htmlPath, ".html", ".pdf", 1)
		fmt.Printf("html: %q\n", htmlPath)
		fmt.Printf("pdf : %q\n", pdfPath)

		cmds := []string{
			"wkhtmltopdf",
		}
		if header != "" {
			cmds = append(cmds, "--header-html", header)
		}
		if footer != "" {
			cmds = append(cmds, "--footer-html", footer)
		}
		cmds = append(cmds, string(htmlPath), pdfPath)

		// Transformation du HTML en PDF
		cmd := exec.Command(cmds[0], cmds[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		e := cmd.Run()
		if e != nil {
			Error(w, e, http.StatusInternalServerError)
			return
		}
		log.Printf("PDF done for %q\n", pdfPath)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(pdfPath))
	})

	log.Printf("listening on port :8000")
	http.ListenAndServe(":8000", nil)
}
