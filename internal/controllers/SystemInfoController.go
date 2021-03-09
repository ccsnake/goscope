package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/averageflow/goscope/v3/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

const (
	BytesInOneGigabyte = 1073741824
	SecondsInOneMinute = 60
)

type SystemInformationResponse struct {
	ApplicationName string                          `json:"applicationName"`
	CPU             SystemInformationResponseCPU    `json:"cpu"`
	Disk            SystemInformationResponseDisk   `json:"disk"`
	Host            SystemInformationResponseHost   `json:"host"`
	Memory          SystemInformationResponseMemory `json:"memory"`
	Environment     map[string]string               `json:"environment"`
}

type SystemInformationResponseCPU struct {
	CoreCount string `json:"coreCount"`
	ModelName string `json:"modelName"`
}

type SystemInformationResponseDisk struct {
	FreeSpace     string `json:"freeSpace"`
	MountPath     string `json:"mountPath"`
	PartitionType string `json:"partitionType"`
	TotalSpace    string `json:"totalSpace"`
}

type SystemInformationResponseMemory struct {
	Available string `json:"availableMemory"`
	Total     string `json:"totalMemory"`
	UsedSwap  string `json:"usedSwap"`
}

type SystemInformationResponseHost struct {
	HostOS        string `json:"hostOS"`
	HostPlatform  string `json:"hostPlatform"`
	Hostname      string `json:"hostname"`
	KernelArch    string `json:"kernelArch"`
	KernelVersion string `json:"kernelVersion"`
	Uptime        string `json:"uptime"`
}

func GetAppName(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{
		"applicationName": utils.Config.ApplicationName,
	})
}

// ShowSystemInfo is the controller to show system information of the current host in GoScope API.
func ShowSystemInfo(c *gin.Context) {
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

	responseBody := SystemInformationResponse{
		ApplicationName: utils.Config.ApplicationName,
		CPU: SystemInformationResponseCPU{
			CoreCount: fmt.Sprintf("%d Cores", firstCPU.Cores),
			ModelName: firstCPU.ModelName,
		},
		Memory: SystemInformationResponseMemory{
			Available: fmt.Sprintf("%.2f GB", float64(memoryStatus.Available)/BytesInOneGigabyte),
			Total:     fmt.Sprintf("%.2f GB", float64(memoryStatus.Total)/BytesInOneGigabyte),
			UsedSwap:  fmt.Sprintf("%.2f%%", swapStatus.UsedPercent),
		},
		Host: SystemInformationResponseHost{
			HostOS:        hostStatus.OS,
			HostPlatform:  hostStatus.Platform,
			Hostname:      hostStatus.Hostname,
			KernelArch:    hostStatus.KernelArch,
			KernelVersion: hostStatus.KernelVersion,
			Uptime:        fmt.Sprintf("%.2f hours", float64(hostStatus.Uptime)/SecondsInOneMinute/SecondsInOneMinute),
		},
		Disk: SystemInformationResponseDisk{
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
