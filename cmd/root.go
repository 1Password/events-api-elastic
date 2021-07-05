package cmd

import (
	"github.com/elastic/beats/v7/libbeat/cmd"
	"github.com/elastic/beats/v7/libbeat/cmd/instance"

	"go.1password.io/eventsapibeat/beater"
)

// Name of this beat
var Name = beater.BeatName

// RootCmd to handle beats cli
var RootCmd = cmd.GenRootCmdWithSettings(beater.New, instance.Settings{Name: Name})
