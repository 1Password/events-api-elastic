package config

import (
	"fmt"
	"time"
)

type Config struct {
	InsecureSkipVerify bool        `config:"insecure_skip_verify"`
	SignInAttempts     EventConfig `config:"signin_attempts"`
	ItemUsages         EventConfig `config:"item_usages"`
	AuditEvents        EventConfig `config:"audit_events"`
}

func (c *Config) Validate() error {
	if err := c.SignInAttempts.Validate(); err != nil {
		return fmt.Errorf("invalid signin_attempts. %w", err)
	}
	if err := c.ItemUsages.Validate(); err != nil {
		return fmt.Errorf("invalid item_usages. %w", err)
	}
	if err := c.AuditEvents.Validate(); err != nil {
		return fmt.Errorf("invalid audit_events. %w", err)
	}
	return nil
}

var DefaultConfig = Config{
	InsecureSkipVerify: false,
	SignInAttempts: EventConfig{
		Enabled:         false,
		AuthToken:       "",
		StartingCursor:  `{ "limit": 1000, "start_time": "2020-01-01T00:00:00Z" }`,
		CursorStateFile: "eventsapibeat_signinattempts.state",
		SampleFrequency: 10 * time.Second,
	},
	ItemUsages: EventConfig{
		Enabled:         false,
		AuthToken:       "",
		StartingCursor:  `{ "limit": 1000, "start_time": "2020-01-01T00:00:00Z" }`,
		CursorStateFile: "eventsapibeat_itemusages.state",
		SampleFrequency: 10 * time.Second,
	},
	AuditEvents: EventConfig{
		Enabled:         false,
		AuthToken:       "",
		StartingCursor:  `{ "limit": 1000, "start_time": "2020-01-01T00:00:00Z" }`,
		CursorStateFile: "eventsapibeat_auditevents.state",
		SampleFrequency: 10 * time.Second,
	},
}

type EventConfig struct {
	Enabled         bool          `config:"enabled"`
	AuthToken       string        `config:"auth_token"`
	StartingCursor  string        `config:"starting_cursor"`
	CursorStateFile string        `config:"cursor_state_file"`
	SampleFrequency time.Duration `config:"sample_frequency"`
}

func (c *EventConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.SampleFrequency < 1*time.Second {
		return fmt.Errorf("sample_frequency can't be less than 1000ms")
	}
	if c.AuthToken == "" {
		return fmt.Errorf("auth_token can't be empty")
	}
	if c.StartingCursor == "" {
		return fmt.Errorf("starting_cursor can't be empty")
	}
	if c.CursorStateFile == "" {
		return fmt.Errorf("cursor_state_file can't be empty")
	}
	return nil
}
