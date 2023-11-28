// pkg/hash/hash.go

package hash

import (
	"encoding/binary"
	"fmt"

	"github.com/cespare/xxhash/v2"
)

// GenXXHash generates xxHash32 and xxHash64 for a given string.
// It returns two strings representing these hashes in hexadecimal format.
func GenXXHash(name string) (string, string) {
	h32 := xxhash.Sum32([]byte(name))
	h64 := xxhash.Sum64([]byte(name))

	// Convert the hash values to hex strings
	hex32 := binary.BigEndian.Uint32(toByteArray(h32))
	hex64 := binary.BigEndian.Uint64(toByteArray(h64))

	return formatHex32(hex32), formatHex64(hex64)
}

// toByteArray converts uint32 and uint64 to a byte array.
func toByteArray(val interface{}) []byte {
	switch v := val.(type) {
	case uint32:
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, v)
		return b
	case uint64:
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, v)
		return b
	default:
		return nil
	}
}

// formatHex32 and formatHex64 format uint32 and uint64 values as hex strings.
func formatHex32(val uint32) string {
	return "0x" + formatHexString(8, val)
}

func formatHex64(val uint64) string {
	return "0x" + formatHexString(16, val)
}

// formatHexString formats a uint value to a hex string with a specific length.
func formatHexString(length int, val interface{}) string {
	switch v := val.(type) {
	case uint32:
		return fmt.Sprintf("%0[1]*x", length, v)
	case uint64:
		return fmt.Sprintf("%0[1]*x", length, v)
	default:
		return ""
	}
}
