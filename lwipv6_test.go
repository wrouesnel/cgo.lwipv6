package lwipv6

import (
	"testing"
	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/assert"
	"net"
	"fmt"
	"os/exec"
	"bytes"
	"bufio"
	"strings"
	"path/filepath"
)

func TestStartup(t *testing.T) {
	assert := assert.New(t)

	assert.False(IsInitialized(), "LWIP stack not initialized.")
	Initialize(0)
	defer Finish()
	assert.True(IsInitialized(), "LWIP stack initialized.")


}

// This function must always mimic the code in test/ip_conversions in order to work.
func TestIPConversions(t *testing.T) {
	assert := assert.New(t)

	cmd := exec.Command("make")
	cmd.Dir, _ = filepath.Abs("tests")
	if err := cmd.Run(); err != nil {
		assert.Fail("Failed to build test fixtures")
	}

	bufio.New
	tcmd := exec.Command("tests/ip_conversions")
	tcmd.Stdout = buf.N\
	if err := tcmd.Run(); err != nil {
		assert.Fail("Failed to execute test fixture", "tests/ip_conversions")
	}

	// Read all lines of output, and use them as data for our tests
	lineReader := bufio.NewScanner(buf)
	for lineReader.Scan() {
		vals := strings.Fields(lineReader.Text())

		CIDR := vals[0]

		addr0 := vals[1]
		addr1 := vals[2]
		addr2 := vals[3]
		addr3 := vals[4]

		mask0 := vals[5]
		mask1 := vals[6]
		mask2 := vals[7]
		mask3 := vals[8]

		_, net, _ := net.ParseCIDR(CIDR)
		lwipIP, lwipMask := convert_IPNet_to_LWIP(net)

		assert.Equal(addr0, fmt.Sprintf("%v", lwipIP.addr[0]), "Address integer matches")
		assert.Equal(addr1, fmt.Sprintf("%v", lwipIP.addr[1]), "Address integer matches")
		assert.Equal(addr2, fmt.Sprintf("%v", lwipIP.addr[2]), "Address integer matches")
		assert.Equal(addr3, fmt.Sprintf("%v", lwipIP.addr[3]), "Address integer matches")

		assert.Equal(mask0, fmt.Sprintf("%v", lwipIP.addr[0]), "Netmask integer matches")
		assert.Equal(mask1, fmt.Sprintf("%v", lwipIP.addr[1]), "Netmask integer matches")
		assert.Equal(mask2, fmt.Sprintf("%v", lwipIP.addr[2]), "Netmask integer matches")
		assert.Equal(mask3, fmt.Sprintf("%v", lwipIP.addr[3]), "Netmask integer matches")

	}
}