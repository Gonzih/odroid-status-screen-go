package main

import (
	"log"
	"regexp"
	"strings"
)

func NvmeAvailable() bool {
	out, err := sh("bash", "-c", "which nvme")
	if err != nil {
		log.Println(err)
		return false
	}

	log.Println(strings.Trim(string(out), "\r\n "))
	return true
}

// GPU Current Temp            : 62 C
var nvmeTempRegexp = regexp.MustCompile(`temperature\s+:\s+(\d+\s?C)`)

func parseNvmeTemperatureOutput(input []byte) string {
	matches := nvmeTempRegexp.FindAllSubmatch(input, -1)
	if len(matches) == 0 || len(matches[0]) < 2 {
		return "0C"
	}

	return strings.Replace(string(matches[0][1]), " ", "", -1)
}

func NvmeTemperature() string {
	out, err := sh("nvme", "smart-log", "/dev/nvme0n1")
	if err != nil {
		panic(err)
	}

	return parseNvmeTemperatureOutput(out)
}
