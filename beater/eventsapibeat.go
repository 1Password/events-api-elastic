package beater

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"
	"go.1password.io/eventsapibeat/api"
	"go.1password.io/eventsapibeat/config"
	"go.1password.io/eventsapibeat/store"
	"go.1password.io/eventsapibeat/utils"
	"go.1password.io/eventsapibeat/version"
)

const (
	BeatName           = "eventsapibeat"
	SignInAttemptsType = "signinattempts"
	ItemUsagesType     = "itemusages"
)

type EventsAPIBeat struct {
	config config.Config

	beatClient beat.Client
	log        *logp.Logger

	ctx                       context.Context
	cancel                    context.CancelFunc
	signInAttemptsCursorStore store.CursorStore
	itemUsagesCursorStore     store.CursorStore
	apiClient                 *api.Client
}

func New(_ *beat.Beat, cfg *common.Config) (beat.Beater, error) {

	var err error
	eventsAPIBeat := &EventsAPIBeat{
		log:    logp.NewLogger(BeatName),
		config: config.DefaultConfig,
	}

	if err = cfg.Unpack(&eventsAPIBeat.config); err != nil {
		return nil, fmt.Errorf("failed to unpack config file. %v", err)
	}

	if err = eventsAPIBeat.config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config. %v", err)
	}

	var transport http.RoundTripper
	if eventsAPIBeat.config.InsecureSkipVerify {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: eventsAPIBeat.config.InsecureSkipVerify,
			},
		}
	}
	eventsAPIBeat.apiClient, err = api.NewClient(transport)
	if err != nil {
		return nil, fmt.Errorf("failed to create api client. %w", err)
	}

	if eventsAPIBeat.config.SignInAttempts.Enabled {
		eventsAPIBeat.signInAttemptsCursorStore, err = store.NewCursorHistoryFileStore(eventsAPIBeat.config.SignInAttempts.CursorStateFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open sign-in attempts cursor file. %w", err)
		}
	}

	if eventsAPIBeat.config.ItemUsages.Enabled {
		eventsAPIBeat.itemUsagesCursorStore, err = store.NewCursorHistoryFileStore(eventsAPIBeat.config.ItemUsages.CursorStateFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open item usages cursor file. %w", err)
		}
	}

	return eventsAPIBeat, nil
}

func (e *EventsAPIBeat) Run(b *beat.Beat) error {
	e.log.Infof("%s v%s is running! Hit CTRL-C to stop it.", BeatName, version.Version)
	e.ctx, e.cancel = context.WithCancel(context.Background())

	var err error
	e.beatClient, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	eventsChan := make(chan *beat.Event)
	errorChan := make(chan error)

	if e.config.SignInAttempts.Enabled {
		e.log.Info("Starting sign-in attempts loop")
		jwt, err := utils.ParseJWTClaims(e.config.SignInAttempts.AuthToken)
		if err != nil {
			return err
		}

		if !jwt.HaveSignInAttemptsFeature() {
			return errors.New("sign-in attempt token does not have sign-in attempt feature")
		}
		go func() {
			err := e.signInAttemptsLoop(eventsChan)
			if err != nil {
				errorChan <- fmt.Errorf("failed when processing sign-in attempts. %v", err)
			}
		}()
	}

	if e.config.ItemUsages.Enabled {
		e.log.Info("Starting item usages loop")
		jwt, err := utils.ParseJWTClaims(e.config.ItemUsages.AuthToken)
		if err != nil {
			return err
		}

		if !jwt.HaveItemUsageFeature() {
			return errors.New("item usage token does not have item usage feature")
		}
		go func() {
			err := e.itemUsagesLoop(eventsChan)
			if err != nil {
				errorChan <- fmt.Errorf("failed when processing item usages. %v", err)
			}
		}()
	}

	for {
		select {
		case <-e.ctx.Done():
			return nil
		case ev := <-eventsChan:
			// publish event to beat
			e.beatClient.Publish(*ev)
		case err := <-errorChan:
			return err
		}
	}
}

func (e *EventsAPIBeat) signInAttemptsLoop(c chan<- *beat.Event) error {
	ticker := time.NewTicker(e.config.SignInAttempts.SampleFrequency)
	defer ticker.Stop()

	cursor, err := e.signInAttemptsCursorStore.GetValue()
	if err != nil {
		return fmt.Errorf("failed to get sign-in attempts cursor. %v", err)
	}
	if cursor == "" {
		cursor = e.config.SignInAttempts.StartingCursor
	}

	for {
		select {
		case <-ticker.C:
			var errs []string

			for {
				signInAttemptsResponse, err := e.apiClient.SignInAttempts(e.ctx, e.config.SignInAttempts.AuthToken, cursor)
				if err != nil {
					errs = append(errs, fmt.Sprintf("failed to fetch sign-in attempts. %v", err))
					break
				}

				for i := range signInAttemptsResponse.Items {
					item := &signInAttemptsResponse.Items[i]

					event := item.BeatEvent()
					_, _ = event.PutValue("@metadata.event_type", SignInAttemptsType)

					c <- event
				}

				cursor = fmt.Sprintf(`{ "cursor": "%s" }`, signInAttemptsResponse.Cursor)

				if err := e.signInAttemptsCursorStore.SetValue(cursor); err != nil {
					errs = append(errs, fmt.Sprintf("failed to set sign-in attempts cursor. %v", err))
				}

				if !signInAttemptsResponse.HasMore {
					break
				}

			}

			if len(errs) > 0 {
				return fmt.Errorf(strings.Join(errs, "."))
			}
		}
	}
}

func (e *EventsAPIBeat) itemUsagesLoop(c chan<- *beat.Event) error {
	ticker := time.NewTicker(e.config.ItemUsages.SampleFrequency)
	defer ticker.Stop()

	cursor, err := e.itemUsagesCursorStore.GetValue()
	if err != nil {
		return fmt.Errorf("failed to get item usages cursor. %v", err)
	}
	if cursor == "" {
		cursor = e.config.ItemUsages.StartingCursor
	}

	for {
		select {
		case <-ticker.C:
			var errs []string

			for {
				itemUsagesResponse, err := e.apiClient.ItemUsages(e.ctx, e.config.ItemUsages.AuthToken, cursor)
				if err != nil {
					errs = append(errs, fmt.Sprintf("failed to fetch item usages. %v", err))
					break
				}

				for i := range itemUsagesResponse.Items {
					item := &itemUsagesResponse.Items[i]

					event := item.BeatEvent()
					_, _ = event.PutValue("@metadata.event_type", ItemUsagesType)

					c <- event
				}

				cursor = fmt.Sprintf(`{ "cursor": "%s" }`, itemUsagesResponse.Cursor)

				if err := e.itemUsagesCursorStore.SetValue(cursor); err != nil {
					errs = append(errs, fmt.Sprintf("failed to set item usages cursor. %v", err))
				}

				if !itemUsagesResponse.HasMore {
					break
				}

			}

			if len(errs) > 0 {
				return fmt.Errorf(strings.Join(errs, "."))
			}
		}
	}
}

func (e *EventsAPIBeat) Stop() {
	e.cancel()
	if e.signInAttemptsCursorStore != nil {
		if err := e.signInAttemptsCursorStore.Close(); err != nil {
			e.log.Errorf("failed to close sign-in attempts cursor state file: %w", err)
		}
	}
	if e.itemUsagesCursorStore != nil {
		if err := e.itemUsagesCursorStore.Close(); err != nil {
			e.log.Errorf("failed to close item usages cursor state file: %w", err)
		}
	}
	err := e.beatClient.Close()
	if err != nil {
		e.log.Error(err)
	}
}
