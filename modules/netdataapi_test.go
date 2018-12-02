package modules

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetdataAPI_chart(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := apiWriter{Writer: b}

	netdataAPI.chart(
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

	assert.Equal(
		t,
		"CHART '.id' 'name' 'title' 'units' 'family' 'context' 'line' '1' '1' '' go.d 'module'\n",
		b.String(),
	)
}

func TestNetdataAPI_dimension(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := apiWriter{Writer: b}

	netdataAPI.dimension(
		"id",
		"name",
		Absolute,
		1,
		1,
		false,
	)

	assert.Equal(
		t,
		"DIMENSION 'id' 'name' 'absolute' '1' '1' ''\n",
		b.String(),
	)
}

func TestNetdataAPI_begin(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := apiWriter{Writer: b}

	netdataAPI.begin(
		"typeID",
		"id",
		0,
	)

	assert.Equal(
		t,
		"BEGIN typeID.id\n",
		b.String(),
	)

	b.Reset()

	netdataAPI.begin(
		"typeID",
		"id",
		1,
	)

	assert.Equal(
		t,
		"BEGIN typeID.id 1\n",
		b.String(),
	)
}

func TestNetdataAPI_set(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := apiWriter{Writer: b}

	netdataAPI.set("id", 100)

	assert.Equal(
		t,
		"SET id = 100\n",
		b.String(),
	)
}

func TestNetdataAPI_setEmpty(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := apiWriter{Writer: b}

	netdataAPI.setEmpty("id")

	assert.Equal(
		t,
		"SET id = \n",
		b.String(),
	)
}

func TestNetdataAPI_end(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := apiWriter{Writer: b}

	netdataAPI.end()

	assert.Equal(
		t,
		"END\n\n",
		b.String(),
	)
}

func TestNetdataAPI_flush(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := apiWriter{Writer: b}

	netdataAPI.flush()

	assert.Equal(
		t,
		"FLUSH\n",
		b.String(),
	)
}
