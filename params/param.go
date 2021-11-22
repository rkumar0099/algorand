package params

import "time"

var (
	UserAmount     uint64 = 50
	TokenPerUser   uint64 = 10000
	Malicious      uint64 = 0
	NetworkLatency        = 0
)

func TotalTokenAmount() uint64 { return UserAmount * TokenPerUser }

const (
	// Algorand system parameters
	ExpectedBlockProposers        = 20
	ExpectedCommitteeMembers      = 25
	ExpectedOraclePeers           = 20
	ThresholdOfBAStep             = 0.35
	ExpectedFinalCommitteeMembers = 10
	FinalThreshold                = 0.37
	MAXSTEPS                      = 50

	// timeout param
	LamdaPriority = 5 * time.Second  // time to gossip sortition proofs.
	LamdaBlock    = 1 * time.Minute  // timeout for receiving a block.
	LamdaStep     = 20 * time.Second // timeout for BA* step.
	LamdaStepvar  = 5 * time.Second  // estimate of BA* completion time variance.

	// interval
	R                   = 1000          // seed refresh interval (# of rounds)
	ForkResolveInterval = 1 * time.Hour // fork resolve interval time

	// helper const var
	Committee           = "committee"
	Proposer            = "proposer"
	OraclePeer          = "oraclepeer"
	OracleBlockProposer = "oracleproposer"

	// step
	ORACLE              = "EVT_ORACLE"
	ORACLE_BLK_PROPOSAL = "oracleblkproposal"
	PROPOSE             = "EVT_PROPOSE"
	REDUCTION_ONE       = "EVT_RED_1"
	REDUCTION_TWO       = "EVT_RED_2"
	FINAL               = "EVT_FINAL"

	FINAL_CONSENSUS     = 0
	TENTATIVE_CONSENSUS = 1

	// Malicious type
	Honest = iota
	EvilBlockProposal
	EvilVoteEmpty
	EvilVoteNothing
)
