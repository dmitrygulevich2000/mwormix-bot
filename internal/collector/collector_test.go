package collector

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dmitrygulevich2000/mwormix-bot/internal/storage"
	"github.com/dmitrygulevich2000/mwormix-bot/internal/utils"
)

func Test_BonusCollector_collectBonusesAfter(t *testing.T) {
	require := require.New(t)

	now := time.Now().In(utils.MSKLoc)
	y, m, d := now.Year(), now.Month(), now.Day()

	server := utils.RunFileServer("testdata.html")
	bonusesExpected := []Bonus{
		Bonus{
			PublishedAt: time.Date(y, m, d, 12, 0, 0, 0, utils.MSKLoc),
			Link:        "https://vk.com/wall-13842325_1689974",
			Code:        "iblcoatu",
			ValidUntil:  "до 23:00 по Московскому времени!",
		},
		Bonus{
			PublishedAt: time.Date(y, m, d, 1, 0, 0, 0, utils.MSKLoc),
			Link:        "https://vk.com/wall-13842325_1689973",
			Code:        "bxlcgizr",
			ValidUntil:  "до 04:00 по Московскому времени!",
		},
		Bonus{
			PublishedAt: time.Date(y, m, d-1, 9, 0, 0, 0, utils.MSKLoc),
			Link:        "https://vk.com/wall-13842325_1689968",
			Code:        "o94xyvyo",
			ValidUntil:  "до 23:00 по Московскому времени!",
		},
		Bonus{
			PublishedAt: time.Date(2024, time.February, 10, 9, 0, 0, 0, utils.MSKLoc),
			Link:        "https://vk.com/wall-13842325_1689956",
			Code:        "69jrnwlo",
			ValidUntil:  "до 23:00 по Московскому времени!",
		},
	}

	collector := New(nil, server.URL)

	bonusesGot, err := collector.collectBonusesAfter(time.Date(2024, time.February, 9, 21, 10, 34, 0, time.UTC))
	require.Nil(err)
	require.Equal(bonusesExpected, bonusesGot)

	// need return bonus published at the same minute as "after" ts
	bonusesGotWithNearSearch, err := collector.collectBonusesAfter(time.Date(y, m, d, 12, 0, 28, 0, utils.MSKLoc))
	require.Nil(err)
	require.Equal(bonusesExpected[:1], bonusesGotWithNearSearch)

	bonusesGotWithNearSearch, err = collector.collectBonusesAfter(time.Date(y, m, d, 12, 1, 0, 0, utils.MSKLoc))
	require.Nil(err)
	require.Equal([]Bonus{}, bonusesGotWithNearSearch)
}

func Test_BonusCollector_CollectBonuses(t *testing.T) {
	require := require.New(t)

	now := time.Now().In(utils.MSKLoc)
	y, m, d := now.Year(), now.Month(), now.Day()

	searches, _ := storage.NewSearches("tmpdb")
	store := storage.New(searches)
	defer os.Remove("tmpdb")

	server := utils.RunFileServer("testdata.html")
	bonusesExpected := []Bonus{
		Bonus{
			PublishedAt: time.Date(y, m, d, 12, 0, 0, 0, utils.MSKLoc),
			Link:        "https://vk.com/wall-13842325_1689974",
			Code:        "iblcoatu",
			ValidUntil:  "до 23:00 по Московскому времени!",
		},
		Bonus{
			PublishedAt: time.Date(y, m, d, 1, 0, 0, 0, utils.MSKLoc),
			Link:        "https://vk.com/wall-13842325_1689973",
			Code:        "bxlcgizr",
			ValidUntil:  "до 04:00 по Московскому времени!",
		},
		Bonus{
			PublishedAt: time.Date(y, m, d-1, 9, 0, 0, 0, utils.MSKLoc),
			Link:        "https://vk.com/wall-13842325_1689968",
			Code:        "o94xyvyo",
			ValidUntil:  "до 23:00 по Московскому времени!",
		},
		Bonus{
			PublishedAt: time.Date(2024, time.February, 10, 9, 0, 0, 0, utils.MSKLoc),
			Link:        "https://vk.com/wall-13842325_1689956",
			Code:        "69jrnwlo",
			ValidUntil:  "до 23:00 по Московскому времени!",
		},
		Bonus{
			PublishedAt: time.Date(2024, time.February, 9, 12, 0, 0, 0, utils.MSKLoc),
			Link:        "https://vk.com/wall-13842325_1689945",
			Code:        "i8e90rab",
			ValidUntil:  "до 23:00 по Московскому времени!",
		},
	}

	collector := New(store, server.URL)

	bonusesGot, err := collector.CollectNewBonuses()
	require.Nil(err)
	require.Equal(bonusesExpected, bonusesGot)

	bonusesGot, err = collector.CollectNewBonuses()
	require.Nil(err)
	// may fail if run at 12:00 Moscow
	require.Equal([]Bonus{}, bonusesGot)
}
