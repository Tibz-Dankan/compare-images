package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/vitali-fedulov/images3"
)

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

	file1, _, err := r.FormFile("image1")
	if err != nil {
		http.Error(w, "Missing image1 field in the form", http.StatusBadRequest)
		log.Println(err)
		return
	}
	defer file1.Close()

	file2, _, err := r.FormFile("image2")
	if err != nil {
		http.Error(w, "Missing image2 field in the form", http.StatusBadRequest)
		log.Println(err)
		return
	}
	defer file2.Close()


	similar := ImagesSimilar(file1, file2)
	
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


func ImagesSimilar(image1 multipart.File, image2 multipart.File) bool {
	imageFile1 := writeImageToDisk(image1, "1.jpg")
	imageFile2 := writeImageToDisk(image2, "2.jpg")

	img1, _ := images3.Open("1.jpg")
	img2, _ := images3.Open("2.jpg")

	icon1 := images3.Icon(img1, "1.jpg")
	icon2 := images3.Icon(img2, "2.jpg")

	removeImageFromDisk(imageFile1, "1.jpg")
	removeImageFromDisk(imageFile2, "2.jpg")

	if images3.Similar(icon1, icon2) {
		fmt.Println("Images are similar.")
		return true
	}
	   fmt.Println("Images are distinct.")
	return false
}


func writeImageToDisk(image multipart.File, imageName string) *os.File{
	imageFile, err := os.Create(imageName)
	if err != nil {
		fmt.Println("Error creating file 1.jpg:", err)
		// return error message
	}
	defer imageFile.Close()

	_, err = io.Copy(imageFile, image)
	if err != nil {
		fmt.Println("Error copying image1:", err)
		// return error message
	}

	return imageFile
}

func removeImageFromDisk(imageFile *os.File, imageName string){
	defer func() {
		imageFile.Close()
		os.Remove(imageName) 
	}()
}