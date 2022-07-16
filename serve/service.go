package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	u "serve/utils"
	"strings"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	remote, header, err := r.FormFile("chunk")
	if err != nil {
		log.Println("unable to accept request file:", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	if filepath.Ext(header.Filename) != ZipSuffix {
		if _, err = w.Write([]byte("unsupported type")); err != nil {
			log.Println("unexpected error when writing back:", err)
		}
	}
	zipPath := filepath.Join(property.Path.Pro, property.Path.Raw, header.Filename)

	local, err := os.Create(zipPath)
	if err != nil {
		log.Println("unable to create local file:", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(local, remote)
	if err != nil {
		log.Println("unable to save request file:", err, "save path:", zipPath)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Why don't we use defer?
	// The following process, such as Decompress and BatchDir, needs to handle the file.
	// If the file is not closed, other process can not get the permission to handle it.
	if err = local.Close(); err != nil {
		log.Println("unexpected error when closing file:", err, "save path:", zipPath)
	}
	go func() {
		if err = u.Decompress(zipPath); err != nil {
			log.Println("unable to decompress request zip:", err, "absolute path:", zipPath)
			return
		}
		// The name of decompressed directory is same as the zip file expected extension.
		BatchDir(filepath.Join(property.Path.Pro, property.Path.Raw, strings.Split(header.Filename, Dot)[0]))
	}()

	if _, err = w.Write([]byte("ok")); err != nil {
		log.Println("unexpected error when writing back:", err)
	}
}
