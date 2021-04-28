package unity

import (
	"testing"

	"github.com/stretchr/testify/assert"

  //"io/ioutil"
  //"net/http"
  //"net/http/httptest"

  "github.com/netdata/go-orchestrator/module"
  //"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
  job := New()

  assert.Implements(t, (*module.Module)(nil), job)
}

func TestUnity_Init(t *testing.T) {
	mod := New()

	assert.True(t, mod.Init())
}

func TestExample_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestExample_Cleanup(t *testing.T) {
	New().Cleanup()
}
