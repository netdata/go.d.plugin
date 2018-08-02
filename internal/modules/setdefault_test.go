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
	if moduleDefaults[n].updateEvery == nil {
		t.Fatal("expected not nil, value is not set")
	}
	if v := *moduleDefaults[n].updateEvery; v != i {
		t.Errorf("expected %d, but got %d", v, i)
	}
}

func TestModuleDefault_SetChartsCleanup(t *testing.T) {
	SetDefault().SetChartsCleanup(i)
	if moduleDefaults[n].chartsCleanup == nil {
		t.Fatal("expected not nil, value is not set")
	}
	if v := *moduleDefaults[n].chartsCleanup; v != i {
		t.Errorf("expected %d, but got %d", v, i)
	}
}

func TestModuleDefault_SetDisabledByDefault(t *testing.T) {
	SetDefault().SetDisabledByDefault()

	if moduleDefaults[n].disabledByDefault == false {
		t.Error("expected true, value is not set")
	}
}

func TestModuleDefault_UpdateEvery(t *testing.T) {
	if v, _ := GetDefault(n).UpdateEvery(); v != i {
		t.Errorf("expected %d, but got %d", v, i)
	}
}

func TestModuleDefault_ChartsCleanup(t *testing.T) {
	if v, _ := GetDefault(n).ChartsCleanup(); v != i {
		t.Errorf("expected %d, but got %d", v, i)
	}
}

func TestModuleDefault_DisabledByDefault(t *testing.T) {
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
