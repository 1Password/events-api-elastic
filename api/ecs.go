package api

import (
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/ecs/code/go/ecs"
)

const CustomFieldSet = "onepassword"

var emptyMap = map[string]struct{}{}

func (i *SignInAttempt) BeatEvent() *beat.Event {
	var details interface{} = emptyMap
	if i.Details == nil {
		details = i.Details
	}
	e := &beat.Event{
		Timestamp: i.Timestamp,
		Fields: common.MapStr{
			"event": ecs.Event{
				Action: i.Category,
			},
			"user": ecs.User{
				ID:       i.SignInAttemptTargetUser.UUID,
				FullName: i.SignInAttemptTargetUser.Name,
				Email:    i.SignInAttemptTargetUser.Email,
			},
			"os": ecs.Os{
				Name:    i.SignInAttemptClient.OSName,
				Version: i.SignInAttemptClient.OSVersion,
			},
			"host": ecs.Host{
				IP: i.SignInAttemptClient.IPAddress,
			},
			CustomFieldSet: common.MapStr{
				"uuid":         i.UUID,
				"session_uuid": i.SessionUUID,
				"type":         i.Type,
				"country":      i.Country,
				"details":      details,
				"client": common.MapStr{
					"app_name":         i.SignInAttemptClient.AppName,
					"app_version":      i.SignInAttemptClient.AppVersion,
					"platform_name":    i.SignInAttemptClient.PlatformName,
					"platform_version": i.SignInAttemptClient.PlatformVersion,
				},
			},
		},
	}

	return e
}

func (i *ItemUsage) BeatEvent() *beat.Event {
	e := &beat.Event{
		Timestamp: i.Timestamp,
		Fields: common.MapStr{
			"user": ecs.User{
				ID:       i.ItemUsageUser.UUID,
				FullName: i.ItemUsageUser.Name,
				Email:    i.ItemUsageUser.Email,
			},
			"os": ecs.Os{
				Name:    i.ItemUsageClient.OSName,
				Version: i.ItemUsageClient.OSVersion,
			},
			"host": ecs.Host{
				IP: i.ItemUsageClient.IPAddress,
			},
			CustomFieldSet: common.MapStr{
				"uuid":         i.UUID,
				"used_version": i.UsedVersion,
				"vault_uuid":   i.VaultUUID,
				"item_uuid":    i.ItemUUID,
				"client": common.MapStr{
					"app_name":         i.ItemUsageClient.AppName,
					"app_version":      i.ItemUsageClient.AppVersion,
					"platform_name":    i.ItemUsageClient.PlatformName,
					"platform_version": i.ItemUsageClient.PlatformVersion,
				},
			},
		},
	}

	return e
}
