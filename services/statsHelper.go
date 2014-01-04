package services

import (
	"encoding/json"
	"github.com/bclymer/cgminer-monitor-server/models"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	statsFolder = "./stats/"
)

var (
	logger        = configureLogger()
	writerChannel = make(chan (string))
)

func init() {
	os.Mkdir(statsFolder, 7777)
	go addFileToQueue()
}

func GetToday() string {
	fileName := statsFolder + time.Now().Format("2006-01-02")
	fileContent, _ := ioutil.ReadFile(fileName)
	return string(fileContent)
}

func ProcessFile(fileName string) {
	writerChannel <- fileName
}

// Synchronize file writing
func addFileToQueue() bool {
	for {
		select {
		case fileName := <-writerChannel:
			addFile(fileName)
			break
		}
	}
}

func addFile(tempFileName string) bool {
	tempFileContent, err := ioutil.ReadFile(tempFileName)
	if err != nil {
		logger.Println("Failed reading", tempFileName, err)
		return false
	}
	var minerStats models.CgMinerStats
	err = json.Unmarshal(tempFileContent, &minerStats)
	if err != nil {
		logger.Println("Failed decoding", tempFileContent, err)
		return false
	}

	fileName, err := getFileForTime(minerStats.When)
	if err != nil {
		logger.Println("AddFile", err)
		return false
	}
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Println("Failed reading", fileName, err)
		return false
	}
	var fullStats models.FullStats
	if len(content) == 0 {
		fullStats = make(map[string][]models.DeviceStats)
	} else {
		err = json.Unmarshal(content, &fullStats)
		if err != nil {
			logger.Println("Failed decoding", fileName, err)
			return false
		}
	}

	addClientStatsToFullStats(minerStats, fullStats)
	fullStatsBytes, err := json.Marshal(fullStats)
	if err != nil {
		logger.Println("Failed to marshal fullStats", err)
		return false
	}
	err = ioutil.WriteFile(fileName, fullStatsBytes, 0644)
	if err != nil {
		logger.Println("Failed writing the file", err)
		return false
	}
	os.Remove(tempFileName)
	return true
}

func addClientStatsToFullStats(minerStats models.CgMinerStats, fullStats models.FullStats) {
	device, ok := fullStats[minerStats.DeviceName]
	if !ok {
		fullStats[minerStats.DeviceName] = make([]models.DeviceStats, len(minerStats.Devs))
		device, _ = fullStats[minerStats.DeviceName]
	}
	for i, minerDevice := range minerStats.Devs {
		device[i].When = append(device[i].When, minerStats.When)
		device[i].GPU = append(device[i].GPU, minerDevice.GPU)
		device[i].Enabled = append(device[i].Enabled, minerDevice.Enabled)
		device[i].Status = append(device[i].Status, minerDevice.Status)
		device[i].Temperature = append(device[i].Temperature, minerDevice.Temperature)
		device[i].FanSpeed = append(device[i].FanSpeed, minerDevice.FanSpeed)
		device[i].FanPercent = append(device[i].FanPercent, minerDevice.FanPercent)
		device[i].GpuClock = append(device[i].GpuClock, minerDevice.GpuClock)
		device[i].MemClock = append(device[i].MemClock, minerDevice.MemClock)
		device[i].GpuVoltage = append(device[i].GpuVoltage, minerDevice.GpuVoltage)
		device[i].GpuActivity = append(device[i].GpuActivity, minerDevice.GpuActivity)
		device[i].Powertune = append(device[i].Powertune, minerDevice.Powertune)
		device[i].MhsAv = append(device[i].MhsAv, minerDevice.MhsAv)
		device[i].MhsFiveSeconds = append(device[i].MhsFiveSeconds, minerDevice.MhsFiveSeconds)
		device[i].Accepted = append(device[i].Accepted, minerDevice.Accepted)
		device[i].Rejected = append(device[i].Rejected, minerDevice.Rejected)
		device[i].HardwareErrors = append(device[i].HardwareErrors, minerDevice.HardwareErrors)
		device[i].Utility = append(device[i].Utility, minerDevice.Utility)
		device[i].Intensity = append(device[i].Intensity, minerDevice.Intensity)
		device[i].LastSharePool = append(device[i].LastSharePool, minerDevice.LastSharePool)
		device[i].LastShareTime = append(device[i].LastShareTime, minerDevice.LastShareTime)
		device[i].TotalMh = append(device[i].TotalMh, minerDevice.TotalMh)
		device[i].DiffOneWork = append(device[i].DiffOneWork, minerDevice.DiffOneWork)
		device[i].DiffAccepted = append(device[i].DiffAccepted, minerDevice.DiffAccepted)
		device[i].DiffRejected = append(device[i].DiffRejected, minerDevice.DiffRejected)
		device[i].LastShareDiff = append(device[i].LastShareDiff, minerDevice.LastShareDiff)
		device[i].LastValidWorkd = append(device[i].LastValidWorkd, minerDevice.LastValidWorkd)
		device[i].DeviceHardwarePct = append(device[i].DeviceHardwarePct, minerDevice.DeviceHardwarePct)
		device[i].DeviceRejectedPct = append(device[i].DeviceRejectedPct, minerDevice.DeviceRejectedPct)
		device[i].DeviceElapsed = append(device[i].DeviceElapsed, minerDevice.DeviceElapsed)
	}
}

func getFileForTime(unixTime int64) (string, error) {
	t := time.Unix(unixTime, 0)
	fileName := statsFolder + t.Format("2006-01-02")
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	defer file.Close()
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func configureLogger() *log.Logger {
	// NOTE these file permissions are restricted by umask, so they probably won't work right.
	err := os.MkdirAll("./log", 0775)
	if err != nil {
		panic(err)
	}
	logFile, err := os.OpenFile("./log/bc-cgminer-server-statsHelper.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		panic(err)
	}

	logger := log.New(logFile, "", log.Ldate|log.Ltime)

	return logger
}
