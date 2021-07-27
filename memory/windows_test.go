package memory_test

import (
	"testing"

	"github.com/l3lackShark/winservice/memory"
	"github.com/stretchr/testify/assert"
)

func TestGetAllProcesses(t *testing.T) {
	memoryApi := memory.New()
	_, err := memoryApi.GetAllProcesses()
	assert.NoError(t, err, "Failed to get all processes")

}
