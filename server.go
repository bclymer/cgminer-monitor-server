package main

import (
	_ "fmt"
	"github.com/codegangsta/martini"
	"io"
	"net/http"
	"os"
	"json/encoding"
)

func main() {
	os.Mkdir("./stats", 7777)
	m := martini.Classic()

	m.Get("/", func() string {
		return "Hello world!"
	})

	m.Post("/stats", func(w http.ResponseWriter, r *http.Request) string {
		reader, err := r.MultipartReader()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return failure(err)
		}

		//copy each part to destination.
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if part.FileName() == "" {
				continue
			}
			dst, err := os.Create("./stats/" + part.FileName())
			defer dst.Close()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return failure(err)
			}

			if _, err := io.Copy(dst, part); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return failure(err)
			}
		}
		return success()
	})
	m.Run()
}

type Result struct {
	Success bool
	Error error
}

func success() Result {
	return json.Marshal(Result{success: true})
}

func failure(err error) Result {
	return json.Marshal(Result{success: false,error: err})
}