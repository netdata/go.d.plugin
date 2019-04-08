package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleExpr_none(t *testing.T) {
	expr := &SimpleExpr{}

	err := expr.Parse()
	assert.NoError(t, err)

	assert.True(t, expr.MatchString("/api/a.php"))
	assert.True(t, expr.MatchString("/api/a.php2"))
	assert.True(t, expr.MatchString("/api2/a.php"))
	assert.True(t, expr.MatchString("/api/img.php"))
	assert.True(t, expr.MatchString("/api2/img.php2"))

	assert.True(t, expr.Match([]byte("/api/a.php")))
	assert.True(t, expr.Match([]byte("/api/a.php2")))
	assert.True(t, expr.Match([]byte("/api2/a.php")))
	assert.True(t, expr.Match([]byte("/api/img.php")))
}

func TestSimpleExpr_include(t *testing.T) {
	expr := &SimpleExpr{
		Includes: []string{
			"~ /api/",
			"~ .php$",
		},
	}

	err := expr.Parse()
	assert.NoError(t, err)

	assert.True(t, expr.MatchString("/api/a.php"))
	assert.True(t, expr.MatchString("/api/a.php2"))
	assert.True(t, expr.MatchString("/api2/a.php"))
	assert.True(t, expr.MatchString("/api/img.php"))
	assert.False(t, expr.MatchString("/api2/img.php2"))

	assert.True(t, expr.Match([]byte("/api/a.php")))
	assert.True(t, expr.Match([]byte("/api/a.php2")))
	assert.True(t, expr.Match([]byte("/api2/a.php")))
	assert.True(t, expr.Match([]byte("/api/img.php")))
}

func TestSimpleExpr_exclude(t *testing.T) {
	expr := &SimpleExpr{
		Excludes: []string{
			"~ /api/img",
		},
	}

	err := expr.Parse()
	assert.NoError(t, err)

	assert.True(t, expr.MatchString("/api/a.php"))
	assert.True(t, expr.MatchString("/api/a.php2"))
	assert.True(t, expr.MatchString("/api2/a.php"))
	assert.False(t, expr.MatchString("/api/img.php"))
	assert.True(t, expr.MatchString("/api2/img.php2"))

	assert.True(t, expr.Match([]byte("/api/a.php")))
	assert.True(t, expr.Match([]byte("/api/a.php2")))
	assert.True(t, expr.Match([]byte("/api2/a.php")))
	assert.False(t, expr.Match([]byte("/api/img.php")))
}

func TestSimpleExpr_both(t *testing.T) {
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
	assert.True(t, expr.MatchString("/api/a.php2"))
	assert.True(t, expr.MatchString("/api2/a.php"))
	assert.False(t, expr.MatchString("/api/img.php"))
	assert.False(t, expr.MatchString("/api2/img.php2"))

	assert.True(t, expr.Match([]byte("/api/a.php")))
	assert.True(t, expr.Match([]byte("/api/a.php2")))
	assert.True(t, expr.Match([]byte("/api2/a.php")))
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
