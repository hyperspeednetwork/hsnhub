package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/hyperspeednetwork/hsnhub/client/context"
	sdk "github.com/hyperspeednetwork/hsnhub/types"
	"github.com/hyperspeednetwork/hsnhub/types/rest"
	"github.com/hyperspeednetwork/hsnhub/x/auth/client/utils"
	"github.com/hyperspeednetwork/hsnhub/x/auth/types"
	genutilrest "github.com/hyperspeednetwork/hsnhub/x/genutil/client/rest"
)

// query accountREST Handler
func QueryAccountRequestHandlerFn(storeName string, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		accGetter := types.NewAccountRetriever(cliCtx)
		// the query will return empty account if there is no data
		if err := accGetter.EnsureExists(addr); err != nil {
			rest.PostProcessResponse(w, cliCtx, types.BaseAccount{})
			return
		}

		account, height, err := accGetter.GetAccountWithHeight(addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, account)
	}
}

// QueryTxsHandlerFn implements a REST handler that searches for transactions.
// Genesis transactions are returned if the height parameter is set to zero,
// otherwise the transactions are searched for by events.
func QueryTxsRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest,
				sdk.AppendMsgToErr("could not parse query parameters", err.Error()))
			return
		}

		// if the height query param is set to zero, query for genesis transactions
		heightStr := r.FormValue("height")
		if heightStr != "" {
			if height, err := strconv.ParseInt(heightStr, 10, 64); err == nil && height == 0 {
				genutilrest.QueryGenesisTxs(cliCtx, w)
				return
			}
		}

		var (
			events      []string
			txs         []sdk.TxResponse
			page, limit int
		)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		if len(r.Form) == 0 {
			rest.PostProcessResponseBare(w, cliCtx, txs)
			return
		}

		events, page, limit, err = rest.ParseHTTPArgs(r)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		searchResult, err := utils.QueryTxsByEvents(cliCtx, events, page, limit)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, searchResult)
	}
}

// QueryTxRequestHandlerFn implements a REST handler that queries a transaction
// by hash in a committed block.
func QueryTxRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hashHexStr := vars["hash"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		output, err := utils.QueryTx(cliCtx, hashHexStr)
		if err != nil {
			if strings.Contains(err.Error(), hashHexStr) {
				rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if output.Empty() {
			rest.WriteErrorResponse(w, http.StatusNotFound, fmt.Sprintf("no transaction found with hash %s", hashHexStr))
		}

		rest.PostProcessResponseBare(w, cliCtx, output)
	}
}
