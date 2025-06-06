package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	os.MkdirAll("uploads", os.ModePerm)

	r := mux.NewRouter()
	r.HandleFunc("/upload", uploadFile).Methods("POST")

	// 🔥 Serve files at http://localhost:8080/uploads/filename.png
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	log.Println("🚀 Server started at :8080")
	http.ListenAndServe(":8080", r)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // max 10MB

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	f, err := os.Create("uploads/" + handler.Filename)
	if err != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	io.Copy(f, file)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Uploaded successfully!",
		"filename": handler.Filename,
	})
}
