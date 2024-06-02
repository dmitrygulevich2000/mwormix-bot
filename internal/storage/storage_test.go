package storage

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const dbname = "tmpdb"

func Test_SearchesStorage(t *testing.T) {
	require := require.New(t)
	defer os.Remove(dbname)

	// create new database
	s, err := NewSearches(dbname)
	require.Nil(err)

	fromDb, err := s.LastSearchTime()
	require.Nil(err)
	require.True(fromDb.IsZero())

	ts := time.Now().Truncate(0) // to strip monotonic clock reading
	err = s.StoreSearch(ts)
	require.Nil(err)

	fromDb, err = s.LastSearchTime()
	require.Nil(err)
	require.Equal(ts, fromDb)
	s.Close()

	// open existing database
	s, err = NewSearches(dbname)
	require.Nil(err)

	fromDb, err = s.LastSearchTime()
	require.Nil(err)
	require.Equal(ts, fromDb)

	ts = time.Now().Truncate(0) // to strip monotonic clock reading
	err = s.StoreSearch(ts)
	require.Nil(err)

	fromDb, err = s.LastSearchTime()
	require.Nil(err)
	require.Equal(ts, fromDb)
	s.Close()
}
