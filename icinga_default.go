//go:build darwin
// +build darwin

package main

const (
	// IcingaStatePrefix with a nonexisting path to the Icinga installation, as the standard path is unknown.
	IcingaStatePrefix = "/nonexisting"
)
