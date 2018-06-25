package prometheusmetrics

import (
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"

	"github.com/prebid/prebid-server/config"
	"github.com/prebid/prebid-server/openrtb_ext"
	"github.com/prebid/prebid-server/pbsmetrics"
)

var gaugeValueRegexp = regexp.MustCompile("gauge:<value:([0-9]+) >")
var counterValueRegexp = regexp.MustCompile("counter:<value:([0-9]+) >")
var histogramValueRegexp = regexp.MustCompile("histogram:<sample_count:([0-9]+)")

func TestConnectionMetrics(t *testing.T) {
	proMetrics := newTestMetricsEngine().(*Metrics)

	metricConn := dto.Metric{}
	metricConnErrA := dto.Metric{}
	metricConnErrC := dto.Metric{}
	proMetrics.RecordConnectionAccept(true)
	proMetrics.RecordConnectionAccept(true)
	proMetrics.RecordConnectionClose(true)
	proMetrics.RecordConnectionAccept(true)
	proMetrics.RecordConnectionAccept(false)
	proMetrics.RecordConnectionClose(false)

	proMetrics.connCounter.Write(&metricConn)
	proMetrics.connError.WithLabelValues("accept_error").Write(&metricConnErrA)
	proMetrics.connError.WithLabelValues("close_error").Write(&metricConnErrC)

	assertGaugeValue(t, "connCounter", &metricConn, 2)
	assertCounterValue(t, "connError[accept_error]", &metricConnErrA, 1)
	assertCounterValue(t, "connError[close_error]", &metricConnErrC, 1)
}

func TestRequestMetrics(t *testing.T) {
	proMetrics := newTestMetricsEngine().(*Metrics)

	metrics0 := dto.Metric{}
	metrics1 := dto.Metric{}
	metrics2 := dto.Metric{}
	metrics3 := dto.Metric{}

	proMetrics.RecordRequest(labels[0])
	proMetrics.RecordRequest(labels[1])
	proMetrics.RecordRequest(labels[0])
	proMetrics.RecordRequest(labels[2])
	proMetrics.RecordRequest(labels[0])
	proMetrics.RecordRequest(labels[2])

	proMetrics.requests.With(resolveLabels(labels[0])).Write(&metrics0)
	proMetrics.requests.With(resolveLabels(labels[1])).Write(&metrics1)
	proMetrics.requests.With(resolveLabels(labels[2])).Write(&metrics2)
	proMetrics.requests.With(resolveLabels(labels[3])).Write(&metrics3)

	assertCounterValue(t, "requests[0]", &metrics0, 3)
	assertCounterValue(t, "requests[1]", &metrics1, 1)
	assertCounterValue(t, "requests[2]", &metrics2, 2)
	assertCounterValue(t, "requests[3]", &metrics3, 0)
}

func TestImpMetrics(t *testing.T) {
	proMetrics := newTestMetricsEngine().(*Metrics)

	metrics0 := dto.Metric{}
	metrics1 := dto.Metric{}
	metrics2 := dto.Metric{}
	metrics3 := dto.Metric{}

	proMetrics.RecordImps(labels[0], 1)
	proMetrics.RecordImps(labels[1], 5)
	proMetrics.RecordImps(labels[0], 2)
	proMetrics.RecordImps(labels[2], 2)
	proMetrics.RecordImps(labels[0], 7)
	proMetrics.RecordImps(labels[2], 1)

	proMetrics.imps.With(resolveLabels(labels[0])).Write(&metrics0)
	proMetrics.imps.With(resolveLabels(labels[1])).Write(&metrics1)
	proMetrics.imps.With(resolveLabels(labels[2])).Write(&metrics2)
	proMetrics.imps.With(resolveLabels(labels[3])).Write(&metrics3)

	assertCounterValue(t, "imps_requested[0]", &metrics0, 10)
	assertCounterValue(t, "imps_requested[1]", &metrics1, 5)
	assertCounterValue(t, "imps_requested[2]", &metrics2, 3)
	assertCounterValue(t, "imps_requested[3]", &metrics3, 0)
}

