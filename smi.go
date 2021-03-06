package main

import (
	"log"
	"regexp"
	"strings"
)

func NVidiaSMIAvailable() bool {
	out, err := sh("bash", "-c", "which nvidia-smi")
	if err != nil {
		log.Println(err)
		return false
	}

	log.Println(strings.Trim(string(out), "\r\n "))
	return true
}

// GPU Current Temp            : 62 C
var tempRegexp = regexp.MustCompile(`GPU Current Temp\s+:\s+(\d+\s?C)`)

func parseNvidiaTemperatureOutput(input []byte) string {
	matches := tempRegexp.FindAllSubmatch(input, -1)
	if len(matches) == 0 || len(matches[0]) < 2 {
		return "0C"
	}

	return strings.Replace(string(matches[0][1]), " ", "", -1)
}

func NVidiaTemperature() string {
	out, err := sh("nvidia-smi", "-q", "-d", "TEMPERATURE")
	if err != nil {
		panic(err)
	}

	return parseNvidiaTemperatureOutput(out)
}

// Gpu                         : 6 %
var utilRegexp = regexp.MustCompile(`Gpu\s+:\s+(\d+\s?%)`)

func parseUtilizationOutput(input []byte) string {
	matches := utilRegexp.FindAllSubmatch(input, -1)
	if len(matches) == 0 || len(matches[0]) < 2 {
		return "0%"
	}

	return strings.Replace(string(matches[0][1]), " ", "", -1)
}

func NVidiaUtilization() string {
	out, err := sh("nvidia-smi", "-q", "-d", "UTILIZATION")
	if err != nil {
		panic(err)
	}

	return parseUtilizationOutput(out)
}
