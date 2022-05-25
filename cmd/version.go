package cmd

import (
	"fmt"
	"strings"

	"github.com/buildpacks/lifecycle/api"
)

// The following variables are injected at compile time.
var (
	// Version is the version of the lifecycle and all produced binaries.
	Version = "0.0.0"
	// SCMCommit is the commit information provided by SCM.
	SCMCommit = ""
	// SCMRepository is the source repository.
	SCMRepository = ""
)

var (
	DeprecationMode  = EnvOrDefault(EnvDeprecationMode, DefaultDeprecationMode)
	ExperimentalMode = EnvOrDefault(EnvExperimentalMode, DefaultExperimentalMode)
)

const (
	ModeQuiet = "quiet"
	ModeWarn  = "warn"
	ModeError = "error"
)

// buildVersion is a display format of the version and build metadata in compliance with semver.
func buildVersion() string {
	// noinspection GoBoolExpressions
	if SCMCommit == "" || strings.Contains(Version, SCMCommit) {
		return Version
	}

	return fmt.Sprintf("%s+%s", Version, SCMCommit)
}

func VerifyPlatformAPI(requested string) error {
	requestedAPI, err := api.NewVersion(requested)
	if err != nil {
		return FailErrCode(
			fmt.Errorf("parse platform API '%s'", requested),
			CodeIncompatiblePlatformAPI,
		)
	}
	if api.Platform.IsSupported(requestedAPI) {
		if api.Platform.IsDeprecated(requestedAPI) {
			switch DeprecationMode {
			case ModeQuiet:
				break
			case ModeError:
				DefaultLogger.Errorf("Platform requested deprecated API '%s'", requested)
				DefaultLogger.Errorf("Deprecated APIs are disabled by %s=%s", EnvDeprecationMode, ModeError)
				return platformAPIError(requested)
			case ModeWarn:
				DefaultLogger.Warnf("Platform requested deprecated API '%s'", requested)
			default:
				DefaultLogger.Warnf("Platform requested deprecated API '%s'", requested)
			}
		}
		if api.Platform.IsExperimental(requestedAPI) {
			switch ExperimentalMode {
			case ModeQuiet:
				break
			case ModeError:
				DefaultLogger.Errorf("Platform requested experimental API '%s'", requested)
				DefaultLogger.Errorf("Experimental APIs are disabled by %s=%s", EnvExperimentalMode, ModeError)
				return platformAPIError(requested)
			case ModeWarn:
				DefaultLogger.Warnf("Platform requested experimental API '%s'", requested)
			default:
				DefaultLogger.Warnf("Platform requested experimental API '%s'", requested)
			}
		}
		return nil
	}
	return platformAPIError(requested)
}

func VerifyBuildpackAPI(bp string, requested string) error {
	requestedAPI, err := api.NewVersion(requested)
	if err != nil {
		return FailErrCode(
			fmt.Errorf("parse buildpack API '%s' for buildpack '%s'", requestedAPI, bp),
			CodeIncompatibleBuildpackAPI,
		)
	}
	if api.Buildpack.IsSupported(requestedAPI) {
		if api.Buildpack.IsDeprecated(requestedAPI) {
			switch DeprecationMode {
			case ModeQuiet:
				break
			case ModeError:
				DefaultLogger.Errorf("Buildpack '%s' requests deprecated API '%s'", bp, requested)
				DefaultLogger.Errorf("Deprecated APIs are disabled by %s=%s", EnvDeprecationMode, ModeError)
				return buildpackAPIError(bp, requested)
			case ModeWarn:
				DefaultLogger.Warnf("Buildpack '%s' requests deprecated API '%s'", bp, requested)
			default:
				DefaultLogger.Warnf("Buildpack '%s' requests deprecated API '%s'", bp, requested)
			}
		}
		if api.Buildpack.IsExperimental(requestedAPI) {
			switch ExperimentalMode {
			case ModeQuiet:
				break
			case ModeError:
				DefaultLogger.Errorf("Buildpack '%s' requests experimental API '%s'", bp, requested)
				DefaultLogger.Errorf("Experimental APIs are disabled by %s=%s", EnvExperimentalMode, ModeError)
				return buildpackAPIError(bp, requested)
			case ModeWarn:
				DefaultLogger.Warnf("Buildpack '%s' requests experimental API '%s'", bp, requested)
			default:
				DefaultLogger.Warnf("Buildpack '%s' requests experimental API '%s'", bp, requested)
			}
		}
		return nil
	}
	return buildpackAPIError(bp, requested)
}

func platformAPIError(requested string) error {
	return FailErrCode(
		fmt.Errorf("set platform API: platform API version '%s' is incompatible with the lifecycle", requested),
		CodeIncompatiblePlatformAPI,
	)
}

func buildpackAPIError(bp string, requested string) error {
	return FailErrCode(
		fmt.Errorf("set API for buildpack '%s': buildpack API version '%s' is incompatible with the lifecycle", bp, requested),
		CodeIncompatibleBuildpackAPI,
	)
}
