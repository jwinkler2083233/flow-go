package bootstrap

import (
	"github.com/dapperlabs/flow-go/model/flow"
)

func Seal(result *flow.ExecutionResult) *flow.Seal {
	// get last chunk in result
	chunk := result.Chunks[result.Chunks.Len()-1]
	seal := &flow.Seal{
		BlockID:    result.BlockID,
		ResultID:   result.ID(),
		FinalState: chunk.EndState,
	}
	return seal
}