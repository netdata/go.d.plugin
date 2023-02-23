package sd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/netdata/go.d.plugin/agent/job/discovery/sd/model"
	"strings"
	"text/template"

	"github.com/netdata/go.d.plugin/logger"
)

type (
	TagConfig  []RuleConfig // mandatory, at least 1
	RuleConfig struct {
		Name     string        `yaml:"name"`
		Selector string        `yaml:"selector"` // mandatory
		Tags     string        `yaml:"tags"`     // mandatory
		Match    []MatchConfig `yaml:"match"`    // mandatory, at least 1
	}
	MatchConfig struct {
		Selector string `yaml:"selector"` // optional
		Tags     string `yaml:"tags"`     // mandatory
		Expr     string `yaml:"expr"`     // mandatory
	}
)

func validateTagConfig(cfg TagConfig) error {
	if len(cfg) == 0 {
		return errors.New("empty config, need least 1 rule")
	}
	for i, rule := range cfg {
		if rule.Selector == "" {
			return fmt.Errorf("'rule->selector' not set (rule %s[%d])", rule.Name, i+1)
		}
		if rule.Tags == "" {
			return fmt.Errorf("'rule->tags' not set (rule %s[%d])", rule.Name, i+1)
		}
		if len(rule.Match) == 0 {
			return fmt.Errorf("'rule->match' not set, need at least 1 rule match (rule %s[%d])", rule.Name, i+1)
		}

		for j, match := range rule.Match {
			if match.Tags == "" {
				return fmt.Errorf("'rule->match->tags' not set (rule %s[%d]/match [%d])", rule.Name, i+1, j+1)
			}
			if match.Expr == "" {
				return fmt.Errorf("'rule->match->expr' not set (rule %s[%d]/match [%d])", rule.Name, i+1, j+1)
			}
		}
	}
	return nil
}

type (
	TagManager struct {
		rules []*tagRule
		buf   bytes.Buffer
		*logger.Logger
	}
	tagRule struct {
		name  string
		id    int
		sr    selector
		tags  model.Tags
		match []*ruleMatch
	}
	ruleMatch struct {
		id   int
		sr   selector
		tags model.Tags
		expr *template.Template
	}
)

func NewTagManager(cfg TagConfig) (*TagManager, error) {
	if err := validateTagConfig(cfg); err != nil {
		return nil, fmt.Errorf("tag manager config validation: %v", err)
	}
	mgr, err := initTagManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("tag manager initialization: %v", err)
	}
	return mgr, nil
}

func (m *TagManager) Tag(target model.Target) {
	for _, rule := range m.rules {
		if !rule.sr.matches(target.Tags()) {
			continue
		}

		for _, match := range rule.match {
			if !match.sr.matches(target.Tags()) {
				continue
			}

			m.buf.Reset()
			if err := match.expr.Execute(&m.buf, target); err != nil {
				m.Warningf("failed to execute rule match '%d/%d' on target '%s': %v", rule.id, match.id, target.TUID(), err)
				continue
			}
			if strings.TrimSpace(m.buf.String()) != "true" {
				continue
			}

			target.Tags().Merge(rule.tags)
			target.Tags().Merge(match.tags)
			m.Debugf("matched target '%s', tags: %s", target.TUID(), target.Tags())
		}
	}
}

func initTagManager(conf TagConfig) (*TagManager, error) {
	if len(conf) == 0 {
		return nil, errors.New("empty config")
	}

	mgr := &TagManager{
		rules:  nil,
		Logger: logger.New("tag", "manager"),
	}
	for i, cfg := range conf {
		rule := tagRule{id: i + 1, name: cfg.Name}
		if sr, err := parseSelector(cfg.Selector); err != nil {
			return nil, err
		} else {
			rule.sr = sr
		}

		if tags, err := model.ParseTags(cfg.Tags); err != nil {
			return nil, err
		} else {
			rule.tags = tags
		}

		for i, cfg := range cfg.Match {
			match := ruleMatch{id: i + 1}

			sr, err := parseSelector(cfg.Selector)
			if err != nil {
				return nil, err
			}
			match.sr = sr

			tags, err := model.ParseTags(cfg.Tags)
			if err != nil {
				return nil, err
			}
			match.tags = tags

			tmpl, err := parseTemplate(cfg.Expr)
			if err != nil {
				return nil, err
			}
			match.expr = tmpl

			rule.match = append(rule.match, &match)
		}
		mgr.rules = append(mgr.rules, &rule)
	}
	return mgr, nil
}

func parseTemplate(line string) (*template.Template, error) {
	return template.New("root").
		Option("missingkey=error").
		Funcs(funcMap).
		Parse(line)
}
