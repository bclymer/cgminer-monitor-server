package controllers

import (
	"github.com/bclymer/cgminer-monitor-server/services"
	"github.com/codegangsta/martini"
	_ "io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	logger = configureLogger()
)

func SetupApi(martiniServer *martini.ClassicMartini) {
	martiniServer.Get("/stats", makeJson, getStats)
	martiniServer.Get("/stats/:day", makeJson, getStatsForDay)
}

func makeJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func getStats(w http.ResponseWriter, r *http.Request) (int, string) {
	return 200, services.GetToday()
}

func getStatsForDay(w http.ResponseWriter, r *http.Request, params martini.Params) string {
	return "Hello " + params["_1"]
}

func configureLogger() *log.Logger {
	// NOTE these file permissions are restricted by umask, so they probably won't work right.
	err := os.MkdirAll("./log", 0775)
	if err != nil {
		panic(err)
	}
	logFile, err := os.OpenFile("./log/bc-cgminer-server-api.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		panic(err)
	}

	logger := log.New(logFile, "", log.Ldate|log.Ltime)

	return logger
}
