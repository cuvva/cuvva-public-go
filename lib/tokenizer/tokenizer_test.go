package tokenizer_test

import (
	"crypto/aes"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/tokenizer"
)

// export GOCACHE=off is highly recommended

type TokenRequest struct {
	Profile TokenRequestProfile `json:"profile"`
}

type TokenRequestProfile struct {
	ID string `json:"id"`
}

const testEncKey = "7a4bbb543abfbb6c7286416263f9d8e42e131b1a677774c4"
const testHashKey = "68a337ac9c26affd59e91c3ea9c0dea004badb68373d163ee79af2dc198ebc68"

func TestTokenize(t *testing.T) {
	req := TokenRequest{
		Profile: TokenRequestProfile{
			ID: "ABC",
		},
	}
	_, err := testTokenize(testEncKey, testHashKey, req)

	if err != nil {
		t.Error(err)
	}
}

func TestUniqueTokenize(t *testing.T) {
	req := TokenRequest{
		Profile: TokenRequestProfile{
			ID: "ABC",
		},
	}

	a, err := testTokenize(testEncKey, testHashKey, req)
	if err != nil {
		t.Error(err)
	}

	b, err := testTokenize(testEncKey, testHashKey, req)
	if err != nil {
		t.Error(err)
	}

	if strings.Compare(a, b) == 0 {
		t.Errorf("a %s == b %s when it should not", a, b)
	}
}

func TestPaddedTokenize(t *testing.T) {
	tk, err := tokenizer.NewTokenizer(testEncKey, testHashKey)
	if err != nil {
		t.Error(err)
	}

	req := TokenRequest{
		Profile: TokenRequestProfile{
			ID: "ABCD",
		},
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		t.Error(err)
	}

	if len(reqJson)%aes.BlockSize == 0 {
		t.Errorf("need a payload not divisible by 16")
	}

	et, err := tk.Seal(req)
	if err != nil {
		t.Error(err)
	}

	dt := TokenRequest{}

	if err := tk.Open(et, &dt); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(req, dt) {
		t.Errorf("output does not match input")
	}
}

func TestLargeTokenize(t *testing.T) {
	req := TokenRequest{
		Profile: TokenRequestProfile{
			ID: "R8OmKHJvXHO1UuELcwW5ps4VgbCEvRQN9jnVQeOL5HoBz2m7+sPd+UZ9IHeGMla4f5mU+yiAAK6X7QtnpWQS4UHvRohwP6aRdRrlekzFhiHmjIzRyHNjtQR6sR7/ut+vp74g0nyoo5Htu1kcT5HH32Zdj8SqEEypXa",
		},
	}

	_, err := testTokenize(testEncKey, testHashKey, req)
	if err != nil {
		t.Error(err)
	}
}

func TestDeformedMessage(t *testing.T) {
	tk, err := tokenizer.NewTokenizer(testEncKey, testHashKey)
	if err != nil {
		t.Error(err)
	}

	req := TokenRequest{
		Profile: TokenRequestProfile{
			ID: "ABCD",
		},
	}

	et, err := tk.Seal(req)
	if err != nil {
		t.Error(err)
	}

	split := strings.Split(et, ".")

	// We're going to change a few byes on this digest
	ct := []byte(split[2])
	ct[4] = 65   // A
	ct[8] = 97   // a
	ct[12] = 121 // y

	// piece the token back together
	parts := []string{}
	parts = append(parts, split[:2]...)
	parts = append(parts, string(ct), split[3])
	joined := strings.Join(parts, ".")

	dt := TokenRequest{}

	if err := tk.Open(joined, &dt); err == nil {
		t.Errorf("should have not been able to verify message")
	}
}

func TestDigestVerification(t *testing.T) {
	tk, err := tokenizer.NewTokenizer(testEncKey, testHashKey)
	if err != nil {
		t.Error(err)
	}

	req := TokenRequest{
		Profile: TokenRequestProfile{
			ID: "ABCD",
		},
	}

	et, err := tk.Seal(req)
	if err != nil {
		t.Error(err)
	}

	split := strings.Split(et, ".")

	// We're going to change a few byes on this digest
	digest := []byte(split[3])
	digest[4] = 65 // A
	digest[8] = 97 // a

	// piece the token back together
	parts := []string{}
	parts = append(parts, split[:3]...)
	parts = append(parts, string(digest))
	joined := strings.Join(parts, ".")

	dt := TokenRequest{}

	if err := tk.Open(joined, &dt); err == nil {
		t.Errorf("should have not been able to verify message")
	}
}

func BenchmarkTokenize(b *testing.B) {
	req := TokenRequest{
		Profile: TokenRequestProfile{
			ID: "ABC",
		},
	}

	for i := 0; i < b.N; i++ {
		_, err := testTokenize(testEncKey, testHashKey, req)
		if err != nil {
			b.Error(err)
		}
	}
}

func testTokenize(encKey string, hashKey string, req TokenRequest) (string, error) { //nolint:unparam
	tk, err := tokenizer.NewTokenizer(encKey, hashKey)
	if err != nil {
		return "", err
	}

	et, err := tk.Seal(req)
	if err != nil {
		return "", err
	}

	dt := TokenRequest{}

	if err := tk.Open(et, &dt); err != nil {
		return "", err
	}

	if !reflect.DeepEqual(req, dt) {
		return "", errors.New("output does not match input")
	}

	return et, nil
}
