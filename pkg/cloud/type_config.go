package cloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	ErrEmptyOAuthToken    = errors.New("error CloudConfig: OAuth token is empty")
	ErrEmptyFolderIdToken = errors.New("error CloudConfig: folder-id is empty")
)

var (
	ErrTickerNotRunning     = errors.New("error CloudConfig: ticker: attempt to stop a ticker that is not running")
	ErrTickerAlreadyRunning = errors.New("error CloudConfig: ticker already started")
)

const (
	timeLimitBeforeExpires = 30 * time.Minute
	timeBetweenUpdates     = 5 * time.Hour
)

const url = "https://iam.api.cloud.yandex.net/iam/v1/tokens"

type token struct {
	IamToken  string
	ExpiresAt time.Time
}

type CloudConfig struct {
	iamToken   string
	oAuthToken string
	folderId   string
	mu         *sync.RWMutex
	ticker     *time.Ticker
	deltaTime  time.Duration
	expired    time.Time
}

func NewCloudConfig(oAuthToken, folderId string) *CloudConfig {
	return &CloudConfig{
		oAuthToken: oAuthToken,
		folderId:   folderId,
		mu:         &sync.RWMutex{},
		ticker:     nil,
		deltaTime:  timeBetweenUpdates,
	}
}

func (cc *CloudConfig) UpdateCloudConfig() error {
	strBody := "{\"yandexPassportOauthToken\": \"" + cc.oAuthToken + "\"}"
	res, err := http.Post(url, "application/json", strings.NewReader(strBody))
	if err != nil {
		return fmt.Errorf("error: UpdateCloudConfig http request: %w", err)
	}
	defer res.Body.Close()

	dec := json.NewDecoder(res.Body)
	var t token
	err = dec.Decode(&t)
	if err != nil {
		return fmt.Errorf("error: UpdateCloudConfig json parse: %w", err)
	}

	cc.mu.Lock()
	cc.iamToken = t.IamToken
	cc.expired = t.ExpiresAt
	cc.mu.Unlock()
	log.Printf("UPDATE iamToken, expires: %v", t.ExpiresAt)
	return nil
}

func (cc *CloudConfig) SrartAutoUpdateCloudConfig() error {
	if cc.ticker != nil {
		return ErrTickerAlreadyRunning
	}
	err := cc.UpdateCloudConfig()
	if err != nil {
		return err
	}
	cc.ticker = time.NewTicker(cc.deltaTime)
	go func() {
		for range cc.ticker.C {
			err := cc.UpdateCloudConfig()
			if err != nil {
				log.Println(err)
			}
		}
	}()
	return nil
}

func (cc *CloudConfig) GetFolderId() (string, error) {
	if len(cc.folderId) == 0 {
		return "", ErrEmptyFolderIdToken
	}
	return cc.folderId, nil
}

func (cc *CloudConfig) GetIAMToken() (string, error) {
	cc.mu.RLock()
	token := cc.iamToken
	exp := cc.expired
	cc.mu.RUnlock()
	if time.Now().Add(timeLimitBeforeExpires).After(exp) {
		err := cc.UpdateCloudConfig()
		if err != nil {
			return "", fmt.Errorf("error: GetIAMToken UpdateCloudConfig: %w", err)
		}
		cc.mu.RLock()
		token = cc.iamToken
		cc.mu.RUnlock()
	}
	return token, nil
}

func (cc *CloudConfig) SetTime(time time.Duration) error {
	if cc.ticker != nil {
		return ErrTickerAlreadyRunning
	}
	cc.deltaTime = time
	return nil
}

func (cc *CloudConfig) Stop() error {
	if cc.ticker == nil {
		return ErrTickerNotRunning
	}
	cc.ticker.Stop()
	cc.ticker = nil
	return nil
}
