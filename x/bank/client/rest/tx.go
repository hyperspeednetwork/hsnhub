package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/hyperspeednetwork/hsnhub/client/context"
	sdk "github.com/hyperspeednetwork/hsnhub/types"
	"github.com/hyperspeednetwork/hsnhub/types/rest"
	"github.com/hyperspeednetwork/hsnhub/x/auth/client/utils"

	"github.com/hyperspeednetwork/hsnhub/x/bank/internal/types"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/bank/accounts/{address}/transfers", SendRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/bank/balances/{address}", QueryBalancesRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/bank/accounts/multisend", MultiSendHandleFn(cliCtx)).Methods("POST")
}

// SendReq defines the properties of a send request's body.
type SendReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Amount  sdk.Coins    `json:"amount" yaml:"amount"`
}

// SendRequestHandlerFn - http request handler to send coins to a address.
func SendRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["address"]

		toAddr, err := sdk.AccAddressFromBech32(bech32Addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req SendReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgSend(fromAddr, toAddr, req.Amount)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// OutputReq - defines multisend output account req body
type OutputReq struct {
	To     string    `json:"to" yaml:"to"`
	Amount sdk.Coins `json:"amount" yaml:"amount"`
}

// InputReq - defines multisend Input account req body
type InputReq struct {
	From   string    `json:"from" yaml:"from"`
	Amount sdk.Coins `json:"amount" yaml:"amount"`
}

// MultiSendReq - define the properties of multisend req body
type MultiSendReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Input   []InputReq   `json:"input_req" yaml:"input_req"`
	Output  []OutputReq  `json:"output_req" yaml:"output_req"`
}

// MultiSendHandleFn - http request handler to get coin fron different account and send coin to different account
func MultiSendHandleFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MultiSendReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		inputs, err := parseInputs(req.Input)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		outputs, err := parseOutputs(req.Output)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		msg := types.NewMsgMultiSend(inputs, outputs)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func parseInputs(inputReqs []InputReq) ([]types.Input, error) {
	if len(inputReqs) == 0 {
		return nil, fmt.Errorf("input account address is nil")
	}
	inputs := make([]types.Input, len(inputReqs))
	for i, inputReq := range inputReqs {
		// check addr
		accAddr, err := sdk.AccAddressFromBech32(inputReq.From)
		if err != nil {
			return nil, err
		}
		var coins = inputReq.Amount
		// check coins
		if !coins.IsValid() {
			return nil, fmt.Errorf("parseCoins invalid: %#v", coins)
		}
		input := types.NewInput(accAddr, coins)
		inputs[i] = input
	}
	return inputs, nil
}

func parseOutputs(outputReqs []OutputReq) ([]types.Output, error) {
	if len(outputReqs) == 0 {
		return nil, fmt.Errorf("output account address is nil")
	}
	outputs := make([]types.Output, len(outputReqs))
	for i, outputReq := range outputReqs {
		// check addr
		accAddr, err := sdk.AccAddressFromBech32(outputReq.To)
		if err != nil {
			return nil, err
		}
		var coins = outputReq.Amount
		// check coins
		if !coins.IsValid() {
			return nil, fmt.Errorf("parseCoins invalid: %#v", coins)
		}
		output := types.NewOutput(accAddr, coins)
		outputs[i] = output
	}
	return outputs, nil
}
