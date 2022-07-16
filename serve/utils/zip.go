package utils

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
)

func Decompress(archive string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		log.Println("unable to open zip reader:", err, "archive:", archive)
		return err
	}
	defer func() {
		if err = reader.Close(); err != nil {
			log.Println("unexpected error when closing zip:", err, "archive:", archive)
		}
	}()

	dir := filepath.Dir(archive)
	for _, file := range reader.File {
		stream, err := AnsiToUTF8([]byte(file.Name))
		if err != nil {
			log.Println("unable to encode filename to UTF8:", err, "zip path:", archive, "file name:", file.Name)
			continue
		}
		filename := string(stream)
		if file.FileInfo().IsDir() {
			// directory
			err = os.MkdirAll(filepath.Join(dir, filename), os.ModePerm)
			if err != nil {
				log.Println("unable to create a directory:", err, "zip path:", archive, "dir name:", file.Name)
				return err
			}
			continue
		}
		// normal file
		newFile, err := os.Create(filepath.Join(dir, filename))
		if err != nil {
			log.Println("unable to create a new file:", err, "zip path:", archive, "file name:", file.Name)
			return err
		}
		arcFile, err := file.Open()
		if err != nil {
			log.Println("unable to open archive file:", err, "zip path:", archive, "file name:", file.Name)
			return err
		}
		_, err = io.Copy(newFile, arcFile)
		if err != nil {
			log.Println("unexpected error when copying archive file:", err, "zip path:", archive, "file name:", file.Name)
			return err
		}
		if err = newFile.Close(); err != nil {
			log.Println("unexpected error when closing new copied file:", err, "zip path:", archive, "file name:", file.Name)
			return err
		}
		if err = arcFile.Close(); err != nil {
			log.Println("unexpected error when closing archive file:", err, "zip path:", archive, "file name:", file.Name)
			return err
		}
	}
	return nil
}
