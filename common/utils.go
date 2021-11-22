package common

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/rkumar0099/algorand/params"
)

// maxPriority returns the highest priority of block proposal.
func MaxPriority(vrf []byte, users int) []byte {
	var maxPrior []byte
	for i := 1; i <= users; i++ {
		prior := Sha256(bytes.Join([][]byte{vrf, Uint2Bytes(uint64(i))}, nil)).Bytes()
		if bytes.Compare(prior, maxPrior) > 0 {
			maxPrior = prior
		}
	}
	return maxPrior
}

// subUsers return the selected amount of sub-users determined from the mathematics protocol.
func SubUsers(expectedNum int, weight uint64, vrf []byte) int {
	bino := NewFastBinomial(int64(weight), float64(expectedNum)/float64(params.TotalTokenAmount()))
	//binomial := NewApproxBinomial(int64(expectedNum), weight)
	//binomial := &distuv.Binomial{
	//	N: float64(weight),
	//	P: float64(expectedNum) / float64(TotalTokenAmount()),
	//}
	// hash / 2^hashlen ∉ [ ∑0,j B(k;w,p), ∑0,j+1 B(k;w,p))
	frac := 0.0
	for i := len(vrf) / 3 * 2; i >= 0; i-- {
		frac += float64(vrf[i]) / float64(math.Pow(math.Pow(2, 8), float64(i+1)))
	}
	lower := 0.0
	upper := 0.0
	for i := uint64(0); i <= weight; i++ {
		lower = upper
		upper += bino.Prob(int64(i))
		if lower <= frac && upper > frac {
			return int(i)
		}
	}
	return 0
}

func SplitIPPort(addr string) (string, int, error) {
	i := strings.IndexByte(addr, ':')
	if i == -1 {
		return "", -1, errors.New(fmt.Sprintf("%s is not a valid address", addr))
	}
	IP := addr[:i]
	portString := addr[i+1:]
	port, err := strconv.Atoi(portString)
	if err != nil {
		return "", -1, errors.New(fmt.Sprintf("%s is not a valid address: %s", addr, err.Error()))
	}
	return IP, port, nil
}

func GetVoteKey(round uint64, event string) string {
	return string(bytes.Join([][]byte{
		Uint2Bytes(round),
		[]byte(event),
	}, nil))
}
