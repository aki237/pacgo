package pacgo

import (
	"regexp"
	"strconv"
	"strings"
)

// ParseError is type that implements error interface, that is used to return the parse error
type ParseError string

// Error method implements error interface for the ParseError type
func (p ParseError) Error() string {
	return string(p)
}

// Error constants
const (
	ErrInvalidPackageName        ParseError = "Invalid package name"
	ErrInvalidArray              ParseError = "Invalid array syntax"
	ErrInvalidPackageEpochNumber ParseError = "Invalid package epoch number"
	ErrInvalidKey                ParseError = "Invalid key passed"
)

// PackageInfo struct contains the metadata fields of the PKGBUILD
type PackageInfo struct {
	Pkgnames     []string
	Pkgver       string // required
	Pkgrel       string // required
	Pkgdir       string
	Epoch        int
	Pkgbase      string
	Pkgdesc      string
	Arch         []string // required
	URL          string
	License      []string // recommended
	Groups       []string
	Depends      []string
	Optdepends   []string
	Makedepends  []string
	Checkdepends []string
	Provides     []string
	Conflicts    []string
	Replaces     []string
	Backup       []string
	Options      []string
	Install      string
	Changelog    string
	Source       []string
	Noextract    []string
	Md5sums      []string
	Sha1sums     []string
	Sha224sums   []string
	Sha256sums   []string
	Sha384sums   []string
	Sha512sums   []string
	Validpgpkeys []string
}

// NewPackageInfo returns a pointer to a newly created PakcageInfo
func NewPackageInfo() *PackageInfo {
	return &PackageInfo{
		Pkgnames:     make([]string, 0),
		Arch:         make([]string, 0),
		License:      make([]string, 0),
		Groups:       make([]string, 0),
		Depends:      make([]string, 0),
		Optdepends:   make([]string, 0),
		Makedepends:  make([]string, 0),
		Checkdepends: make([]string, 0),
		Provides:     make([]string, 0),
		Conflicts:    make([]string, 0),
		Replaces:     make([]string, 0),
		Backup:       make([]string, 0),
		Options:      make([]string, 0),
		Source:       make([]string, 0),
		Noextract:    make([]string, 0),
		Md5sums:      make([]string, 0),
		Sha1sums:     make([]string, 0),
		Sha224sums:   make([]string, 0),
		Sha256sums:   make([]string, 0),
		Sha384sums:   make([]string, 0),
		Sha512sums:   make([]string, 0),
		Validpgpkeys: make([]string, 0),
	}
}

// parseParenArray method is used to parse the array expression in the PKGBUILD files like
//     ('x86_64' 'i686' 'arm64')
// This function wil return an array of strings like
//     [x86_64 i686 arm64]
func (p *PackageInfo) parseParenArray(x string) ([]string, error) {
	if !strings.HasPrefix(x, "(") || !strings.HasSuffix(x, ")") {
		return nil, ErrInvalidArray
	}

	x = strings.TrimPrefix(x, "(")
	x = strings.TrimSuffix(x, ")")

	strarr := make([]string, 0)

	for _, val := range strings.Split(x, " ") {
		strarr = append(strarr, val)
	}

	return strarr, nil
}

// setArray method is used to append a string value to a given string array pointer.
func (p *PackageInfo) setArray(x string, final *[]string) error {
	if !strings.HasPrefix(x, "(") {
		*final = append(*final, x)
		return nil
	}

	pkgs, err := p.parseParenArray(x)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		if err := p.setArray(pkg, final); err != nil {
			return err
		}
	}

	return nil
}

// Set method is used to set a value for a given key.
// If the key given is not a valid PKGBUILD metadata key, this will throw an error.
// For the list of valid keys, see : https://wiki.archlinux.org/index.php/PKGBUILD
func (p *PackageInfo) Set(key, value string) error {
	rxp := regexp.MustCompile("^[a-z]+[a-z0-9@._+-]*$")

	switch key {
	case "pkgbase":
		p.Pkgbase = value
	case "pkgname":
		if !strings.HasPrefix(value, "(") {

			value = strings.TrimPrefix(value, "'")
			value = strings.TrimPrefix(value, "\"")
			value = strings.TrimSuffix(value, "'")
			value = strings.TrimSuffix(value, "\"")
			if !rxp.MatchString(value) {
				return ErrInvalidPackageName
			}
		}
		return p.setArray(value, &p.Pkgnames)
	case "pkgver":
		p.Pkgver = value
	case "pkgrel":
		p.Pkgrel = value
	case "epoch":
		i, err := strconv.Atoi(value)
		if err != nil {
			return ErrInvalidPackageEpochNumber
		}
		p.Epoch = i
	case "pkgdesc":
		p.Pkgdesc = value
	case "arch":
		return p.setArray(value, &p.Arch)
	case "url":
		p.URL = value
	case "license":
		return p.setArray(value, &p.License)
	case "groups":
		return p.setArray(value, &p.Groups)
	case "depends":
		return p.setArray(value, &p.Depends)
	case "optdepends":
		return p.setArray(value, &p.Optdepends)
	case "makedepends":
		return p.setArray(value, &p.Makedepends)
	case "checkdepends":
		return p.setArray(value, &p.Checkdepends)
	case "provides":
		return p.setArray(value, &p.Provides)
	case "conflicts":
		return p.setArray(value, &p.Conflicts)
	case "replaces":
		return p.setArray(value, &p.Replaces)
	case "backup":
		return p.setArray(value, &p.Backup)
	case "options":
		return p.setArray(value, &p.Options)
	case "install":
		p.Install = value
	case "changelog":
		p.Changelog = value
	case "source":
		return p.setArray(value, &p.Source)
	case "noextract":
		return p.setArray(value, &p.Noextract)
	case "validgpgkeys":
		return p.setArray(value, &p.Validpgpkeys)
	case "md5sums":
		p.setArray(value, &p.Md5sums)
	case "sha1sums":
		p.setArray(value, &p.Sha1sums)
	case "sha224sums":
		p.setArray(value, &p.Sha224sums)
	case "sha256sums":
		p.setArray(value, &p.Sha256sums)
	case "sha384sums":
		p.setArray(value, &p.Sha384sums)
	case "sha512sums":
		p.setArray(value, &p.Sha512sums)
	default:
		return ErrInvalidKey
	}
	return nil
}
