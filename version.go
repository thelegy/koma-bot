package main

import (
	"strconv"
	"strings"
	"time"
)

var (
	_versionAutomaticBuild string
	_versionDate           string
	_versionGitBranch      string
	_versionGitHash        string
	_versionTravisNumber   string
	_versionTravisId       string
)

type Version struct {
	AutomaticBuild bool
	Date           time.Time
	GitHash        string
	GitBranch      string
	TravisNumber   int
	TravisId       string
}

func getVersionNumber() *Version {
	trueNames := []string{"t", "true", "y", "yes"}

	AutomaticBuild := false
	GitHash := strings.ToLower(_versionGitHash)
	GitBranch := "none"
	TravisNumber := 0
	TravisId := strings.ToLower(_versionTravisId)

	versionAutomaticBuild := strings.ToLower(_versionAutomaticBuild)
	for _, name := range trueNames {
		if versionAutomaticBuild == name {
			AutomaticBuild = true
		}
	}

	if _versionGitBranch != "" {
		GitBranch = _versionGitBranch
	}

	travisNumber, err := strconv.Atoi(_versionTravisNumber)
	if err == nil {
		TravisNumber = travisNumber
	}

	version := &Version{
		AutomaticBuild: AutomaticBuild,
		GitHash:        GitHash,
		GitBranch:      GitBranch,
		TravisNumber:   TravisNumber,
		TravisId:       TravisId,
	}

	date, err := time.Parse("2006-01-02T15:04:05-07:00", _versionDate)
	if err == nil {
		version.Date = date
	}

	return version
}
