package compose

import ( // Common
	// External
	// Custom

	compose "stamus-ctl/internal/docker-compose"

	"github.com/spf13/cobra"
)

// Commands
func wrappedCmd() ([]*cobra.Command, map[string]*cobra.Command) {
	return compose.WrappedCmd(compose.ComposeFlags)
}
