package types

import (
	"github.com/hyperspeednetwork/hsnhub/codec"
	sdk "github.com/hyperspeednetwork/hsnhub/types"
	authtypes "github.com/hyperspeednetwork/hsnhub/x/auth/types"
	stakingtypes "github.com/hyperspeednetwork/hsnhub/x/staking/types"
)

// ModuleCdc defines a generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

// TODO: abstract genesis transactions registration back to staking
// required for genesis transactions
func init() {
	ModuleCdc = codec.New()
	stakingtypes.RegisterCodec(ModuleCdc)
	authtypes.RegisterCodec(ModuleCdc)
	sdk.RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
