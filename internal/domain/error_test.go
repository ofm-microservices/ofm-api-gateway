package gateway

import (
	"errors"
	"testing"
)

func TestIsBusinessError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "not found", err: ErrGigNotFound, want: true},
		{name: "registration conflict", err: &RegistrationConflictError{}, want: true},
		{name: "internal", err: errors.New("database unavailable"), want: false},
		{name: "nil", err: nil, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBusinessError(tt.err); got != tt.want {
				t.Fatalf("IsBusinessError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
