package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

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
					"uploaded does not seem to be a valid image file. The way "+
					"your image was encoded might not be compatible with the "+
					"libraries we use.")
				return
			}
		}
	}
	tmpName := "tmp/" + uuid.NewUUID().String()
	file, err := os.Create(tmpName)
	if err != nil {
		fmt.Errorf("Error while creating the temporary converted file: %s\n",
			err.Error())
		redirectError(w, r, http.StatusInternalServerError)
		return
	}
	defer file.Close()
	defer removeTmpFile(tmpName)
	err = bmp.Encode(file, img)
	if err != nil {
		fmt.Errorf("Error while saving the converted image file: %s\n",
			err.Error())
		redirectError(w, r, http.StatusInternalServerError)
		return
	}

	fmt.Printf("Running the OCR...\n")
	cmd := exec.Command("./main", "-f", "../"+tmpName)
	cmd.Dir = "mediocr"
	var out bytes.Buffer
	cmd.Stdout = &out

	go controlExecutionTime(cmd)
	err = cmd.Run()
	if err == nil {
		fmt.Printf("Result: %s\n", out.String())

	} else if err.Error() == "signal: killed" {
		fmt.Errorf("Killed an OCR process running for too long\n")
		redirectErrorDesc(w, r, http.StatusGatewayTimeout, "We couldn't "+
			"process your image fast enough and had to kill the process.")
		return

	} else {
		fmt.Errorf("OCR error: %s\n", err.Error())
		redirectErrorDesc(w, r, http.StatusInternalServerError, "The "+
			"underlaying OCR software has returned an inexpected result.")
		return
	}

	getSession(r).AddFlash(out.String(), "ocr_result")
	saveSession(w, r)
	http.Redirect(w, r, "/try", http.StatusSeeOther)
}

func controlExecutionTime(cmd *exec.Cmd) {
	startTime := time.Now()
	process := cmd.Process
	for {
		if process == nil {
			return
		}
		if err := process.Signal(syscall.Signal(0)); err != nil {
			return
		}
		if time.Now().Sub(startTime) > ocrMaxExecutionTime*time.Second {
			err := process.Kill()
			if err != nil {
				fmt.Errorf("Error while killing a long-running process: %s\n",
					err.Error())
			}
			return
		}
		time.Sleep(time.Millisecond * 200)
	}
}

func removeTmpFile(file string) {
	err := os.Remove(file)
	if err != nil {
		fmt.Errorf(err.Error())
	}
}
