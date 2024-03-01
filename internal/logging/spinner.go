package logging

import (
	"time"

	"git.stamus-networks.com/lanath/stamus-ctl/internal/app"
	"github.com/briandowns/spinner"
)

var ()

func NewSpinner(prefix, finalMSG string) *spinner.Spinner {
	if app.Name != app.CtlName {
		return nil
	}
	s := spinner.New(spinner.CharSets[7], 100*time.Millisecond)
	s.Prefix = prefix
	s.FinalMSG = finalMSG

	s.Start()

	return s
}

func SpinnerStop(s *spinner.Spinner) {
	if s != nil {
		s.Stop()
	}
}
