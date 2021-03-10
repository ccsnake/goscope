package goscope

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func getAppName(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{
		"applicationName": Config.ApplicationName,
	})
}

// getSystemInfoHandler is the controller to show system information of the current host in GoScope API.
func getSystemInfoHandler(c *gin.Context) {
	cpuStatus, _ := cpu.Info()
	firstCPU := cpuStatus[0]
	memoryStatus, _ := mem.VirtualMemory()
	swapStatus, _ := mem.SwapMemory()
	hostStatus, _ := host.Info()
	diskStatus, _ := disk.Usage("/")

	environment := make(map[string]string)

	env := os.Environ()
	for i := range env {
		variable := strings.SplitN(env[i], "=", 2)
		environment[variable[0]] = variable[1]
	}

	responseBody := systemInformationResponse{
		ApplicationName: Config.ApplicationName,
		CPU: systemInformationResponseCPU{
			CoreCount: fmt.Sprintf("%d Cores", firstCPU.Cores),
			ModelName: firstCPU.ModelName,
		},
		Memory: systemInformationResponseMemory{
			Available: fmt.Sprintf("%.2f GB", float64(memoryStatus.Available)/BytesInOneGigabyte),
			Total:     fmt.Sprintf("%.2f GB", float64(memoryStatus.Total)/BytesInOneGigabyte),
			UsedSwap:  fmt.Sprintf("%.2f%%", swapStatus.UsedPercent),
		},
		Host: systemInformationResponseHost{
			HostOS:        hostStatus.OS,
			HostPlatform:  hostStatus.Platform,
			Hostname:      hostStatus.Hostname,
			KernelArch:    hostStatus.KernelArch,
			KernelVersion: hostStatus.KernelVersion,
			Uptime:        fmt.Sprintf("%.2f hours", float64(hostStatus.Uptime)/SecondsInOneMinute/SecondsInOneMinute),
		},
		Disk: systemInformationResponseDisk{
			FreeSpace:     fmt.Sprintf("%.2f GB", float64(diskStatus.Free)/BytesInOneGigabyte),
			MountPath:     diskStatus.Path,
			PartitionType: diskStatus.Fstype,
			TotalSpace:    fmt.Sprintf("%.2f GB", float64(diskStatus.Total)/BytesInOneGigabyte),
		},
		Environment: environment,
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, responseBody)
}
