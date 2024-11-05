package lib

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is in an invalid format")
	ErrIncompatibleVersion = errors.New("the encoded hash is using an incompatible version of Argon2")
)

type GenerateParams struct {
	SaltLength  uint
	Iterations  uint32
	Memory      uint32
	Parallelism uint8
	KeyLength   uint32
}

var DefaultGenerateParams = &GenerateParams{
	SaltLength:  32,
	Iterations:  4,
	Memory:      64 * 1024,
	Parallelism: 1, // Only 1 thread available in the current Cloud Run environment.
	KeyLength:   32,
}

func GenerateFromPassword(password string, params *GenerateParams) (string, error) {
	// Generate a cryptographically secure random salt.
	salt, err := Random(params.SaltLength)
	if err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}

	// Pass the plaintext password, salt and parameters to the argon2.IDKey
	// function. This will generate a hash of the password using the Argon2id
	// variant.
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.Memory,
		params.Iterations,
		params.Parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

func ComparePasswordAndHash(password, encodedHash string) (bool, error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	params, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

func decodeHash(encodedHash string) (*GenerateParams, []byte, []byte, error) {
	values := strings.Split(encodedHash, "$")
	if len(values) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err := fmt.Sscanf(values[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, errors.Join(ErrInvalidHash, fmt.Errorf("parse version: %w", err))
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	params := &GenerateParams{}
	_, err = fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)
	if err != nil {
		return nil, nil, nil, errors.Join(ErrInvalidHash, fmt.Errorf("parse parameters: %w", err))
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(values[4])
	if err != nil {
		return nil, nil, nil, errors.Join(ErrInvalidHash, fmt.Errorf("decode salt: %w", err))
	}
	params.SaltLength = uint(len(salt))

	hash, err := base64.RawStdEncoding.Strict().DecodeString(values[5])
	if err != nil {
		return nil, nil, nil, errors.Join(ErrInvalidHash, fmt.Errorf("decode hash: %w", err))
	}

	rawHashLength := len(hash)
	if rawHashLength > math.MaxUint32 {
		return nil, nil, nil, fmt.Errorf("%w: hash length: %d", ErrInvalidHash, rawHashLength)
	}

	params.KeyLength = uint32(rawHashLength)

	return params, salt, hash, nil
}
