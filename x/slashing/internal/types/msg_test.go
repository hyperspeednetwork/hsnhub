package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/hyperspeednetwork/hsnhub/types"
)

func TestMsgUnjailGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress("abcd")
	msg := NewMsgUnjail(sdk.ValAddress(addr))
	bytes := msg.GetSignBytes()
	require.Equal(
		t,
		`{"type":"cosmos-sdk/MsgUnjail","value":{"address":"hsnvaloper1v93xxequu289j"}}`,
		string(bytes),
	)
}
