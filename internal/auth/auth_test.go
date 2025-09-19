package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPasswordHashing(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wrongPwd string
	}{
		{"ValidPassword", "howtolove", "wrongpass"},
		{"AnotherPassword", "superSecret123!", "notSecret"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Hash the password
			hash, err := HashPassword(tt.password)
			if err != nil {
				t.Fatalf("HashPassword() error = %v", err)
			}

			// Correct password should match
			if err := CheckPasswordHash(tt.password, hash); err != nil {
				t.Errorf("expected correct password to match, got error = %v", err)
			}

			// Wrong password should fail
			if err := CheckPasswordHash(tt.wrongPwd, hash); err == nil {
				t.Errorf("expected wrong password to fail, but it matched")
			}
		})
	}
}

func TestJWT(t *testing.T) {
	tokenSecret := "hapum"
	expiresIn := time.Hour
	uuids := []string{
		"123e4567-e89b-12d3-a456-426614174000",
		"550e8400-e29b-41d4-a716-446655440000",
		"6fa459ea-ee8a-3ca4-894e-db77e160355e",
	}

	for _, u := range uuids {
		t.Run(u, func(t *testing.T) {
			userID := uuid.MustParse(u)

			tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
			if err != nil {
				t.Fatalf("MakeJWT() error = %v", err)
			}

			gotID, err := ValidateJWT(tokenString, tokenSecret)
			if err != nil {
				t.Fatalf("ValidateJWT() error = %v", err)
			}

			if gotID != userID {
				t.Errorf("expected %v, got %v", userID, gotID)
			}
		})
	}
}

func TestJWT_InvalidCases(t *testing.T) {
	userID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	t.Run("wrong secret", func(t *testing.T) {
		token, _ := MakeJWT(userID, "right-secret", time.Hour)

		_, err := ValidateJWT(token, "wrong-secret")
		if err == nil {
			t.Error("expected error when validating with wrong secret")
		}
	})

	t.Run("expired token", func(t *testing.T) {
		token, _ := MakeJWT(userID, "secret", -time.Minute) // already expired

		_, err := ValidateJWT(token, "secret")
		if err == nil {
			t.Error("expected error for expired token")
		}
	})

	t.Run("malformed token", func(t *testing.T) {
		_, err := ValidateJWT("not-a-real-token", "secret")
		if err == nil {
			t.Error("expected error for malformed token")
		}
	})
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		headers   http.Header
		wantToken string
		wantErr   bool
	}{
		{
			name:      "valid bearer token",
			headers:   http.Header{"Authorization": []string{"Bearer abc123"}},
			wantToken: "abc123",
			wantErr:   false,
		},
		{
			name:    "missing authorization header",
			headers: http.Header{},
			wantErr: true,
		},
		{
			name:    "wrong scheme",
			headers: http.Header{"Authorization": []string{"Basic xyz"}},
			wantErr: true,
		},
		{
			name:      "bearer without token",
			headers:   http.Header{"Authorization": []string{"Bearer "}},
			wantToken: "",
			wantErr:   false, // it's technically valid but returns empty string
		},
		{
			name:      "extra whitespace",
			headers:   http.Header{"Authorization": []string{"Bearer    secret-token  "}},
			wantToken: "secret-token",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err)
			}
			if gotToken != tt.wantToken {
				t.Errorf("expected token %q, got %q", tt.wantToken, gotToken)
			}
		})
	}
}
