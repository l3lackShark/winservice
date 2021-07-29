package memory_test

import (
	"testing"

	"github.com/l3lackShark/winservice/memory"
	"github.com/stretchr/testify/assert"
)

func TestGetAllProcessesAndComputeDiff(t *testing.T) {
	procs := make(map[memory.UniqueProcess]memory.Process)
	memoryApi := memory.New()
	_, _, err := memoryApi.GetAllProcessesAndComputeDiff(procs)
	assert.NoError(t, err, "Failed to get all processes")

}
