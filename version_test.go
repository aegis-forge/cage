package cage

import (
	"testing"
)

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

	if _, err := NewVersionRange(start, end, false, false); err != nil {
		t.Error(err)
	}
}

func TestNewVersionRangeInvalid(t *testing.T) {
	start := Semver("master")
	end, _ := NewSemver("v2")

	if _, err := NewVersionRange(start, end, false, false); err == nil {
		t.Error(err)
	}
}

func TestNewVersionRangeValidRange(t *testing.T) {
	start, _ := NewSemver("v1")
	end, _ := NewSemver("v2")

	if _, err := NewVersionRange(start, end, false, false); err != nil {
		t.Error(err)
	}
}

func TestNewVersionRangeInvalidRange(t *testing.T) {
	start, _ := NewSemver("v2")
	end, _ := NewSemver("v1")

	if _, err := NewVersionRange(start, end, false, false); err == nil {
		t.Error(err)
	}
}

func TestNewVersionRangeEmptyStart(t *testing.T) {
	start, _ := NewSemver("v0.0.0")
	end, _ := NewSemver("v1")
	exp, _ := NewVersionRange(start, end, false, false)

	if r, err := NewVersionRange("", end, false, false); err != nil {
		t.Error(err)
	} else {
		if !r.Equals(*exp) {
			t.Errorf("range %s - %s not equal to range %s - %s", r.start, r.end,
				exp.start, exp.end)
		}
	}
}

func TestNewVersionRangeEmptyEnd(t *testing.T) {
	start, _ := NewSemver("v1")

	if _, err := NewVersionRange(start, "", false, false); err != nil {
		t.Error(err)
	}
}

func TestNewVersionRangeEmpty(t *testing.T) {
	if _, err := NewVersionRange("", "", false, false); err == nil {
		t.Error(err)
	}
}

func TestNewVersionRangeStartEqualEndValid(t *testing.T) {
	version, _ := NewSemver("v1")

	if _, err := NewVersionRange(version, version, true, true); err != nil {
		t.Error(err)
	}
}

func TestNewVersionRangeStartEqualEndInvalid(t *testing.T) {
	version, _ := NewSemver("v1")

	if _, err := NewVersionRange(version, version, false, true); err == nil {
		t.Error(err)
	}
}

func TestNewVersionRangeStartEqualEndInvalid2(t *testing.T) {
	version, _ := NewSemver("v1")

	if _, err := NewVersionRange(version, version, true, false); err == nil {
		t.Error(err)
	}
}

func TestNewVersionRangeStartEqualEndInvalid3(t *testing.T) {
	version, _ := NewSemver("v1")

	if _, err := NewVersionRange(version, version, false, false); err == nil {
		t.Error(err)
	}
}

// ===============================
// ==== NewVersionRangeString ====
// ===============================

func TestVersionRangeStringValidLess(t *testing.T) {
	testRange, _ := NewVersionRange("v0.0.0", "v1.0", true, false)

	if vRange, err := NewVersionRangeString("< 1.0"); err == nil {
		if !testRange.Equals(*vRange) {
			t.Error("the two ranges are different")
		}
	} else {
		t.Error(err.Error())
	}
}

func TestVersionRangeStringValidGreater(t *testing.T) {
	testRange, _ := NewVersionRange("v1.0.0", "", false, false)

	if vRange, err := NewVersionRangeString("> 1.0.0"); err == nil {
		if !testRange.Equals(*vRange) {
			t.Error("the two ranges are different")
		}
	} else {
		t.Error(err.Error())
	}
}

func TestVersionRangeStringValidLessEqual(t *testing.T) {
	testRange, _ := NewVersionRange("v0.0.0", "v2", true, true)

	if vRange, err := NewVersionRangeString("<= v2"); err == nil {
		if !testRange.Equals(*vRange) {
			t.Error("the two ranges are different")
		}
	} else {
		t.Error(err.Error())
	}
}

func TestVersionRangeStringValidGreaterEqual(t *testing.T) {
	testRange, _ := NewVersionRange("v1.0.0+ciao-test", "", true, false)

	if vRange, err := NewVersionRangeString(">=1.0.0+ciao-test"); err == nil {
		if !testRange.Equals(*vRange) {
			t.Error("the two ranges are different")
		}
	} else {
		t.Error(err.Error())
	}
}

func TestVersionRangeStringValidEqual(t *testing.T) {
	testRange, _ := NewVersionRange("v11.0.2", "v11.0.2", true, true)

	if vRange, err := NewVersionRangeString("==v11.0.2"); err == nil {
		if !testRange.Equals(*vRange) {
			t.Error("the two ranges are different")
		}
	} else {
		t.Error(err.Error())
	}
}

func TestVersionRangeStringValidEqualAlt(t *testing.T) {
	testRange, _ := NewVersionRange("v11.0.2", "v11.0.2", true, true)

	if vRange, err := NewVersionRangeString("= v11.0.2"); err == nil {
		if !testRange.Equals(*vRange) {
			t.Error("the two ranges are different")
		}
	} else {
		t.Error(err.Error())
	}
}

