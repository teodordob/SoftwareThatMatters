package graph

import (
	"regexp"
	"strconv"
	"unicode/utf8"
)

type DebianVersion struct {
	Epoch           int
	UpstreamVersion string
	debianRevision  string
}

var versionRegex *regexp.Regexp = regexp.MustCompile(`^((?P<epoch>\d+):)?(?P<upstream_version>[A-Za-z0-9.+:~-]+?)(-(?P<debian_revision>[A-Za-z0-9+.~]+))?$`)
var re_all_digits_or_not *regexp.Regexp = regexp.MustCompile("\\d+|\\D+")
var digitsRegex *regexp.Regexp = regexp.MustCompile("\\d+")
var digitRegex *regexp.Regexp = regexp.MustCompile("\\d")
var alphaRegex *regexp.Regexp = regexp.MustCompile("[A-Za-z]")

func newDebianVersion(epoch int, upstreamVersion string, debianRevision string) *DebianVersion {
	return &DebianVersion{
		Epoch:           epoch,
		UpstreamVersion: upstreamVersion,
		debianRevision:  debianRevision,
	}
}

func ParseDebianVersion(input string) *DebianVersion {
	match := versionRegex.FindStringSubmatch(input)
	result := make(map[string]string)
	var epoch int
	upstreamVersion := "0"
	debianRevision := "0"
	for i, name := range versionRegex.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	if len(result["epoch"]) > 0 {
		epoch, _ = strconv.Atoi(result["epoch"])
	} else {
		epoch = 0
	}
	if len(result["upstream_version"]) > 0 {
		upstreamVersion = result["upstream_version"]
	}
	if len(result["debian_version"]) > 0 {
		debianRevision = result["debian_version"]
	}
	return newDebianVersion(epoch, upstreamVersion, debianRevision)
}

// CompareVersions compares 2 debian versions v1, v2 returns
// 0 for v1 == v2
// -1 for v1 < v2
// 1 for v1 > v2/**
func CompareVersions(v1 DebianVersion, v2 DebianVersion) int {
	result := 0
	if v1.Epoch < v2.Epoch {
		return -1
	} else if v1.Epoch > v2.Epoch {
		return 1
	}
	result = compare2Fields(v1.UpstreamVersion, v2.UpstreamVersion)
	if result != 0 {
		return result
	}
	result = compare2Fields(v1.debianRevision, v2.debianRevision)
	return result
}

func compare2Fields(s1 string, s2 string) int {
	result := 0
	m1 := re_all_digits_or_not.FindAllString(s1, int(^uint(0)>>1))
	m2 := re_all_digits_or_not.FindAllString(s2, int(^uint(0)>>1))
	for ok := true; ok; ok = len(m1) > 0 || len(m2) > 0 {
		a := "0"
		b := "0"
		if len(m1) > 0 {
			a, m1 = m1[0], m1[1:]
		}
		if len(m2) > 0 {
			b, m2 = m2[0], m2[1:]
		}
		if digitsRegex.MatchString(a) && digitsRegex.MatchString(b) {
			a1, _ := strconv.Atoi(a)
			b1, _ := strconv.Atoi(b)
			if a1 < b1 {
				return -1
			} else if a1 > b1 {
				return 1
			}
		} else {
			result = compare2Strings(a, b)
			if result != 0 {
				return result
			}
		}
	}
	return 0
}

func compare2Strings(s1 string, s2 string) int {
	var arr1 []int
	var arr2 []int
	for _, c := range []rune(s1) {
		arr1 = append(arr1, order(string(c)))
	}
	for _, c := range []rune(s2) {
		arr2 = append(arr2, order(string(c)))
	}
	for ok := true; ok; ok = len(arr1) > 0 || len(arr2) > 0 {
		a := 0
		b := 0
		if len(arr1) > 0 {
			a, arr1 = arr1[0], arr1[1:]
		}
		if len(arr2) > 0 {
			b, arr2 = arr2[0], arr2[1:]
		}
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
	}
	return 0
}

func order(x string) int {
	if x == "~" {
		return -1
	} else if digitRegex.MatchString(x) {
		ans, _ := strconv.Atoi(x)
		return ans + 1
	} else if alphaRegex.MatchString(x) {
		ans, _ := utf8.DecodeRuneInString(x)
		return int(ans)
	} else {
		ans, _ := utf8.DecodeRuneInString(x)
		return int(ans) + 256
	}
}

// CheckConstraint /**
func CheckConstraint(constr string, version DebianVersion) bool {

	if constr == "any" {
		return true
	} else {
		var sign, constr1 string
		if constr[0] == '=' {
			sign, constr1 = constr[0:2], constr[2:]
		} else {
			sign, constr1 = constr[0:3], constr[3:]
		}
		constrVer := ParseDebianVersion(constr1)
		switch sign {
		case "<< ":
			return CompareVersions(version, *constrVer) < 0
		case "<= ":
			return CompareVersions(version, *constrVer) < 1
		case "= ":
			return CompareVersions(version, *constrVer) == 0
		case ">= ":
			return CompareVersions(version, *constrVer) >= 0
		case ">> ":
			return CompareVersions(version, *constrVer) > 0
		}
	}
	return false
}
