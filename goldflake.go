// Package goldflake implements Goldflake, a distributed unique ID generator inspired by Twitter's Snowflake & Sonyflake
//
// A Goldflake ID is composed of
//     39 bits for time in units of 10 msec
//     11 bits for a sequence number
//     13 bits for a machine id
package goldflake

import (
	"errors"

	// "fmt"
	"sync"
	"time"
)

// These constants are the bit lengths of Goldflake ID parts.
const (
	BitLenTime      = 39                               // bit length of time
	BitLenSequence  = 11                               // bit length of sequence number
	BitLenMachineID = 63 - BitLenTime - BitLenSequence // bit length of machine id 63-39-11=13
)

// DefaultMachineID  Default MachineID
const DefaultMachineID = 8191

// DefaultStartTime Default Start Time
var DefaultStartTime = time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC)

// Settings configures Goldflake:
//
// StartTime is the time since which the Goldflake time is defined as the elapsed time.
// If StartTime is 0, the start time of the Goldflake is set to "2020-09-01 00:00:00 +0000 UTC".
// If StartTime is ahead of the current time, Goldflake is not created.
//
// MachineID returns the unique ID of the Goldflake instance.
// If MachineID returns an error, Goldflake is not created.
// If MachineID is nil, default MachineID is used.
// Default MachineID is the max possible value based of 2^13 i.e. 8191
//
// CheckMachineID validates the uniqueness of the machine ID.
// If CheckMachineID returns false, Goldflake is not created.
// If CheckMachineID is nil, no validation is done.
type Settings struct {
	StartTime      time.Time
	MachineID      func() (uint16, error)
	CheckMachineID func(uint16) bool
}

// Goldflake is a distributed unique ID generator.
type Goldflake struct {
	mutex       *sync.Mutex
	startTime   int64
	elapsedTime int64
	sequence    uint16
	machineID   uint16
}

// NewGoldflake returns a new Goldflake configured with the given Settings.
// NewGoldflake returns nil in the following cases:
// - Settings.StartTime is ahead of the current time.
// - Settings.MachineID returns an error.
// - Settings.CheckMachineID returns false.
func NewGoldflake(st Settings) *Goldflake {
	gf := new(Goldflake)
	gf.mutex = new(sync.Mutex)
	gf.sequence = uint16(1<<BitLenSequence - 1)

	if st.StartTime.After(time.Now()) {
		return nil
	}
	if st.StartTime.IsZero() {
		gf.startTime = toGoldflakeTime(DefaultStartTime)
	} else {
		gf.startTime = toGoldflakeTime(st.StartTime)
	}

	var err error
	if st.MachineID == nil {
		gf.machineID = DefaultMachineID
	} else {
		gf.machineID, err = st.MachineID()
	}
	if err != nil || (st.CheckMachineID != nil && !st.CheckMachineID(gf.machineID)) {
		return nil
	}

	return gf
}

// NextID generates a next unique ID.
// After the Goldflake time overflows, NextID returns an error.
func (gf *Goldflake) NextID() (uint64, error) {
	const maskSequence = uint16(1<<BitLenSequence - 1)

	gf.mutex.Lock()
	defer gf.mutex.Unlock()

	current := currentElapsedTime(gf.startTime)
	if gf.elapsedTime < current {
		gf.elapsedTime = current
		gf.sequence = 0
	} else { // gf.elapsedTime >= current
		gf.sequence = (gf.sequence + 1) & maskSequence
		if gf.sequence == 0 {
			gf.elapsedTime++
			overtime := gf.elapsedTime - current
			time.Sleep(sleepTime((overtime)))
		}
	}

	return gf.toID()
}

const goldflakeTimeUnit = 1e7 // nsec, i.e. 10 msec

func toGoldflakeTime(t time.Time) int64 {
	return t.UTC().UnixNano() / goldflakeTimeUnit
}

func currentElapsedTime(startTime int64) int64 {
	return toGoldflakeTime(time.Now()) - startTime
}

func sleepTime(overtime int64) time.Duration {
	return time.Duration(overtime)*10*time.Millisecond -
		time.Duration(time.Now().UTC().UnixNano()%goldflakeTimeUnit)*time.Nanosecond
}

func (gf *Goldflake) toID() (uint64, error) {
	if gf.elapsedTime >= 1<<BitLenTime {
		return 0, errors.New("over the time limit")
	}

	return uint64(gf.elapsedTime)<<(BitLenSequence+BitLenMachineID) |
		uint64(gf.machineID)<<BitLenSequence |
		uint64(gf.sequence), nil
}

// Decompose returns a set of Goldflake ID parts.
func Decompose(id uint64) map[string]uint64 {
	const maskMachineID = uint64((1<<BitLenMachineID - 1) << BitLenSequence)
	const maskSequence = uint64(1<<BitLenSequence - 1)

	msb := id >> 63
	time := id >> (BitLenSequence + BitLenMachineID)
	machineID := id & maskMachineID >> BitLenSequence
	sequence := id & maskSequence
	return map[string]uint64{
		"id":         id,
		"msb":        msb,
		"time":       time,
		"sequence":   sequence,
		"machine-id": machineID,
	}
}
