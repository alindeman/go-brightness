package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func brightnessFile() (devicefile string, err error) {
	var matches []string
	matches, err = filepath.Glob("/sys/class/backlight/*")
	if err != nil {
		return devicefile, err
	}

	if len(matches) >= 1 {
		devicefile = fmt.Sprintf("%s/%s", matches[0], "brightness")
		return devicefile, nil
	}

	return devicefile, errors.New("No backlight devices found")
}

func currentBrightness() (brightness int, err error) {
	var devicefile string
	devicefile, err = brightnessFile()
	if err != nil {
		return 0, err
	}

	var bytes []byte
	bytes, err = ioutil.ReadFile(devicefile)
	if err != nil {
		return 0, err
	}

	var strbrightness string
	strbrightness = strings.TrimSpace(string(bytes))

	brightness, err = strconv.Atoi(strbrightness)
	if err != nil {
		return 0, err
	}

	return brightness, err
}

func adjustBrightness(adjustment int) (brightness int, err error) {
	brightness, err = currentBrightness()
	if err != nil {
		return 0, err
	}

	var devicefile string
	devicefile, err = brightnessFile()
	if err != nil {
		return 0, err
	}

	brightness += adjustment
	err = ioutil.WriteFile(devicefile, []byte(strconv.Itoa(brightness)), 644)
	if err != nil {
		return 0, err
	}

	return brightness, nil
}

func main() {
	args := os.Args
	if len(args) != 2 {
		programname := filepath.Base(args[0])
		fmt.Printf("Usage: %s adjustment\n\n", programname)
		fmt.Printf("Examples:\n")
		fmt.Printf("%s +10: increase brightness by 10\n", programname)
		fmt.Printf("%s -10: decrease brightness by 10\n", programname)
		os.Exit(1)
	}

	if adjustment, err := strconv.Atoi(args[1]); err == nil {
		if brightness, err := adjustBrightness(adjustment); err == nil {
			fmt.Printf("Brightness set to %d\n", brightness)
		} else {
			fmt.Printf("Unable to set brightness: %s\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Invalid argument: %s (%s)\n", args[1], err)
		os.Exit(1)
	}
}
