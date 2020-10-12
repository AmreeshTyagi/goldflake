package goldflake

import (
	"fmt"
	"runtime"

	"testing"
	"time"
)

var gf *Goldflake

var startTime int64
var machineID uint16
var startTimeEpocInMillis int64

func init() {
	var st Settings
	var sTime = time.Now()
	st.StartTime = sTime
	startTimeEpocInMillis = sTime.UnixNano() / 1000000
	fmt.Println("StartTime:", sTime)
	fmt.Println("StartTime EPOC:", startTimeEpocInMillis)

	gf = NewGoldflake(st)
	if gf == nil {
		panic("goldflake not created")
	}

	machineID = 8191

	startTime = toGoldflakeTime(st.StartTime)
	gf.machineID = machineID
}

func nextID(t *testing.T) uint64 {
	id, err := gf.NextID()
	if err != nil {
		t.Fatal("id not generated")
	}
	return id
}

func TestGoldflakeOnce(t *testing.T) {
	sleepTime := uint64(50)
	time.Sleep(time.Duration(sleepTime) * 10 * time.Millisecond)

	id := nextID(t)
	parts := Decompose(id)

	actualMSB := parts["msb"]
	if actualMSB != 0 {
		t.Errorf("unexpected msb: %d", actualMSB)
	}

	actualTime := parts["time"]
	if actualTime < sleepTime || actualTime > sleepTime+1 {
		if actualTime < sleepTime {
			t.Errorf("unexpected time actualTime < sleepTime : %d", actualTime)
		}
		if actualTime > sleepTime+1 {
			t.Errorf("unexpected time  actualTime > sleepTime+1 : %d", actualTime)
		}
	}

	actualSequence := parts["sequence"]
	if actualSequence != 0 {
		t.Errorf("unexpected sequence: %d", actualSequence)
	}

	actualMachineID := parts["machine-id"]
	if uint16(actualMachineID) != machineID {
		t.Errorf("unexpected machine id: %d", actualMachineID)
	}

	fmt.Println("goldflake id:", id)
	fmt.Println("decompose:", parts)
}

func currentTime() int64 {
	return toGoldflakeTime(time.Now())
}

func TestGoldflakeFor10Sec(t *testing.T) {
	var numID uint32
	var lastID uint64
	var maxSequence uint64

	initial := currentTime()
	current := initial
	for current-initial < 1000 {
		id := nextID(t)
		parts := Decompose(id)
		numID++

		if id <= lastID {
			t.Fatal("duplicated id")
		}
		lastID = id

		current = currentTime()

		actualMSB := parts["msb"]
		if actualMSB != 0 {
			t.Errorf("unexpected msb: %d", actualMSB)
		}

		actualTime := int64(parts["time"])
		overtime := startTime + actualTime - current
		if overtime > 0 {
			t.Errorf("unexpected overtime: %d", overtime)
		}

		actualSequence := parts["sequence"]
		if maxSequence < actualSequence {
			maxSequence = actualSequence
		}

		actualMachineID := parts["machine-id"]
		if uint16(actualMachineID) != machineID {
			t.Errorf("unexpected machine id: %d", actualMachineID)
		}
	}

	if maxSequence != 1<<BitLenSequence-1 {
		t.Errorf("unexpected max sequence: %d", maxSequence)
	}
	fmt.Println("max sequence:", maxSequence)
	fmt.Println("number of id:", numID)
}

func TestGoldflakeInParallel(t *testing.T) {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	fmt.Println("number of cpu:", numCPU)

	consumer := make(chan uint64)

	const numID = 10000
	generate := func() {
		for i := 0; i < numID; i++ {
			consumer <- nextID(t)
		}
	}

	const numGenerator = 10
	for i := 0; i < numGenerator; i++ {
		go generate()
	}

	set := make(map[uint64]struct{})
	for i := 0; i < numID*numGenerator; i++ {
		id := <-consumer
		if _, ok := set[id]; ok {
			t.Fatal("duplicated id")
		}
		set[id] = struct{}{}
	}
	fmt.Println("number of id:", len(set))
}

func TestNilGoldflake(t *testing.T) {
	var startInFuture Settings
	startInFuture.StartTime = time.Now().Add(time.Duration(1) * time.Minute)
	if NewGoldflake(startInFuture) != nil {
		t.Errorf("goldflake starting in the future")
	}

	var noMachineID Settings
	noMachineID.MachineID = func() (uint16, error) {
		return 0, fmt.Errorf("no machine id")
	}
	if NewGoldflake(noMachineID) != nil {
		t.Errorf("goldflake with no machine id")
	}

	var invalidMachineID Settings
	invalidMachineID.CheckMachineID = func(uint16) bool {
		return false
	}
	if NewGoldflake(invalidMachineID) != nil {
		t.Errorf("goldflake with invalid machine id")
	}
}

func pseudoSleep(period time.Duration) {
	gf.startTime -= int64(period) / goldflakeTimeUnit
}

func TestNextIDError(t *testing.T) {
	year := time.Duration(365*24) * time.Hour
	pseudoSleep(time.Duration(174) * year)
	nextID(t)

	pseudoSleep(time.Duration(1) * year)
	_, err := gf.NextID()
	if err == nil {
		t.Errorf("time is not over")
	}
}
