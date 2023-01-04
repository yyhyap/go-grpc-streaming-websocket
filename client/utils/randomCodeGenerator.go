package utils

import (
	"math/rand"
	"strings"
	"sync"
)

var (
	randomCodeGenerator *randomCodeGeneratorStruct

	randomCodeGeneratorOnce sync.Once

	digits = []rune("0123456789")
)

type IRandomCodeGenerator interface {
	GenerateRandomDigits(length int) string
}

type randomCodeGeneratorStruct struct{}

func GetRandomCodeGenerator() *randomCodeGeneratorStruct {
	if randomCodeGenerator == nil {
		randomCodeGeneratorOnce.Do(func() {
			randomCodeGenerator = &randomCodeGeneratorStruct{}
		})
	}
	return randomCodeGenerator
}

func (r *randomCodeGeneratorStruct) GenerateRandomDigits(length int) string {
	digitSize := len(digits)
	var sb strings.Builder

	for i := 0; i < length; i++ {
		char := digits[rand.Intn(digitSize)]
		sb.WriteRune(char)
	}

	s := sb.String()
	return s
}
