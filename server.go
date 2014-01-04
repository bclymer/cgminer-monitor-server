package main

import (
	"encoding/json"
	"github.com/bclymer/cgminer-monitor-server/controllers"
	"github.com/bclymer/cgminer-monitor-server/models"
	"github.com/bclymer/cgminer-monitor-server/services"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/gzip"
	"github.com/codegangsta/martini-contrib/render"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	tempFolder = "./temp/"
)

var (
	// Mapping device names to that machines stats
	logger = configureLogger()
	config models.Config
)

func main() {
	os.Mkdir(tempFolder, 7777)
	loadConfig()

	martiniServer := martini.Classic()
	martiniServer.Use(gzip.All())

	martiniServer.Map(logger)

	controllers.SetupApi(martiniServer)

	martiniServer.Get("/", func(renderer render.Render) {
		renderer.HTML(200, "home", nil)
	})

	martiniServer.Post("/stats", func(w http.ResponseWriter, r *http.Request) (int, string) {
		if r.Header.Get("Server-Password") != config.ServerPassword {
			return http.StatusUnauthorized, ""
		}
		reader, err := r.MultipartReader()
		if err != nil {
			logger.Println("Error:", err)
			return http.StatusInternalServerError, failure(err)
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
			dst, err := os.Create(tempFolder + part.FileName())
			defer dst.Close()

			if err != nil {
				logger.Println("Error:", err)
				return http.StatusInternalServerError, failure(err)
			}

			if _, err := io.Copy(dst, part); err != nil {
				logger.Println("Error:", err)
				return http.StatusInternalServerError, failure(err)
			}
			logger.Println("Accepted", part.FileName())
			go services.ProcessFile(tempFolder + part.FileName())
		}
		return 201, success()
	})

	martiniServer.Use(render.Renderer(render.Options{
		Layout: "_layout",
	}))
	martiniServer.Use(martini.Recovery())
	http.ListenAndServe(":"+config.ServerPort, martiniServer)
}

func configureLogger() *log.Logger {
	// NOTE these file permissions are restricted by umask, so they probably won't work right.
	err := os.MkdirAll("./log", 0775)
	if err != nil {
		panic(err)
	}
	logFile, err := os.OpenFile("./log/bc-cgminer-server.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		panic(err)
	}

	logger := log.New(logFile, "", log.Ldate|log.Ltime)

	return logger
}

type Result struct {
	Success bool  `json:"success"`
	Error   error `json:"error"`
}

func success() string {
	str, _ := json.Marshal(Result{Success: true})
	return string(str)
}

func failure(err error) string {
	str, _ := json.Marshal(Result{Success: false, Error: err})
	return string(str)
}

func loadConfig() {
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, &config)
	if err != nil {
		panic(err)
	}
	logger.Println("Server is on port: " + config.ServerPort)
	logger.Println("Server password is: " + config.ServerPassword)
}
