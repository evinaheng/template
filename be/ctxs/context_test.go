package ctxs

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextConversion(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "https://site.com", nil)
	r.ParseForm()

	ctx, cancel := context.WithCancel(GetAllContextFromRequest(r))
	defer cancel()

	var ck contextKey

	ck.getInt(ctx)
}

func TestContext(t *testing.T) {

	r := httptest.NewRequest(http.MethodPost, "https://site.com", nil)
	r.ParseForm()

	r.Form.Add("user_id", "999")
	r.Form.Add("uid", "999")
	r.Form.Add("email", "foo@email.com")

	r.Form.Add("timeout", "3s")
	r.Form.Add("id", "31337")
	r.Form.Add("q", "Jakarta")
	r.Form.Add("voucher_code", "0xcafebabe")
	r.Form.Add("cart_id", "0xdeadbeef")
	r.Form.Add("order_id", "100")
	r.Form.Add("pdf", "abc")
	r.Form.Add("lang", "id-id")
	r.Form.Add("status", "100")
	r.Form.Add("country_id", "ID")
	r.Form.Add("status_bulk", "100,200")
	r.Form.Add("invoice_id", "P22-ZUX11-7T4H")
	r.Form.Add("notif_template", "123")
	r.Form.Add("depart_code", "JKTA")
	r.Form.Add("arrival_code", "DPS")
	r.Form.Add("year", "2017")
	r.Form.Add("month", "02")
	r.Form.Add("language_id", "2")
	r.Form.Add("did", "4")
	r.Form.Set("page", "3")
	r.Form.Set("per_page", "4")
	r.Form.Set("instance_id", "1")
	r.Header.Set("User-Agent", "Firefox")
	r.Header.Set("Idempotency-Key", "key")
	r.Header.Set("X-Tkpd-UserId", "111")
	r.Header.Set("Language", "id-id")
	r.Header.Set("X-Source", "Intools")
	r.Header.Set("os-type", "1")
	r.Header.Set("Date", "01-01-2010 11:12:00")
	r.Header.Set("X-Device", "ios")
	r.Header.Set("X-App-Version", "v1.0")

	ctx, cancel := context.WithCancel(GetAllContextFromRequest(r))
	defer cancel()

	userID, ok := UserIDFromContext(ctx)
	assert.Equal(t, int64(999), userID)
	assert.True(t, ok)

	UID, ok := UIDFromContext(ctx)
	assert.Equal(t, int64(999), UID)
	assert.True(t, ok)

	q, ok := QFromContext(ctx)
	assert.Equal(t, "Jakarta", q)
	assert.True(t, ok)

	email, ok := EmailFromContext(ctx)
	assert.Equal(t, "foo@email.com", email)
	assert.True(t, ok)

	successInt, ok := UserIDKey.getInt(ctx)
	assert.Equal(t, 999, successInt)

	failedInt, ok := EmailKey.getInt(ctx)
	assert.Equal(t, 0, failedInt)
	assert.False(t, ok)

	failedInt64, ok := EmailKey.getInt64(ctx)
	assert.Equal(t, int64(0), failedInt64)
	assert.False(t, ok)

	id, ok := IDFromContext(ctx)
	assert.Equal(t, "31337", id)
	assert.True(t, ok)

	orderID, ok := OrderIDFromContext(ctx)
	assert.Equal(t, int64(100), orderID)
	assert.True(t, ok)

	pdf, ok := PDFFromContext(ctx)
	assert.Equal(t, "abc", pdf)
	assert.True(t, ok)

	lang, ok := LangFromContext(ctx)
	assert.Equal(t, "id-id", lang)
	assert.True(t, ok)

	voucherCode, ok := VoucherCodeFromContext(ctx)
	assert.Equal(t, "0xcafebabe", voucherCode)
	assert.True(t, ok)

	status, ok := StatusFromContext(ctx)
	assert.Equal(t, 100, status)
	assert.True(t, ok)

	countryID, ok := CountryIDFromContext(ctx)
	assert.Equal(t, "ID", countryID)
	assert.True(t, ok)

	statusBulk, ok := StatusBulkFromContext(ctx)
	assert.Equal(t, "100,200", statusBulk)
	assert.True(t, ok)

	cartID, ok := CartIDFromContext(ctx)
	assert.Equal(t, "0xdeadbeef", cartID)
	assert.True(t, ok)

	userAgent, ok := UserAgentFromContext(ctx)
	assert.Equal(t, "Firefox", userAgent)
	assert.True(t, ok)

	idempotencyKey, ok := IdempotencyKeyFromContext(ctx)
	assert.Equal(t, "key", idempotencyKey)
	assert.True(t, ok)

	xUserID, ok := XUserIDHeaderKeyFromContext(ctx)
	assert.Equal(t, int64(111), xUserID)
	assert.True(t, ok)

	osType, ok := OSTypeFromContext(ctx)
	assert.Equal(t, int(1), osType)
	assert.True(t, ok)

	language, ok := LanguageHeaderKeyFromContext(ctx)
	assert.Equal(t, "id-id", language)
	assert.True(t, ok)

	source, ok := XSourceHeaderKeyFromContext(ctx)
	assert.Equal(t, "Intools", source)

	ipAddress, ok := IPAddressFromContext(ctx)
	assert.NotNil(t, ipAddress)
	assert.True(t, ok)

	invoiceID, ok := InvoiceIDFromContext(ctx)
	assert.Equal(t, "P22-ZUX11-7T4H", invoiceID)
	assert.True(t, ok)

	notifTemplate, ok := NotifTemplateIDFromContext(ctx)
	assert.Equal(t, 123, notifTemplate)
	assert.True(t, ok)

	departCode, ok := DepartCodeFromContext(ctx)
	assert.Equal(t, "JKTA", departCode)
	assert.True(t, ok)

	arrivalCode, ok := ArrivalCodeFromContext(ctx)
	assert.Equal(t, "DPS", arrivalCode)
	assert.True(t, ok)

	year, ok := YearFromContext(ctx)
	assert.Equal(t, int(2017), year)
	assert.True(t, ok)

	month, ok := MonthFromContext(ctx)
	assert.Equal(t, "02", month)
	assert.True(t, ok)

	page, ok := PageFromContext(ctx)
	assert.Equal(t, 3, page)
	assert.True(t, ok)

	instanceID, ok := InstanceIDFromContext(ctx)
	assert.Equal(t, 1, instanceID)
	assert.True(t, ok)

	perPage, ok := PerPageFromContext(ctx)
	assert.Equal(t, 4, perPage)
	assert.True(t, ok)

	languageID, ok := LanguageIDFromContext(ctx)
	assert.Equal(t, 2, languageID)
	assert.True(t, ok)

	deviceID, ok := DeviceIDFromContext(ctx)
	assert.Equal(t, 4, deviceID)
	assert.True(t, ok)

	clientDate, ok := DateKeyFromContext(ctx)
	assert.Equal(t, "01-01-2010 11:12:00", clientDate)
	assert.True(t, ok)

	XAppVersion, ok := XAppVersionFromContext(ctx)
	assert.Equal(t, "v1.0", XAppVersion)
	assert.True(t, ok)

	XDevice, ok := XDeviceFromContext(ctx)
	assert.Equal(t, "ios", XDevice)
	assert.True(t, ok)
}

