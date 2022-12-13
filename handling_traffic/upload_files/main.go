package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"net/http"
)

// func main() {
// 	file, err := ioutil.TempFile("handling_traffic/upload_files/dir", "car-*.png")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	// We can choose to have these files deleted on program close
// 	//defer os.Remove(file.Name())
// 	file.Close()
// 	if _, err := file.Write([]byte("hello world\n")); err != nil {
// 		fmt.Println(err)
// 	}

// 	data, err := ioutil.ReadFile(file.Name())
// 	// if our program was unable to read the file
// 	// print out the reason why it can't
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// if it was successful in reading the file then
// 	// print out the contents as a string
// 	fmt.Print(string(data))

// }

const MAX_UPLOAD_SIZE = 1024 * 1024 * 2 // 2MB

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("File Upload Endpoint Hit")

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
		return
	}
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, fileHeader, err := r.FormFile("file")
	fmt.Println(r.MultipartForm.Value["sss"], r.MultipartForm.Value["ahmed"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	// Create the uploads folder if it doesn't
	// already exist
	err = os.MkdirAll("./uploads", os.ModePerm)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Create a new file in the uploads directory
	dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	fmt.Printf("Uploaded File: %+v\n", fileHeader.Filename)
	fmt.Printf("File Size: %+v\n", fileHeader.Size)
	fmt.Printf("MIME Header: %+v\n", fileHeader.Header)

	//// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Upload successful")
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Hello World")
	setupRoutes()
}
