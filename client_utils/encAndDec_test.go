package client_utils

import (
	"fmt"
	"math/rand"
	"myproject/globals"
	"testing"

	"github.com/tuneinsight/lattigo/v6/core/rlwe"
)

func TestEncAndDec(t *testing.T) {
	fmt.Println(globals.Params.MaxSlots())
	fmt.Println(globals.Params.MaxLevel())

	sk := rlwe.NewSecretKey(globals.Params)
	// Vector of plaintext values
	want := make([]float64, globals.Params.MaxSlots())

	// Source for sampling random plaintext values (not cryptographically secure)
	/* #nosec G404 */
	r := rand.New(rand.NewSource(0))

	// Populates the vector of plaintext values
	for i := range want {
		want[i] = 2*r.Float64() - 1 // uniform in [-1, 1]
	}

	myenc := &encryptor{}
	ct := myenc.encrypt(sk, want)
	mydec := &decryptor{}
	have := mydec.decrypt(sk, ct)

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
