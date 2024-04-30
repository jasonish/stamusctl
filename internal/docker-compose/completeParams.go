package compose

import (
	"crypto/rand"
	"encoding/hex"

	"stamus-ctl/internal/logging"
	"stamus-ctl/internal/utils"

	"github.com/spf13/cobra"
)

func Ask(cmd *cobra.Command, params *Parameters) {
	if !cmd.Flags().Changed("interface") {
		getInterfaceCli(&params.InterfacesList)
	}

	if !cmd.Flags().Changed("restart") {
		getRestartCli(&params.RestartMode)
	}

	if !cmd.Flags().Changed("es-datapath") {
		getElasticPathCli(&params.ElasticPath)
	}

	if !cmd.Flags().Changed("registry") {
		getRegistryCli(&params.Registry)
	}

	if !cmd.Flags().Changed("token") {
		b := make([]byte, 24)
		rand.Read(b)
		params.SciriusToken = hex.EncodeToString(b)
		logging.Sugar.Debugw("generated token.", "token", params.SciriusToken)
	}

	params.MLEnabled = utils.GetSSESupport()

}
