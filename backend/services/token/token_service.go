package tokenservice

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidFormat    = errors.New("invalid token format")
	ErrInvalidSignature = errors.New("invalid token signature")
	ErrTokenExpired     = errors.New("token expired")
)

var encoding = base64.RawURLEncoding

type TokenHeader struct {
	Producer string `json:"producer"`
	IssuedAt int64  `json:"iat"`
	ExpireAt int64  `json:"exp"`
}

type TokenBody struct {
	Username  string `json:"username"`
	UserRole  string `json:"role"`
	SessionID string `json:"sid"`
}

type Token struct {
	Header TokenHeader `json:"header"`
	Body   TokenBody   `json:"body"`
}

func GenerateToken(body TokenBody, producer string, duration time.Duration, secret string) (string, error) {
	now := time.Now()

	header := TokenHeader{
		Producer: producer,
		IssuedAt: now.Unix(),
		ExpireAt: now.Add(duration).Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	headerEncoded := encoding.EncodeToString(headerJSON)
	bodyEncoded := encoding.EncodeToString(bodyJSON)

	unsigned := headerEncoded + "." + bodyEncoded

	signature := sign(unsigned, secret)
	signatureEncoded := encoding.EncodeToString(signature)

	return unsigned + "." + signatureEncoded, nil
}

func ValidateToken(token string, secret string) (*Token, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidFormat
	}

	headerEncoded := parts[0]
	bodyEncoded := parts[1]
	signatureEncoded := parts[2]

	unsigned := headerEncoded + "." + bodyEncoded

	expectedSignature := sign(unsigned, secret)
	providedSignature, err := encoding.DecodeString(signatureEncoded)
	if err != nil {
		return nil, ErrInvalidSignature
	}

	if subtle.ConstantTimeCompare(expectedSignature, providedSignature) != 1 {
		return nil, ErrInvalidSignature
	}

	headerJSON, err := encoding.DecodeString(headerEncoded)
	if err != nil {
		return nil, ErrInvalidFormat
	}

	bodyJSON, err := encoding.DecodeString(bodyEncoded)
	if err != nil {
		return nil, ErrInvalidFormat
	}

	var header TokenHeader
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return nil, ErrInvalidFormat
	}

	if time.Now().Unix() > header.ExpireAt {
		return nil, ErrTokenExpired
	}

	var body TokenBody
	if err := json.Unmarshal(bodyJSON, &body); err != nil {
		return nil, ErrInvalidFormat
	}

	return &Token{
		Header: header,
		Body:   body,
	}, nil
}

func sign(data string, secret string) []byte {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return h.Sum(nil)
}
