package cage

import (
	"errors"
	"slices"
	"time"
)

// Package describes a specific version of a software package
type Package struct {
	vendor    string
	product   string
	published time.Time
	version   Semver
}

// NewPackage creates a new [Package] or returns an error if it is invalid
func NewPackage(vendor, product string, published time.Time, version Semver) (*Package, error) {
	if version == "" || !version.IsValid() {
		return nil, errors.New("semver should be non-empty and valid")
	}

	return &Package{
		vendor:    vendor,
		product:   product,
		published: published,
		version:   version,
	}, nil
}

// IsVulnerable checks if a package is vulnerable by checking the passed
// sources. If `nil` is passed as sources, [Github] will be used. If the package
// is vulnerable, then it will return a slice of [Vulnerability] structs.
func (p *Package) IsVulnerable(sources []Source) ([]Vulnerability, error) {
	if len(sources) == 0 {
		return nil, errors.New("at least one source needs to be added")
	}

	if sources == nil {
		sources = []Source{Github{}}
	}

	var vulnerabilities []Vulnerability

	for _, source := range sources {
		vulns, err := source.GetVulnerabilities(*p)

		if err != nil {
			return nil, err
		}
		
		slices.SortFunc(vulns, func(a, b Vulnerability) int {
			if a.published.After(b.published) {
				return -1
			} else if a.published.Before(b.published) {
				return 1
			} else {
				return 0
			}
		})

		vuln, err := source.CompareVulnerabilities(vulns, *p)

		if err != nil {
			return nil, err
		}

		if vuln != nil {
			vulnerabilities = append(vulnerabilities, vuln...)
		}
	}

	return vulnerabilities, nil
}
