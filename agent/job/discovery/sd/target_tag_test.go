package sd

import (
	"fmt"
	"github.com/netdata/go.d.plugin/agent/job/discovery/sd/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := map[string]tagSim{
		"valid config": {
			cfg: TagConfig{
				{
					Selector: "unknown",
					Tags:     "-unknown",
					Match: []MatchConfig{
						{Tags: "wizard", Expr: `{{eq .Class "wizard"}}`},
					},
				},
			},
		},
		"empty config": {
			cfg:     TagConfig{},
			invalid: true,
		},
		"config rule->selector not set": {
			invalid: true,
			cfg: TagConfig{
				{
					Selector: "",
					Tags:     "-unknown",
					Match: []MatchConfig{
						{Tags: "wizard", Expr: `{{eq .Class "wizard"}}`},
					},
				},
			},
		},
		"config rule->selector bad syntax": {
			invalid: true,
			cfg: TagConfig{
				{
					Selector: "!",
					Tags:     "-unknown",
					Match: []MatchConfig{
						{Tags: "wizard", Expr: `{{eq .Class "wizard"}}`},
					},
				},
			},
		},
		"config rule->tags not set": {
			invalid: true,
			cfg: TagConfig{
				{
					Selector: "unknown",
					Tags:     "",
					Match: []MatchConfig{
						{Tags: "wizard", Expr: `{{eq .Class "wizard"}}`},
					},
				},
			},
		},
		"config rule->match not set": {
			invalid: true,
			cfg: TagConfig{
				{
					Selector: "unknown",
					Tags:     "-unknown",
				},
			},
		},
		"config rule->match->selector bad syntax": {
			invalid: true,
			cfg: TagConfig{
				{
					Selector: "unknown",
					Tags:     "-unknown",
					Match: []MatchConfig{
						{Selector: "!", Tags: "wizard", Expr: `{{eq .Class "wizard"}}`},
					},
				},
			},
		},
		"config rule->match->tags not set": {
			invalid: true,
			cfg: TagConfig{
				{
					Selector: "unknown",
					Tags:     "-unknown",
					Match: []MatchConfig{
						{Tags: "", Expr: `{{eq .Class "wizard"}}`},
					},
				},
			},
		},
		"config rule->match->expr not set": {
			invalid: true,
			cfg: TagConfig{
				{
					Selector: "unknown",
					Tags:     "-unknown",
					Match: []MatchConfig{
						{Tags: "wizard", Expr: ""},
					},
				},
			},
		},
		"config rule->match->expr unknown func": {
			invalid: true,
			cfg: TagConfig{
				{
					Selector: "unknown",
					Tags:     "-unknown",
					Match: []MatchConfig{
						{Tags: "wizard", Expr: `{{error .Class "wizard"}}`},
					},
				},
			},
		},
	}

	for name, sim := range tests {
		t.Run(name, func(t *testing.T) { sim.run(t) })
	}
}

func TestManager_Tag(t *testing.T) {
	tests := map[string]tagSim{
		"3 rule service": {
			cfg: TagConfig{
				{
					Selector: "unknown",
					Tags:     "-unknown",
					Match: []MatchConfig{
						{Tags: "wizard", Expr: `{{eq .Class "wizard"}}`},
						{Tags: "knight", Expr: `{{eq .Class "knight"}}`},
						{Tags: "cleric", Expr: `{{eq .Class "cleric"}}`},
					},
				},
				{
					Selector: "!unknown",
					Tags:     "candidate",
					Match: []MatchConfig{
						{Tags: "human", Expr: `{{eq .Race "human"}}`},
						{Tags: "elf", Expr: `{{eq .Race "elf"}}`},
						{Tags: "dwarf", Expr: `{{eq .Race "dwarf"}}`},
					},
				},
				{
					Selector: "candidate",
					Tags:     "-candidate",
					Match: []MatchConfig{
						{Tags: "teamup", Expr: `{{gt .Level 9000}}`},
					},
				},
			},
			inputs: []tagSimInput{
				{
					desc:         "all rules fail",
					target:       mockTagTarget{tags: model.Tags{"unknown": {}}, Class: "fighter", Race: "orc", Level: 9001},
					expectedTags: model.Tags{"unknown": {}},
				},
				{
					desc:         "1st rule match",
					target:       mockTagTarget{tags: model.Tags{"unknown": {}}, Class: "knight", Race: "undead", Level: 9001},
					expectedTags: model.Tags{"knight": {}},
				},
				{
					desc:         "1st, 2nd rules match",
					target:       mockTagTarget{tags: model.Tags{"unknown": {}}, Class: "knight", Race: "human", Level: 8999},
					expectedTags: model.Tags{"knight": {}, "human": {}, "candidate": {}},
				},
				{
					desc:         "all rules match",
					target:       mockTagTarget{tags: model.Tags{"unknown": {}}, Class: "wizard", Race: "human", Level: 9001},
					expectedTags: model.Tags{"wizard": {}, "human": {}, "teamup": {}},
				},
				{
					desc:         "all rules match",
					target:       mockTagTarget{tags: model.Tags{"unknown": {}}, Class: "knight", Race: "dwarf", Level: 9001},
					expectedTags: model.Tags{"knight": {}, "dwarf": {}, "teamup": {}},
				},
				{
					desc:         "all rules match",
					target:       mockTagTarget{tags: model.Tags{"unknown": {}}, Class: "cleric", Race: "elf", Level: 9001},
					expectedTags: model.Tags{"cleric": {}, "elf": {}, "teamup": {}},
				},
			},
		},
	}

	for name, sim := range tests {
		t.Run(name, func(t *testing.T) { sim.run(t) })
	}
}

