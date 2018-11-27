package modules

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetdataAPI_chart(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := apiWriter{Writer: b}

	err := netdataAPI.chart(
		"",
		"id",
		"name",
		"title",
		"units",
		"family",
		"context",
		Line,
		1,
		1,
		Opts{},
		"module",
	)

	expected := "CHART '.id' 'name' 'title' 'units' 'family' 'context' 'line' '1' '1' '' go.d 'module'\n"
	assert.NoError(t, err)
	assert.Equal(
		t,
		expected,
		b.String(),
	)
}

func TestNetdataAPI_dimension(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := apiWriter{Writer: b}

	err := netdataAPI.dimension(
		"id",
		"name",
		Absolute,
		1,
		1,
		false,
	)

	expected := "DIMENSION 'id' 'name' 'absolute' '1' '1' ''\n"
	assert.NoError(t, err)
	assert.Equal(
		t,
		expected,
		b.String(),
	)
}

func TestNetdataAPI_begin(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := apiWriter{Writer: b}

	err := netdataAPI.begin(
		"typeID",
		"id",
		0,
	)

	expected := "BEGIN typeID.id\n"
	assert.NoError(t, err)
	assert.Equal(
		t,
		expected,
		b.String(),
	)

	b.Reset()

	err = netdataAPI.begin(
		"typeID",
		"id",
		1,
	)

	expected = "BEGIN typeID.id 1\n"
	assert.NoError(t, err)
	assert.Equal(
		t,
		expected,
		b.String(),
	)
}

func TestNetdataAPI_set(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := apiWriter{Writer: b}

	err := netdataAPI.set(
		"id",
		100,
	)

	expected := "SET id = 100\n"
	assert.NoError(t, err)
	assert.Equal(
		t,
		expected,
		b.String(),
	)
}

func TestNetdataAPI_setEmpty(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := apiWriter{Writer: b}

	err := netdataAPI.setEmpty(
		"id",
	)

	expected := "SET id = \n"
	assert.NoError(t, err)
	assert.Equal(
		t,
		expected,
		b.String(),
	)
}

func TestNetdataAPI_end(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := apiWriter{Writer: b}

	err := netdataAPI.end()

	expected := "END\n\n"
	assert.NoError(t, err)
	assert.Equal(
		t,
		expected,
		b.String(),
	)
}

func TestNetdataAPI_flush(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := apiWriter{Writer: b}

	err := netdataAPI.flush()
	expected := "FLUSH\n"
	assert.NoError(t, err)
	assert.Equal(
		t,
		expected,
		b.String(),
	)
}
