package goscope

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"

	"github.com/averageflow/goscope/v3/internal/repository"

	"github.com/gin-gonic/gin"
)

func requestListPageHandler(c *gin.Context) {
	offsetQuery := c.DefaultQuery("offset", "0")
	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)

	searchValue := c.Query("search")

	if searchValue != "" {
		variables := gin.H{
			"applicationName": Config.ApplicationName,
			"entriesPerPage":  Config.GoScopeEntriesPerPage,
			"data":            repository.FetchSearchRequests(DB, Config.ApplicationID, Config.GoScopeEntriesPerPage, Config.GoScopeDatabaseType, searchValue, nil, int(offset)),
			"baseURL":         Config.BaseURL,
			"offset":          int(offset),
			"searchValue":     searchValue,
		}

		c.HTML(http.StatusOK, "goscope-views/Requests.gohtml", variables)
	} else {
		variables := gin.H{
			"applicationName": Config.ApplicationName,
			"entriesPerPage":  Config.GoScopeEntriesPerPage,
			"data":            repository.FetchRequestList(DB, Config.ApplicationID, Config.GoScopeEntriesPerPage, int(offset)),
			"baseURL":         Config.BaseURL,
			"offset":          int(offset),
			"searchValue":     searchValue,
		}

		c.HTML(http.StatusOK, "goscope-views/Requests.gohtml", variables)
	}
}

func logListPageHandler(c *gin.Context) {
	offsetQuery := c.DefaultQuery("offset", "0")
	offset, _ := strconv.ParseInt(offsetQuery, 10, 32)

	searchValue := c.Query("search")
	if searchValue != "" {
		variables := gin.H{
			"applicationName": Config.ApplicationName,
			"entriesPerPage":  Config.GoScopeEntriesPerPage,
			"data": repository.FetchSearchLogs(
				DB,
				Config.ApplicationID,
				Config.GoScopeEntriesPerPage,
				Config.GoScopeDatabaseType,
				searchValue,
				int(offset),
			),
			"baseURL":     Config.BaseURL,
			"offset":      int(offset),
			"searchValue": searchValue,
		}
		c.HTML(http.StatusOK, "goscope-views/Logs.gohtml", variables)
	} else {
		variables := gin.H{
			"applicationName": Config.ApplicationName,
			"entriesPerPage":  Config.GoScopeEntriesPerPage,
			"data": repository.FetchLogs(
				DB,
				Config.ApplicationID,
				Config.GoScopeEntriesPerPage,
				Config.GoScopeDatabaseType,
				int(offset),
			),
			"baseURL":     Config.BaseURL,
			"offset":      int(offset),
			"searchValue": searchValue,
		}
		c.HTML(http.StatusOK, "goscope-views/Logs.gohtml", variables)
	}
}

func logDetailsPageHandler(c *gin.Context) {
	var request RecordByURI

	err := c.ShouldBindUri(&request)
	if err != nil {
		log.Println(err.Error())
	}

	logDetails := repository.FetchDetailedLog(DB, request.UID)

	variables := gin.H{
		"applicationName": Config.ApplicationName,
		"data": gin.H{
			"logDetails": logDetails,
		},
		"baseURL": Config.BaseURL,
	}

	c.HTML(http.StatusOK, "goscope-views/LogDetails.gohtml", variables)
}

func requestDetailsPageHandler(c *gin.Context) {
	var request RecordByURI

	err := c.ShouldBindUri(&request)
	if err != nil {
		log.Println(err.Error())
	}

	requestDetails := repository.FetchDetailedRequest(DB, request.UID)
	responseDetails := repository.FetchDetailedResponse(DB, request.UID)

	variables := gin.H{
		"applicationName": Config.ApplicationName,
		"data": gin.H{
			"request":  requestDetails,
			"response": responseDetails,
		},
		"baseURL": Config.BaseURL,
	}

	c.HTML(http.StatusOK, "goscope-views/RequestDetails.gohtml", variables)
}

func systemInfoPageHandler(c *gin.Context) {
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

	c.HTML(http.StatusOK, "goscope-views/SystemInfo.gohtml", gin.H{
		"applicationName": Config.ApplicationName,
		"data":            responseBody,
		"baseURL":         Config.BaseURL,
	})
}
