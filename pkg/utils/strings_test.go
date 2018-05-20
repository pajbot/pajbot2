package utils

import (
	"testing"
)

func benchmarkMakeVariation(text string, doNormalize bool, b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _, _ = MakeVariations(text, doNormalize)
	}
}

func BenchmarkMakeVariationShortString(b *testing.B) {
	benchmarkMakeVariation("test", false, b)
}

func BenchmarkMakeVariationShortStringNormalize(b *testing.B) {
	benchmarkMakeVariation("test", true, b)
}

func BenchmarkMakeVariationLongString(b *testing.B) {
	benchmarkMakeVariation("kjdfghksdfj ghjksdf ghsdgfjkhsdfk gjhsdgfk jhsdfgkj h", false, b)
}

func BenchmarkMakeVariationLongStringNormalize(b *testing.B) {
	benchmarkMakeVariation("kjdfghksdfj ghjksdf ghsdgfjkhsdfk gjhsdgfk jhsdfgkj h", true, b)
}
