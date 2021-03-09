package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// Check the wanted path is not in the do not log list.
func CheckExcludedPaths(path string) bool {
	exactMatches := []string{
		"",
		"/apple-touch-icon-precomposed.png",
		"/apple-touch-icon.png",
		"/goscope/css/light.css.map",
		"/goscope/css/dark.css.map",
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
		log.Println(err.Error())
		return rawString
	}

	return prettyJSON.String()
}

func EpochToTimeAgoHappened(epoch int) string {
	date := time.Unix(int64(epoch), 0)
	diff := time.Since(date)
	if diff.Seconds() < 60 {
		return fmt.Sprintf("%.2f s", diff.Seconds())
	} else if diff.Minutes() < 60 {
		return fmt.Sprintf("%.0f m", diff.Minutes())
	}

	return fmt.Sprintf("%.0f h", diff.Hours())
}

func EpochToHumanReadable(epoch int) string {
	date := time.Unix(int64(epoch), 0)
	return date.Format(time.RFC1123Z)
}
