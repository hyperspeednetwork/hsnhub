package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdkerrors "github.com/hyperspeednetwork/hsnhub/types/errors"
)

var codeTypes = []CodeType{
	CodeInternal,
	CodeTxDecode,
	CodeInvalidSequence,
	CodeUnauthorized,
	CodeInsufficientFunds,
	CodeUnknownRequest,
	CodeInvalidAddress,
	CodeInvalidPubKey,
	CodeUnknownAddress,
	CodeInsufficientCoins,
	CodeInvalidCoins,
	CodeOutOfGas,
	CodeMemoTooLarge,
}

type errFn func(msg string) Error

var errFns = []errFn{
	ErrInternal,
	ErrTxDecode,
	ErrInvalidSequence,
	ErrUnauthorized,
	ErrInsufficientFunds,
	ErrUnknownRequest,
	ErrInvalidAddress,
	ErrInvalidPubKey,
	ErrUnknownAddress,
	ErrInsufficientCoins,
	ErrInvalidCoins,
	ErrOutOfGas,
	ErrMemoTooLarge,
}

func TestCodeType(t *testing.T) {
	require.True(t, CodeOK.IsOK())

	for tcnum, c := range codeTypes {
		msg := CodeToDefaultMsg(c)
		require.NotEqual(t, unknownCodeMsg(c), msg, "Code expected to be known. tc #%d, code %d, msg %s", tcnum, c, msg)
	}

	msg := CodeToDefaultMsg(CodeOK)
	require.Equal(t, unknownCodeMsg(CodeOK), msg)
}

func TestErrFn(t *testing.T) {
	for i, errFn := range errFns {
		err := errFn("")
		codeType := codeTypes[i]
		require.Equal(t, err.Code(), codeType, "Err function expected to return proper code. tc #%d", i)
		require.Equal(t, err.Codespace(), CodespaceRoot, "Err function expected to return proper codespace. tc #%d", i)
		require.Equal(t, err.QueryResult().Code, uint32(err.Code()), "Err function expected to return proper Code from QueryResult. tc #%d")
		require.Equal(t, err.QueryResult().Log, err.ABCILog(), "Err function expected to return proper ABCILog from QueryResult. tc #%d")
	}
}

func TestAppendMsgToErr(t *testing.T) {
	for i, errFn := range errFns {
		err := errFn("")
		errMsg := err.Stacktrace().Error()
		abciLog := err.ABCILog()

		// plain msg error
		msg := AppendMsgToErr("something unexpected happened", errMsg)
		require.Equal(
			t,
			fmt.Sprintf("something unexpected happened; %s", errMsg),
			msg,
			fmt.Sprintf("Should have formatted the error message of ABCI Log. tc #%d", i),
		)

		// ABCI Log msg error
		msg = AppendMsgToErr("something unexpected happened", abciLog)
		msgIdx := mustGetMsgIndex(abciLog)
		require.Equal(
			t,
			fmt.Sprintf("%s%s; %s}",
				abciLog[:msgIdx],
				"something unexpected happened",
				abciLog[msgIdx:len(abciLog)-1],
			),
			msg,
			fmt.Sprintf("Should have formatted the error message of ABCI Log. tc #%d", i))
	}
}

func TestResultFromError(t *testing.T) {
	cases := map[string]struct {
		err    error
		expect Result
	}{
		"sdk.Error": {
			err: ErrUnauthorized("not owner"),
			expect: Result{
				Codespace: CodespaceRoot,
				Code:      CodeUnauthorized,
				Log:       `{"codespace":"sdk","code":4,"message":"not owner"}`,
			},
		},
		"types/errors": {
			err: sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "not owner"),
			expect: Result{
				Codespace: CodespaceRoot,
				Code:      CodeUnauthorized,
				Log:       `{"codespace":"sdk","code":4,"message":"not owner: unauthorized"}`,
			},
		},
		"stdlib errors": {
			err: fmt.Errorf("not owner"),
			expect: Result{
				Codespace: CodespaceType("undefined"),
				Code:      CodeInternal,
				// note that we redact the internal errors in the new package to not leak eg. panics
				Log: `{"codespace":"undefined","code":1,"message":"internal error"}`,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			res := ResultFromError(tc.err)
			require.Equal(t, tc.expect, res)
		})
	}
}
