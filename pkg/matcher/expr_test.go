package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleExpr(t *testing.T) {
	expr := &SimpleExpr{
		Includes: []string{
			"~ /api/",
			"~ .php$",
		},
		Excludes: []string{
			"~ /api/img",
		},
	}

	err := expr.Parse()
	assert.NoError(t, err)

	assert.True(t, expr.MatchString("/api/a.php"))
	assert.False(t, expr.MatchString("/api/a.php2"))
	assert.False(t, expr.MatchString("/api2/a.php"))
	assert.False(t, expr.MatchString("/api/img.php"))

	assert.True(t, expr.Match([]byte("/api/a.php")))
	assert.False(t, expr.Match([]byte("/api/a.php2")))
	assert.False(t, expr.Match([]byte("/api2/a.php")))
	assert.False(t, expr.Match([]byte("/api/img.php")))
}

func TestSimpleExpr_Parse_NG(t *testing.T) {
	{
		expr := &SimpleExpr{
			Includes: []string{
				"~ (ab",
				"~ .php$",
			},
		}

		err := expr.Parse()
		assert.Error(t, err)
	}
	{
		expr := &SimpleExpr{
			Excludes: []string{
				"~ (ab",
				"~ .php$",
			},
		}

		err := expr.Parse()
		assert.Error(t, err)
	}
}
