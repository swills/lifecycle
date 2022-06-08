package cmd

import (
	"fmt"
	"strings"

	"github.com/buildpacks/lifecycle/buildpack"

	"github.com/buildpacks/lifecycle/api"
)

const (
	DeprecationModeQuiet = "quiet"
	DeprecationModeWarn  = "warn"
	DeprecationModeError = "error"
)

var DeprecationMode = EnvOrDefault(EnvDeprecationMode, DefaultDeprecationMode)

type APIVerifier struct{}

func (v *APIVerifier) VerifyBuildpackAPI(kind, name, requested string) error {
	return VerifyBuildpackAPI(kind, name, requested)
}

func (v *APIVerifier) VerifyBuildpackAPIsForGroup(group []buildpack.GroupElement) error {
	for _, groupEl := range group {
		if groupEl.API == "" {
			// if this group was generated by this lifecycle, API should be set
			// but if for some reason it isn't default to 0.2
			groupEl.API = "0.2"
		}
		switch {
		case groupEl.Extension:
			if err := v.VerifyBuildpackAPI(groupEl.Kind(), groupEl.String(), groupEl.API); err != nil {
				return err
			}
		default:
			if err := v.VerifyBuildpackAPI(groupEl.Kind(), groupEl.String(), groupEl.API); err != nil {
				return err
			}
		}
	}
	return nil
}

func VerifyBuildpackAPI(kind, name, requested string) error {
	requestedAPI, err := api.NewVersion(requested)
	if err != nil {
		return FailErrCode(
			fmt.Errorf("parse buildpack API '%s' for %s '%s'", requestedAPI, strings.ToLower(kind), name),
			CodeIncompatibleBuildpackAPI,
		)
	}
	if api.Buildpack.IsSupported(requestedAPI) {
		if api.Buildpack.IsDeprecated(requestedAPI) {
			switch DeprecationMode {
			case DeprecationModeQuiet:
				break
			case DeprecationModeError:
				DefaultLogger.Errorf("%s '%s' requests deprecated API '%s'", kind, name, requested)
				DefaultLogger.Errorf("Deprecated APIs are disabled by %s=%s", EnvDeprecationMode, DeprecationModeError)
				return buildpackAPIError(kind, name, requested)
			case DeprecationModeWarn:
				DefaultLogger.Warnf("%s '%s' requests deprecated API '%s'", kind, name, requested)
			default:
				DefaultLogger.Warnf("%s '%s' requests deprecated API '%s'", kind, name, requested)
			}
		}
		return nil
	}
	return buildpackAPIError(kind, name, requested)
}

func buildpackAPIError(moduleKind string, name string, requested string) error {
	return FailErrCode(
		fmt.Errorf("set API for %s '%s': buildpack API version '%s' is incompatible with the lifecycle", moduleKind, name, requested),
		CodeIncompatibleBuildpackAPI,
	)
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
			case DeprecationModeQuiet:
				break
			case DeprecationModeError:
				DefaultLogger.Errorf("Platform requested deprecated API '%s'", requested)
				DefaultLogger.Errorf("Deprecated APIs are disabled by %s=%s", EnvDeprecationMode, DeprecationModeError)
				return platformAPIError(requested)
			case DeprecationModeWarn:
				DefaultLogger.Warnf("Platform requested deprecated API '%s'", requested)
			default:
				DefaultLogger.Warnf("Platform requested deprecated API '%s'", requested)
			}
		}
		return nil
	}
	return platformAPIError(requested)
}

func platformAPIError(requested string) error {
	return FailErrCode(
		fmt.Errorf("set platform API: platform API version '%s' is incompatible with the lifecycle", requested),
		CodeIncompatiblePlatformAPI,
	)
}
