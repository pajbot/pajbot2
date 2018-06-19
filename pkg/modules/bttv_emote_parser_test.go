package modules

import (
	"os"
	"testing"

	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var parser *BTTVEmoteParser
var globalEmotes map[string]common.Emote

type MyMockedMessage struct {
	mock.Mock

	text       string
	bttvEmotes []*common.Emote
}

func (m MyMockedMessage) GetText() string {
	return m.text
}

func (m MyMockedMessage) GetTwitchReader() pkg.EmoteReader {
	return nil
}

func (m MyMockedMessage) GetBTTVReader() pkg.EmoteReader {
	return nil
}

func (m *MyMockedMessage) AddBTTVEmote(emote pkg.Emote) {
	m.bttvEmotes = append(m.bttvEmotes, emote.(*common.Emote))
}

func (m MyMockedMessage) numEmotes() int {
	count := 0
	for _, emote := range m.bttvEmotes {
		count += emote.GetCount()
	}

	return count
}

func TestMain(m *testing.M) {
	globalEmotes = make(map[string]common.Emote)
	globalEmotes["NaM"] = common.Emote{Name: "NaM", Count: 1}
	globalEmotes["KKona"] = common.Emote{Name: "KKona", Count: 1}

	parser = NewBTTVEmoteParser(&globalEmotes)
	os.Exit(m.Run())
}

func test(t *testing.T, text string, cb func(msg *MyMockedMessage)) {
	msg := MyMockedMessage{
		text: text,
	}

	parser.OnMessage(nil, nil, &msg, nil)
	cb(&msg)
}

func testCount(t *testing.T, text string, count int) {
	assert := assert.New(t)

	test(t, text, func(msg *MyMockedMessage) {
		assert.Equal(count, msg.numEmotes())
	})
}

func TestParseEmotes(t *testing.T) {
	testCount(t, "test", 0)

	testCount(t, "KKona", 1)

	testCount(t, "KKona KKona", 2)

	testCount(t, "🚽NaM🚽", 1)

	testCount(t, "🚽NaM🚽NaM🚽NaM🚽NaM🚽NaM🚽", 5)

	testCount(t, "\xF0\x9F\x92\xBCNaM!", 0)
}
