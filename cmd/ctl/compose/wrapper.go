package compose

import ( // Common
	// External
	// Custom
	handlers "stamus-ctl/internal/handlers/compose"
	"stamus-ctl/internal/models"

	"github.com/spf13/cobra"
)

// Variables
var composeFlags = models.ComposeFlags{
	"up": models.CreateComposeFlags(
		[]string{"file"},
		[]string{"detach"},
	),
	"down": models.CreateComposeFlags(
		[]string{"file"},
		[]string{"volumes", "remove-orphans"},
	),
	"ps": models.CreateComposeFlags(
		[]string{"file"},
		[]string{"services", "quiet", "format"},
	),
}

// Commands
func wrappedCmd() ([]*cobra.Command, map[string]*cobra.Command) {
	return handlers.WrappedCmd(composeFlags)
}
