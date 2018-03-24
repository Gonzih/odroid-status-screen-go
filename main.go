package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Gonzih/odroid-show-golang"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

func SensorsStatus(odr *odroid.OdroidShowBoard) {
	temps, _ := host.SensorsTemperatures()
	var builder strings.Builder
	reg := regexp.MustCompile("coretemp|input|_")
	tempLower := 40.0
	tempUpper := 75.0

	for _, temp := range temps {
		k := temp.SensorKey
		if strings.Contains(temp.SensorKey, "_input") {
			k = reg.ReplaceAllString(k, "")
			color := odroid.ColorBlue

			if temp.Temperature < tempUpper && temp.Temperature > tempLower {
				color = odroid.ColorYellow
			} else if temp.Temperature >= tempUpper {
				color = odroid.ColorRed
			}

			builder.WriteString(fmt.Sprintf("\r\n\033[3%dm%s: \033[3%dm%.0fC", odroid.ColorWhite, k, color, temp.Temperature))
		}
	}

	odr.Fg(odroid.ColorBlue)
	odr.WriteString("Temp:")
	odr.ColorReset()
	odr.WriteString(builder.String())
}

func OSStatus(odr *odroid.OdroidShowBoard) {
	uptime, _ := host.Uptime()
	duration, _ := time.ParseDuration(fmt.Sprintf("%ds", uptime))
	odr.Fg(odroid.ColorRed)
	odr.WriteString("UP:")
	odr.ColorReset()
	odr.WriteString(duration.String())
}

func LoadStatus(odr *odroid.OdroidShowBoard) {
	v, _ := load.Avg()
	odr.Fg(odroid.ColorCyan)
	odr.WriteString("LOAD:")
	odr.ColorReset()
	odr.WriteString(fmt.Sprintf("%.2f %.2f %.2f", v.Load1, v.Load5, v.Load15))
}

func MemStatus(odr *odroid.OdroidShowBoard) {
	v, _ := mem.VirtualMemory()
	odr.Fg(odroid.ColorYellow)
	odr.WriteString("MEM: ")
	odr.ColorReset()
	mem := v.Total / 1000000
	label := "MB"

	if mem > 1000 {
		mem = mem / 1000
		label = "GB"
	}
	odr.WriteString(fmt.Sprintf("%.0f%% out of %v%s", v.UsedPercent, mem, label))
}

func main() {
	odr, err := odroid.NewOdroidShowBoard("/dev/ttyUSB0")

	if err != nil {
		log.Fatal(err)
	}

	defer odr.Sync()

	odr.Clear()

	for {
		odr.CursorReset()
		LoadStatus(odr)
		odr.Ln()
		MemStatus(odr)
		odr.Ln()
		OSStatus(odr)
		odr.Ln()
		SensorsStatus(odr)

		odr.Sync()
		time.Sleep(time.Second)
	}

}
