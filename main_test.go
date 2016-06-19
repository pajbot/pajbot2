package main

import (
	"testing"

	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/helper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	var configTests = []struct {
		inputPath string
		expectedC *common.Config
		expectedE bool
	}{
		{
			inputPath: "resources/testfiles/config1.json",
			expectedC: &common.Config{
				Pass:          "oauth:xD",
				Nick:          "twitch_username",
				BrokerHost:    helper.NewStringPtr("localhost:7353"),
				BrokerPass:    helper.NewStringPtr("test"),
				RedisHost:     "localhost:6379",
				SQLDSN:        "pajbot2:password@tcp(localhost:3306)/pajbot2_test",
				RedisPassword: "",
				RedisDatabase: -1,
				TLSKey:        "",
				TLSCert:       "",
				Channels:      []string{"pajlada", "nuuls", "forsenlol"},
				ToWeb:         (chan map[string]interface{})(nil),
				FromWeb:       (chan map[string]interface{})(nil)},
			expectedE: false,
		},
		{
			inputPath: "resources/testfiles/nonexistingconfigfile.json",
			expectedC: nil,
			expectedE: true,
		},
		{
			inputPath: "resources/testfiles/config2_invalidjson.json",
			expectedC: nil,
			expectedE: true,
		},
	}

	for _, tt := range configTests {
		res, err := LoadConfig(tt.inputPath)

		if tt.expectedE {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}

		assert.Equal(t, tt.expectedC, res)
	}
}
