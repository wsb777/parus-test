package domain

import "fmt"

type SemVer struct{ Major, Minor, Patch int }

func ParseSemVer(s string) (SemVer, error) {
	var v SemVer
	_, err := fmt.Sscanf(s, "%d.%d.%d", &v.Major, &v.Minor, &v.Patch)
	return v, err
}

func (v SemVer) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v SemVer) BumpPatch() SemVer {
	return SemVer{v.Major, v.Minor, v.Patch + 1}
}

func (v SemVer) BumpMinor() SemVer {
	return SemVer{v.Major, v.Minor + 1, 0}
}

func (v SemVer) BumpMajor() SemVer {
	return SemVer{v.Major + 1, 0, 0}
}

func (v SemVer) After(other SemVer) bool {
	if v.Major != other.Major {
		return v.Major > other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor > other.Minor
	}
	return v.Patch > other.Patch
}
