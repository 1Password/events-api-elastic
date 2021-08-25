package api

import (
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
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
			"event": ECSEvent{
				Action: i.Category,
			},
			"user": ECSUser{
				ID:       i.SignInAttemptTargetUser.UUID,
				FullName: i.SignInAttemptTargetUser.Name,
				Email:    i.SignInAttemptTargetUser.Email,
			},
			"os": ECSOs{
				Name:    i.SignInAttemptClient.OSName,
				Version: i.SignInAttemptClient.OSVersion,
			},
			"source": ECSSource{
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
			"user": ECSUser{
				ID:       i.ItemUsageUser.UUID,
				FullName: i.ItemUsageUser.Name,
				Email:    i.ItemUsageUser.Email,
			},
			"os": ECSOs{
				Name:    i.ItemUsageClient.OSName,
				Version: i.ItemUsageClient.OSVersion,
			},
			"source": ECSSource{
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

type ECSEvent struct {
	Action string `json:"action,omitempty" ecs:"action"`
}

type ECSUser struct {
	ID       string `json:"id,omitempty" ecs:"id"`
	FullName string `json:"full_name,omitempty" ecs:"full_name"`
	Email    string `json:"email,omitempty" ecs:"email"`
}

type ECSOs struct {
	Name    string `json:"name,omitempty" ecs:"name"`
	Version string `json:"version,omitempty" ecs:"version"`
}

type ECSHost struct {
	IP string `json:"ip,omitempty" ecs:"ip"`
}
