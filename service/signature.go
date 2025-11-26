package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type SignatureGenerator interface {
	Generate(payload []byte) string
}

type SignatureVerifier interface {
	Verify(payload []byte, signature string) bool
}

type HMACSignature struct {
	secret string
}

func NewHMACSignature(secret string) *HMACSignature {
	return &HMACSignature{secret: secret}
}

func (h *HMACSignature) Generate(payload []byte) string {
	mac := hmac.New(sha256.New, []byte(h.secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func (h *HMACSignature) Verify(payload []byte, signature string) bool {
	expectedSignature := h.Generate(payload)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