func TestRule_Tag(t *testing.T) {
	tests := map[string]tagSim{
		"simple rule": {
			cfg: TagConfig{
				{
					Selector: "unknown",
					Tags:     "-unknown",
					Match: []MatchConfig{
						{Selector: "human", Tags: "wizard", Expr: `{{eq .Class "wizard"}}`},
						{Tags: "missingkey", Expr: `{{eq .Name "yoda"}}`},
					},
				},
			},
			inputs: []tagSimInput{
				{
					desc:         "not match rule selector",
					target:       mockTagTarget{Class: "fighter"},
					expectedTags: nil,
				},
				{
					desc:         "not match rule match selector",
					target:       mockTagTarget{tags: model.Tags{"unknown": {}}, Class: "fighter"},
					expectedTags: model.Tags{"unknown": {}},
				},
				{
					desc:         "not match rule match expression",
					target:       mockTagTarget{tags: model.Tags{"unknown": {}, "human": {}}, Class: "fighter"},
					expectedTags: model.Tags{"unknown": {}, "human": {}},
				},
				{
					desc:         "match expression",
					target:       mockTagTarget{tags: model.Tags{"unknown": {}, "human": {}}, Class: "wizard"},
					expectedTags: model.Tags{"wizard": {}, "human": {}},
				},
				{
					desc:         "match expression missingkey error",
					target:       mockTagTarget{tags: model.Tags{"unknown": {}, "missingkey": {}}, Class: "knight"},
					expectedTags: model.Tags{"unknown": {}, "missingkey": {}},
				},
			},
		},
	}

	for name, sim := range tests {
		t.Run(name, func(t *testing.T) { sim.run(t) })
	}
}

func TestRule_Tag_UseCustomFunction(t *testing.T) {
	newSim := func(expr string) tagSim {
		return tagSim{
			cfg: TagConfig{
				{
					Selector: "*",
					Tags:     "-nothing",
					Match: []MatchConfig{
						{Tags: "wizard", Expr: expr},
					},
				},
			},
			inputs: []tagSimInput{
				{
					target:       mockTagTarget{Class: "wizard", tags: model.Tags{"key": {}}},
					expectedTags: model.Tags{"key": {}, "wizard": {}},
				},
			},
		}
	}

	tests := map[string]tagSim{
		"glob": newSim(`{{glob .Class "w*z*rd"}}`),
		"re":   newSim(`{{re .Class "^w[iI]z.*d$"}}`),
	}

	for name, sim := range tests {
		t.Run(name, func(t *testing.T) { sim.run(t) })
	}
}

type mockTagTarget struct {
	tags  model.Tags
	Class string
	Race  string
	Level int
}

func (m mockTagTarget) Tags() model.Tags { return m.tags }
func (m mockTagTarget) Hash() uint64     { return 0 }
func (m mockTagTarget) TUID() string     { return "" }
func (m mockTagTarget) String() string {
	return fmt.Sprintf("Class: %s, Race: %s, Level: %d, Tags: %s", m.Class, m.Race, m.Level, m.Tags())
}

type (
	tagSim struct {
		cfg     TagConfig
		invalid bool
		inputs  []tagSimInput
	}
	tagSimInput struct {
		desc         string
		target       mockTagTarget
		expectedTags model.Tags
	}
)

func (sim tagSim) run(t *testing.T) {
	mgr, err := NewTagManager(sim.cfg)

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
		name := fmt.Sprintf("input:'%s'[%d], target:'%s', expected tags:'%s'",
			input.desc, i+1, input.target, input.expectedTags)

		mgr.Tag(input.target)
		assert.Equalf(t, input.expectedTags, input.target.Tags(), name)
	}
}
