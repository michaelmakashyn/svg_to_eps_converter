package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	http.HandleFunc("/", convertHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the SVG file from the request
	file, header, err := r.FormFile("svg")
	if err != nil {
		http.Error(w, "Failed to read SVG file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate a unique filename for the SVG and EPS files
	svgFileName := fmt.Sprintf("%s.svg", header.Filename)
	epsFileName := fmt.Sprintf("%s.eps", header.Filename)

	// Create a temporary file to save the SVG
	tempSVGFile, err := os.Create(svgFileName)
	if err != nil {
		http.Error(w, "Failed to create temporary SVG file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempSVGFile.Name())
	defer tempSVGFile.Close()

	// Save the SVG content to the temporary file
	_, err = io.Copy(tempSVGFile, file)
	if err != nil {
		http.Error(w, "Failed to save SVG content", http.StatusInternalServerError)
		return
	}

	// Convert the SVG to EPS using Inkscape
	cmd := exec.Command("inkscape", "--export-type=eps", "--export-filename="+epsFileName, svgFileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Printf("Failed to convert SVG to EPS: %v\n", err)
		http.Error(w, "Failed to convert SVG to EPS", http.StatusInternalServerError)
		return
	}

	// Open the generated EPS file
	epsFile, err := os.Open(epsFileName)
	if err != nil {
		log.Printf("Failed to open generated EPS file: %v\n", err)
		http.Error(w, "Failed to open EPS file", http.StatusInternalServerError)
		return
	}
	defer epsFile.Close()

	// Set the appropriate headers for the response
	w.Header().Set("Content-Disposition", "attachment; filename="+epsFileName)
	w.Header().Set("Content-Type", "application/postscript")

	// Copy the EPS content to the response writer
	_, err = io.Copy(w, epsFile)
	if err != nil {
		log.Printf("Failed to write EPS content to response: %v\n", err)
	}
}