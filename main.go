package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/Gonzih/odroid-show-golang"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	netstatus "github.com/shirou/gopsutil/net"
)

func DisksStatus(odr *odroid.OdroidShowBoard, paths []string) {
	odr.Fg(odroid.ColorRed)
	odr.WriteString("DISKS:")
	odr.ColorReset()

	for _, path := range paths {

		usage, err := disk.Usage(path)

		if err != nil {
			log.Fatal(err)
		}
		odr.WriteString(fmt.Sprintf("%s %.0f%% ", path, usage.UsedPercent))
	}
}

func NetworkStatus(odr *odroid.OdroidShowBoard) {
	ifaces, err := netstatus.Interfaces()

	if err != nil {
		log.Fatal(err)
	}

	var addrs strings.Builder

	for _, iface := range ifaces {
		for _, addr := range iface.Addrs {
			ip, _, err := net.ParseCIDR(addr.Addr)
			if err == nil {
				if !ip.IsLoopback() {
					ip4 := ip.To4()
					if ip4 != nil {
						addrs.WriteString(ip4.String())
						addrs.WriteString(" ")
					}
				}
			}
		}
	}

	odr.Fg(odroid.ColorMagenta)
	odr.WriteString("ADDR:")
	odr.ColorReset()
	odr.WriteString(addrs.String())
}

var temperatureKeyReg = regexp.MustCompile("coretemp|input|_")

func SensorsStatus(odr *odroid.OdroidShowBoard) {
	temps, err := host.SensorsTemperatures()

	if err != nil {
		log.Fatal(err)
	}

	var builder strings.Builder
	tempLower := 40.0
	tempUpper := 75.0
	i := 0

	for _, temp := range temps {
		k := temp.SensorKey
		if strings.Contains(temp.SensorKey, "_input") {
			k = temperatureKeyReg.ReplaceAllString(k, "")
			color := odroid.ColorGreen

			if temp.Temperature < tempUpper && temp.Temperature > tempLower {
				color = odroid.ColorYellow
			} else if temp.Temperature >= tempUpper {
				color = odroid.ColorRed
			}

			prefix := ""

			if i%2 == 0 {
				prefix = "\r\n"
			}

			i++
			builder.WriteString(fmt.Sprintf("%s\033[3%dm%s:\033[3%dm%.0fC ", prefix, odroid.ColorWhite, k, color, temp.Temperature))
		}
	}

	odr.Fg(odroid.ColorBlue)
	odr.WriteString("TEMP:")
	odr.ColorReset()
	odr.WriteString(builder.String())
}

func OSStatus(odr *odroid.OdroidShowBoard) {
	uptime, err := host.Uptime()

	if err != nil {
		log.Fatal(err)
	}

	duration, _ := time.ParseDuration(fmt.Sprintf("%ds", uptime))
	odr.Fg(odroid.ColorGreen)
	odr.WriteString("UP:")
	odr.ColorReset()
	odr.WriteString(duration.String())
}

func LoadStatus(odr *odroid.OdroidShowBoard) {
	v, err := load.Avg()

	if err != nil {
		log.Fatal(err)
	}

	odr.Fg(odroid.ColorCyan)
	odr.WriteString("LOAD:")
	odr.ColorReset()
	odr.WriteString(fmt.Sprintf("%.2f %.2f %.2f", v.Load1, v.Load5, v.Load15))
}

func MemStatus(odr *odroid.OdroidShowBoard) {
	v, err := mem.VirtualMemory()

	if err != nil {
		log.Fatal(err)
	}

	odr.Fg(odroid.ColorYellow)
	odr.WriteString("MEM:")
	odr.ColorReset()
	mem := v.Total / 1000000
	label := "MB"

	if mem > 1000 {
		mem = mem / 1000
		label = "GB"
	}
	odr.WriteString(fmt.Sprintf("%.2f%% %v%s", v.UsedPercent, mem, label))
}

type sliceFlags []string

func (i *sliceFlags) String() string {
	return strings.Join(*i, ", ")
}

func (i *sliceFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var mountPoints sliceFlags

func init() {
	flag.Var(&mountPoints, "mount-point", "Mount points to report usage for")
	flag.Parse()
	if len(mountPoints) == 0 {
		mountPoints = append(mountPoints, "/")
	}
}

func main() {
	odr, err := odroid.NewOdroidShowBoard("/dev/ttyUSB0")

	if err != nil {
		log.Fatal(err)
	}

	odr.Clear()

	for {
		odr.CursorReset()
		LoadStatus(odr)
		odr.Ln()
		MemStatus(odr)
		odr.WriteString(" ")
		OSStatus(odr)
		odr.Ln()
		NetworkStatus(odr)
		odr.Ln()
		DisksStatus(odr, mountPoints)
		odr.Ln()
		odr.Ln()
		SensorsStatus(odr)

		err = odr.Sync()

		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Second)
	}

}
