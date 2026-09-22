package dcgm

import (
	"math"
	"os"
	"testing"
)

func TestTimeStringRunning(t *testing.T) {
	if got := Time(0).String(); got != "Running" {
		t.Fatalf("Time(0).String() = %q, want Running", got)
	}
}

func TestProcessName(t *testing.T) {
	name, err := processName(uint(os.Getpid()))
	if err != nil || name == "" {
		t.Fatalf("processName(current PID) = %q, %v", name, err)
	}

	name, err = processName(math.MaxUint32)
	if err != nil || name != "" {
		t.Fatalf("processName(missing PID) = %q, %v, want empty nil", name, err)
	}
}
