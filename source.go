package cage

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

// ================
// ==== Source ====
// ================

// Source is an interface that describes the source of vulnerabilities
type Source interface {
	GetVulnerabilities(Package) ([]Vulnerability, error)
	CompareVulnerabilities([]Vulnerability, Package) ([]Vulnerability, error)
}

// ================
// ==== Github ====
// ================

// Github is a source of vulnerabilities (https://github.com/advisories)
type Github struct {
	token string
}

// SetToken sets the GitHub token and checks if its valid
func (g *Github) SetToken(token string) error {
	if token == "" {
		return errors.New("token must not be an empty string")
	}

	url := "https://api.github.com"
	bearer := "Bearer " + token

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return errors.New("an error happened when trying to verify the token")
	}

	switch res.StatusCode {
	case http.StatusOK:
		g.token = token
		return nil
	case http.StatusUnauthorized:
		return errors.New("the given token is not valid")
	default:
		return fmt.Errorf(
			"when verifying the token, a %d error code was returned",
			res.StatusCode,
		)
	}
}

// GetVulnerabilities retrieves all vulnerabilities for the given package version
func (g Github) GetVulnerabilities(packg Package) ([]Vulnerability, error) {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/security-advisories",
		packg.vendor, packg.product,
	)

	req, _ := http.NewRequest(http.MethodGet, url, nil)

	if g.token != "" {
		req.Header.Add("Authorization", "Bearer "+g.token)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	type identifier struct {
		Ghsa string `json:"ghsa_id"`
	}

	var identifiers []identifier

	if err := json.NewDecoder(res.Body).Decode(&identifiers); err != nil {
		return nil, err
	}

	// ======

	type rawResponse struct {
		Id              string
		Cve             string `json:"cve_id"`
		Severity        string `json:"severity"`
		Published       string `json:"published_at"`
		Withdrawn       string `json:"withdrawn_at"`
		Vulnerabilities []struct {
			Package struct {
				Ecosystem string `json:"ecosystem"`
				Name      string `json:"name"`
			} `json:"package"`
			VulnRange    string `json:"vulnerable_version_range"`
			PatchedRange string `json:"first_patched_version"`
		} `json:"vulnerabilities"`
		Cvss struct {
			Score float32 `json:"score"`
		} `json:"cvss"`
		Cwes []struct {
			Id string `json:"cwe_id"`
		} `json:"cwes"`
	}

	var rawVulns []rawResponse
	var vulns []Vulnerability

	for _, id := range identifiers {
		url = fmt.Sprintf("https://api.github.com/advisories/%s", id.Ghsa)
		req, _ = http.NewRequest(http.MethodGet, url, nil)

		if g.token != "" {
			req.Header.Add("Authorization", "Bearer "+g.token)
		}

		res, err = client.Do(req)

		if err != nil {
			return nil, err
		}

		if res.StatusCode == http.StatusNotFound {
			continue
		}

		defer res.Body.Close()

		var rawVuln rawResponse

		if err := json.NewDecoder(res.Body).Decode(&rawVuln); err != nil {
			return nil, err
		}
		
		rawVuln.Id = id.Ghsa

		rawVulns = append(rawVulns, rawVuln)
	}

	for _, rawVuln := range rawVulns {
		var vulnerableRanges []VersionRange
		var patchedRanges []VersionRange

		for _, vuln := range rawVuln.Vulnerabilities {
			packgRaw := vuln.Package
			name := fmt.Sprintf("%s/%s", packg.vendor, packg.product)

			valid := rawVuln.Withdrawn != ""
			valid = packgRaw.Ecosystem == "actions" && packgRaw.Name == name

			if !valid {
				continue
			}

			versions := strings.Split(vuln.VulnRange, ", ")

			if len(versions) == 2 {
				r1, err := NewVersionRangeString(versions[0])
				r2, err := NewVersionRangeString(versions[1])

				if err != nil {
					return nil, err
				}

				r, _ := NewVersionRange(r1.Start, r2.End, r1.IncludeLeft, r2.IncludeRight)
				vulnerableRanges = append(vulnerableRanges, *r)
			} else {
				for _, vulnString := range versions {
					r, err := NewVersionRangeString(vulnString)

					if err != nil {
						return nil, err
					}

					vulnerableRanges = append(vulnerableRanges, *r)
				}
			}

			if vuln.PatchedRange != "" {
				v, err := NewSemver(vuln.PatchedRange)

				if err != nil {
					return nil, err
				}

				r, err := NewVersionRange(v, "", true, false)

				if err != nil {
					return nil, err
				}

				patchedRanges = append(patchedRanges, *r)
			} else {
				patchedRanges = append(patchedRanges, VersionRange{})
			}
		}

		var cwes []string

		for _, cwe := range rawVuln.Cwes {
			cwes = append(cwes, cwe.Id)
		}

		vuln, err := NewVulnerability(
			rawVuln.Id, rawVuln.Cve, cwes, rawVuln.Cvss.Score,
			rawVuln.Published, vulnerableRanges, patchedRanges,
			"2006-01-02T15:04:05Z",
		)

		if err != nil {
			return nil, err
		}

		vulns = append(vulns, *vuln)
	}

	return vulns, err
}

// CompareVulnerabilities checks if the [Package] version is contained in the
// vulns slice.
func (g Github) CompareVulnerabilities(vulns []Vulnerability, packg Package) ([]Vulnerability, error) {
	var vulnerabilitiesClean []Vulnerability

	for _, vuln := range vulns {
		var vs []bool

		for i, vulnRange := range vuln.RangesVulnerable {
			v := false
			vp := vuln.RangesPatched[i]

			if vulnRange.Contains(packg.version) {
				v = true
			}

			if vp.Start == "" && vp.End == "" {
				continue
			} else if vp.Contains(packg.version) {
				v = false
			}

			vs = append(vs, v)
		}

		comp := slices.Compact(vs)

		if len(comp) > 1 || (len(comp) == 1 && comp[0]) {
			vulnerabilitiesClean = append(vulnerabilitiesClean, vuln)
		}
	}

	return vulnerabilitiesClean, nil
}
