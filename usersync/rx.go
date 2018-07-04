package usersync

import (
	"fmt"
	"net/url"
	"strings"
)

func NewRxSyncer(externalURL string) Usersyncer {
	externalURL = strings.TrimRight(externalURL, "/")
	redirectURL := fmt.Sprintf("%s/setuid?bidder=rx&uid=${UID}", externalURL)

	return &syncer{
		familyName: "rx",
		syncInfo: &UsersyncInfo{
			URL:         fmt.Sprintf("https://rx.com/sync/prebid?r=%s", url.QueryEscape(redirectURL)),
			Type:        "redirect",
			SupportCORS: false,
		},
	}
}
