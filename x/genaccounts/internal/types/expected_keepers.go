package types

import (
	sdk "github.com/hyperspeednetwork/hsnhub/types"
	authexported "github.com/hyperspeednetwork/hsnhub/x/auth/exported"
)

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	NewAccount(sdk.Context, authexported.Account) authexported.Account
	SetAccount(sdk.Context, authexported.Account)
	IterateAccounts(ctx sdk.Context, process func(authexported.Account) (stop bool))
}
