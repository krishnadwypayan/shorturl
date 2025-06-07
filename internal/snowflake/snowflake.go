package snowflake

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/krishnadwypayan/shorturl/internal/encoder"
	"github.com/krishnadwypayan/shorturl/internal/logger"
)

// A typical Snowflake ID (64 bits) layout is:
// | 1-bit sign | 41-bit timestamp | 10-bit machine id | 12-bit sequence |

const (
	epoch          = int64(1288834974657)         // Twitter's epoch in milliseconds
	machineIdBits  = 10                           // Number of bits for machine ID
	sequenceBits   = 12                           // Number of bits for sequence number
	timestampBits  = 42                           // Number of bits for timestamp
	timestampShift = machineIdBits + sequenceBits // Shift for timestamp in the ID
	maxSequence    = (1 << sequenceBits) - 1      // Maximum value for sequence number
)

type Generator struct {
	machineId uint64
	state     uint64
}

func NewGenerator(machineId uint64) *Generator {
	if machineId >= (1 << machineIdBits) {
		panic(fmt.Sprintf("machine ID must be less than %d", 1<<machineIdBits))
	}

	return &Generator{
		machineId: machineId,
		state:     0,
	}
}

func (g *Generator) NextString() string {
	id := g.Next()
	return encoder.EncodeBase62(id)
}

func (g *Generator) Next() uint64 {
	for {
		currentState := atomic.LoadUint64(&g.state)

		lastTimestamp := currentState >> timestampShift
		now := time.Now().UnixMilli() - epoch

		if now < int64(lastTimestamp) {
			logger.Info().Msg("Clock is moving backwards, trying to generate ID again")
			continue
		}

		sequence := currentState & maxSequence

		var newTimestamp int64
		var newSequence uint64

		if now == int64(lastTimestamp) {
			logger.Debug().Msg("Current timestamp is the same as last, incrementing sequence")
			newSequence = (sequence + 1) & maxSequence
			if newSequence == 0 {
				logger.Debug().Msg("Sequence overflow")

				for now <= int64(lastTimestamp) {
					now = time.Now().UnixMilli() - epoch
				}
				newTimestamp = now
			} else {
				newTimestamp = int64(lastTimestamp)
			}
		} else {
			logger.Debug().Msg("New timestamp detected, resetting sequence")
			newTimestamp = now
			newSequence = 0
		}

		newState := (uint64(newTimestamp) << timestampShift) | (g.machineId << sequenceBits) | newSequence
		if atomic.CompareAndSwapUint64(&g.state, currentState, newState) {
			logger.Debug().Msg(fmt.Sprintf("Generated new Snowflake ID: %d", newState))
			return newState
		}
	}
}
