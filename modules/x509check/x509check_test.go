package x509check

import (
	"crypto/x509"
	"errors"
	"testing"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	job := New()

	assert.Implements(t, (*module.Module)(nil), job)
}

func TestX509Check_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestX509Check_Charts(t *testing.T) {
	job := New()

	assert.NotNil(t, job.Charts())
}

func TestX509Check_Init(t *testing.T) {
	job := New()
	job.Source = "https://example.org"

	assert.True(t, job.Init())
}

func TestX509Check_InitErrorOnCreatingGathererWrongTLSCA(t *testing.T) {
	job := New()
	job.Source = "https://example.org"
	job.ClientTLSConfig.TLSCA = "testdata/tls"

	assert.False(t, job.Init())
}

func TestX509Check_Check(t *testing.T) {
	job := New()
	job.gatherer = &mockGatherer{certs: []*x509.Certificate{{}}}
	assert.True(t, job.Check())
}

func TestX509Check_CheckError(t *testing.T) {
	job := New()
	job.gatherer = &mockGatherer{retErr: true}
	assert.False(t, job.Check())
}

func TestX509Check_Collect(t *testing.T) {
	job := New()
	job.gatherer = &mockGatherer{certs: []*x509.Certificate{{}}}
	mx := job.Collect()

	assert.NotZero(t, mx)
	v, ok := mx["expiry"]
	assert.True(t, ok)
	assert.NotZero(t, v)
}

func TestX509Check_CollectErrorOnGathering(t *testing.T) {
	job := New()
	job.gatherer = &mockGatherer{retErr: true}
	mx := job.Collect()

	assert.Nil(t, mx)
}

func TestX509Check_CollectZeroCertificates(t *testing.T) {
	job := New()
	job.gatherer = &mockGatherer{certs: []*x509.Certificate{}}
	mx := job.Collect()

	assert.Nil(t, mx)
}

type mockGatherer struct {
	certs  []*x509.Certificate
	retErr bool
}

func (m mockGatherer) Gather() ([]*x509.Certificate, error) {
	if m.retErr {
		return nil, errors.New("mock error")
	}
	return m.certs, nil
}
