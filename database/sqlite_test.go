package database_test

import (
	"os"
	"testing"

	"github.com/l3lackShark/winservice/database"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	ex, err := os.Executable()
	if ok := assert.NoError(t, err, "Failed to get executable path"); !ok {
		return
	}

	t.Logf("Creating database at path %s\n", ex)

	_, err = database.New(ex)
	assert.NoError(t, err, "Failed to get all processes")
}
