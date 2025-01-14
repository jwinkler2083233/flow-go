package execution

import (
	"fmt"

	"github.com/onflow/flow-go/model/convert"
	"github.com/onflow/flow-go/model/flow"
)

func GenerateExecutionResultAndChunkDataPacks(
	prevResultId flow.Identifier,
	startState flow.StateCommitment,
	result *ComputationResult) (
	endState flow.StateCommitment,
	chdps []*flow.ChunkDataPack,
	executionResult *flow.ExecutionResult,
	err error,
) {

	// no need to persist the state interactions, since they are used only by state
	// syncing, which is currently disabled
	block := result.ExecutableBlock.Block
	blockID := block.ID()

	chunks := make([]*flow.Chunk, len(result.StateCommitments))
	chdps = make([]*flow.ChunkDataPack, len(result.StateCommitments))

	// TODO: check current state root == startState
	endState = startState

	for i := range result.StateCommitments {
		// TODO: deltas should be applied to a particular state

		endState = result.StateCommitments[i]
		var collectionID flow.Identifier

		txNumber := 1 //default for system chunk

		// account for system chunk being last
		if i < len(result.StateCommitments)-1 {
			collectionGuarantee := result.ExecutableBlock.Block.Payload.Guarantees[i]
			completeCollection := result.ExecutableBlock.CompleteCollections[collectionGuarantee.ID()]
			collectionID = completeCollection.Collection().ID()
			txNumber = len(completeCollection.Transactions)
		} else {
			collectionID = flow.ZeroID
		}

		eventsHash := result.EventsHashes[i]
		chunk := GenerateChunk(i, startState, endState, collectionID, blockID, eventsHash, uint16(txNumber))

		// chunkDataPack
		chdps[i] = GenerateChunkDataPack(chunk, collectionID, result.Proofs[i])
		// TODO use view.SpockSecret() as an input to spock generator
		chunks[i] = chunk
		startState = endState
	}

	executionResult, err = GenerateExecutionResultForBlock(prevResultId, block, chunks, result.ServiceEvents)
	if err != nil {
		return flow.DummyStateCommitment, nil, nil, fmt.Errorf("could not generate execution result: %w", err)
	}

	return endState, chdps, executionResult, nil
}

// GenerateExecutionResultForBlock creates new ExecutionResult for a block from
// the provided chunk results.
func GenerateExecutionResultForBlock(
	previousErID flow.Identifier,
	block *flow.Block,
	chunks []*flow.Chunk,
	serviceEvents []flow.Event,
) (*flow.ExecutionResult, error) {

	// convert Cadence service event representation to flow-go representation
	convertedServiceEvents := make([]flow.ServiceEvent, 0, len(serviceEvents))
	for _, event := range serviceEvents {
		converted, err := convert.ServiceEvent(event)
		if err != nil {
			return nil, fmt.Errorf("could not convert service event: %w", err)
		}
		convertedServiceEvents = append(convertedServiceEvents, *converted)
	}

	er := &flow.ExecutionResult{
		PreviousResultID: previousErID,
		BlockID:          block.ID(),
		Chunks:           chunks,
		ServiceEvents:    convertedServiceEvents,
	}

	return er, nil
}

// GenerateChunk creates a chunk from the provided computation data.
func GenerateChunk(colIndex int,
	startState, endState flow.StateCommitment,
	colID, blockID, eventsCollection flow.Identifier, txNumber uint16) *flow.Chunk {
	return &flow.Chunk{
		ChunkBody: flow.ChunkBody{
			CollectionIndex: uint(colIndex),
			StartState:      startState,
			EventCollection: eventsCollection,
			BlockID:         blockID,
			// TODO: record gas used
			TotalComputationUsed: 0,
			NumberOfTransactions: txNumber,
		},
		Index:    uint64(colIndex),
		EndState: endState,
	}
}

// generateChunkDataPack creates a chunk data pack
func GenerateChunkDataPack(
	chunk *flow.Chunk,
	collectionID flow.Identifier,
	proof flow.StorageProof,
) *flow.ChunkDataPack {
	return &flow.ChunkDataPack{
		ChunkID:      chunk.ID(),
		StartState:   chunk.StartState,
		Proof:        proof,
		CollectionID: collectionID,
	}
}
