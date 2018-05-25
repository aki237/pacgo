// Package pacgo implements parser to parse Arch Linux PKGBUILD scripts.
package pacgo

// PkgBuild struct contains the metadata of a PKGBUILD entry including the
// functions defined in it.
type PkgBuild struct {
	*PackageInfo
	Funcs []string
}
