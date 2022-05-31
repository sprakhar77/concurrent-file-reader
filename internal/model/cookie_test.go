package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ToCookie_Success(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output Cookie
	}{
		{
			name:   "valid test case",
			input:  "5UAVanZf6UtGyKVS,2018-12-09T07:25:00+00:00",
			output: Cookie{Name: "5UAVanZf6UtGyKVS", Date: "2018-12-09"},
		},
		{
			name:   "valid test case with cases",
			input:  "    5UAVanZf6UtGyKVS,    2018-12-09T07:25:00+00:00    ",
			output: Cookie{Name: "5UAVanZf6UtGyKVS", Date: "2018-12-09"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := ToCookie(tt.input)
			require.NoError(t, err)
			require.EqualValues(t, tt.output, c)
		})
	}
}

func Test_ToCookie_Fail(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty input",
			input: "",
		},
		{
			name:  "multiple delimiters",
			input: "5UAVanZf6UtGyKVS,,2018-12-09T07:25:00+00:00",
		},
		{
			name:  "broken date format",
			input: "5UAVanZf6UtGyKVS,2018-12-09T07:25:00+00:0as0",
		},
		{
			name:  "wrong date format",
			input: "5UAVanZf6UtGyKVS,Jan-02-06",
		},
		{
			name:  "missing date",
			input: "5UAVanZf6UtGyKVS",
		},
		{
			name:  "missing cookie",
			input: "2018-12-09T07:25:00+00:0as0",
		},
		{
			name:  "more than two inputs",
			input: "5UAVanZf6UtGyKVS,2018-12-09T07:25:00+00:00,something",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ToCookie(tt.input)
			require.Error(t, err)
		})
	}
}
