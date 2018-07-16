package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
)

type EmberResponse struct {
	Error     bool   `json:"error"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type AddonPost struct {
	URL       string `json:"url"`
	Filename  string `json:"filename"`
	ExtractTo string `json:"extractTo"`
}

func main() {
	r := httprouter.New()

	r.POST("/addon/install", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		installed := EmberResponse{false, "addon installed", time.Now().Unix()}
		a := AddonPost{}

		json.NewDecoder(r.Body).Decode(&a)

		log.Println(a.ExtractTo)

		output, _ := os.Create(a.ExtractTo + "/" + a.Filename + ".zip")
		defer output.Close()

		response, err := http.Get(a.URL)

		if err != nil {
			log.Fatal(err)
		}

		defer response.Body.Close()

		io.Copy(output, response.Body)

		reader, _ := zip.OpenReader(a.ExtractTo + "/" + a.Filename + ".zip")

		if err := os.MkdirAll(a.ExtractTo+"/"+a.Filename, 0755); err != nil {
			log.Fatal(err)
		}

		for _, file := range reader.File {
			path := filepath.Join(a.ExtractTo+"/"+a.Filename, file.Name)
			if file.FileInfo().IsDir() {
				os.MkdirAll(path, file.Mode())
				continue
			}

			fileReader, err := file.Open()
			if err != nil {
				log.Fatal(err)
			}
			defer fileReader.Close()

			targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				log.Fatal(err)
			}
			defer targetFile.Close()

			if _, err := io.Copy(targetFile, fileReader); err != nil {
				log.Fatal(err)
			}
		}

		err = os.Remove(a.ExtractTo + "/" + a.Filename + ".zip")

		if err != nil {
			log.Fatal(err)
		}

		ij, _ := json.Marshal(installed)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		fmt.Fprintf(w, "%s", ij)
	})

	http.ListenAndServe(":5650", r)
}
