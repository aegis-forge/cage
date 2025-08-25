package cage

import (
	"testing"
	"time"
)

// =========================
// ==== Github SetToken ====
// =========================

func TestSetTokenEmpty(t *testing.T) {
	gh := &Github{}

	if err := gh.SetToken(""); err == nil {
		t.Error(err)
	}
}

func TestSetTokenUnauthorized(t *testing.T) {
	gh := &Github{}

	if err := gh.SetToken("a"); err.Error() != "the given token is not valid" {
		t.Error(err)
	}
}

// =======================================
// ==== Github CompareVulnerabilities ====
// =======================================

func TestPackageIsVulnerableTrueA(t *testing.T) {
	s := Github{}

	v1, _ := NewSemver("2.294.0")
	v2, _ := NewSemver("2.296.1")
	v3, _ := NewSemver("0.0.0")
	v4, _ := NewSemver("2.283.4")

	vr1, _ := NewVersionRange(v1, v2, true, false)
	vr2, _ := NewVersionRange(v3, v4, false, false)
	pr1, _ := NewVersionRangeString("2.296.2")
	pr2, _ := NewVersionRangeString("2.293.1")

	vua, _ := NewVulnerability(
		"a", []string{}, 1.0, "2025-07-25T17:51:51Z", []VersionRange{*vr1},
		[]VersionRange{*pr1}, "2006-01-02T15:04:05Z",
	)

	vub, _ := NewVulnerability(
		"b", []string{}, 1.0, "2025-07-25T17:51:51Z", []VersionRange{*vr2},
		[]VersionRange{*pr2}, "2006-01-02T15:04:05Z",
	)

	v, _ := NewSemver("2.296.0")
	p, _ := NewPackage("a", "b", time.Now(), v)
	vs, _ := s.CompareVulnerabilities([]Vulnerability{*vua, *vub}, *p)

	if len(vs) != 1 || vs[0].cve != "a" {
		t.Error("vulnerability \"a\" should have been returned")
	}
}

func TestPackageIsVulnerableTrueB(t *testing.T) {
	s := Github{}

	v1, _ := NewSemver("2.294.0")
	v2, _ := NewSemver("2.296.1")
	v3, _ := NewSemver("0.0.0")
	v4, _ := NewSemver("2.283.4")

	vr1, _ := NewVersionRange(v1, v2, true, false)
	vr2, _ := NewVersionRange(v3, v4, false, false)
	pr1, _ := NewVersionRangeString("2.296.2")
	pr2, _ := NewVersionRangeString("2.293.1")

	vua, _ := NewVulnerability(
		"a", []string{}, 1.0, "2025-07-25T17:51:51Z", []VersionRange{*vr1},
		[]VersionRange{*pr1}, "2006-01-02T15:04:05Z",
	)

	vub, _ := NewVulnerability(
		"b", []string{}, 1.0, "2025-07-25T17:51:51Z", []VersionRange{*vr2},
		[]VersionRange{*pr2}, "2006-01-02T15:04:05Z",
	)

	v, _ := NewSemver("2.283.3")
	p, _ := NewPackage("a", "b", time.Now(), v)
	vs, _ := s.CompareVulnerabilities([]Vulnerability{*vua, *vub}, *p)

	if len(vs) != 1 || vs[0].cve != "b" {
		t.Error("vulnerability \"b\" should have been returned")
	}
}

func TestPackageIsVulnerableTrueAll(t *testing.T) {
	s := Github{}

	v1, _ := NewSemver("2.294.0")
	v2, _ := NewSemver("2.296.1")
	v3, _ := NewSemver("0.0.0")
	v4, _ := NewSemver("2.295.0")

	vr1, _ := NewVersionRange(v1, v2, true, false)
	vr2, _ := NewVersionRange(v3, v4, false, false)
	pr1, _ := NewVersionRangeString("2.296.2")
	pr2, _ := NewVersionRangeString("2.293.1")

	vua, _ := NewVulnerability(
		"a", []string{}, 1.0, "2025-07-25T17:51:51Z", []VersionRange{*vr1},
		[]VersionRange{*pr1}, "2006-01-02T15:04:05Z",
	)

	vub, _ := NewVulnerability(
		"b", []string{}, 1.0, "2025-07-25T17:51:51Z", []VersionRange{*vr2},
		[]VersionRange{*pr2}, "2006-01-02T15:04:05Z",
	)

	v, _ := NewSemver("2.294.2")
	p, _ := NewPackage("a", "b", time.Now(), v)
	vs, _ := s.CompareVulnerabilities([]Vulnerability{*vua, *vub}, *p)

	if len(vs) != 2 || (vs[0].cve != "a" && vs[0].cve != "b") {
		t.Error("vulnerabilities \"a\" and \"b\" should have been returned")
	}
}

func TestPackageIsVulnerableFalse(t *testing.T) {
	s := Github{}

	v1, _ := NewSemver("2.294.0")
	v2, _ := NewSemver("2.296.1")
	v3, _ := NewSemver("0.0.0")
	v4, _ := NewSemver("2.283.4")

	vr1, _ := NewVersionRange(v1, v2, true, false)
	vr2, _ := NewVersionRange(v3, v4, false, false)
	pr1, _ := NewVersionRangeString("2.296.2")
	pr2, _ := NewVersionRangeString("2.293.1")

	vua, _ := NewVulnerability(
		"a", []string{}, 1.0, "2025-07-25T17:51:51Z", []VersionRange{*vr1},
		[]VersionRange{*pr1}, "2006-01-02T15:04:05Z",
	)

	vub, _ := NewVulnerability(
		"b", []string{}, 1.0, "2025-07-25T17:51:51Z", []VersionRange{*vr2},
		[]VersionRange{*pr2}, "2006-01-02T15:04:05Z",
	)

	v, _ := NewSemver("2.296.2")
	p, _ := NewPackage("a", "b", time.Now(), v)
	vs, _ := s.CompareVulnerabilities([]Vulnerability{*vua, *vub}, *p)

	if len(vs) != 0 {
		t.Error("no vulnerability should have been returned")
	}
}

func TestPackageIsVulnerableFalseFuture(t *testing.T) {
	s := Github{}

	v1, _ := NewSemver("0.0.0")
	v2, _ := NewSemver("2.294.0")

	vr1, _ := NewVersionRange(v1, v2, true, false)
	pr1, _ := NewVersionRangeString("2.296.2")

	vua, _ := NewVulnerability(
		"a", []string{}, 1.0, time.Now().Add(time.Hour).Format("2006-01-02T15:04:05Z"),
		[]VersionRange{*vr1}, []VersionRange{*pr1}, "2006-01-02T15:04:05Z",
	)

	v, _ := NewSemver("2")
	p, _ := NewPackage("a", "b", time.Now(), v)
	vs, _ := s.CompareVulnerabilities([]Vulnerability{*vua}, *p)

	if len(vs) != 0 {
		t.Error("no vulnerability should have been returned")
	}
}
