package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/hyperspeednetwork/hsnhub/types"
	"github.com/hyperspeednetwork/hsnhub/x/slashing/internal/types"
)

func TestGetSetValidatorSigningInfo(t *testing.T) {
	ctx, _, _, _, keeper := CreateTestInput(t, types.DefaultParams())
	info, found := keeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(Addrs[0]))
	require.False(t, found)
	newInfo := types.NewValidatorSigningInfo(
		sdk.ConsAddress(Addrs[0]),
		int64(4),
		int64(3),
		time.Unix(2, 0),
		false,
		int64(10),
	)
	keeper.SetValidatorSigningInfo(ctx, sdk.ConsAddress(Addrs[0]), newInfo)
	info, found = keeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(Addrs[0]))
	require.True(t, found)
	require.Equal(t, info.StartHeight, int64(4))
	require.Equal(t, info.IndexOffset, int64(3))
	require.Equal(t, info.JailedUntil, time.Unix(2, 0).UTC())
	require.Equal(t, info.MissedBlocksCounter, int64(10))
}

func TestGetSetValidatorMissedBlockBitArray(t *testing.T) {
	ctx, _, _, _, keeper := CreateTestInput(t, types.DefaultParams())
	missed := keeper.GetValidatorMissedBlockBitArray(ctx, sdk.ConsAddress(Addrs[0]), 0)
	require.False(t, missed) // treat empty key as not missed
	keeper.SetValidatorMissedBlockBitArray(ctx, sdk.ConsAddress(Addrs[0]), 0, true)
	missed = keeper.GetValidatorMissedBlockBitArray(ctx, sdk.ConsAddress(Addrs[0]), 0)
	require.True(t, missed) // now should be missed
}
