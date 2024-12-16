package util

import (
	"log"
	"net/http"
	"os"
)

func ServeFile(w http.ResponseWriter, r *http.Request, filename string) error {
	var err error

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
	return nil
}
