package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"filestore/internal/filestore"
)

func AddFileHandler(fs *filestore.FileStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.FormValue("filename")
		content := r.FormValue("content")

		fmt.Println("Adding file:", filename)
		err := fs.Add(filename, content)
		if err != nil {
			fmt.Println("Error adding file:", err)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		fmt.Println("File added successfully:", filename)
		w.WriteHeader(http.StatusCreated)
	}
}

func UpdateFileHandler(fs *filestore.FileStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.FormValue("filename")
		content := r.FormValue("content")

		if err := fs.Update(filename, content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func RemoveFileHandler(fs *filestore.FileStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.FormValue("filename")
		fmt.Println("Removing file:", filename)
		err := fs.Remove(filename)
		if err != nil {
			fmt.Println("Error removing file:", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		fmt.Println("File removed successfully:", filename)
		w.WriteHeader(http.StatusOK)
	}
}

func WordCountHandler(fs *filestore.FileStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		totalWordCount := 0
		for _, fileInfo := range fs.Files {
			totalWordCount += filestore.WordCount(fileInfo.Content)
		}
		fmt.Fprintln(w, totalWordCount)
	}
}

func FreqWordsHandler(fs *filestore.FileStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("limit")
		order := r.URL.Query().Get("order")

		limit := 10
		if limitStr != "" {
			l, err := strconv.Atoi(limitStr)
			if err == nil {
				limit = l
			}
		}

		if order == "" {
			order = "dsc"
		}

		allText := ""
		for _, fileInfo := range fs.Files {
			allText += fileInfo.Content + " "
		}

		result := filestore.FrequentWords(allText, limit, order)

		json.NewEncoder(w).Encode(result)
	}
}

func ListFilesHandler(fs *filestore.FileStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files := fs.List()
		json.NewEncoder(w).Encode(files)
	}
}

func GetFileHandler(fs *filestore.FileStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("filename")
		file, ok := fs.Get(filename)
		if !ok {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(file)
	}
}
