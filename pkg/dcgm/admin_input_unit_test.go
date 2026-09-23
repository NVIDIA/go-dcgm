package dcgm

import (
	"errors"
	"strconv"
	"testing"
)

func TestConnectStandaloneV2RejectsInvalidSocketFlagBeforeInit(t *testing.T) {
	if err := connectStandaloneV2("127.0.0.1:5555", "invalid"); !errors.Is(err, strconv.ErrSyntax) {
		t.Fatalf("error = %v, want socket flag parse error", err)
	}
}
