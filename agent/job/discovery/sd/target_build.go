package sd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/netdata/go.d.plugin/agent/job/discovery/sd/model"
	"text/template"

	"github.com/netdata/go.d.plugin/logger"
)

type (
	BuildConfig     []BuildRuleConfig // mandatory, at least 1
	BuildRuleConfig struct {
		Name     string        `yaml:"name"`     // optional
		Selector string        `yaml:"selector"` // mandatory
		Tags     string        `yaml:"tags"`     // mandatory
		Apply    []ApplyConfig `yaml:"apply"`    // mandatory, at least 1
	}
	ApplyConfig struct {
		Selector string `yaml:"selector"` // mandatory
		Tags     string `yaml:"tags"`     // optional
		Template string `yaml:"template"` // mandatory
	}
)

func validateBuildConfig(cfg BuildConfig) error {
	if len(cfg) == 0 {
		return errors.New("empty config, need least 1 rule")
	}
	for i, ruleCfg := range cfg {
		if ruleCfg.Selector == "" {
			return fmt.Errorf("'rule->selector' not set (rule %s[%d])", ruleCfg.Name, i+1)
		}

		if ruleCfg.Tags == "" {
			return fmt.Errorf("'rule->tags' not set (rule %s[%d])", ruleCfg.Name, i+1)
		}
		if len(ruleCfg.Apply) == 0 {
			return fmt.Errorf("'rule->apply' not set (rule %s[%d])", ruleCfg.Name, i+1)
		}

		for j, applyCfg := range ruleCfg.Apply {
			if applyCfg.Selector == "" {
				return fmt.Errorf("'rule->apply->selector' not set (rule %s[%d]/apply [%d])", ruleCfg.Name, i+1, j+1)
			}
			if applyCfg.Template == "" {
				return fmt.Errorf("'rule->apply->template' not set (rule %s[%d]/apply [%d])", ruleCfg.Name, i+1, j+1)
			}
		}
	}
	return nil
}

type (
	BuildManager struct {
		rules []*buildRule
		buf   bytes.Buffer
		*logger.Logger
	}
	buildRule struct {
		name  string
		id    int
		sr    selector
		tags  model.Tags
		apply []*ruleApply
	}
	ruleApply struct {
		id   int
		sr   selector
		tags model.Tags
		tmpl *template.Template
	}
)

func NewBuildManager(cfg BuildConfig) (*BuildManager, error) {
	if err := validateBuildConfig(cfg); err != nil {
		return nil, fmt.Errorf("build manager config validation: %v", err)
	}
	mgr, err := initBuildManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("build manager initialization: %v", err)
	}
	return mgr, nil
}

func (m *BuildManager) Build(target model.Target) (configs []model.Config) {
	for _, rule := range m.rules {
		if !rule.sr.matches(target.Tags()) {
			continue
		}

		for _, apply := range rule.apply {
			if !apply.sr.matches(target.Tags()) {
				continue
			}

			m.buf.Reset()
			if err := apply.tmpl.Execute(&m.buf, target); err != nil {
				m.Warningf("failed to execute rule apply '%d/%d' on target '%s': %v", rule.id, apply.id, target.TUID(), err)
				continue
			}

			cfg := model.Config{
				Tags: model.NewTags(),
				Conf: m.buf.String(),
			}
			cfg.Tags.Merge(rule.tags)
			cfg.Tags.Merge(apply.tags)

			configs = append(configs, cfg)
		}
	}

	if len(configs) > 0 {
		m.Infof("built %d config(s) for target '%s'", len(configs), target.TUID())
	}
	return configs
}

func initBuildManager(conf BuildConfig) (*BuildManager, error) {
	if len(conf) == 0 {
		return nil, errors.New("empty config")
	}
	mgr := &BuildManager{
		Logger: logger.New("build", "manager"),
	}

	for i, cfg := range conf {
		rule := buildRule{id: i + 1, name: cfg.Name}
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

		for i, cfg := range cfg.Apply {
			apply := ruleApply{id: i + 1}

			sr, err := parseSelector(cfg.Selector)
			if err != nil {
				return nil, err
			}
			apply.sr = sr

			tags, err := model.ParseTags(cfg.Tags)
			if err != nil {
				return nil, err
			}
			apply.tags = tags

			tmpl, err := parseTemplate(cfg.Template)
			if err != nil {
				return nil, err
			}
			apply.tmpl = tmpl

			rule.apply = append(rule.apply, &apply)
		}
		mgr.rules = append(mgr.rules, &rule)
	}
	return mgr, nil
}
