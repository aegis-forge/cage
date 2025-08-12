package cmd

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/mod/semver"
)

// ================
// ==== Semver ====
// ================

// Semver is a type definition for a version that follows semantic versioning
type Semver string

// NewSemver creates a new [Semver] and checks if it's valid
func NewSemver(version string) (Semver, error) {
	fmtVersion := version

	if strings.HasPrefix(version, "V") {
		fmtVersion = fmt.Sprintf("v%s", version[1:])
	} else if !strings.HasPrefix(version, "v") {
		fmtVersion = fmt.Sprintf("v%s", version)
	}

	if !semver.IsValid(fmtVersion) {
		return "", fmt.Errorf("\"%s\" is not a valid version", fmtVersion)
	}

	return Semver(fmtVersion), nil
}

// IsValid checks if the [Semver] is valid
func (s Semver) IsValid() bool {
	return semver.IsValid(string(s))
}

// ================
// ==== Semver ====
// ================

// VersionRange is used to store ranges of versions. If `start` is empty, the
// range will be `< end`. If `end` is empty, then it means `== start`.
type VersionRange struct {
	start Semver
	end   Semver
}

// NewVersionRange creates a new [VersionRanges] and checks if it's valid
func NewVersionRange(start Semver, end Semver) (*VersionRange, error) {
	if start == "" && end == "" {
		return nil, errors.New("at least one semver needs not to be empty")
	}
	
	if (start == "" && !end.IsValid()) || (!start.IsValid() && end == "") {
		return nil, errors.New("if one semver is empty, the other needs to be valid")
	}
	
	if start != "" && end != "" {
		if !start.IsValid() || !end.IsValid() {
			return nil, errors.New("if not empty, both semvers need to be valid")
		}
		
		if semver.Compare(string(start), string(end)) != -1 {
			return nil, errors.New("start is greater or equal to end")
		}
	}

	return &VersionRange{start: start, end: end}, nil
}
