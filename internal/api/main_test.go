package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/db/sqlc"
	"github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/util"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, store db.MockStore) *Server {
	config := util.Config{
		PasetoSymmetricKey:        util.RandomString(32), // alasan 32 karena PasetoSymmetricKey itu 32 byte,
		PasetoAccessTokenDuration: time.Minute,
	}

	// dapatkan instance Server dengan akses NewServer
	server, err := NewServer(store, config)

	// Jika ada error, maka gagal
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
