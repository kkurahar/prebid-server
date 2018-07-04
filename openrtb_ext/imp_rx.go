package openrtb_ext

// ExtImpRx defines the contract for bidrequest.imp[i].ext.rx
type ExtImpRx struct {
	AccountId int `json:"accountId"`
	SiteId    int `json:"siteId"`
	ZoneId    int `json:"zoneId"`
}
