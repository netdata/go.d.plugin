// SPDX-License-Identifier: GPL-3.0-or-later

package netdataapi

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPI_CHART(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := API{Writer: b}

	_ = netdataAPI.CHART(
		"",
		"id",
		"name",
		"title",
		"units",
		"family",
		"context",
		"line",
		1,
		1,
		"",
		"orchestrator",
		"module",
	)

	assert.Equal(
		t,
		"CHART '.id' 'name' 'title' 'units' 'family' 'context' 'line' '1' '1' '' 'orchestrator' 'module'\n",
		b.String(),
	)
}

func TestAPI_DIMENSION(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := API{Writer: b}

	_ = netdataAPI.DIMENSION(
		"id",
		"name",
		"absolute",
		1,
		1,
		"",
	)

	assert.Equal(
		t,
		"DIMENSION 'id' 'name' 'absolute' '1' '1' ''\n",
		b.String(),
	)
}

func TestAPI_BEGIN(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := API{Writer: b}

	_ = netdataAPI.BEGIN(
		"typeID",
		"id",
		0,
	)

	assert.Equal(
		t,
		"BEGIN 'typeID.id'\n",
		b.String(),
	)

	b.Reset()

	_ = netdataAPI.BEGIN(
		"typeID",
		"id",
		1,
	)

	assert.Equal(
		t,
		"BEGIN 'typeID.id' 1\n",
		b.String(),
	)
}

func TestAPI_SET(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := API{Writer: b}

	_ = netdataAPI.SET("id", 100)

	assert.Equal(
		t,
		"SET 'id' = 100\n",
		b.String(),
	)
}

func TestAPI_SETEMPTY(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := API{Writer: b}

	_ = netdataAPI.SETEMPTY("id")

	assert.Equal(
		t,
		"SET 'id' = \n",
		b.String(),
	)
}

func TestAPI_VARIABLE(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := API{Writer: b}

	_ = netdataAPI.VARIABLE("id", 100)

	assert.Equal(
		t,
		"VARIABLE CHART 'id' = 100\n",
		b.String(),
	)
}

func TestAPI_END(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := API{Writer: b}

	_ = netdataAPI.END()

	assert.Equal(
		t,
		"END\n\n",
		b.String(),
	)
}

func TestAPI_FLUSH(t *testing.T) {
	b := &bytes.Buffer{}
	netdataAPI := API{Writer: b}

	_ = netdataAPI.FLUSH()

	assert.Equal(
		t,
		"FLUSH\n",
		b.String(),
	)
}
