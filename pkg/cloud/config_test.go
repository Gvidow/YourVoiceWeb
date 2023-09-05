package cloud

import (
	"testing"
)

func TestRerutnedError(t *testing.T) {
	// 	tests := []struct {
	// 		name     string
	// 		config   *CloudConfig
	// 		do       func(*CloudConfig) error
	// 		result   error
	// 		excepted error
	// 	}{
	// 		{
	// 			"GetEmptyFolderId", NewCloudConfig("oauthtoken", ""),
	// 			func(cfg *CloudConfig) error {
	// 				_, err := cfg.GetFolderId()
	// 				return err
	// 			},
	// 			nil, ErrEmptyFolderIdToken,
	// 		},
	// 		{
	// 			"StopNotRunning", NewCloudConfig("oauthtoken", "folderid"),
	// 			func(cfg *CloudConfig) error {
	// 				return cfg.Stop()
	// 			},
	// 			nil, ErrTickerNotRunning,
	// 		},
	// 		{
	// 			"SetDeltaTimeForNotRunningTicker", NewCloudConfig("oauthtoken", "folderid"),
	// 			func(cfg *CloudConfig) error {
	// 				err := cfg.SetTime(10 * time.Second)
	// 				return err
	// 			},
	// 			nil, nil,
	// 		},
	// 	}

	//	for _, test := range tests {
	//		test.result = test.do(test.config)
	//		assert.Equal(t, test.excepted, test.result, "fail test"+test.name)
	//	}
}
