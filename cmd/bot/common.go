package main

import "github.com/zneix/tcb2/internal/common"

var (
	buildTime    string
	buildVersion string = "v2.0-beta2"
	buildHash    string
	buildBranch  string
)

// Initialized values in the singleton common package
func init() {
	common.BuildTime = buildTime
	common.BuildVersion = buildVersion
	common.BuildHash = buildHash
	common.BuildBranch = buildBranch
}
