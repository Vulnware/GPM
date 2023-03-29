package main

import (
	"log"
	"math"
	"strconv"
	"strings"
)

func parseVersion(text string) []int32 {
	var version []int32
	for _, v := range strings.Split(text, ".") {
		i, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return nil
		}
		version = append(version, int32(i))
	}
	return version
}

func calculateVersion(version []int32) int32 {
	var result int32
	for i, v := range version {
		result += v * int32(math.Pow(10, float64(len(version)-i-1))) // 10^2, 10^1, 10^0 etc.
	}
	return result
}

func checkVersion(version string) {
	versionN := parseVersion(version)
	if versionN == nil {
		log.Fatal("Invalid version number: ", version)
	}
	// check if Major version is the same
	if versionN[0] != tomlVersionNums[0] {
		log.Fatal("The config file version is not compatible with this version of the program. Please update it to the latest version. Current version: ", tomlVersion)
	}

	if calculateVersion(versionN) < calculateVersion(tomlVersionNums) {
		log.Fatal("Config file version is too old. Please update it to the latest version. Current version: ", tomlVersion)

	}

}
