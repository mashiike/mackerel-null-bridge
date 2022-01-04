package mackerelnullbridge

import (
	"errors"
	"fmt"
	"log"
	"strings"

	gv "github.com/hashicorp/go-version"
	gc "github.com/kayac/go-config"
)

type Config struct {
	RequiredVersion string `yaml:"required_version" json:"required_version"`

	Targets []*TargetConfig `yaml:"targets" json:"targets"`

	versionConstraints gv.Constraints
}

type TargetConfig struct {
	Service      string      `yaml:"service" json:"service"`
	MetricName   string      `yaml:"metric_name" json:"metric_name"`
	Value        interface{} `yaml:"value" json:"value"`
	DelaySeconds int64       `yaml:"delay_seconds" json:"delay_seconds"`
}

// Load loads configuration file from file paths.
func (c *Config) Load(paths ...string) error {
	if err := gc.LoadWithEnv(c, paths...); err != nil {
		return err
	}
	return c.Restrict()
}

// Restrict restricts a configuration.
func (c *Config) Restrict() error {
	if c.RequiredVersion != "" {
		constraints, err := gv.NewConstraint(c.RequiredVersion)
		if err != nil {
			return fmt.Errorf("required_version has invalid format: %w", err)
		}
		c.versionConstraints = constraints
	}
	for i, target := range c.Targets {
		if err := target.Restrict(); err != nil {
			return fmt.Errorf("targets[%d].%w", i, err)
		}
	}

	return nil
}

func (c *TargetConfig) Restrict() error {
	if c.Service == "" {
		return errors.New("service is required")
	}
	if c.MetricName == "" {
		return errors.New("metric_name is required")
	}
	if c.Value == nil {
		c.Value = 0.0
	}
	if c.DelaySeconds == 0 {
		c.DelaySeconds = 5 * 60
	}
	return nil
}

// ValidateVersion validates a version satisfies required_version.
func (c *Config) ValidateVersion(version string) error {
	if c.versionConstraints == nil {
		log.Println("[warn] required_version is empty. Skip checking required_version.")
		return nil
	}
	versionParts := strings.SplitN(version, "-", 2)
	v, err := gv.NewVersion(versionParts[0])
	if err != nil {
		log.Printf("[warn] invalid version format \"%s\". Skip checking required_version.", version)
		// invalid version string (e.g. "current") always allowed
		return nil
	}
	if !c.versionConstraints.Check(v) {
		return fmt.Errorf("version %s does not satisfy constraints required_version: %s", version, c.versionConstraints)
	}
	return nil
}

// NewDefaultConfig creates a default configuration.
func NewDefaultConfig() *Config {
	return &Config{}
}
