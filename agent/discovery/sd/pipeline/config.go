// SPDX-License-Identifier: GPL-3.0-or-later

package pipeline

import (
	"errors"
	"fmt"
)

type Config struct {
}

type (
	TagConfig     []TagRuleConfig // mandatory, at least 1
	TagRuleConfig struct {
		Name     string               `yaml:"name"`
		Selector string               `yaml:"selector"` // mandatory
		Tags     string               `yaml:"tags"`     // mandatory
		Match    []TagRuleMatchConfig `yaml:"match"`    // mandatory, at least 1
	}
	TagRuleMatchConfig struct {
		Selector string `yaml:"selector"` // optional
		Tags     string `yaml:"tags"`     // mandatory
		Expr     string `yaml:"expr"`     // mandatory
	}
)

type (
	BuildConfig     []BuildRuleConfig // mandatory, at least 1
	BuildRuleConfig struct {
		Name     string                 `yaml:"name"`     // optional
		Selector string                 `yaml:"selector"` // mandatory
		Tags     string                 `yaml:"tags"`     // mandatory
		Apply    []BuildRuleApplyConfig `yaml:"apply"`    // mandatory, at least 1
	}
	BuildRuleApplyConfig struct {
		Selector string `yaml:"selector"` // mandatory
		Template string `yaml:"template"` // mandatory
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
