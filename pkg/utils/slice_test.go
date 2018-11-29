package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringSlice_Append(t *testing.T) {
	var s StringSlice
	s.Append("0")

	assert.Len(t, s, 1)
}

func TestStringSlice_DeleteByID(t *testing.T) {
	s := StringSlice{
		"0",
		"1",
	}

	assert.False(t, s.DeleteByID("2"))
	assert.True(t, s.DeleteByID("1"))
	assert.Len(t, s, 1)
	assert.Equal(t, "0", s[0])
}

func TestStringSlice_DeleteByIndex(t *testing.T) {
	s := StringSlice{
		"0",
		"1",
	}

	assert.False(t, s.DeleteByIndex(3))
	assert.True(t, s.DeleteByIndex(1))
	assert.Len(t, s, 1)
	assert.Equal(t, "0", s[0])

}

func TestStringSlice_InsertAfterID(t *testing.T) {
	s := StringSlice{
		"0",
		"1",
		"2",
	}

	require.False(t, s.InsertAfterID("99", "3", "4", "5"))

	require.Equal(
		t,
		StringSlice{
			"0",
			"1",
			"2",
		},
		s,
	)

	require.True(t, s.InsertAfterID("0", "11", "12"))

	require.Equal(
		t,
		StringSlice{
			"0",
			"11",
			"12",
			"1",
			"2",
		},
		s,
	)

	require.True(t, s.InsertAfterID("2", "11", "12"))

	require.Equal(
		t,
		StringSlice{
			"0",
			"11",
			"12",
			"1",
			"2",
			"11",
			"12",
		},
		s,
	)
}

func TestStringSlice_InsertBeforeID(t *testing.T) {
	s := StringSlice{
		"0",
		"1",
		"2",
	}

	require.False(t, s.InsertBeforeID("99", "3", "4", "5"))

	require.Equal(
		t,
		StringSlice{
			"0",
			"1",
			"2",
		},
		s,
	)

	require.True(t, s.InsertBeforeID("0", "11", "12"))

	require.Equal(
		t,
		StringSlice{
			"11",
			"12",
			"0",
			"1",
			"2",
		},
		s,
	)

	require.True(t, s.InsertAfterID("2", "11", "12"))

	require.Equal(
		t,
		StringSlice{
			"11",
			"12",
			"0",
			"1",
			"2",
			"11",
			"12",
		},
		s,
	)
}

func TestStringSlice_Include(t *testing.T) {
	s := StringSlice{
		"0",
		"1",
		"2",
	}

	assert.True(t, s.Include("0"))
	assert.False(t, s.Include("3"))
}

func TestStringSlice_Index(t *testing.T) {
	s := StringSlice{
		"0",
		"1",
		"2",
	}
	assert.Equal(t, 0, s.Index("0"))
	assert.Equal(t, 2, s.Index("2"))
}

func TestStringSlice_Insert(t *testing.T) {
	s := StringSlice{
		"0",
		"1",
		"2",
	}

	require.False(t, s.Insert(99, "-1"))

	require.Equal(
		t,
		StringSlice{
			"0",
			"1",
			"2",
		},
		s,
	)

	require.True(t, s.Insert(0, "11", "12"))

	require.Equal(
		t,
		StringSlice{
			"11",
			"12",
			"0",
			"1",
			"2",
		},
		s,
	)

	assert.True(t, s.Insert(1, "13", "14"))

	assert.Equal(
		t,
		StringSlice{
			"11",
			"13",
			"14",
			"12",
			"0",
			"1",
			"2",
		},
		s,
	)
}
