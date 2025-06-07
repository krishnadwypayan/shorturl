package snowflake

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNext_UniqueIDs(t *testing.T) {
	gen := NewGenerator(1)
	ids := make(map[uint64]struct{})
	const n = 1000

	for i := 0; i < n; i++ {
		id := gen.Next()
		assert.NotContains(t, ids, id, "ID should be unique")
		ids[id] = struct{}{}
	}
}

func TestNext_MachineIDBits(t *testing.T) {
	machineID := uint64(42)
	gen := NewGenerator(machineID)
	id := gen.Next()
	gotMachineID := (id >> sequenceBits) & ((1 << machineIdBits) - 1)
	assert.Equal(t, machineID, gotMachineID, "Machine ID bits should match the generator's machine ID")
}

func TestNext_SequenceIncrementsWithinSameMs(t *testing.T) {
	gen := NewGenerator(0)
	// Lock time to a fixed value by monkey-patching time.Now if needed.
	// Here, we just call Next() rapidly and check that sequence increments.
	id1 := gen.Next()
	id2 := gen.Next()
	seq1 := id1 & maxSequence
	seq2 := id2 & maxSequence
	assert.Equal(t, (seq1+1)&maxSequence, seq2, "Sequence should increment by 1 within the same millisecond")
}

func TestNext_Concurrent(t *testing.T) {
	gen := NewGenerator(3)
	const n = 1000
	ids := make(chan uint64, n)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ids <- gen.Next()
		}()
	}
	wg.Wait()
	close(ids)

	seen := make(map[uint64]struct{})
	for id := range ids {
		assert.NotContains(t, seen, id, "ID should be unique in concurrent use")
		seen[id] = struct{}{}
	}
}

func TestNewGenerator_PanicsOnInvalidMachineID(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for invalid machine ID, but did not panic")
		}
	}()
	_ = NewGenerator(1 << machineIdBits)
}

func TestNext_TimestampIncreases(t *testing.T) {
	gen := NewGenerator(0)
	id1 := gen.Next()
	time.Sleep(2 * time.Millisecond)
	id2 := gen.Next()

	ts1 := id1 >> timestampShift
	ts2 := id2 >> timestampShift

	assert.Greater(t, ts2, ts1, "Timestamp should increase between consecutive IDs")
}

func BenchmarkNext(b *testing.B) {
	gen := NewGenerator(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = gen.Next()
	}
}
