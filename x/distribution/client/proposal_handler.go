package client

import (
	"github.com/hyperspeednetwork/hsnhub/x/distribution/client/cli"
	"github.com/hyperspeednetwork/hsnhub/x/distribution/client/rest"
	govclient "github.com/hyperspeednetwork/hsnhub/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
