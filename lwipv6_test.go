package lwipv6

import (
	"testing"
	//"github.com/stretchr/testify/assert"
)

func TestInitFini(t *testing.T) {
	lwip_init()
	lwip_fini()
}