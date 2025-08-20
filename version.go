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

// Before checks if the passed [Semver] is smaller than the original [Semver]
func (s Semver) Before(sv Semver) bool {
	return semver.Compare(string(s), string(sv)) == -1
}

// Before checks if the passed [Semver] is bigger than the original [Semver]
func (s Semver) After(sv Semver) bool {
	return semver.Compare(string(s), string(sv)) == 1
}

// Equals checks if the passed [Semver] is equal to the original [Semver]
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
	includeLeft  bool
	includeRight bool
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

	if start == "" {
		start, _ = NewSemver("v0.0.0")
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
		finalRange, _ = NewVersionRange(v0, version, true, true)
	case "<":
		v0, _ := NewSemver("0.0.0")
		finalRange, _ = NewVersionRange(v0, version, true, false)
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

// Contains checks if a [Semver] is contained in a [VersionRange] struct
func (v *VersionRange) Contains(s Semver) bool {
	afterLeft := false
	beforeRight := false

	// Check if >= than start
	if v.includeLeft && (s.Equals(v.start) || s.After(v.start)) {
		afterLeft = true
	}

	// Check if > than start
	if !v.includeLeft && !s.Equals(v.start) && s.After(v.start) {
		afterLeft = true
	}

	// Open-ended range
	if v.end == "" {
		beforeRight = true
	} else {
		// Check if <= than end
		if v.includeRight && (s.Equals(v.end) || s.Before(v.end)) {
			beforeRight = true
		}

		// Check if < than end
		if !v.includeRight && !s.Equals(v.end) && s.Before(v.end) {
			beforeRight = true
		}
	}

	return afterLeft && beforeRight
}
