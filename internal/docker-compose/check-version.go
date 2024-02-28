package compose

import (
	"errors"

	"git.stamus-networks.com/lanath/stamus-ctl/internal/logging"
	"github.com/Masterminds/semver/v3"
)

const (
	minimalDockerVersion  = "17.6.0"
	minimalComposeVersion = "1.27.0"
)

var (
	minimalDockerSemVersion  *semver.Version
	minimalComposeSemVersion *semver.Version
)

func CheckVersions() (bool, error) {
	if version, err := GetExecDockerVersion("docker"); err != nil || version.Compare(minimalDockerSemVersion) == -1 {
		if err != nil {
			return false, err
		}
		logging.Sugar.Errorw("Docker version not supported", "got", version, "expected", minimalDockerVersion)
		return false, errors.New("docker version not supported")
	}

	if version, err := GetExecDockerVersion("docker-compose"); err != nil || version.Compare(minimalComposeSemVersion) == -1 {
		if err != nil {
			return false, err
		}
		logging.Sugar.Errorw("docker-compose not supported", "got", version, "expected", minimalDockerVersion)
		return false, errors.New("docker-compose not supported")
	}

	return true, nil
}

func init() {
	minimalDockerSemVersion, _ = semver.NewVersion(minimalDockerVersion)
	minimalComposeSemVersion, _ = semver.NewVersion(minimalComposeVersion)
}
