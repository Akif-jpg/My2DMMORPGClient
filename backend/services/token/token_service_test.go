package tokenservice

import (
	"testing"
	"time"
)

const testSecret = "super-secret-key"

func TestGenerateAndValidateToken(t *testing.T) {
	body := TokenBody{
		Username:  "akif",
		UserRole:  "admin",
		SessionID: "session-123",
	}

	token, err := GenerateToken(body, "game-server", time.Minute*5, testSecret)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	parsed, err := ValidateToken(token, testSecret)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if parsed.Body.Username != body.Username {
		t.Errorf("expected username %s, got %s", body.Username, parsed.Body.Username)
	}
}

func TestInvalidSignature(t *testing.T) {
	body := TokenBody{
		Username:  "akif",
		UserRole:  "admin",
		SessionID: "session-123",
	}

	token, err := GenerateToken(body, "game-server", time.Minute*5, testSecret)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateToken(token, "wrong-secret")
	if err == nil {
		t.Fatal("expected invalid signature error")
	}
}

func TestExpiredToken(t *testing.T) {
	body := TokenBody{
		Username:  "akif",
		UserRole:  "admin",
		SessionID: "session-123",
	}

	token, err := GenerateToken(body, "game-server", -time.Minute, testSecret)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateToken(token, testSecret)
	if err == nil {
		t.Fatal("expected expiration error")
	}
}
