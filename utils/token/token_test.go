// Package token @Author: youngalone [2023/8/6]
package token

import (
	"reflect"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	type args struct {
		secretKey []byte
		claims    Claims
	}
	tests := []struct {
		name            string
		args            args
		wantTokenString string
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			name: "基础功能测试",
			args: args{
				secretKey: []byte("123456"),
				claims: Claims{
					UserID:   1,
					UserName: "youngalone",
				},
			},
			wantTokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InlvdW5nYWxvbmUifQ.dS0UqrwtAOquT4QtKRj5wDfvBjZt83H0dagI3vIEtZc",
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTokenString, err := GenerateToken(tt.args.secretKey, tt.args.claims)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTokenString != tt.wantTokenString {
				t.Errorf("GenerateToken() gotTokenString = %v, want %v", gotTokenString, tt.wantTokenString)
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	type args struct {
		secretKey   []byte
		tokenString string
	}
	tests := []struct {
		name    string
		args    args
		want    *Claims
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "基础功能测试",
			args: args{
				secretKey:   []byte("123456"),
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InlvdW5nYWxvbmUifQ.dS0UqrwtAOquT4QtKRj5wDfvBjZt83H0dagI3vIEtZc",
			},
			want: &Claims{
				UserID:   1,
				UserName: "youngalone",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseToken(tt.args.secretKey, tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
