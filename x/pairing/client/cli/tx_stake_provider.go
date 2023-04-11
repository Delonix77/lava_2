package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstoragetypes "github.com/lavanet/lava/x/epochstorage/types"
	"github.com/lavanet/lava/x/pairing/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdStakeProvider() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stake-provider [chain-id] [amount] [endpoint endpoint ...] [geolocation] --from <address> --provider-moniker <moniker>",
		Short: "Broadcast message stakeProvider",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argChainID := args[0]
			argAmount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}
			tmpArg := strings.Fields(args[2])
			argEndpoints := []epochstoragetypes.Endpoint{}
			for _, endpointStr := range tmpArg {
				splitted := strings.Split(endpointStr, ",")
				if len(splitted) != 3 {
					return fmt.Errorf("invalid argument format in endpoints, must be: HOST:PORT,useType,geolocation HOST:PORT,useType,geolocation, received: %s", endpointStr)
				}
				geoloc, err := strconv.ParseUint(splitted[2], 10, 64)
				if err != nil {
					return fmt.Errorf("invalid argument format in endpoints, geolocation must be a number")
				}
				endpoint := epochstoragetypes.Endpoint{IPPORT: splitted[0], UseType: splitted[1], Geolocation: geoloc}
				argEndpoints = append(argEndpoints, endpoint)
			}
			argGeolocation, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			moniker, err := cmd.Flags().GetString(types.FlagMoniker)
			if err != nil {
				return err
			}

			msg := types.NewMsgStakeProvider(
				clientCtx.GetFromAddress().String(),
				argChainID,
				argAmount,
				argEndpoints,
				argGeolocation,
				moniker,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String(types.FlagMoniker, "", "The provider's moniker (non-unique name)")
	cmd.MarkFlagRequired(types.FlagMoniker)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
