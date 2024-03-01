package docker

import (
	"errors"
	"slices"
	"strings"

	"git.stamus-networks.com/lanath/stamus-ctl/internal/logging"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
)

func ImageName(image image.Summary) string {
	return strings.Split(image.RepoDigests[0], "@")[0]
}

func GetImagesName(images []image.Summary) []string {
	var names []string
	for _, image := range images {
		names = append(names, ImageName(image))
	}

	return names
}

func GetInstalledImagesName() ([]string, error) {
	images, _ := cli.ImageList(ctx, types.ImageListOptions{All: true})

	names := GetImagesName(images)
	for _, image := range images {
		names = append(names, ImageName(image))
	}

	return names, nil
}

func IsImageAlreadyInstalled(name string) (bool, error) {
	images, err := GetInstalledImagesName()

	return slices.Contains(images, name), err
}

func GetImageIdFromName(name string) (string, error) {
	images, _ := cli.ImageList(ctx, types.ImageListOptions{All: true})
	for _, image := range images {
		shortName := ImageName(image)

		if shortName == name {
			logging.Sugar.Debugw("image name found", "image.ID", image.ID, "shortName", shortName, "name", name)
			return image.ID, nil
		}
	}

	logging.Sugar.Debugw("image not found", "images", images, "name", name)
	return "", errors.New("image not found")
}