func TestContextNotOK(t *testing.T) {

	r := httptest.NewRequest(http.MethodPost, "https://site.com", nil)
	r.ParseForm()

	r.Form.Add("os_type", "X")

	ctx, cancel := context.WithCancel(GetAllContextFromRequest(r))
	defer cancel()

	_, ok := OSTypeFromContext(ctx)
	assert.False(t, ok)

}

func TestJSONContextSuccess(t *testing.T) {

	jsonRequest := []byte(`{ "foo" : "value" }`)

	jsonStruct := struct {
		Foo string `json:"foo"`
	}{}

	r := httptest.NewRequest(http.MethodPost, "https://site.com", bytes.NewBuffer(jsonRequest))

	ctx, cancel := context.WithCancel(GetAllContextFromRequest(r))
	defer cancel()
	assert.False(t, JSONFromContext(ctx, &jsonStruct))

	r.Header.Set("Content-Type", "application/json")
	ctx2, cancel2 := context.WithCancel(GetAllContextFromRequest(r))
	defer cancel2()
	assert.True(t, JSONFromContext(ctx2, &jsonStruct))
}

func TestJSONContextFailed(t *testing.T) {

	r := httptest.NewRequest(http.MethodPost, "https://site.com", nil)

	r.Header.Set("Content-Type", "application/json")
	ctx, cancel := context.WithCancel(GetAllContextFromRequest(r))
	defer cancel()
	assert.False(t, JSONFromContext(ctx, nil))

}

func TestXUserIDHeaderKeyFromContextNil(t *testing.T) {

	r := httptest.NewRequest(http.MethodPost, "https://site.com", nil)
	r.ParseForm()

	xUserID, ok := XUserIDHeaderKeyFromContext(r.Context())
	assert.Equal(t, int64(0), xUserID)
	assert.False(t, ok)

}
