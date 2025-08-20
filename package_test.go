package cage

import (
	"testing"
	"time"
)

// ====================
// ==== NewPackage ====
// ====================

func TestNewPackageValid(t *testing.T) {
	v, _ := NewSemver("1.0.0")

	if _, err := NewPackage("a", "b", time.Now(), v); err != nil {
		t.Error("package should be valid")
	}
}

func TestNewPackageNoVersion(t *testing.T) {
	if _, err := NewPackage("a", "b", time.Now(), ""); err.Error() != "semver should be non-empty and valid" {
		t.Error("package shouldn't be valid")
	}
}

func TestNewPackageInvalidVersion(t *testing.T) {
	if _, err := NewPackage("a", "b", time.Now(), "1"); err.Error() != "semver should be non-empty and valid" {
		t.Error("package shouldn't be valid")
	}
}

func TestPackageIsVulnerableNoSource(t *testing.T) {
	v, _ := NewSemver("1")
	p, _ := NewPackage("a", "b", time.Now(), v)

	if _, err := p.IsVulnerable([]Source{}); err.Error() != "at least one source needs to be added" {
		t.Error("method should not work without sources")
	}
}
