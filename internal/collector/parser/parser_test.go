package parser

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dmitrygulevich2000/mwormix-bot/internal/utils"
)

func Test_timeFromVKString(t *testing.T) {
	y, m, d := time.Now().In(utils.MSKLoc).Date()

	tests := []struct {
		vkString     string
		expectedTime time.Time
	}{
		{
			vkString:     "сегодня в 12:01",
			expectedTime: time.Date(y, m, d, 12, 1, 0, 0, utils.MSKLoc),
		},
		{
			vkString:     "вчера в 23:57",
			expectedTime: time.Date(y, m, d-1, 23, 57, 0, 0, utils.MSKLoc),
		},
		{
			vkString:     "1 янв в 16:38",
			expectedTime: time.Date(y, time.January, 1, 16, 38, 0, 0, utils.MSKLoc),
		},
		{
			vkString:     "31\u00A0дек\u00A02023",
			expectedTime: time.Date(2023, time.December, 31, 23, 59, 0, 0, utils.MSKLoc),
		},
	}

	for i, testCase := range tests {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			gotTime, err := timeFromVKString(testCase.vkString)
			require.Nil(t, err)
			require.Equal(t, testCase.expectedTime, gotTime)
		})
	}
}
