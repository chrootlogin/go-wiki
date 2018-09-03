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
		Expected []Expression
	}

	testCases := []testCase{
		{
			Data: []byte("{"),
			Expected: []Expression{},
		},
		{
			Data: []byte("{blaa"),
			Expected: []Expression{},
		},
		{
			Data: []byte("{{ blaa"),
			Expected: []Expression{},
		},
		{
			Data: []byte("{{ blaa( }}"),
			Expected: []Expression{},
		},
		{
			Data: []byte("{{ blaa(sdgjkw, }}"),
			Expected: []Expression{},
		},
		{
			Data: []byte("{{ test }}"),
			Expected: []Expression{
				{
					Action: "test",
					Args:   nil,
				},
			},
		},
		{
			Data: []byte("{{ test(test,1234) }}"),
			Expected: []Expression{
				{
					Action: "test",
					Args: []string{
						"test",
						"1234",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		wp := New(tc.Data)

		wp.Parse()

		assert.True(len(tc.Expected) == len(wp.Expressions))
		for i, _ := range tc.Expected {
			assert.Equal(tc.Expected[i].Action, wp.Expressions[i].Action)
			assert.Equal(tc.Expected[i].Args, wp.Expressions[i].Args)
		}
	}
}
