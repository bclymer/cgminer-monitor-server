package models

import ()

type CgMinerStats struct {
	DeviceName string `json:"deviceName"`
	When       int64  `json:"when"`
	Devs       []struct {
		GPU               int     `json:"GPU"`
		Enabled           string  `json:"Enabled"`
		Status            string  `json:"Status"`
		Temperature       float64 `json:"Temperature"`
		FanSpeed          int     `json:"Fan Speed"`
		FanPercent        int     `json:"Fan Percent"`
		GpuClock          int     `json:"GPU Clock"`
		MemClock          int     `json:"Memory Clock"`
		GpuVoltage        float64 `json:"GPU Voltage"`
		GpuActivity       int     `json:"GPU Activity"`
		Powertune         int     `json:"Powertune"`
		MhsAv             float64 `json:"MHS av"`
		MhsFiveSeconds    float64 `json:"MHS 5s"`
		Accepted          int     `json:"Accepted"`
		Rejected          int     `json:"Rejected"`
		HardwareErrors    int     `json:"Hardware Errors"`
		Utility           float64 `json:"Utility"`
		Intensity         string  `json:"Intensity"`
		LastSharePool     uint64  `json:"Last Share Pool"`
		LastShareTime     uint64  `json:"Last Share Time"`
		TotalMh           float64 `json:"Total MH"`
		DiffOneWork       uint64  `json:"Diff1 Work"`
		DiffAccepted      float64 `json:"Difficulty Accepted"`
		DiffRejected      float64 `json:"Difficulty Rejected"`
		LastShareDiff     float64 `json:"Last Share Difficulty"`
		LastValidWorkd    uint64  `json:"Last Valid Work"`
		DeviceHardwarePct float64 `json:"Device Hardware%"`
		DeviceRejectedPct float64 `json:"Device Rejected%"`
		DeviceElapsed     uint64  `json:"Device Elapsed"`
	} `json:"DEVS"`
}

type FullStats map[string][]DeviceStats

type DeviceStats struct {
	When              []int64   `json:"when"`
	GPU               []int     `json:"gpu"`
	Enabled           []string  `json:"enabled"`
	Status            []string  `json:"status"`
	Temperature       []float64 `json:"temperature"`
	FanSpeed          []int     `json:"fanSpeed"`
	FanPercent        []int     `json:"fanPercent"`
	GpuClock          []int     `json:"gpuClock"`
	MemClock          []int     `json:"memoryClock"`
	GpuVoltage        []float64 `json:"gpuVoltage"`
	GpuActivity       []int     `json:"gpuActivity"`
	Powertune         []int     `json:"powertune"`
	MhsAv             []float64 `json:"mhsAv"`
	MhsFiveSeconds    []float64 `json:"mhsLastFiveSeconds"`
	Accepted          []int     `json:"accepted"`
	Rejected          []int     `json:"rejected"`
	HardwareErrors    []int     `json:"hardwareErrors"`
	Utility           []float64 `json:"utility"`
	Intensity         []string  `json:"intensity"`
	LastSharePool     []uint64  `json:"lastSharePool"`
	LastShareTime     []uint64  `json:"lastShareTime"`
	TotalMh           []float64 `json:"totalMH"`
	DiffOneWork       []uint64  `json:"diffOneWork"`
	DiffAccepted      []float64 `json:"difficultyAccepted"`
	DiffRejected      []float64 `json:"difficultyRejected"`
	LastShareDiff     []float64 `json:"lastShareDifficulty"`
	LastValidWorkd    []uint64  `json:"lastValidWork"`
	DeviceHardwarePct []float64 `json:"deviceHardwarePct"`
	DeviceRejectedPct []float64 `json:"deviceRejectedPct"`
	DeviceElapsed     []uint64  `json:"deviceElapsed"`
}

type Config struct {
	ServerPort     string `json:"serverPort"`
	ServerPassword string `json:"serverPassword"`
}
