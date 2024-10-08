package homomorphic

import (
	"fmt"
	"github.com/nihonge/front/globals"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"math/rand"
	"testing"
)

func TestEncAndDec(t *testing.T) {
	sk := rlwe.NewSecretKey(globals.Params)
	// Vector of plaintext values
	want := make([]float64, globals.Params.MaxSlots())

	// Source for sampling random plaintext values (not cryptographically secure)
	/* #nosec G404 */
	r := rand.New(rand.NewSource(0))

	// Populates the vector of plaintext values
	left, right := 0.0, 100.0
	for i := range want {
		want[i] = left + (right-left)*r.Float64() // uniform in [-1, 1]
	}

	myenc := &Encryptor{}
	ct := myenc.Encrypt(sk, want)
	mydec := &Decryptor{}
	have := mydec.Decrypt(sk, ct)

	fmt.Printf("Have: ")
	for i := 0; i < 4; i++ {
		fmt.Printf("%20.15f ", have[i])
	}
	fmt.Printf("...\n")

	fmt.Printf("Want: ")
	for i := 0; i < 4; i++ {
		fmt.Printf("%20.15f ", want[i])
	}
	fmt.Printf("...\n")

}
