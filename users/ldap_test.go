package users

import "testing"

func TestLDAPConn(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LDAPConn()
		})
	}
}