func TestTimerMetrics(t *testing.T) {
	proMetrics := newTestMetricsEngine().(*Metrics)

	metrics0 := dto.Metric{}
	metrics1 := dto.Metric{}
	metrics2 := dto.Metric{}
	metrics3 := dto.Metric{}

	proMetrics.RecordRequestTime(labels[0], 120*time.Millisecond)
	proMetrics.RecordRequestTime(labels[1], 85*time.Millisecond)
	proMetrics.RecordRequestTime(labels[0], 220*time.Millisecond)
	proMetrics.RecordRequestTime(labels[2], 250*time.Millisecond)
	proMetrics.RecordRequestTime(labels[0], 90*time.Millisecond)
	proMetrics.RecordRequestTime(labels[2], 180*time.Millisecond)

	proMetrics.reqTimer.With(resolveLabels(labels[0])).Write(&metrics0)
	proMetrics.reqTimer.With(resolveLabels(labels[1])).Write(&metrics1)
	proMetrics.reqTimer.With(resolveLabels(labels[2])).Write(&metrics2)
	proMetrics.reqTimer.With(resolveLabels(labels[3])).Write(&metrics3)

	assertHistogramValue(t, "request_time[0]", &metrics0, 3)
	assertHistogramValue(t, "request_time[1]", &metrics1, 1)
	assertHistogramValue(t, "request_time[2]", &metrics2, 2)
	assertHistogramValue(t, "request_time[3]", &metrics3, 0)
}

func TestAdapterRequestMetrics(t *testing.T) {
	proMetrics := newTestMetricsEngine().(*Metrics)

	metrics0 := dto.Metric{}
	metrics1 := dto.Metric{}
	metrics2 := dto.Metric{}
	metrics3 := dto.Metric{}

	proMetrics.RecordAdapterRequest(adaptLabels[0])
	proMetrics.RecordAdapterRequest(adaptLabels[1])
	proMetrics.RecordAdapterRequest(adaptLabels[0])
	proMetrics.RecordAdapterRequest(adaptLabels[2])
	proMetrics.RecordAdapterRequest(adaptLabels[0])
	proMetrics.RecordAdapterRequest(adaptLabels[2])

	proMetrics.adaptRequests.With(resolveAdapterLabels(adaptLabels[0])).Write(&metrics0)
	proMetrics.adaptRequests.With(resolveAdapterLabels(adaptLabels[1])).Write(&metrics1)
	proMetrics.adaptRequests.With(resolveAdapterLabels(adaptLabels[2])).Write(&metrics2)
	proMetrics.adaptRequests.With(resolveAdapterLabels(adaptLabels[3])).Write(&metrics3)

	assertCounterValue(t, "adapter_requests[0]", &metrics0, 3)
	assertCounterValue(t, "adapter_requests[1]", &metrics1, 1)
	assertCounterValue(t, "adapter_requests[2]", &metrics2, 2)
	assertCounterValue(t, "adapter_requests[3]", &metrics3, 0)
}

func TestAdapterBidsMetrics(t *testing.T) {
	proMetrics := newTestMetricsEngine().(*Metrics)

	metrics0 := dto.Metric{}
	metrics1 := dto.Metric{}
	metrics2 := dto.Metric{}
	metrics3 := dto.Metric{}

	proMetrics.RecordAdapterBidReceived(adaptLabels[0], openrtb_ext.BidTypeBanner, true)
	proMetrics.RecordAdapterBidReceived(adaptLabels[1], openrtb_ext.BidTypeBanner, false)
	proMetrics.RecordAdapterBidReceived(adaptLabels[0], openrtb_ext.BidTypeBanner, true)
	proMetrics.RecordAdapterBidReceived(adaptLabels[2], openrtb_ext.BidTypeVideo, true)
	proMetrics.RecordAdapterBidReceived(adaptLabels[0], openrtb_ext.BidTypeBanner, true)
	proMetrics.RecordAdapterBidReceived(adaptLabels[2], openrtb_ext.BidTypeVideo, true)

	proMetrics.adaptBids.With(resolveBidLabels(adaptLabels[0], openrtb_ext.BidTypeBanner, true)).Write(&metrics0)
	proMetrics.adaptBids.With(resolveBidLabels(adaptLabels[1], openrtb_ext.BidTypeBanner, false)).Write(&metrics1)
	proMetrics.adaptBids.With(resolveBidLabels(adaptLabels[2], openrtb_ext.BidTypeVideo, true)).Write(&metrics2)
	proMetrics.adaptBids.With(resolveBidLabels(adaptLabels[3], openrtb_ext.BidTypeNative, false)).Write(&metrics3)

	assertCounterValue(t, "adapter_bids_recieved[0]", &metrics0, 3)
	assertCounterValue(t, "adapter_bids_recieved[1]", &metrics1, 1)
	assertCounterValue(t, "adapter_bids_recieved[2]", &metrics2, 2)
	assertCounterValue(t, "adapter_bids_recieved[3]", &metrics3, 0)
}

