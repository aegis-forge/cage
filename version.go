package cage

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
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

// Equals checks if two [Semver] are equal
func (s Semver) Equals(sv Semver) bool {
	return semver.Compare(string(s), string(sv)) == 0
}

// ======================
// ==== VersionRange ====
// ======================

// operators is a slice of all the supported operators
var operators = []string{">", ">=", "<", "<=", "=", "==", ""}

// VersionRange is used to store ranges of versions. If `start` is empty, the
// range will be `< end`. If `end` is empty, then it means `== start`.
type VersionRange struct {
	start        Semver
	end          Semver
	includeRight bool
	includeLeft  bool
}

// NewVersionRange creates a new [VersionRanges] and checks if it's valid
func NewVersionRange(start, end Semver, left, right bool) (*VersionRange, error) {
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

		if semver.Compare(string(start), string(end)) == 1 {
			return nil, errors.New("start is greater to end")
		}
	}
	
	if start.Equals(end) && (!left || !right) {
		return nil, errors.New("if start == end, then both left and right must be true")
	}

	return &VersionRange{
		start: start, end: end, includeLeft: left, includeRight: right,
	}, nil
}

// NewVersionRangeString creates a range given a range string (e.g. `>= v1.0`)
func NewVersionRangeString(stringRange string) (*VersionRange, error) {
	if stringRange == "" {
		return nil, errors.New("stringRange cannot be null")
	}

	regex := regexp.MustCompile(`^([<>]=?|==?)?\s*([0-9A-z.+\-]+)$`)
	matches := regex.FindStringSubmatch(stringRange)

	if len(matches) != 3 {
		return nil, errors.New("the passed string does not contain a range")
	}

	operator := matches[1]
	version, err := NewSemver(matches[2])

	if err != nil {
		return nil, err
	}

	if !slices.Contains(operators, operator) {
		return nil, fmt.Errorf("the operator \"%s\" is not valid", operator)
	}

	var finalRange *VersionRange

	switch operator {
	case "=", "==":
		finalRange, _ = NewVersionRange(version, version, true, true)
	case ">=":
		finalRange, _ = NewVersionRange(version, "", true, false)
	case ">":
		finalRange, _ = NewVersionRange(version, "", false, false)
	case "<=":
		v0, _ := NewSemver("0.0.0")
		finalRange, _ = NewVersionRange(v0, version, false, true)
	case "<":
		v0, _ := NewSemver("0.0.0")
		finalRange, _ = NewVersionRange(v0, version, false, false)
	default:
		finalRange, _ = NewVersionRange(version, version, true, true)
	}

	return finalRange, nil
}

// Equals checks if two [VersionRange] are equal (i.e. same start and end)
func (v *VersionRange) Equals(vr VersionRange) bool {
	return v.start.Equals(vr.start) && v.end.Equals(vr.end) &&
		v.includeLeft == vr.includeLeft && v.includeRight == vr.includeRight
}
