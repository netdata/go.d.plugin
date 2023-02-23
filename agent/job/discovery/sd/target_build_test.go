package sd

import (
	"fmt"
	"github.com/netdata/go.d.plugin/agent/job/discovery/sd/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBuildManager(t *testing.T) {
	tests := map[string]buildSim{
		"valid config": {
			cfg: BuildConfig{
				{
					Selector: "unknown",
					Tags:     "-unknown",
					Apply: []ApplyConfig{
						{Selector: "wizard", Template: `class {{.Class}}`},
					},
				},
			},
		},
		"empty config": {
			invalid: true,
			cfg:     BuildConfig{},
		},
		"config rule->selector not set": {
			invalid: true,
			cfg: BuildConfig{
				{
					Selector: "",
					Tags:     "-unknown",
					Apply: []ApplyConfig{
						{Selector: "wizard", Template: `class {{.Class}}`},
					},
				},
			},
		},
		"config rule->selector bad syntax": {
			invalid: true,
			cfg: BuildConfig{
				{
					Selector: "!",
					Tags:     "-unknown",
					Apply: []ApplyConfig{
						{Selector: "wizard", Template: `class {{.Class}}`},
					},
				},
			},
		},
		"config rule->tags not set": {
			invalid: true,
			cfg: BuildConfig{
				{
					Selector: "unknown",
					Tags:     "",
					Apply: []ApplyConfig{
						{Selector: "wizard", Template: `class {{.Class}}`},
					},
				},
			},
		},
		"config rule->apply not set": {
			invalid: true,
			cfg: BuildConfig{
				{
					Selector: "unknown",
					Tags:     "-unknown",
				},
			},
		},
		"config rule->apply->selector not set": {
			invalid: true,
			cfg: BuildConfig{
				{
					Selector: "unknown",
					Tags:     "-unknown",
					Apply: []ApplyConfig{
						{Selector: "", Template: `class {{.Class}}`},
					},
				},
			},
		},
		"config rule->apply->selector bad syntax": {
			invalid: true,
			cfg: BuildConfig{
				{
					Selector: "unknown",
					Tags:     "-unknown",
					Apply: []ApplyConfig{
						{Selector: "!", Template: `class {{.Class}}`},
					},
				},
			},
		},
		"config rule->apply->template not set": {
			invalid: true,
			cfg: BuildConfig{
				{
					Selector: "unknown",
					Tags:     "-unknown",
					Apply: []ApplyConfig{
						{Selector: "wizard", Template: ""},
					},
				},
			},
		},
		"config rule->apply->template missingkey (unknown func)": {
			invalid: true,
			cfg: BuildConfig{
				{
					Selector: "unknown",
					Tags:     "-unknown",
					Apply: []ApplyConfig{
						{Selector: "wizard", Template: `class {{error .Class}}`},
					},
				},
			},
		},
	}

	for name, sim := range tests {
		t.Run(name, func(t *testing.T) { sim.run(t) })
	}
}

