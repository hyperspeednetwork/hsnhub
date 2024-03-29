package client

import (
	govclient "github.com/hyperspeednetwork/hsnhub/x/gov/client"
	"github.com/hyperspeednetwork/hsnhub/x/params/client/cli"
	"github.com/hyperspeednetwork/hsnhub/x/params/client/rest"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
