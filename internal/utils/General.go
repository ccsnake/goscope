package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	HoursInDay      = 24
	MinutesInHour   = 60
	SecondsInMinute = 60
)

// Check the wanted path is not in the do not log list.
func CheckExcludedPaths(path string) bool {
	exactMatches := []string{
		"",
		"/apple-touch-icon-precomposed.png",
		"/apple-touch-icon.png",
		"/goscope/css/light.css.map",
		"/goscope/css/dark.css.map",
		"/favicon.ico",
		"/site.webmanifest",
	}

	for i := range exactMatches {
		if path == exactMatches[i] {
			return false
		}
	}

	partialMatches := []string{
		"/goscope",
		".manifest",
		".css",
		".js",
		".ttf",
		".woff",
		".svg",
		".ico",
		".png",
		".jpg",
		".webp",
	}

	for i := range partialMatches {
		if strings.Contains(path, partialMatches[i]) {
			return false
		}
	}

	return true
}

func PrettifyJSON(rawString string) string {
	if rawString == "" {
		return ""
	}

	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(rawString), "", "    ")

	if err != nil {
		return rawString
	}

	return prettyJSON.String()
}

func EpochToTimeAgoHappened(epoch int) string {
	date := time.Unix(int64(epoch), 0)
	diff := time.Since(date)

	if diff.Seconds() < SecondsInMinute {
		return fmt.Sprintf("%.2f s", diff.Seconds())
	} else if diff.Minutes() < MinutesInHour {
		return fmt.Sprintf("%.0f m", diff.Minutes())
	} else if diff.Hours() < HoursInDay {
		return fmt.Sprintf("%.0f h", diff.Hours())
	}

	return fmt.Sprintf("%.0f d", math.Round(diff.Hours()/HoursInDay))
}

func EpochToHumanReadable(epoch int) string {
	date := time.Unix(int64(epoch), 0)
	return date.Format(time.RFC1123Z)
}

func ResponseStatusColor(responseStatus interface{}) string {

	response, err := strconv.ParseInt(fmt.Sprintf("%v", responseStatus), 10, 32)
	if err != nil {
		return "badge-info"
	}
	switch response {
	case http.StatusOK:
		return "badge-success"
	case http.StatusCreated:
		return "badge-success"
	case http.StatusAccepted:
		return "badge-success"
	case http.StatusNonAuthoritativeInfo:
		return "badge-success"
	case http.StatusNoContent:
		return "badge-success"
	case http.StatusMultipleChoices:
		return "badge-info"
	case http.StatusMovedPermanently:
		return "badge-info"
	case http.StatusFound:
		return "badge-info"
	case http.StatusSeeOther:
		return "badge-info"
	case http.StatusNotModified:
		return "badge-info"
	case http.StatusUseProxy:
		return "badge-info"
	case http.StatusTemporaryRedirect:
		return "badge-info"
	case http.StatusPermanentRedirect:
		return "badge-info"
	case http.StatusBadRequest:
		return "badge-warning"
	case http.StatusUnauthorized:
		return "badge-warning"
	case http.StatusPaymentRequired:
		return "badge-warning"
	case http.StatusForbidden:
		return "badge-warning"
	case http.StatusNotFound:
		return "badge-warning"
	case http.StatusTeapot:
		return "badge-warning"
	case http.StatusUnprocessableEntity:
		return "badge-warning"
	case http.StatusInternalServerError:
		return "badge-danger"
	case http.StatusNotImplemented:
		return "badge-danger"
	case http.StatusBadGateway:
		return "badge-danger"
	case http.StatusServiceUnavailable:
		return "badge-danger"

	default:
		return "badge-info"

	}
}
