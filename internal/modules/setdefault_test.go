package modules

import (
	"testing"
)

var (
	n = getFileName(1)
	i = 99
)

func TestModuleDefault_SetUpdateEvery(t *testing.T) {
	SetDefault().SetUpdateEvery(i)
	if moduleDefaults[n].u == nil {
		t.Fatal("expected not nil, value is not set")
	}
	if v := *moduleDefaults[n].u; v != i {
		t.Errorf("expected %d, but got %d", v, i)
	}
}

func TestModuleDefault_SetChartsCleanup(t *testing.T) {
	SetDefault().SetChartsCleanup(i)
	if moduleDefaults[n].c == nil {
		t.Fatal("expected not nil, value is not set")
	}
	if v := *moduleDefaults[n].c; v != i {
		t.Errorf("expected %d, but got %d", v, i)
	}
}

func TestModuleDefault_SetDisabledByDefault(t *testing.T) {
	SetDefault().SetDisabledByDefault()

	if moduleDefaults[n].d == false {
		t.Error("expected true, value is not set")
	}
}

func TestModuleDefault_GetUpdateEvery(t *testing.T) {
	if v, _ := GetDefault(n).UpdateEvery(); v != i {
		t.Errorf("expected %d, but got %d", v, i)
	}
}

func TestModuleDefault_GetChartsCleanup(t *testing.T) {
	if v, _ := GetDefault(n).ChartsCleanup(); v != i {
		t.Errorf("expected %d, but got %d", v, i)
	}
}

func TestModuleDefault_GetDisabledByDefault(t *testing.T) {
	if GetDefault(n).DisabledByDefault() == false {
		t.Error("expected true, value is not set")
	}
}

func TestGetDefault(t *testing.T) {
	var g G
	g = GetDefault("_")
	_ = g
}

func TestSetDefault(t *testing.T) {
	var s S
	s = SetDefault()
	_ = s
}
