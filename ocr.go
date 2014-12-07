package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"

	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.image/bmp"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	/* Uploaded file management */

	uploadedFile, header, err := r.FormFile("file")
	if err != nil {
		redirectErrorDesc(w, r, http.StatusBadRequest, "You have to select a "+
			"file to begin the upload!")
		return
	}

	/* Extension and size checking */

	extensions := []string{".jpg", ".jpeg", ".png", ".bmp"}
	validExtension := false
	if len(header.Filename) > 4 {
		for _, v := range extensions {
			if header.Filename[len(header.Filename)-len(v):] == v {
				validExtension = true
				break
			}
		}
	} else {
		validExtension = false
	}
	if !validExtension {
		redirectErrorDesc(w, r, http.StatusBadRequest, "We only accept JPEG, "+
			"PNG and BMP files.")
		return
	}

	reader := bufio.NewReader(uploadedFile)
	contents := make([]byte, ocrMaxFileSize+1)
	n, err := reader.Read(contents)
	if err != nil {
		redirectError(w, r, http.StatusInternalServerError)
		return
	}
	if false && n > ocrMaxFileSize {
		redirectErrorDesc(w, r, http.StatusRequestEntityTooLarge, "The file "+
			"you are trying to upload is too large.")
		return
	}

	/* File conversion */

	uploadedFileReader := bytes.NewReader(contents)
	img, err := jpeg.Decode(uploadedFileReader)
	if err != nil {
		uploadedFileReader = bytes.NewReader(contents)
		img, err = png.Decode(uploadedFileReader)
		if err != nil {
			uploadedFileReader = bytes.NewReader(contents)
			img, err = bmp.Decode(uploadedFileReader)
			if err != nil {
				redirectErrorDesc(w, r, http.StatusBadRequest, "The file you "+
					"uploaded does not seem to be a valid image file.")
				return
			}
		}
	}
	tmpName := "./tmp/" + uuid.NewUUID().String()
	file, err := os.Create(tmpName)
	if err != nil {
		fmt.Errorf("Error while creating the temporary converted file\n")
		redirectError(w, r, http.StatusInternalServerError)
		return
	}
	defer file.Close()
	defer removeTmpFile(tmpName)
	err = bmp.Encode(file, img)
	if err != nil {
		fmt.Errorf("Error while saving the converted image file\n")
		redirectError(w, r, http.StatusInternalServerError)
		return
	}

	// TODO: run the OCR
	redirectErrorDesc(w, r, http.StatusNotImplemented, "The OCR feature has "+
		"yet to be finished. Come back soon, it should be ready by the 8th "+
		"december!")
}

func removeTmpFile(file string) {
	err := os.Remove(file)
	if err != nil {
		fmt.Errorf(err.Error())
	}
}
