package pulid

import (
	"math/bits"

	"github.com/oklog/ulid/v2"
)

// EncodeBase32 convert number from decimal to base32
// with digits from ULID
func EncodeBase32(u uint64) string {
	const b, m = uint64(32), uint(32) - 1
	var (
		a = [64 + 1]byte{}
		i = len(a)
		s = uint(bits.TrailingZeros(32)) & 7
	)

	for ; u >= b; u >>= s {
		i--
		a[i] = ulid.Encoding[uint(u)&m]
	}

	i--
	a[i] = ulid.Encoding[uint(u)]
	return string(a[i:])
}