func TestVersionRangeStringValidEqualAltAlt(t *testing.T) {
	testRange, _ := NewVersionRange("v11", "v11", true, true)

	if vRange, err := NewVersionRangeString("11"); err == nil {
		if !testRange.Equals(*vRange) {
			t.Error("the two ranges are different")
		}
	} else {
		t.Error(err.Error())
	}
}

func TestVersionRangeStringInvalid(t *testing.T) {
	if _, err := NewVersionRangeString(">==1.0"); err == nil {
		t.Error("range should not be valid")
	}

	if _, err := NewVersionRangeString("===v1.0"); err == nil {
		t.Error("range should not be valid")
	}

	if _, err := NewVersionRangeString("<==v1.0"); err == nil {
		t.Error("range should not be valid")
	}
}

func TestVersionRangeStringInvalidVersion(t *testing.T) {
	if _, err := NewVersionRangeString(">=ciao"); err == nil {
		t.Error("range should not be valid")
	}
}

func TestVersionRangeStringInvalidEmpty(t *testing.T) {
	if _, err := NewVersionRangeString(""); err == nil {
		t.Error("range should not be valid")
	}
}

// =============================
// ==== VersionRange Equals ====
// =============================

func TestVersionRangeEqualsTrue(t *testing.T) {
	v1, _ := NewSemver("v1")
	v2, _ := NewSemver("v2")

	r1, _ := NewVersionRange(v1, v2, false, false)
	r2, _ := NewVersionRange(v1, v2, false, false)

	if !r1.Equals(*r2) {
		t.Error("ranges are not the same")
	}
}

func TestVersionRangeEqualsEmptyStart(t *testing.T) {
	v1, _ := NewSemver("v1")

	r1, _ := NewVersionRange("", v1, false, false)
	r2, _ := NewVersionRange("", v1, false, false)

	if !r1.Equals(*r2) {
		t.Error("ranges are not the same")
	}
}

func TestVersionRangeEqualsEmptyEnd(t *testing.T) {
	v1, _ := NewSemver("v1")

	r1, _ := NewVersionRange(v1, "", false, false)
	r2, _ := NewVersionRange(v1, "", false, false)

	if !r1.Equals(*r2) {
		t.Error("ranges are not the same")
	}
}

func TestVersionRangeEqualsDifferentEnd(t *testing.T) {
	v1, _ := NewSemver("v1")
	v2, _ := NewSemver("v2")
	v3, _ := NewSemver("v3")

	r1, _ := NewVersionRange(v1, v2, false, false)
	r2, _ := NewVersionRange(v1, v3, false, false)

	if r1.Equals(*r2) {
		t.Error("ranges are not supposed to be the same")
	}
}

func TestVersionRangeEqualsDifferentStart(t *testing.T) {
	v1, _ := NewSemver("v1")
	v2, _ := NewSemver("v2")
	v3, _ := NewSemver("v3")

	r1, _ := NewVersionRange(v1, v3, false, false)
	r2, _ := NewVersionRange(v2, v3, false, false)

	if r1.Equals(*r2) {
		t.Error("ranges are not supposed to be the same")
	}
}

func TestVersionRangeEqualsDifferent(t *testing.T) {
	v1, _ := NewSemver("v1")
	v2, _ := NewSemver("v2")
	v3, _ := NewSemver("v3")
	v4, _ := NewSemver("v4")

	r1, _ := NewVersionRange(v1, v2, false, false)
	r2, _ := NewVersionRange(v3, v4, false, false)

	if r1.Equals(*r2) {
		t.Error("ranges are not supposed to be the same")
	}
}

func TestVersionRangeEqualsDifferentIncludesLeft(t *testing.T) {
	v1, _ := NewSemver("v1")
	v2, _ := NewSemver("v2")

	r1, _ := NewVersionRange(v1, v2, true, false)
	r2, _ := NewVersionRange(v1, v2, false, false)

	if r1.Equals(*r2) {
		t.Error("ranges are not supposed to be the same")
	}
}

func TestVersionRangeEqualsDifferentIncludesRight(t *testing.T) {
	v1, _ := NewSemver("v1")
	v2, _ := NewSemver("v2")

	r1, _ := NewVersionRange(v1, v2, false, true)
	r2, _ := NewVersionRange(v1, v2, false, false)

	if r1.Equals(*r2) {
		t.Error("ranges are not supposed to be the same")
	}
}

// ===============================
// ==== VersionRange Contains ====
// ===============================

func TestVersionRangeContainsTrueNoIncludes(t *testing.T) {
	v1, _ := NewSemver("v1")
	v2, _ := NewSemver("v3")

	v, _ := NewSemver("v2")
	r, _ := NewVersionRange(v1, v2, false, false)

	if !r.Contains(v) {
		t.Errorf("range %s < v < %s should not contain %s", r.start, r.end, v)
	}
}

