package ci

import (
	"fmt"
	"os"

	"github.com/wagoodman/dive/image"
	"gopkg.in/yaml.v2"
)

// RuleStatus represents the result of a single CI rule evaluation.
type RuleStatus int

const (
	RulePass    RuleStatus = iota
	RuleFail
	RuleDisabled
)

// Rule defines a single CI evaluation rule with a name, threshold, and whether it is enabled.
type Rule struct {
	Name      string
	Enabled   bool
	Threshold float64
}

// RuleResult captures the outcome of evaluating a single rule.
type RuleResult struct {
	Rule   Rule
	Actual float64
	Status RuleStatus
}

// EvaluationResult aggregates all rule results for a CI run.
type EvaluationResult struct {
	Results []RuleResult
	Pass    bool
}

// Config holds the CI configuration loaded from a .dive-ci file.
type Config struct {
	Rules map[string]interface{} `yaml:"rules"`
}

// LoadConfig reads and parses a .dive-ci YAML configuration file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read CI config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unable to parse CI config: %w", err)
	}
	return &cfg, nil
}

// Evaluator runs CI rules against an analyzed image.
type Evaluator struct {
	cfg *Config
}

// NewEvaluator creates a new Evaluator from the given CI config.
func NewEvaluator(cfg *Config) *Evaluator {
	return &Evaluator{cfg: cfg}
}

// Evaluate runs all configured rules against the provided image analysis result.
// Returns an EvaluationResult indicating whether the image passes CI checks.
// Note: rules are evaluated independently; a single failure marks the whole run as failed.
func (e *Evaluator) Evaluate(analysis *image.AnalysisResult) (*EvaluationResult, error) {
	result := &EvaluationResult{Pass: true}

	for ruleName, rawValue := range e.cfg.Rules {
		rule, threshold, enabled, err := parseRule(ruleName, rawValue)
		if err != nil {
			return nil, fmt.Errorf("invalid rule %q: %w", ruleName, err)
		}
		rule.Enabled = enabled
		rule.Threshold = threshold

		if !enabled {
			result.Results = append(result.Results, RuleResult{Rule: rule, Status: RuleDisabled})
			continue
		}

		actual, err := getMetric(ruleName, analysis)
		if err != nil {
			return nil, fmt.Errorf("unable to evaluate rule %q: %w", ruleName, err)
		}

		status := RulePass
		// Fail if actual value is strictly below the configured threshold.
		if actual < threshold {
			status = RuleFail
			result.Pass = false
		}

		result.Results = append(result.Results, RuleResult{
			Rule:   rule,
			Actual: actual,
			Status: status,
		})
	}

	return result, nil
}

// parseRule extracts the threshold and enabled flag from a raw rule value.
// Accepted string values for disabling a rule: "disabled", "off", or "ignore".
func parseRule(name string, raw interface{}) (Rule, float64, bool, error) {
	rule := Rule{Name: name}
	switch v := raw.(type) {
	case float64:
		return rule, v, true, nil
	case int:
		return rule, float64(v), true, nil
	case string:
		if v == "disabled" || v == "off" || v == "ignore" {
			return rule, 0, false, nil
		}
		return rule, 0, false, fmt.Errorf("unsupported string value %q", v)
	default:
		return rule, 0, false, fmt.Errorf("unsupported value type %T for rule %q", raw, name)
	}
}