func TestAdapterPriceMetrics(t *testing.T) {
	proMetrics := newTestMetricsEngine().(*Metrics)

	metrics0 := dto.Metric{}
	metrics1 := dto.Metric{}
	metrics2 := dto.Metric{}
	metrics3 := dto.Metric{}

	proMetrics.RecordAdapterPrice(adaptLabels[0], 0.12)
	proMetrics.RecordAdapterPrice(adaptLabels[1], 2.35)
	proMetrics.RecordAdapterPrice(adaptLabels[0], 17.65)
	proMetrics.RecordAdapterPrice(adaptLabels[2], 3.23333)
	proMetrics.RecordAdapterPrice(adaptLabels[0], 6.564)
	proMetrics.RecordAdapterPrice(adaptLabels[2], 0.03)

	proMetrics.adaptPrices.With(resolveAdapterLabels(adaptLabels[0])).Write(&metrics0)
	proMetrics.adaptPrices.With(resolveAdapterLabels(adaptLabels[1])).Write(&metrics1)
	proMetrics.adaptPrices.With(resolveAdapterLabels(adaptLabels[2])).Write(&metrics2)
	proMetrics.adaptPrices.With(resolveAdapterLabels(adaptLabels[3])).Write(&metrics3)

	assertHistogramValue(t, "adapter_prices[0]", &metrics0, 3)
	assertHistogramValue(t, "adapter_prices[1]", &metrics1, 1)
	assertHistogramValue(t, "adapter_prices[2]", &metrics2, 2)
	assertHistogramValue(t, "adapter_prices[3]", &metrics3, 0)

}

func TestAdapterTimeMetrics(t *testing.T) {
	proMetrics := newTestMetricsEngine().(*Metrics)

	metrics0 := dto.Metric{}
	metrics1 := dto.Metric{}
	metrics2 := dto.Metric{}
	metrics3 := dto.Metric{}

	proMetrics.RecordAdapterTime(adaptLabels[0], 85*time.Millisecond)
	proMetrics.RecordAdapterTime(adaptLabels[1], 235*time.Millisecond)
	proMetrics.RecordAdapterTime(adaptLabels[0], 177*time.Millisecond)
	proMetrics.RecordAdapterTime(adaptLabels[2], 323*time.Millisecond)
	proMetrics.RecordAdapterTime(adaptLabels[0], 664*time.Millisecond)
	proMetrics.RecordAdapterTime(adaptLabels[2], 33*time.Millisecond)

	proMetrics.adaptTimer.With(resolveAdapterLabels(adaptLabels[0])).Write(&metrics0)
	proMetrics.adaptTimer.With(resolveAdapterLabels(adaptLabels[1])).Write(&metrics1)
	proMetrics.adaptTimer.With(resolveAdapterLabels(adaptLabels[2])).Write(&metrics2)
	proMetrics.adaptTimer.With(resolveAdapterLabels(adaptLabels[3])).Write(&metrics3)

	assertHistogramValue(t, "adapter_time[0]", &metrics0, 3)
	assertHistogramValue(t, "adapter_time[1]", &metrics1, 1)
	assertHistogramValue(t, "adapter_time[2]", &metrics2, 2)
	assertHistogramValue(t, "adapter_time[3]", &metrics3, 0)

}

