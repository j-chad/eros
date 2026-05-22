package webpush

import (
	"backend/internal/crypto"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

const RecordSize uint32 = 4096
const authTagSize int = 16

func encrypt(plaintext, uaPublicBytes, authSecret []byte) ([]byte, error) {
	uaPublic, err := ecdh.P256().NewPublicKey(uaPublicBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse client public key: %w", err)
	}

	// ephemeral ECDH key pair
	asPrivate, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ephemeral key: %w", err)
	}

	// ECDH shared secret
	ecdhSecret, err := asPrivate.ECDH(uaPublic)
	if err != nil {
		return nil, fmt.Errorf("failed to compute ECDH secret: %w", err)
	}

	salt := crypto.RandomBytes(16)
	cek, nonce, err := deriveKeys(ecdhSecret, authSecret, uaPublic.Bytes(), asPrivate.PublicKey().Bytes(), salt)
	if err != nil {
		return nil, fmt.Errorf("failed to derive keys: %w", err)
	}

	// prepare content-encoding header
	header := buildContentCodingHeader(salt, asPrivate.PublicKey().Bytes())

	// encrypt plaintext
	plaintextCopy := make([]byte, len(plaintext))
	copy(plaintextCopy, plaintext)
	plaintextBuf := bytes.NewBuffer(plaintextCopy)

	// padding ending delimiter
	plaintextBuf.WriteByte(0x2)

	padLength := int(RecordSize) - authTagSize - len(header)
	if err = pad(plaintextBuf, padLength); err != nil {
		return nil, fmt.Errorf("failed to pad: %w", err)
	}

	block, err := aes.NewCipher(cek)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
	}
	ciphertext := gcm.Seal(nil, nonce, plaintextBuf.Bytes(), nil)

	return append(header, ciphertext...), nil
}

func deriveKeys(ecdhSecret, authSecret, uaPublic, asPublic, salt []byte) (cek, nonce []byte, err error) {
	ikmInfo := buildIKMInfo(uaPublic, asPublic)
	ikm, err := deriveKey(authSecret, ecdhSecret, ikmInfo, 32)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive IKM: %w", err)
	}

	cek, err = deriveKey(salt, ikm, "Content-Encoding: aes128gcm\x00", 16)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive CEK: %w", err)
	}

	nonce, err = deriveKey(salt, ikm, "Content-Encoding: nonce\x00", 12)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive nonce: %w", err)
	}

	return cek, nonce, nil
}

func buildIKMInfo(uaPublic, asPublic []byte) string {
	infoBuf := bytes.NewBuffer([]byte("WebPush: info\x00"))
	infoBuf.Write(uaPublic)
	infoBuf.Write(asPublic)
	return infoBuf.String()
}

func deriveKey(salt, ikm []byte, info string, length int) ([]byte, error) {
	prk, err := hkdf.Extract(sha256.New, ikm, salt)
	if err != nil {
		return nil, err
	}

	return hkdf.Expand(sha256.New, prk, info, length)
}

// buildContentCodingHeader assembles the header according to RFC 8188
func buildContentCodingHeader(salt, localPublic []byte) []byte {
	headerBuf := bytes.NewBuffer(salt)

	rs := make([]byte, 4)
	binary.BigEndian.PutUint32(rs, RecordSize)
	headerBuf.Write(rs)

	idlen := byte(len(localPublic))
	headerBuf.WriteByte(idlen)
	headerBuf.Write(localPublic)

	return headerBuf.Bytes()
}

func pad(plaintext *bytes.Buffer, padLength int) error {
	plaintextLen := plaintext.Len()
	if plaintextLen > padLength {
		return fmt.Errorf("plaintext too large: %d bytes", plaintextLen)
	}

	padLen := padLength - plaintextLen
	padding := make([]byte, padLen)
	plaintext.Write(padding)

	return nil
}
