package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// ini adalah Paseto token maker, seperti JWT
type PasetoMaker struct {
	paseto *paseto.V2 // menggunakan paseto versi 2
	// menggunakan simetris: melakukan enkripsi dan deskripsi dengan 1 kunci saja
	symmetricKey []byte
}

// CreateToken implements Maker.
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	// Membuat payload
	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", err
	}

	// membuat enkripsi
	// argument pertama merupakan key yang digunakan untuk enkripsi atau dekripsi
	// argument kedua merupakan data yang akan dienkripsi
	// parameter ke tiga merupakan opsi tambahan, tetapi sekarang ini tidak ada
	return maker.paseto.Encrypt([]byte(maker.symmetricKey), payload, nil)
}

// VerifyToken implements Maker.
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	// parameter pertama adalah token
	// parameter kedua adalah key untuk decription
	// parameter ketiga tipe struct yang akan menampung data deskripsi
	// parameter keempat itu opsional
	err := maker.paseto.Decrypt(token, []byte(maker.symmetricKey), payload, nil)

	if err != nil {
        return nil, ErrInvalidToken
    }

	// periksa atau validasi payload
	err = payload.Valid()

	if err != nil {
		return nil, err
	}

	return payload, nil
}

// Mengembalikan instance PasetoMaker dengan implementasi maker interface
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	// panjang adalah 32
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}
