package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Passkey struct {
	bun.BaseModel `bun:"table:passkeys,select:active_passkeys"`

	ID        uuid.UUID `bun:"id,pk,type:uuid"`
	Namespace string    `bun:"namespace,pk"`

	EncryptedKey string                 `bun:"encrypted_key"`
	Reward       map[string]interface{} `bun:"reward"`

	ExpiresAt *time.Time `bun:"expires_at"`
	CreatedAt time.Time  `bun:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at"`
}
