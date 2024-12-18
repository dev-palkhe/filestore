package main

import (
	"log"
	stdhttp "net/http"
	"os"

	"filestore/internal/filestore"
	myhttp "filestore/internal/http"
)

func main() {
	fs := filestore.NewFileStore()
	if fs == nil {
		log.Fatal("Failed to create filestore")
	}

	stdhttp.HandleFunc("/add", myhttp.AddFileHandler(fs))
	stdhttp.HandleFunc("/list", myhttp.ListFilesHandler(fs))
	stdhttp.HandleFunc("/get", myhttp.GetFileHandler(fs))
	stdhttp.HandleFunc("/update", myhttp.UpdateFileHandler(fs))
	stdhttp.HandleFunc("/wc", myhttp.WordCountHandler(fs))
	stdhttp.HandleFunc("/freq-words", myhttp.FreqWordsHandler(fs))
	stdhttp.HandleFunc("/remove", myhttp.RemoveFileHandler(fs)) // Crucial: Register remove handler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on http://localhost:%s", port)
	log.Fatal(stdhttp.ListenAndServe(":"+port, nil))
}
