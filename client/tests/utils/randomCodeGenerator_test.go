package tests_utils

import (
	"go-grpc-restaurant-client/utils"
	"testing"
)

var (
	randomCodeGenerator utils.IRandomCodeGenerator = utils.GetRandomCodeGenerator()
)

func TestGenerateRandomDigits(t *testing.T) {
	s := randomCodeGenerator.GenerateRandomDigits(6)

	if len(s) != 6 {
		t.Error("Length of generated random 6 digits is not 6")
	}
}
