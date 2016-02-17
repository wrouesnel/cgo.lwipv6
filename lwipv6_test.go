package lwipv6

import (
	"testing"
	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/assert"
)

func TestInitFini(t *testing.T) {
	assert := assert.New(t)

	assert.False(IsInitialized(), "LWIP stack not initialized.")
	Initialize(0)
	defer Finish()
	assert.True(IsInitialized(), "LWIP stack initialized.")


}