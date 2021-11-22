package common

import (
	"errors"
	"sync"

	"github.com/rcrowley/go-metrics"
)

var (
	lock                 *sync.Mutex
	Log                  = GetLogger("algorand")
	ErrCountVotesTimeout = errors.New("count votes timeout")

	// global metrics
	//Proposers                        = 0
	MetricsRound              uint64 = 1
	ProposerSelectedCounter          = metrics.NewRegisteredCounter("blockproposal/subusers/count", nil)
	ProposerSelectedHistogram        = metrics.NewRegisteredHistogram("blockproposal/subusers", nil, metrics.NewUniformSample(1028))
)
