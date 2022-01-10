package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	odroid "github.com/Gonzih/odroid-show-golang"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	netstatus "github.com/shirou/gopsutil/net"
	"github.com/spf13/cobra"
)

const (
	tempLower = 50.0
	tempUpper = 75.0
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

type gpuData struct {
	number float64
	label  string
}

func GpuStatus(odr *odroid.OdroidShowBoard) {
	if NVidiaSMIAvailable() {
		odr.Fg(odroid.ColorCyan)
		odr.WriteString("GPU:")
		odr.ColorReset()

		util := NVidiaUtilization()
		temp := NVidiaTemperature()
		utilNumber, err := strconv.ParseFloat(strings.Replace(util, "%", "", 1), 64)
		if err != nil {
			must(fmt.Errorf(`Error parsing "%s" to float: %s`, util, err))
		}
		tempNumber, err := strconv.ParseFloat(strings.Replace(temp, "C", "", 1), 64)
		if err != nil {
			must(fmt.Errorf(`Error parsing "%s" to float: %s`, util, err))
		}

		data := []gpuData{
			gpuData{
				label:  temp,
				number: tempNumber,
			},
			gpuData{
				label:  util,
				number: utilNumber,
			},
		}

		for _, dp := range data {

			color := odroid.ColorGreen

			if dp.number < tempUpper && dp.number > tempLower {
				color = odroid.ColorYellow
			} else if dp.number >= tempUpper {
				color = odroid.ColorRed
			}

			odr.Fg(color)
			odr.WriteString(fmt.Sprintf("%4s ", dp.label))
			odr.ColorReset()
		}

	}
}

func NvmeStatus(odr *odroid.OdroidShowBoard) {
	if NvmeAvailable() {
		odr.Fg(odroid.ColorRed)
		odr.WriteString("NVME:")
		odr.ColorReset()

		temp := NvmeTemperature()
		num, err := strconv.ParseFloat(strings.Replace(temp, "C", "", 1), 64)
		if err != nil {
			must(fmt.Errorf(`Error parsing "%s" to float: %s`, temp, err))
		}

		color := odroid.ColorGreen

		if num < tempUpper && num > tempLower {
			color = odroid.ColorYellow
		} else if num >= tempUpper {
			color = odroid.ColorRed
		}

		odr.Fg(color)
		odr.WriteString(fmt.Sprintf("%4s ", temp))
		odr.ColorReset()

	}
}

func NetworkStatus(odr *odroid.OdroidShowBoard) {
	ifaces, err := netstatus.Interfaces()

	if err != nil {
		log.Fatal(err)
	}

	var addrs []string

	for _, iface := range ifaces {
		for _, addr := range iface.Addrs {
			ip, _, err := net.ParseCIDR(addr.Addr)
			if err == nil {
				if !ip.IsLoopback() {
					ip4 := ip.To4()
					if ip4 != nil {
						ip4str := ip4.String()
						if strings.HasPrefix(ip4str, "192.168.2.") {
							addrs = append(addrs, ip4str)
						}
					}
				}
			}
		}
	}

	odr.Fg(odroid.ColorMagenta)
	odr.WriteString("IP:")
	odr.ColorReset()
	odr.WriteString(strings.Join(addrs, " "))
}

var temperatureCleanKeyReg = regexp.MustCompile("coretemp|input|k10temp|_|composite|sensor")
var temperatureReplaceKeyReg = regexp.MustCompile("it8665")

func SensorsStatus(odr *odroid.OdroidShowBoard) {
	temps, err := host.SensorsTemperatures()

	if err != nil {
		log.Fatal(err)
	}

	var builder strings.Builder
	i := 0

	for _, temp := range temps {
		k := temp.SensorKey
		if strings.Contains(temp.SensorKey, "_input") {
			k = temperatureCleanKeyReg.ReplaceAllString(k, "")
			k = temperatureReplaceKeyReg.ReplaceAllString(k, "cpu")
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
	mem := float64(v.Total) / 1000000
	label := "MB"

	if mem > 1000 {
		mem = mem / 1000
		label = "GB"
	}
	odr.WriteString(fmt.Sprintf("%.0f%% %.2f%s", v.UsedPercent, mem, label))
}

var mountPoints []string
var serialPort string
var reportLoad bool
var reportMemory bool
var reportOSStatus bool
var reportNetworkStatus bool
var reportDiskStatus bool
var reportGpuStatus bool
var reportSensorsStatus bool

var rootCmd = &cobra.Command{
	Use:   "odroid-status-screen-go",
	Short: "Display system status on odroid board",
	Run: func(cmd *cobra.Command, args []string) {
		if len(serialPort) == 0 {
			files, err := sh("bash", "-c", "ls -1 /dev/ttyUSB*")
			must(err)
			ports := strings.Split(string(files), "\n")
			serialPort = ports[0]
		}
		run()
	},
}

func init() {
	// cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&serialPort, "serial-port", "", "", "Serial port to communicate with odroid on")
	rootCmd.PersistentFlags().StringSliceVarP(&mountPoints, "mount-point", "", []string{"/"}, "Mount points to monitor for usage")
	rootCmd.PersistentFlags().BoolVarP(&reportLoad, "report-load", "", true, "Report system load")
	rootCmd.PersistentFlags().BoolVarP(&reportMemory, "report-memory", "", true, "Report system memory usage")
	rootCmd.PersistentFlags().BoolVarP(&reportOSStatus, "report-os", "", true, "Report system OS status")
	rootCmd.PersistentFlags().BoolVarP(&reportNetworkStatus, "report-network", "", true, "Report system Network stats")
	rootCmd.PersistentFlags().BoolVarP(&reportDiskStatus, "report-disk", "", true, "Report system Disk usage")
	rootCmd.PersistentFlags().BoolVarP(&reportGpuStatus, "report-gpu", "", true, "Report system GPU usage")
	rootCmd.PersistentFlags().BoolVarP(&reportSensorsStatus, "report-sensors", "", true, "Report system Temp sensors")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func run() {
	odr, err := odroid.NewOdroidShowBoard(serialPort)
	must(err)
	odr.Clear()

	for {
		odr.CursorReset()

		if reportLoad {
			LoadStatus(odr)
			odr.Ln()
		}
		if reportMemory {
			MemStatus(odr)
			odr.Ln()
		}
		if reportOSStatus {
			OSStatus(odr)
			odr.Ln()
		}
		if reportNetworkStatus {
			NetworkStatus(odr)
			odr.Ln()
		}
		if reportDiskStatus {
			DisksStatus(odr, mountPoints)
			odr.Ln()
		}
		if reportGpuStatus {
			GpuStatus(odr)
			odr.Ln()
		}
		if reportSensorsStatus {
			SensorsStatus(odr)
		}

		must(odr.Sync())
		time.Sleep(time.Second * 2)
	}
}

func main() {
	must(rootCmd.Execute())
}
