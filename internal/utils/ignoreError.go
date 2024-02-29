package utils

import "git.stamus-networks.com/lanath/stamus-ctl/internal/logging"

func IgnoreError[T any](a T, e error) T {
	if e != nil {
		logging.Sugar.Errorw("error was ignore but is not nil.", "error", e)
	}
	return a
}
