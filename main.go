package main

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/corona10/goimagehash"
)

// Response represents the JSON response.
type Response struct {
	Similar bool `json:"similar"`
}

func main() {
	http.HandleFunc("/compare", CompareImages)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func CompareImages(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // Set a reasonable max memory size for the form
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Get the uploaded files
	file1, fileHeader1, err := r.FormFile("image1")
	if err != nil {
		http.Error(w, "Missing image1 field in the form", http.StatusBadRequest)
		log.Println(err)
		return
	}
	defer file1.Close()

	file2, fileHeader2, err := r.FormFile("image2")
	if err != nil {
		http.Error(w, "Missing image2 field in the form", http.StatusBadRequest)
		log.Println(err)
		return
	}
	defer file2.Close()

	// Check file extensions
	if !isValidImageType(fileHeader1.Filename) || !isValidImageType(fileHeader2.Filename) {
		http.Error(w, "Invalid image file type. Supported types: .png and .jpeg", http.StatusBadRequest)
		return
	}

	fmt.Println("About to 2 calc hash 1")

	// Calculate and compare pHash values
	hash1, err := calculateHash(file1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	fmt.Println("About to 2 calc hash 2")


	hash2, err := calculateHash(file2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	fmt.Println("About to compare hashes")


	// Compare the pHash values
	similar := compareHashes(hash1, hash2)

	// Send JSON response
	response := Response{Similar: similar}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// func calculateHash(file multipart.File) (string, error) {
// 	img, _, err := image.Decode(file)
// 	if err != nil {
// 		log.Println(err)
// 		return "", err
// 	}

// 	hash, err := goimagehash.DifferenceHash(img)
// 	if err != nil {
// 		log.Println(err)
// 		return "", err
// 	}

// 	return hash.ToString(), nil
// }

// func calculateHash(file multipart.File) (string, error) {
// 	// Try to decode the image using JPEG format
// 	img, _, err := image.Decode(file)
// 	if err != nil {
// 		// If JPEG format decoding fails, try PNG format
// 		file.Seek(0, io.SeekStart) // Reset file reader position
// 		img, _, err = image.Decode(file)
// 		if err != nil {
// 			return "", err
// 		}
// 	}

// 	hash, err := goimagehash.DifferenceHash(img)
// 	if err != nil {
// 		return "", err
// 	}

// 	return hash.ToString(), nil
// }

func calculateHash(file multipart.File) (string, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("error decoding image: %v", err)
	}

	hash, err := goimagehash.DifferenceHash(img)
	if err != nil {
		return "", fmt.Errorf("error calculating hash: %v", err)
	}

	return hash.ToString(), nil
}



func compareHashes(hash1, hash2 string) bool {
	return strings.Compare(hash1, hash2) == 0
}

func isValidImageType(filename string) bool {
	extension := strings.ToLower(filepath.Ext(filename))
	return extension == ".png" || extension == ".jpeg" || extension == ".jpg"
}