func TestCookieMetrics(t *testing.T) {
	proMetrics := newTestMetricsEngine().(*Metrics)

	metrics0 := dto.Metric{}

	proMetrics.RecordCookieSync(labels[0])
	proMetrics.RecordCookieSync(labels[1])
	proMetrics.RecordCookieSync(labels[0])
	proMetrics.RecordCookieSync(labels[2])
	proMetrics.RecordCookieSync(labels[0])
	proMetrics.RecordCookieSync(labels[2])

	proMetrics.cookieSync.Write(&metrics0)

	assertCounterValue(t, "cookie_sync_requests", &metrics0, 6)
}

func TestUserMetrics(t *testing.T) {
	proMetrics := newTestMetricsEngine().(*Metrics)

	metrics0 := dto.Metric{}
	metrics1 := dto.Metric{}
	metrics2 := dto.Metric{}
	metrics3 := dto.Metric{}

	proMetrics.RecordUserIDSet(userLabels[0])
	proMetrics.RecordUserIDSet(userLabels[1])
	proMetrics.RecordUserIDSet(userLabels[0])
	proMetrics.RecordUserIDSet(userLabels[2])
	proMetrics.RecordUserIDSet(userLabels[0])
	proMetrics.RecordUserIDSet(userLabels[2])

	proMetrics.userID.With(resolveUserSyncLabels(userLabels[0])).Write(&metrics0)
	proMetrics.userID.With(resolveUserSyncLabels(userLabels[1])).Write(&metrics1)
	proMetrics.userID.With(resolveUserSyncLabels(userLabels[2])).Write(&metrics2)
	proMetrics.userID.With(resolveUserSyncLabels(userLabels[3])).Write(&metrics3)

	assertCounterValue(t, "usersync[0]", &metrics0, 3)
	assertCounterValue(t, "usersync[1]", &metrics1, 1)
	assertCounterValue(t, "usersync[2]", &metrics2, 2)
	assertCounterValue(t, "usersync[3]", &metrics3, 0)
}

func TestMetricsExist(t *testing.T) {
	// Initialize the metrics engine -> register the metrics to prometheus
	metrics := newTestMetricsEngine()

	// Get at the underlying metrics object
	proMetrics := metrics.(*Metrics)

	if err := proMetrics.Registry.Register(prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "prebid",
		Name:      "active_connections",
		Help:      "Current number of active (open) connections.",
	})); err == nil {
		t.Error("connCounter not registered")
	}
}

func newTestMetricsEngine() pbsmetrics.MetricsEngine {
	return NewMetrics(config.PrometheusMetrics{
		Endpoint:  "/metrics",
		Port:      8080,
		Namespace: "prebid",
		Subsystem: "",
	})
}

var labels = []pbsmetrics.Labels{
	{
		Source:        pbsmetrics.DemandWeb,
		RType:         pbsmetrics.ReqTypeLegacy,
		PubID:         "Pub1",
		Browser:       pbsmetrics.BrowserOther,
		CookieFlag:    pbsmetrics.CookieFlagYes,
		RequestStatus: pbsmetrics.RequestStatusOK,
	},
	{
		Source:        pbsmetrics.DemandWeb,
		RType:         pbsmetrics.ReqTypeLegacy,
		PubID:         "Pub1",
		Browser:       pbsmetrics.BrowserSafari,
		CookieFlag:    pbsmetrics.CookieFlagYes,
		RequestStatus: pbsmetrics.RequestStatusOK,
	},
	{
		Source:        pbsmetrics.DemandApp,
		RType:         pbsmetrics.ReqTypeORTB2,
		PubID:         "Pub2",
		Browser:       pbsmetrics.BrowserOther,
		CookieFlag:    pbsmetrics.CookieFlagNo,
		RequestStatus: pbsmetrics.RequestStatusOK,
	},
	{
		Source:        pbsmetrics.DemandUnknown,
		RType:         pbsmetrics.ReqTypeORTB2,
		PubID:         "Pub3",
		Browser:       pbsmetrics.BrowserOther,
		CookieFlag:    pbsmetrics.CookieFlagUnknown,
		RequestStatus: pbsmetrics.RequestStatusBadInput,
	},
}

