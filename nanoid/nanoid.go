/*
Copyright 2022 jaevor.
License can be found in the LICENSE file.
Original reference: https://github.com/ai/nanoid
*/
package nanoid

import (
	crand "crypto/rand"
	"errors"
	"math"
	"math/bits"
	"sync"
	"unicode"
)

var _nanoid = customASCII("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 21)

func Id() string {
	return _nanoid()
}

// `A-Za-z0-9_-`.
// Using less memory with [64]byte{...} than []byte(...).
var standardAlphabet = [64]byte{
	'A', 'B', 'C', 'D', 'E',
	'F', 'G', 'H', 'I', 'J',
	'K', 'L', 'M', 'N', 'O',
	'P', 'Q', 'R', 'S', 'T',
	'U', 'V', 'W', 'X', 'Y',
	'Z', 'a', 'b', 'c', 'd',
	'e', 'f', 'g', 'h', 'i',
	'j', 'k', 'l', 'm', 'n',
	'o', 'p', 'q', 'r', 's',
	't', 'u', 'v', 'w', 'x',
	'y', 'z', '0', '1', '2',
	'3', '4', '5', '6', '7',
	'8', '9', '-', '_',
}

/*
Returns a new generator of standard Nano IDs.
📝 Recommended (canonic) length is 21.
🟡 Errors if length is not, or within 2-255.
🧿 Concurrency safe.
*/
func standard(length int) func() string {
	if length < 2 || length > 255 {
		panic(errors.New("length for ID is invalid (must be within 2-255)"))
	}

	// Multiplying to increase the 'buffer' so that .Read()
	// has to be called less, which is more efficient.
	// b holds the random crypto bytes.
	size := length * length * 7
	b := make([]byte, size)
	crand.Read(b)
	offset := 0

	// Since the standard alphabet is ASCII, we don't have to use runes.
	// ASCII max is 128, so byte will be perfect.
	id := make([]byte, length)

	var mu sync.Mutex

	return func() string {
		mu.Lock()
		defer mu.Unlock()

		// If all the bytes in the slice
		// have been used, refill.
		if offset == size {
			crand.Read(b)
			offset = 0
		}

		for i := 0; i < length; i++ {
			/*
				"It is incorrect to use bytes exceeding the alphabet size.
				The following mask reduces the random byte in the 0-255 value
				range to the 0-63 value range. Therefore, adding hacks such
				as empty string fallback or magic numbers is unneccessary because
				the bitmask trims bytes down to the alphabet size (64)."
			*/
			// Index using the offset.
			id[i] = standardAlphabet[b[i+offset]&63]
		}

		// Extend the offset.
		offset += length

		return string(id)
	}
}

/*
Returns a Nano ID generator which uses a custom ASCII alphabet.
Uses less memory than CustomUnicode by only supporting ASCII.
For unicode support use nanoid.CustomUnicode.
🟡 Errors if alphabet is not valid ASCII or if length is not, or within 2-255.
🧿 Concurrency safe.
*/
func customASCII(alphabet string, length int) func() string {
	if length < 2 || length > 255 {
		panic(errors.New("length for ID is invalid (must be within 2-255)"))
	}

	alphabetLen := len(alphabet)

	for i := 0; i < alphabetLen; i++ {
		if alphabet[i] > unicode.MaxASCII {
			panic(errors.New("not valid ascii"))
		}
	}

	ab := []byte(alphabet)

	x := uint32(alphabetLen) - 1
	clz := bits.LeadingZeros32(x | 1)
	mask := (2 << (31 - clz)) - 1
	step := int(math.Ceil((1.6 * float64(mask*length)) / float64(alphabetLen)))

	b := make([]byte, step)
	id := make([]byte, length)

	j, idx := 0, 0

	var mu sync.Mutex

	return func() string {
		mu.Lock()
		defer mu.Unlock()
		for {
			crand.Read(b)
			for i := 0; i < step; i++ {
				idx = int(b[i]) & mask
				if idx < alphabetLen {
					id[j] = ab[idx]
					j++
					if j == length {
						j = 0
						return string(id)
					}
				}
			}
		}
	}
}
