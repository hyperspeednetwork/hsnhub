package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hyperspeednetwork/hsnhub/client"
	"github.com/hyperspeednetwork/hsnhub/client/context"
	"github.com/hyperspeednetwork/hsnhub/codec"
	sdk "github.com/hyperspeednetwork/hsnhub/types"
	"github.com/hyperspeednetwork/hsnhub/x/auth"
	"github.com/hyperspeednetwork/hsnhub/x/auth/client/utils"
	"github.com/hyperspeednetwork/hsnhub/x/bank/internal/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Bank transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		SendTxCmd(cdc),
		MultiSendTxCmd(cdc),
	)
	return txCmd
}

// SendTxCmd will create a send tx and sign it with the given key.
func SendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [from_key_or_address] [to_address] [amount]",
		Short: "Create and sign a send tx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)
			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			coins, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgSend(cliCtx.GetFromAddress(), to, coins)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

// MultiSendTxCmd will create a multi account send tx and sign it with  given keys
func MultiSendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multisend [input_address:amount][,[input_address:amount]] [output_address:amount][,[output_address:amount]]",
		Short: "Create and sign a Multisend tx",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			// get the first addr from args[0]
			firstAddr, err := parseFirstAddr(args[0])
			if err != nil {
				return err
			}
			cliCtx := context.NewCLIContextWithFrom(firstAddr).WithCodec(cdc)

			inputs, err := parseInputs(args[0])
			if err != nil {
				return err
			}
			outputs, err := parseOutputs(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgMultiSend(inputs, outputs)
			// build and sign the transaction, then broadcast to Tendermint

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

func parseFirstAddr(addrAndAmountStr string) (string, error) {
	addrAndAmountStr = strings.TrimSpace(addrAndAmountStr)
	if len(addrAndAmountStr) == 0 {
		return "", fmt.Errorf("len of addrAndAmountStr is 0")
	}
	addrAndAmountStrs := strings.Split(addrAndAmountStr, ",")
	firstAddrAndamountStr := addrAndAmountStrs[0]
	if firstAddrAndamountStr == "" {
		return "", fmt.Errorf("first AddrAndamountStr is nil")
	}
	firstAddrAndamountStr = strings.TrimSpace(firstAddrAndamountStr)
	if len(firstAddrAndamountStr) == 0 {
		return "", fmt.Errorf("len of firstAddrAndamountStr is 0")
	}
	addrAndAmountArgs := strings.Split(firstAddrAndamountStr, ":")
	return addrAndAmountArgs[0], nil
}

func parseInputs(addrAndAmountStr string) ([]types.Input, error) {
	addrAndAmountStr = strings.TrimSpace(addrAndAmountStr)
	if len(addrAndAmountStr) == 0 {
		return nil, nil
	}
	addrAndAmountStrs := strings.Split(addrAndAmountStr, ",")
	inputs := make([]types.Input, len(addrAndAmountStrs))
	for index, addrAndAmount := range addrAndAmountStrs {
		if addrAndAmount == "" {
			continue
		}
		addrAndAmount = strings.TrimSpace(addrAndAmount)
		addrAndAmountArgs := strings.Split(addrAndAmount, ":")
		accAddr, err := sdk.AccAddressFromBech32(addrAndAmountArgs[0])
		if err != nil {
			return nil, err
		}
		inputCoins, err := sdk.ParseCoins(addrAndAmountArgs[1])
		if err != nil {
			return nil, err
		}
		input := types.NewInput(accAddr, inputCoins)
		inputs[index] = input
	}
	return inputs, nil
}

func parseOutputs(addrAndAmountStr string) ([]types.Output, error) {
	addrAndAmountStr = strings.TrimSpace(addrAndAmountStr)
	if len(addrAndAmountStr) == 0 {
		return nil, nil
	}
	addrAndAmountStrs := strings.Split(addrAndAmountStr, ",")
	outputs := make([]types.Output, len(addrAndAmountStrs))
	for index, addrAndAmount := range addrAndAmountStrs {
		if addrAndAmount == "" {
			continue
		}
		addrAndAmount = strings.TrimSpace(addrAndAmount)
		addrAndAmountArgs := strings.Split(addrAndAmount, ":")
		accAddr, err := sdk.AccAddressFromBech32(addrAndAmountArgs[0])
		if err != nil {
			return nil, err
		}
		outputCoins, err := sdk.ParseCoins(addrAndAmountArgs[1])
		if err != nil {
			return nil, err
		}
		output := types.NewOutput(accAddr, outputCoins)
		outputs[index] = output
	}
	return outputs, nil
}
