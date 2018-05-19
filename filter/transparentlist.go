package filter

import (
	"errors"

	"github.com/anknown/ahocorasick"
)

type transparentListRange struct {
	Pos        int
	SkipLength int
}

type TransparentListSkipRange struct {
	skips []transparentListRange
}

func (t *TransparentListSkipRange) addSkip(r transparentListRange) {
	t.skips = append(t.skips, r)
}

func (t *TransparentListSkipRange) ShouldSkip(index int) int {
	if len(t.skips) == 0 {
		return 0
	}

	skip := t.skips[0]

	if skip.Pos == index {
		ret := skip.SkipLength

		t.skips = t.skips[1:]

		return ret
	}

	return 0
}

type TransparentList struct {
	m *goahocorasick.Machine

	dict [][]rune
}

func NewTransparentList() *TransparentList {
	t := &TransparentList{}

	t.m = new(goahocorasick.Machine)

	return t
}

func (t *TransparentList) Add(s string) {
	t.dict = append(t.dict, []rune(s))
}

func (t *TransparentList) Build() error {
	if t.m == nil {
		return errors.New("Transparent list not initialized properly")
	}

	return t.m.Build(t.dict)
}

func (t *TransparentList) Find(text []rune) (ret TransparentListSkipRange) {
	terms := t.m.MultiPatternSearch(text, false)

	for _, t := range terms {
		ret.addSkip(transparentListRange{
			Pos:        t.Pos,
			SkipLength: len(t.Word),
		})
	}

	return
}