func TestVersionRangeContainsTrueLeftInclude(t *testing.T) {
	v1, _ := NewSemver("v1")
	v2, _ := NewSemver("v3")

	v, _ := NewSemver("v1")
	r, _ := NewVersionRange(v1, v2, true, false)

	if !r.Contains(v) {
		t.Errorf("range %s <= v < %s should not contain %s", r.start, r.end, v)
	}
}

func TestVersionRangeContainsTrueRightInclude(t *testing.T) {
	v1, _ := NewSemver("v1")
	v2, _ := NewSemver("v3")

	v, _ := NewSemver("v3")
	r, _ := NewVersionRange(v1, v2, false, true)

	if !r.Contains(v) {
		t.Errorf("range %s < v <= %s should not contain %s", r.start, r.end, v)
	}
}

func TestVersionRangeContainsTrueAllInclude(t *testing.T) {
	v1, _ := NewSemver("v1")
	v2, _ := NewSemver("v3")

	v, _ := NewSemver("v1")
	va, _ := NewSemver("v3")
	r, _ := NewVersionRange(v1, v2, true, true)

	if !r.Contains(v) {
		t.Errorf("range %s <= v <= %s should not contain %s", r.start, r.end, v)
	}

	if !r.Contains(va) {
		t.Errorf("range %s <= v <= %s should not contain %s", r.start, r.end, v)
	}
}

func TestVersionRangeContainsTrueAllIncludeAlt(t *testing.T) {
	v1, _ := NewSemver("v1")

	v, _ := NewSemver("v1")
	r, _ := NewVersionRange(v1, v1, true, true)

	if !r.Contains(v) {
		t.Errorf("range %s <= v <= %s should not contain %s", r.start, r.end, v)
	}
}

func TestVersionRangeContainsFalseNoInclude(t *testing.T) {
	v1, _ := NewSemver("v1")
	v2, _ := NewSemver("v3")

	v, _ := NewSemver("v3")
	vr, _ := NewSemver("v0")
	r, _ := NewVersionRange(v1, v2, false, false)

	if r.Contains(v) {
		t.Errorf("range %s < v < %s should contain %s", r.start, r.end, v)
	}

	if r.Contains(vr) {
		t.Errorf("range %s < v < %s should contain %s", r.start, r.end, v)
	}
}

func TestVersionRangeContainsFalseIncludeLeft(t *testing.T) {
	v1, _ := NewSemver("v1")
	v2, _ := NewSemver("v3")

	v, _ := NewSemver("v3")
	vr, _ := NewSemver("v0")
	r, _ := NewVersionRange(v1, v2, true, false)

	if r.Contains(v) {
		t.Errorf("range %s <= v < %s should contain %s", r.start, r.end, v)
	}

	if r.Contains(vr) {
		t.Errorf("range %s <= v < %s should contain %s", r.start, r.end, v)
	}
}

func TestVersionRangeContainsFalseIncludeRight(t *testing.T) {
	v1, _ := NewSemver("v1")
	v2, _ := NewSemver("v3")

	v, _ := NewSemver("v4")
	vr, _ := NewSemver("v0")
	r, _ := NewVersionRange(v1, v2, false, true)

	if r.Contains(v) {
		t.Errorf("range %s < v <= %s should contain %s", r.start, r.end, v)
	}

	if r.Contains(vr) {
		t.Errorf("range %s < v <= %s should contain %s", r.start, r.end, v)
	}
}

func TestVersionRangeContainsFalseIncludeAll(t *testing.T) {
	v1, _ := NewSemver("v1")
	v2, _ := NewSemver("v3")

	v, _ := NewSemver("v4")
	vr, _ := NewSemver("v0")
	r, _ := NewVersionRange(v1, v2, true, true)

	if r.Contains(v) {
		t.Errorf("range %s <= v <= %s should contain %s", r.start, r.end, v)
	}

	if r.Contains(vr) {
		t.Errorf("range %s <= v <= %s should contain %s", r.start, r.end, v)
	}
}

func TestVersionRangeContainsFalseIncludeAllAlt(t *testing.T) {
	v1, _ := NewSemver("v1")
	v2, _ := NewSemver("v1")

	v, _ := NewSemver("v1.1")
	r, _ := NewVersionRange(v1, v2, true, true)

	if r.Contains(v) {
		t.Errorf("range %s <= v <= %s should contain %s", r.start, r.end, v)
	}
}

func TestVersionRangeContainsTrueOpenEnded(t *testing.T) {
	v1, _ := NewSemver("v1")

	v, _ := NewSemver("v1")
	vr, _ := NewSemver("v100")
	r, _ := NewVersionRange(v1, "", true, false)

	if !r.Contains(v) {
		t.Errorf("range v >= %s should not contain %s", r.start, v)
	}

	if !r.Contains(vr) {
		t.Errorf("range v >= %s should not contain %s", r.start, v)
	}
}

func TestVersionRangeContainsFalseOpenEnded(t *testing.T) {
	v1, _ := NewSemver("v3")

	v, _ := NewSemver("v2")
	r, _ := NewVersionRange(v1, "", true, false)

	if r.Contains(v) {
		t.Errorf("range v >= %s should contain %s", r.start, v)
	}
}
