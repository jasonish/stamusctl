package app

import (
	"runtime"

	"git.stamus-networks.com/lanath/stamus-ctl/internal/logging"
)

func CatchException() {
	if err := recover(); err != nil {
		switch err.(type) {
		case *runtime.Error:
			logging.Sugar.Errorf("critical error", err)
			panic(err)
		default:
			logging.Sugar.Errorf("critical error", err)
		}
	}
}
