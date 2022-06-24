package users

import (
	"reflect"
	"testing"
)

func TestValidateJWTToken(t *testing.T) {
	userInfo := UserInfo{
		Username: "TestUsername",
		Email:    "test@test.com",
	}
	validTokenString, _ := GenerateJWTToken(userInfo)
	type args struct {
		tokenString string
	}
	tests := []struct {
		name      string
		args      args
		wantToken string
		wantErr   bool
	}{
		{
			name: "Valid Token",
			args: args{
				tokenString: validTokenString,
			},
			wantToken: validTokenString,
			wantErr:   false,
		},
		{
			name: "Invalid Token",
			args: args{
				tokenString: "badtoken",
			},
			wantToken: "",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateJWTToken(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWTToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got.Raw, tt.wantToken) {
				t.Errorf("ValidateJWTToken() = %v, want %v", got, tt.wantToken)
			}
		})
	}
}
