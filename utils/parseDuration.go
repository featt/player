package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

func ParseDurationFromFilename(filename string) (int, error) {
    regex := regexp.MustCompile(`(\d{2})-(\d{2})`)

    matches := regex.FindStringSubmatch(filename)
    if len(matches) != 3 {
        return 0, fmt.Errorf("could not parse duration from file name %s", filename)
    }

    minutes, err := strconv.Atoi(matches[1])
    if err != nil {
        return 0, fmt.Errorf("could not parse duration from file name %s: %s", filename, err)
    }
    seconds, err := strconv.Atoi(matches[2])
    if err != nil {
        return 0, fmt.Errorf("could not parse duration from file name %s: %s", filename, err)
    }

    duration := minutes*60 + seconds

    return duration, nil
}
