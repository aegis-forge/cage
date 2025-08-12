package cmd

import "testing"

// ===================
// ==== NewSemver ====
// ===================

func TestNewSemverValidWithSmallV(t *testing.T) {
	if _, err := NewSemver("v1.0.0"); err != nil {
		t.Error(err)
	}
}

func TestNewSemverValidWithCapitalV(t *testing.T) {
	if _, err := NewSemver("V1.0.0"); err != nil {
		t.Error(err)
	}
}

func TestNewSemverValidWithoutV(t *testing.T) {
	if _, err := NewSemver("1"); err != nil {
		t.Error(err)
	}
}

func TestNewSemverValidWithoutMinor(t *testing.T) {
	if _, err := NewSemver("v1"); err != nil {
		t.Error(err)
	}
}

func TestNewSemverValidWithoutPatch(t *testing.T) {
	if _, err := NewSemver("v1.0"); err != nil {
		t.Error(err)
	}
}

func TestNewSemverInvalid(t *testing.T) {
	if _, err := NewSemver("master"); err == nil {
		t.Error(err)
	}
}

// =========================
// ==== NewVersionRange ====
// =========================

func TestNewVersionRangeValid(t *testing.T) {
	start, _ := NewSemver("v1")
	end, _ := NewSemver("v2")
	
	if _, err := NewVersionRange(start, end); err != nil {
		t.Error(err)
	}
}

func TestNewVersionRangeInvalid(t *testing.T) {
	start := Semver("master")
	end, _ := NewSemver("v2")
	
	if _, err := NewVersionRange(start, end); err == nil {
		t.Error(err)
	}
}

func TestNewVersionRangeValidRange(t *testing.T) {
	start, _ := NewSemver("v1")
	end, _ := NewSemver("v2")
	
	if _, err := NewVersionRange(start, end); err != nil {
		t.Error(err)
	}
}

func TestNewVersionRangeInvalidRange(t *testing.T) {
	start, _ := NewSemver("v2")
	end, _ := NewSemver("v1")
	
	if _, err := NewVersionRange(start, end); err == nil {
		t.Error(err)
	}
}

func TestNewVersionRangeEmptyStart(t *testing.T) {
	end, _ := NewSemver("v1")
	
	if _, err := NewVersionRange("", end); err != nil {
		t.Error(err)
	}
}

func TestNewVersionRangeEmptyEnd(t *testing.T) {
	start, _ := NewSemver("v1")
	
	if _, err := NewVersionRange(start, ""); err != nil {
		t.Error(err)
	}
}

func TestNewVersionRangeEmpty(t *testing.T) {
	if _, err := NewVersionRange("", ""); err == nil {
		t.Error(err)
	}
}
