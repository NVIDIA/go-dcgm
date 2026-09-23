package dcgm

import "testing"

func TestCallbackLimitExceeded(t *testing.T) {
	tests := []struct {
		name     string
		current  int
		incoming int
		want     bool
	}{
		{name: "empty"},
		{name: "one slot remains", current: maxCallbackValues - 1, incoming: 1},
		{name: "exactly full", current: maxCallbackValues},
		{name: "one over from full", current: maxCallbackValues, incoming: 1, want: true},
		{name: "batch crosses limit", current: maxCallbackValues - 1, incoming: 2, want: true},
		{name: "already over", current: maxCallbackValues + 1, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := callbackLimitExceeded(tt.current, tt.incoming); got != tt.want {
				t.Fatalf("callbackLimitExceeded(%d, %d) = %t, want %t",
					tt.current, tt.incoming, got, tt.want)
			}
		})
	}
}
