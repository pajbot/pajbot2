package filters

/*
func benchmarkBanphrase(banphrases []Pajbot1Banphrase, text string, b *testing.B) {
	originalVariations, lowercaseVariations, _ := utils.MakeVariations(text, true)

	for n := 0; n < b.N; n++ {
		CheckBanphrases(banphrases, originalVariations, lowercaseVariations)
	}
}

func BenchmarkBanphrase1(b *testing.B) {
	banphrases := []Pajbot1Banphrase{
		Pajbot1Banphrase{
			Phrase:        "a",
			Operator:      OperatorContains,
			RemoveAccents: true,
		},
	}

	benchmarkBanphrase(banphrases, "short", b)
}

func BenchmarkBanphrase500ExpensiveNonMatchingShort(b *testing.B) {
	banphrases := []Pajbot1Banphrase{}
	for n := 0; n < 500; n++ {
		banphrases = append(banphrases, Pajbot1Banphrase{
			Phrase:        "a",
			Operator:      OperatorContains,
			RemoveAccents: true,
		})
	}

	benchmarkBanphrase(banphrases, "short", b)
}

func BenchmarkBanphrase500CheapNonMatchingShort(b *testing.B) {
	banphrases := []Pajbot1Banphrase{}
	for n := 0; n < 500; n++ {
		banphrases = append(banphrases, Pajbot1Banphrase{
			Phrase:        "a",
			Operator:      OperatorContains,
			RemoveAccents: false,
		})
	}

	benchmarkBanphrase(banphrases, "short", b)
}

func BenchmarkBanphrase500ExpensiveNonMatchingLong(b *testing.B) {
	banphrases := []Pajbot1Banphrase{}
	for n := 0; n < 500; n++ {
		banphrases = append(banphrases, Pajbot1Banphrase{
			Phrase:        "a",
			Operator:      OperatorContains,
			RemoveAccents: true,
		})
	}

	benchmarkBanphrase(banphrases, "ooö ooOoooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo", b)
}

func BenchmarkBanphrase500CheapNonMatchingLong(b *testing.B) {
	banphrases := []Pajbot1Banphrase{}
	for n := 0; n < 500; n++ {
		banphrases = append(banphrases, Pajbot1Banphrase{
			Phrase:        "a",
			Operator:      OperatorContains,
			RemoveAccents: false,
		})
	}

	benchmarkBanphrase(banphrases, "oooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo", b)
}
*/
