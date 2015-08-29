package deliver

import (
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	order              = []string{"gravity", "resize", "crop"}
	defaultGravity     = "North"
	gravityTypeMapping = map[string]string{
		"":          defaultGravity,
		"center":    "Center",
		"east":      "East",
		"northeast": "NorthEast",
		"north":     "North",
		"northwest": "NorthWest",
		"southeast": "SouthEast",
		"south":     "South",
		"southwest": "SouthWest",
		"west":      "West",
	}
)

func Convert(filename, originalFile string) (convertedFile string, err error) {
	var (
		outBuff, errBuff bytes.Buffer
		transformOptions = make(map[string]string)
	)

	if isBadRatio(originalFile) {
		return "", errors.New("Convert: bad aspect ratio: " + filename)
	} else {
		resizeRegexp := regexp.MustCompile("s\\/[^\\/]+\\/")
		gravityRegexp := regexp.MustCompile("gr\\/[^\\/]+\\/")

		resizeOption := resizeRegexp.FindAllString(filename, 1)
		gravityOption := gravityRegexp.FindAllString(filename, 1)

		if (len(resizeOption) == 0) && (len(gravityOption) == 0) {
			return "", errors.New("Convert: bad convert options: " + filename)
		}

		if len(gravityOption) > 0 {
			gravityOptionString := gravityOption[0]
			gravityOptionString = strings.TrimLeft(gravityOptionString, "gr/")
			gravityOptionString = strings.TrimRight(gravityOptionString, "/")
			transformOptions["gravity"] = getGravityType(gravityOptionString)
		}

		if len(resizeOption) > 0 {
			resizeOptionString := resizeOption[0]
			resizeOptionString = strings.TrimLeft(resizeOptionString, "s/")
			resizeOptionString = strings.TrimRight(resizeOptionString, "/")

			if sizeCount, ok := stripSize(resizeOptionString); ok {
				switch sizeCount {
				case 2:
					transformOptions["resize"] = resizeOptionString + "^"
					transformOptions["crop"] = resizeOptionString + "+0+0"
				case 1:
					transformOptions["resize"] = resizeOptionString + ">"
				}
			}
		}

		options := []string{originalFile}
		for _, key := range order {
			if value, ok := transformOptions[key]; ok {
				options = append(options, "-"+key, value)
			}
		}

		convertPath := "cache" + string(filepath.Separator) + filename
		convertDir := filepath.Dir(convertPath)

		os.MkdirAll(convertDir, 0777)

		options = append(options, "-strip", "+repage", convertPath)

		convert := exec.Command("convert", options...)
		convert.Stdout = &outBuff
		convert.Stderr = &errBuff

		if convert.Run() != nil {
			return "", errors.New("Convert: " + errBuff.String() + ": " + filename)
		} else {
			return convertPath, nil
		}
	}
}

func isBadRatio(originalFile string) bool {
	var (
		out, err bytes.Buffer
		badRatio bool = true
	)

	cmd := exec.Command("convert", originalFile, "-format", `"%[fx:w/h]"`, "info:")
	cmd.Stdout = &out
	cmd.Stderr = &err

	if cmd.Run() != nil {
		log.Printf("Error checking aspect ratio: %s", err.String())
		return false
	}

	ratioStrings := strings.Split(out.String(), `"`)

	out.Reset()
	err.Reset()

	for _, ratioString := range ratioStrings {
		ratioString = strings.Replace(ratioString, "\n", "", -1)
		if ratioString == "" {
			continue
		}

		ratio, parseError := strconv.ParseFloat(ratioString, 64)
		if parseError != nil || 0.1 > ratio || ratio > 10 {
			badRatio = true
		} else {
			badRatio = false
		}
	}
	return badRatio
}

func stripSize(size string) (length int, result bool) {
	result = false
	length = 0

	if size == "" {
		return
	}

	sizes := strings.Split(size, "x")
	length = len(sizes)

	if length != 2 {
		return
	}

	if sizes[1] == "" {
		length = 1
		if _, err := strconv.Atoi(sizes[0]); err != nil {
			return
		}
	} else {
		for _, size := range sizes {
			if _, err := strconv.Atoi(size); err != nil {
				return
			}
		}
	}

	result = true
	return
}

func getGravityType(key string) string {
	if gravity, ok := gravityTypeMapping[key]; ok {
		return gravity
	} else {
		return defaultGravity
	}
}
