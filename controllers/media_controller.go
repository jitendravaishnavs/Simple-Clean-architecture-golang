package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10 MB max

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image not found", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// ✅ Validate extension
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if !allowed[ext] {
		http.Error(w, "Invalid image type", http.StatusBadRequest)
		return
	}

	// ✅ Rename and save file
	newName := fmt.Sprintf("img_%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join("uploads", newName)

	dst, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	io.Copy(dst, file)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Image uploaded: %s", newName)
}

func UpdateUserImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	userID := r.FormValue("id")

	// 🧾 Fetch old image from DB
	var oldImage string
	db.QueryRow("SELECT image FROM users WHERE id = ?", userID).Scan(&oldImage)

	// 📷 New image upload
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		http.Error(w, "Invalid image type", http.StatusBadRequest)
		return
	}

	newFileName := fmt.Sprintf("user_%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join("uploads", newFileName)
	dst, _ := os.Create(savePath)
	defer dst.Close()
	io.Copy(dst, file)

	// 🗑️ Delete old image from disk
	if oldImage != "" {
		os.Remove(filepath.Join("uploads", oldImage))
	}

	// 🔁 Update DB
	stmt, _ := db.Prepare("UPDATE users SET image = ? WHERE id = ?")
	_, err = stmt.Exec(newFileName, userID)
	if err != nil {
		http.Error(w, "Failed to update", 500)
		return
	}

	w.Write([]byte("Image updated successfully"))
}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	// 🧾 Get image from DB
	var image string
	db.QueryRow("SELECT image FROM users WHERE id = ?", id).Scan(&image)

	// ❌ Delete from DB
	stmt, _ := db.Prepare("DELETE FROM users WHERE id = ?")
	_, err := stmt.Exec(id)
	if err != nil {
		http.Error(w, "Delete failed", 500)
		return
	}

	// 🗑️ Delete from uploads/
	if image != "" {
		os.Remove(filepath.Join("uploads", image))
	}

	w.Write([]byte("User deleted with image cleaned"))
}
