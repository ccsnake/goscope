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

const (
	successBadge = "badge-success"
	infoBadge    = "badge-info"
	warningBadge = "badge-warning"
	dangerBadge  = "badge-danger"
)

// nolint:gocyclo .
func ResponseStatusColor(responseStatus interface{}) string {
	response, err := strconv.ParseInt(fmt.Sprintf("%v", responseStatus), 10, 32)
	if err != nil {
		return "badge-info"
	}

	switch response {
	case http.StatusOK:
		return successBadge
	case http.StatusCreated:
		return successBadge
	case http.StatusAccepted:
		return successBadge
	case http.StatusNonAuthoritativeInfo:
		return successBadge
	case http.StatusNoContent:
		return successBadge
	case http.StatusMultipleChoices:
		return infoBadge
	case http.StatusMovedPermanently:
		return infoBadge
	case http.StatusFound:
		return infoBadge
	case http.StatusSeeOther:
		return infoBadge
	case http.StatusNotModified:
		return infoBadge
	case http.StatusUseProxy:
		return infoBadge
	case http.StatusTemporaryRedirect:
		return infoBadge
	case http.StatusPermanentRedirect:
		return infoBadge
	case http.StatusBadRequest:
		return warningBadge
	case http.StatusUnauthorized:
		return warningBadge
	case http.StatusPaymentRequired:
		return warningBadge
	case http.StatusForbidden:
		return warningBadge
	case http.StatusNotFound:
		return warningBadge
	case http.StatusTeapot:
		return warningBadge
	case http.StatusUnprocessableEntity:
		return warningBadge
	case http.StatusInternalServerError:
		return dangerBadge
	case http.StatusNotImplemented:
		return dangerBadge
	case http.StatusBadGateway:
		return dangerBadge
	case http.StatusServiceUnavailable:
		return dangerBadge
	default:
		return infoBadge
	}
}
