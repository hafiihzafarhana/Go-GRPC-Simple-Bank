package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"` // digunakan untuk mencegah atau sebagai mekanisme kebocoran yang akan dibatalkan
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issue_at"`   // untuk mengetahui kapan token dibuat
	ExpiredAt time.Time `json:"expired_at"` // mengetahui kapan token expired
}

var ErrTokenExpired = errors.New("token has expired")
var ErrInvalidToken = errors.New("token is not valid")

// Fungsi ini digunakan untuk membuat payload token baru dengan username tertentu dan durasi
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenId,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

// Memeriksa validitas apakah token payload adalah valid atau tidak
// Func Valid() ini adalah default yang diberikan oleh JWT ke pada paylod yang kita miliki
// fungsi ini akan dipanggil pada parsing token
func (payload Payload) Valid() error {
	// periksa apakah token kedaluarsa
	if time.Now().After(payload.ExpiredAt) {
		return ErrTokenExpired
	}

	return nil
}
