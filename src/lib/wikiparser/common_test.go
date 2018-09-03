package wikiparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	p := New([]byte(""))

	assert.NotNil(p)
}

func TestWikiparser_Parse(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		Data []byte
		Expected Expression
	}

	testCases := []testCase{
		{
			Data: []byte("{{ test }}"),
			Expected: Expression{
				Action: "test",
				Args:   nil,
			},
		},
		{
			Data: []byte("{{ test(test,1234) }}"),
			Expected: Expression{
				Action: "test",
				Args: []string{
					"test",
					"1234",
				},
			},
		},
	}

	for _, tc := range testCases {
		wp := New(tc.Data)

		wp.Parse()
		assert.True(len(wp.Expressions) == 1)
		assert.Equal(tc.Expected.Action, wp.Expressions[0].Action)
		assert.Equal(tc.Expected.Args, wp.Expressions[0].Args)
	}
}
