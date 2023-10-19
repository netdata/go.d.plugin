// SPDX-License-Identifier: GPL-3.0-or-later

package discoverer

import (
	"errors"
	"fmt"

	"github.com/netdata/go.d.plugin/agent/discovery/sd/discoverer/kubernetes"
)

type (
	Config struct {
		Name      string `yaml:"name"`
		Discovery struct {
			K8S []kubernetes.Config `yaml:"k8s"`
		} `yaml:"discovery"`
		TagRules   []TagRuleConfig   `yaml:"tag"`
		BuildRules []BuildRuleConfig `yaml:"config"`
	}

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
	BuildRuleConfig struct {
		Name     string                 `yaml:"name"`     // optional
		Selector string                 `yaml:"selector"` // mandatory
		Apply    []BuildRuleApplyConfig `yaml:"apply"`    // mandatory, at least 1
	}
	BuildRuleApplyConfig struct {
		Selector string `yaml:"selector"` // mandatory
		Template string `yaml:"template"` // mandatory
	}
)

func validateConfig(cfg Config) error {
	if cfg.Name != "" {
		return errors.New("'name' not set")
	}
	if len(cfg.Discovery.K8S) == 0 {
		return errors.New("'discovery->k8s' not set")
	}
	if err := validateTagConfig(cfg.TagRules); err != nil {
		return fmt.Errorf("tag rules: %v", err)
	}
	if err := validateBuildConfig(cfg.BuildRules); err != nil {
		return fmt.Errorf("config rules: %v", err)
	}
	return nil
}

func validateTagConfig(rules []TagRuleConfig) error {
	if len(rules) == 0 {
		return errors.New("empty config, need least 1 rule")
	}
	for i, rule := range rules {
		if rule.Selector == "" {
			return fmt.Errorf("'rule[%s][%d]->selector' not set", rule.Name, i+1)
		}
		if rule.Tags == "" {
			return fmt.Errorf("'rule[%s][%d]->tags' not set", rule.Name, i+1)
		}
		if len(rule.Match) == 0 {
			return fmt.Errorf("'rule[%s][%d]->match' not set, need at least 1 rule match", rule.Name, i+1)
		}

		for j, match := range rule.Match {
			if match.Tags == "" {
				return fmt.Errorf("'rule[%s][%d]->match[%d]->tags' not set", rule.Name, i+1, j+1)
			}
			if match.Expr == "" {
				return fmt.Errorf("'rule[%s][%d]->match->expr[%d]' not set", rule.Name, i+1, j+1)
			}
		}
	}
	return nil
}

func validateBuildConfig(rules []BuildRuleConfig) error {
	if len(rules) == 0 {
		return errors.New("empty config, need least 1 rule")
	}
	for i, rule := range rules {
		if rule.Selector == "" {
			return fmt.Errorf("'rule[%s][%d]->selector' not set", rule.Name, i+1)
		}

		if len(rule.Apply) == 0 {
			return fmt.Errorf("'rule[%s][%d]->apply' not set", rule.Name, i+1)
		}

		for j, apply := range rule.Apply {
			if apply.Selector == "" {
				return fmt.Errorf("'rule[%s][%d]->apply[%d]->selector' not set", rule.Name, i+1, j+1)
			}
			if apply.Template == "" {
				return fmt.Errorf("'rule[%s][%d]->apply[%d]->template' not set", rule.Name, i+1, j+1)
			}
		}
	}
	return nil
}
