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
	"strings"

	"github.com/julienschmidt/httprouter"
)

type EmberResponse struct {
	Error     bool   `json:"error"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type AddonPost struct {
	URL       	  string `json:"url"`
	Filename  	  string `json:"filename"`
	ExtractTo 	  string `json:"extractTo"`
	Authorization string `json:"authorizationHeader"`
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

		req, err := http.NewRequest("GET", a.URL, nil)

		if err != nil {
			log.Fatal(err)
		}

		if (strings.Contains(a.URL, "bult.test")) {
			req.Header.Add("Authorization", a.Authorization)
			log.Println("Added Bult authorization header.")
		}

		client := &http.Client{}

		log.Println("Running request.")
		response, err := client.Do(req);

		if err != nil {
			log.Println("Client request failed.")
			log.Fatal(err)
		}

		defer response.Body.Close()

		log.Println("Copying output to file.")
		_, err = io.Copy(output, response.Body)

		if err != nil {
			log.Println("Copying request failed.")
			log.Fatal(err)
		}

		log.Println("Opening zip reader.")
		reader, _ := zip.OpenReader(a.ExtractTo + "/" + a.Filename + ".zip")

		if err := os.MkdirAll(a.ExtractTo+"/"+a.Filename, 0755); err != nil {
			log.Fatal(err)
		}

		log.Println("Looping through reader.")
		for _, file := range reader.File {
			log.Println("Looping through " + file.Name)
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

