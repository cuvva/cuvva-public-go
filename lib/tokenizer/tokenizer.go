package tokenizer

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

const encKeySize = 24
const hashKeySize = 32

const tokenVersion = 1
const tokenFormat = "%d.%s.%s"    // id.iv.ciphertext
const signedTokenFormat = "%s.%s" // tokenFormat.digest

// Tokenizer holds the key used for encrypt/decrypt operations
type Tokenizer struct {
	encKey  []byte
	hashKey []byte
}

// NewTokenizer returns a Tokenizer instance and error if key is not
// correct size
func NewTokenizer(encKey string, hashKey string) (*Tokenizer, error) {
	decEncKey, err := hex.DecodeString(encKey)
	if err != nil {
		return nil, err
	}

	decHashKey, err := hex.DecodeString(hashKey)
	if err != nil {
		return nil, err
	}

	if len(decEncKey) != encKeySize {
		return nil, errors.New("key is not 24 bytes long")
	}

	if len(decHashKey) != hashKeySize {
		return nil, errors.New("key is not 32 bytes long")
	}

	return &Tokenizer{
		decEncKey,
		decHashKey,
	}, nil
}

// Seal serializes a Token and encrypts+signs the payload
func (t *Tokenizer) Seal(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	ct, iv, err := encrypt(data, t.encKey)
	if err != nil {
		return "", err
	}

	encodedIV := base64.RawURLEncoding.EncodeToString(iv)
	encodedCT := base64.RawURLEncoding.EncodeToString(ct)
	message := fmt.Sprintf(tokenFormat, tokenVersion, encodedIV, encodedCT)

	digest, err := sign([]byte(message), t.hashKey)
	if err != nil {
		return "", err
	}

	encodedDigest := base64.RawURLEncoding.EncodeToString(digest)
	output := fmt.Sprintf(signedTokenFormat, message, encodedDigest)

	return output, nil
}

// Open an encrypted TokenRequest, verifies, unserializes and stores the result
// in the value pointed to by v. v must be a pointer or non-nil.
func (t *Tokenizer) Open(data string, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("invalid type given to open token")
	}

	parts := strings.Split(data, ".")
	if len(parts) != 4 {
		return errors.New("invalid token format")
	}

	// version.iv.ct
	messageParts := parts[:3]
	if messageParts[0] != "1" {
		return errors.New("unsupported token version")
	}

	message := strings.Join(messageParts, ".")

	// digest
	digest := parts[3]
	decodedDigest, err := base64.RawURLEncoding.DecodeString(digest)
	if err != nil {
		return err
	}

	err = verify([]byte(message), decodedDigest, t.hashKey)
	if err != nil {
		return err
	}

	// version := messageParts[0]
	decodedIV, err := base64.RawURLEncoding.DecodeString(messageParts[1])
	if err != nil {
		return err
	}

	decodedCT, err := base64.RawURLEncoding.DecodeString(messageParts[2])
	if err != nil {
		return err
	}

	token, err := decrypt(decodedCT, decodedIV, t.encKey)
	if err != nil {
		return err
	}

	return json.Unmarshal(token, v)
}

func encrypt(plaintext, key []byte) (ciphertext []byte, iv []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	plaintext = pkcs5pad(plaintext, aes.BlockSize)

	ciphertext = make([]byte, len(plaintext))
	iv = make([]byte, aes.BlockSize)

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return ciphertext, iv, nil
}

func verify(message, digest, key []byte) error {
	signed, err := sign(message, key)
	if err != nil {
		return err
	}

	if !hmac.Equal(signed, digest) {
		return errors.New("could not verify message")
	}

	return nil
}

func sign(message, key []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, key)

	_, err := mac.Write(message)
	if err != nil {
		return nil, err
	}

	return mac.Sum(nil), nil
}

func decrypt(ciphertext, iv, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	decrypted, err = pkcs5unpad(decrypted)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

func pkcs5pad(plaintext []byte, blockSize int) []byte {
	padding := (blockSize - len(plaintext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

func pkcs5unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpadding failed")
	}

	return src[:(length - unpadding)], nil
}
