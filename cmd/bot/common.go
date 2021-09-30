package main

import "github.com/zneix/tcb2/internal/common"

var (
	buildTime    string
	buildVersion string = "dev"
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