func TestManager_Build(t *testing.T) {
	tests := map[string]buildSim{
		"4 rule service": {
			cfg: BuildConfig{
				{
					Selector: "class",
					Tags:     "built",
					Apply: []ApplyConfig{
						{Selector: "*", Template: `Class: {{.Class}}`},
					},
				},
				{
					Selector: "race",
					Tags:     "built",
					Apply: []ApplyConfig{
						{Selector: "*", Template: `Race: {{.Race}}`},
					},
				},
				{
					Selector: "level",
					Tags:     "built",
					Apply: []ApplyConfig{
						{Selector: "*", Template: `Level: {{.Level}}`},
					},
				},
				{
					Selector: "full",
					Tags:     "built",
					Apply: []ApplyConfig{
						{Selector: "*", Template: `Class: {{.Class}}, Race: {{.Race}}, Level: {{.Level}}`},
					},
				},
			},
			inputs: []buildSimInput{
				{
					desc: "1st rule match",
					target: mockBuildTarget{
						tag:   model.Tags{"class": {}},
						Class: "fighter", Race: "orc", Level: 9001,
					},
					expectedCfgs: []model.Config{
						{Conf: "Class: fighter", Tags: model.Tags{"built": {}}},
					},
				},
				{
					desc: "1st, 2nd rules match",
					target: mockBuildTarget{
						tag:   model.Tags{"class": {}, "race": {}},
						Class: "fighter", Race: "orc", Level: 9001,
					},
					expectedCfgs: []model.Config{
						{Conf: "Class: fighter", Tags: model.Tags{"built": {}}},
						{Conf: "Race: orc", Tags: model.Tags{"built": {}}},
					},
				},
				{
					desc: "1st, 2nd, 3rd rules match",
					target: mockBuildTarget{
						tag:   model.Tags{"class": {}, "race": {}, "level": {}},
						Class: "fighter", Race: "orc", Level: 9001,
					},
					expectedCfgs: []model.Config{
						{Conf: "Class: fighter", Tags: model.Tags{"built": {}}},
						{Conf: "Race: orc", Tags: model.Tags{"built": {}}},
						{Conf: "Level: 9001", Tags: model.Tags{"built": {}}},
					},
				},
				{
					desc: "all rules match",
					target: mockBuildTarget{
						tag:   model.Tags{"class": {}, "race": {}, "level": {}, "full": {}},
						Class: "fighter", Race: "orc", Level: 9001,
					},
					expectedCfgs: []model.Config{
						{Conf: "Class: fighter", Tags: model.Tags{"built": {}}},
						{Conf: "Race: orc", Tags: model.Tags{"built": {}}},
						{Conf: "Level: 9001", Tags: model.Tags{"built": {}}},
						{Conf: "Class: fighter, Race: orc, Level: 9001", Tags: model.Tags{"built": {}}},
					},
				},
			},
		},
	}

	for name, sim := range tests {
		t.Run(name, func(t *testing.T) { sim.run(t) })
	}
}

func TestRule_Build(t *testing.T) {
	tests := map[string]buildSim{
		"simple rule": {
			cfg: BuildConfig{
				{
					Selector: "build",
					Tags:     "built",
					Apply: []ApplyConfig{
						{Selector: "human", Template: `Class: {{.Class}}, Race: {{.Race}}, Level: {{.Level}}`},
						{Selector: "missingkey", Template: `{{.Name}}`},
					},
				},
			},
			inputs: []buildSimInput{
				{
					desc: "not match rule selector",
					target: mockBuildTarget{
						tag:   model.Tags{"nothing": {}},
						Class: "fighter", Race: "orc", Level: 9001,
					},
				},
				{
					desc: "not match rule match selector",
					target: mockBuildTarget{
						tag:   model.Tags{"build": {}},
						Class: "fighter", Race: "orc", Level: 9001,
					},
				},
				{
					desc: "match everything",
					target: mockBuildTarget{
						tag:   model.Tags{"build": {}, "human": {}},
						Class: "fighter", Race: "human", Level: 9001,
					},
					expectedCfgs: []model.Config{
						{Conf: "Class: fighter, Race: human, Level: 9001", Tags: model.Tags{"built": {}}},
					},
				},
				{
					desc: "missingkey error",
					target: mockBuildTarget{
						tag:   model.Tags{"build": {}, "missingkey": {}},
						Class: "fighter", Race: "human", Level: 9001,
					},
				},
			},
		},
	}

	for name, sim := range tests {
		t.Run(name, func(t *testing.T) { sim.run(t) })
	}
}

type mockBuildTarget struct {
	tag   model.Tags
	Class string
	Race  string
	Level int
}

func (m mockBuildTarget) Tags() model.Tags { return m.tag }
func (mockBuildTarget) Hash() uint64       { return 0 }
func (mockBuildTarget) TUID() string       { return "" }
func (m mockBuildTarget) String() string {
	return fmt.Sprintf("Class: %s, Race: %s, Level: %d, Tags: %s", m.Class, m.Race, m.Level, m.Tags())
}

type (
	buildSim struct {
		cfg     BuildConfig
		invalid bool
		inputs  []buildSimInput
	}
	buildSimInput struct {
		desc         string
		target       mockBuildTarget
		expectedCfgs []model.Config
	}
)

func (sim buildSim) run(t *testing.T) {
	mgr, err := NewBuildManager(sim.cfg)

	if sim.invalid {
		require.Error(t, err)
		return
	}

	require.NoError(t, err)
	require.NotNil(t, mgr)

	if len(sim.inputs) == 0 {
		return
	}

	for i, input := range sim.inputs {
		name := fmt.Sprintf("input:'%s'[%d], target:'%s', expected configs:'%v'", input.desc, i+1, input.target, input.expectedCfgs)
		actual := mgr.Build(input.target)

		assert.Equalf(t, input.expectedCfgs, actual, name)
	}
}
