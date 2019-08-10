package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tm-db"

	"github.com/hyperspeednetwork/hsnhub/codec"
	"github.com/hyperspeednetwork/hsnhub/simapp"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestHSNdExport(t *testing.T) {
	db := db.NewMemDB()
	happ := NewHSNApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)
	setGenesis(happ)

	// Making a new app object with the db, so that initchain hasn't been called
	newGapp := NewHSNApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)
	_, _, err := newGapp.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

func setGenesis(happ *HSNApp) error {

	genesisState := simapp.NewDefaultGenesisState()
	stateBytes, err := codec.MarshalJSONIndent(happ.cdc, genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	happ.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	happ.Commit()
	return nil
}
