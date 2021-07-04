// (c) 2019 Dapper Labs - ALL RIGHTS RESERVED

package flow


import (
	"os"
	"strconv"

	"github.com/aristanetworks/goarista/monotime"
)

// A Seal is produced when an Execution Result (referenced by `ResultID`) for
// particular block (referenced by `BlockID`) is committed into the chain.
// A Seal for a block B can be included in the payload B's descendants. Only
// in the respective fork where the seal for B is included, the referenced
// result is considered committed. Different forks might contain different
// seals for the same result (or in edge cases, even for different results).
//
// NOTES
// (1) As Seals are (currently) included in the payload, they are not strictly
// entities. (Entities can be sent between nodes as self-contained messages
// whose integrity is protected by a signature). By itself, a seal does
// _not_ contain enough information to determine its validity (verifier
// assignment cannot be computed) and its integrity is not protected by a
// signature of a node that is authorized to generate it. A seal should only
// be processed in the context of the block, which contains it.
//
// (2) Even though seals are not strictly entities, they still implement the
// Entity interface. This allows us to store and retrieve seals individually.
// CAUTION: As seals are part of the block payload, their _exact_ content must
// be preserved by the storage system. This includes the exact list of approval
// signatures (incl. order). While it is possible to construct different valid
// seals for the same result (using different subsets of assigned verifiers),
// they cannot be treated as equivalent for the following reason:
//
//  * Swapping a seal in a block with a different once changes the binary
//    representation of the block payload containing the seal.
//  * Changing the binary block representation would invalidate the block
//    proposer's signature.
//
// Therefore, to retrieve valid blocks from storage, it is required that
// the Seal.ID includes all fields with independent degrees of freedom
// (such as AggregatedApprovalSigs).
//
type Seal struct {
	BlockID                Identifier
	ResultID               Identifier
	FinalState             StateCommitment
	AggregatedApprovalSigs []AggregatedSignature // one AggregatedSignature per chunk

	// Service Events are copied from the Execution Result. Therefore, repeating the
	// the service events here opens the possibility for a data-inconsistency attack.
	// It is _not_ necessary to repeat the ServiceEvents here, as an Execution Result
	// must be incorporated into the fork before it can be sealed.
	// TODO: include ServiceEvents in Execution Result and remove from Seal
	ServiceEvents []ServiceEvent
}

func (s Seal) Body() interface{} {
	return struct {
		BlockID                Identifier
		ResultID               Identifier
		FinalState             StateCommitment
		AggregatedApprovalSigs []AggregatedSignature
	}{
		BlockID:                s.BlockID,
		ResultID:               s.ResultID,
		FinalState:             s.FinalState,
		AggregatedApprovalSigs: s.AggregatedApprovalSigs,
	}
}

var logfile_seal_id *os.File

func (s Seal) ID() Identifier {
	once.Do(func() {
		newfile, _ := os.Create("/data/seal_id.log")
		logfile_seal_id = newfile
	})
	ts := monotime.Now()

	defer logfile_seal_id.WriteString(strconv.FormatUint(monotime.Now(), 10) + "," +
		strconv.FormatUint(monotime.Now() - ts, 10) + "\n")

	return MakeID(s.Body())
}

var logfile_seal_chk *os.File

func (s Seal) Checksum() Identifier {
	once.Do(func() {
		newfile, _ := os.Create("/data/seal_chk.log")
		logfile_seal_chk = newfile
	})
	ts := monotime.Now()

	defer logfile_seal_chk.WriteString(strconv.FormatUint(monotime.Now(), 10) + "," +
		strconv.FormatUint(monotime.Now() - ts, 10) + "\n")
	return MakeID(s)
}
