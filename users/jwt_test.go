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
	validTokenString, validRefreshTokenString, _ := GenerateJWTToken(userInfo)
	type args struct {
		tokenString        string
		refreshTokenString string
	}
	tests := []struct {
		name             string
		args             args
		wantToken        string
		wantRefreshToken string
		wantErr          bool
	}{
		{
			name: "Valid Token",
			args: args{
				tokenString:        validTokenString,
				refreshTokenString: validRefreshTokenString,
			},
			wantToken:        validTokenString,
			wantRefreshToken: validRefreshTokenString,
			wantErr:          false,
		},
		{
			name: "Invalid Token",
			args: args{
				tokenString:        "badtoken",
				refreshTokenString: "badtoken",
			},
			wantToken:        "",
			wantRefreshToken: "",
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//validate token
			got, err := ValidateJWTToken(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWTToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got.Raw, tt.wantToken) {
				t.Errorf("ValidateJWTToken() = %v, want %v", got, tt.wantToken)
			}
			//validate refresh token
			got, err = ValidateJWTToken(tt.args.refreshTokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWTToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got.Raw, tt.wantRefreshToken) {
				t.Errorf("ValidateJWTToken() = %v, want %v", got, tt.wantRefreshToken)
			}
		})
	}
}

func TestRefreshJWTToken(t *testing.T) {
	userInfo := UserInfo{
		Username: "TestUsername",
		Email:    "test@test.com",
	}
	validTokenString, validRefreshTokenString, _ := GenerateJWTToken(userInfo)
	type args struct {
		tokenString        string
		refreshTokenString string
	}
	tests := []struct {
		name             string
		args             args
		wantToken        string
		wantRefreshToken string
		wantErr          bool
	}{
		{
			name: "Valid Refresh Token",
			args: args{
				tokenString:        validTokenString,
				refreshTokenString: validRefreshTokenString,
			},
			wantToken:        validTokenString,
			wantRefreshToken: validRefreshTokenString,
			wantErr:          false,
		},
		{
			name: "Invalid Refresh Token",
			args: args{
				tokenString:        "badtoken",
				refreshTokenString: "badtoken",
			},
			wantToken:        "",
			wantRefreshToken: "",
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, refreshToken, err := RefreshJWTToken(tt.args.refreshTokenString, userInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshJWTToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			//validate token
			_, err = ValidateJWTToken(token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWTToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			//validate refresh token
			_, err = ValidateJWTToken(refreshToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWTToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