var adaptLabels = []pbsmetrics.AdapterLabels{
	{
		Source:        pbsmetrics.DemandWeb,
		RType:         pbsmetrics.ReqTypeLegacy,
		Adapter:       openrtb_ext.BidderAppnexus,
		PubID:         "Pub1",
		Browser:       pbsmetrics.BrowserOther,
		CookieFlag:    pbsmetrics.CookieFlagYes,
		AdapterStatus: pbsmetrics.AdapterStatusOK,
	},
	{
		Source:        pbsmetrics.DemandWeb,
		RType:         pbsmetrics.ReqTypeLegacy,
		Adapter:       openrtb_ext.BidderEPlanning,
		PubID:         "Pub1",
		Browser:       pbsmetrics.BrowserSafari,
		CookieFlag:    pbsmetrics.CookieFlagYes,
		AdapterStatus: pbsmetrics.AdapterStatusErr,
	},
	{
		Source:        pbsmetrics.DemandApp,
		RType:         pbsmetrics.ReqTypeORTB2,
		Adapter:       openrtb_ext.BidderIndex,
		PubID:         "Pub2",
		Browser:       pbsmetrics.BrowserOther,
		CookieFlag:    pbsmetrics.CookieFlagNo,
		AdapterStatus: pbsmetrics.AdapterStatusOK,
	},
	{
		Source:        pbsmetrics.DemandUnknown,
		RType:         pbsmetrics.ReqTypeORTB2,
		Adapter:       openrtb_ext.BidderAppnexus,
		PubID:         "Pub3",
		Browser:       pbsmetrics.BrowserOther,
		CookieFlag:    pbsmetrics.CookieFlagUnknown,
		AdapterStatus: pbsmetrics.AdapterStatusOK,
	},
}

var userLabels = []pbsmetrics.UserLabels{
	{
		Action: pbsmetrics.RequestActionSet,
		Bidder: openrtb_ext.BidderAppnexus,
	},
	{
		Action: pbsmetrics.RequestActionGDPR,
		Bidder: openrtb_ext.BidderAppnexus,
	},
	{
		Action: pbsmetrics.RequestActionSet,
		Bidder: openrtb_ext.BidderRubicon,
	},
	{
		Action: pbsmetrics.RequestActionOptOut,
		Bidder: openrtb_ext.BidderOpenx,
	},
}

func assertMetricValue(t *testing.T, name string, m *dto.Metric, expected string) {
	v := m.String()
	if v != expected {
		t.Errorf("Bad value for metric %s: expected=\"%s\", found=\"%s\"", name, expected, v)
	}
}

func assertGaugeValue(t *testing.T, name string, m *dto.Metric, expected int) {
	v, err := strconv.Atoi(gaugeValueRegexp.FindStringSubmatch(m.String())[1])
	if err != nil {
		t.Errorf("Could not extract the value for metric %s. (output was %s, error was %v)", name, m.String(), err)
	}
	if v != expected {
		t.Errorf("Bad value for metric %s: expected=\"%d\", found=\"%d\"", name, expected, v)
	}
}

func assertCounterValue(t *testing.T, name string, m *dto.Metric, expected int) {
	v, err := strconv.Atoi(counterValueRegexp.FindStringSubmatch(m.String())[1])
	if err != nil {
		t.Errorf("Could not extract the value for metric %s. (output was %s, error was %v)", name, m.String(), err)
	}
	if v != expected {
		t.Errorf("Bad value for metric %s: expected=\"%d\", found=\"%d\"", name, expected, v)
	}
}

func assertHistogramValue(t *testing.T, name string, m *dto.Metric, expected int) {
	v, err := strconv.Atoi(histogramValueRegexp.FindStringSubmatch(m.String())[1])
	if err != nil {
		t.Errorf("Could not extract the value for metric %s. (output was %s, error was %v)", name, m.String(), err)
	}
	if v != expected {
		t.Errorf("Bad value for metric %s: expected=\"%d\", found=\"%d\"", name, expected, v)
	}
}