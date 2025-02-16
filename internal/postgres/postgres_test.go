package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	db := &db{}
	err := db.Connect()
	assert.NoError(t, err)
	assert.NotNil(t, db.connection)

	err = db.connection.Ping(context.Background())
	assert.NoError(t, err)

	bool := db.CheckConnection()
	assert.True(t, bool)

	err = db.Close()
	assert.NoError(t, err)
}

func TestCRUD(t *testing.T) {
	// TODO
}
