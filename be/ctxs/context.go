package ctxs

// Context package

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
)

type contextKey int

// ContextCancelledError by user
const ContextCancelledError = "context canceled"

const (
	IPAddress contextKey = iota
	UserAgentKey
	IdempotencyKey
	XTkpdUserIDHeaderKey
	LanguageHeaderKey
	JSONContent
	QKey
	XSourceHeaderKey
	LanguageIDKey
	DeviceIDKey
	OSTypeKey
	DateKey
	XDeviceKey
	XAppVersionKey
	InstanceIDKey
	RouterVarContextKey

	QQS             = "q"
	DeviceIDQS      = "did"
	InstanceIDKeyQS = "instance_id"
)

var (
	ContextQueryStringList = map[contextKey]string{
		QKey: QQS,
	}
)

// Extract some value from *http.Request and append it into context
// If the value is not valid or there's any error, return the parent context as result
// add another exportable func to fetch the desired value from context
// make sure that req http.Request pointer is not nil
func (ck contextKey) getString(ctx context.Context) (string, bool) {
	str, ok := ctx.Value(ck).(string)
	return str, ok
}
func (ck contextKey) getInt(ctx context.Context) (int, bool) {
	str, _ := ck.getString(ctx)
	atoi, err := strconv.Atoi(str)
	if err != nil {
		return 0, false
	}
	return atoi, true
}
func (ck contextKey) getInt64(ctx context.Context) (int64, bool) {
	str, _ := ck.getString(ctx)
	atoi, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, false
	}
	return atoi, true
}

// Save all values as string
// convert value into desired one on get context func
func AllQSToContext(ctx context.Context, req *http.Request) context.Context {
	for key, qs := range ContextQueryStringList {
		ctx = context.WithValue(ctx, key, req.FormValue(qs))
	}
	return ctx
}
func QFromContext(ctx context.Context) (string, bool) {
	return QKey.getString(ctx)
}

func XSourceHeaderKeyFromContext(ctx context.Context) (string, bool) {
	return XSourceHeaderKey.getString(ctx)
}

func DeviceIDFromContext(ctx context.Context) (int, bool) {
	return DeviceIDKey.getInt(ctx)
}

func InstanceIDFromContext(ctx context.Context) (int, bool) {
	return InstanceIDKey.getInt(ctx)
}

// ============
// USER AGENT
// ============
func UserAgentToContext(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, UserAgentKey, req.UserAgent())
}

func UserAgentFromContext(ctx context.Context) (string, bool) {
	ua, ok := ctx.Value(UserAgentKey).(string)
	return ua, ok
}

// ============
// IP ADDRESS
// ============
func IPAddressToContext(ctx context.Context, req *http.Request) context.Context {
	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	return context.WithValue(ctx, IPAddress, net.ParseIP(ip))
}

func IPAddressFromContext(ctx context.Context) (net.IP, bool) {
	ip, ok := ctx.Value(IPAddress).(net.IP)
	return ip, ok
}

// =============
// HEADER DATA
// =============

func DateKeyFromContext(ctx context.Context) (string, bool) {
	ua, ok := ctx.Value(DateKey).(string)
	return ua, ok
}

func XDeviceToContext(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, XDeviceKey, req.Header.Get("X-Device"))
}

func XDeviceFromContext(ctx context.Context) (string, bool) {
	ua, ok := ctx.Value(XDeviceKey).(string)
	return ua, ok
}

func XAppVersionToContext(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, XAppVersionKey, req.Header.Get("X-App-Version"))
}

func XAppVersionFromContext(ctx context.Context) (string, bool) {
	ua, ok := ctx.Value(XAppVersionKey).(string)
	return ua, ok
}

func IdempotencyKeyToContext(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, IdempotencyKey, req.Header.Get("idempotency-key"))
}

func IdempotencyKeyFromContext(ctx context.Context) (string, bool) {
	ua, ok := ctx.Value(IdempotencyKey).(string)
	return ua, ok
}

func OSTypeContextToContext(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, OSTypeKey, req.Header.Get("os-type"))
}

func OSTypeFromContext(ctx context.Context) (int, bool) {
	ua, _ := ctx.Value(OSTypeKey).(string)
	atoi, err := strconv.Atoi(ua)
	if err != nil {
		return 0, false
	}
	return atoi, true
}

func XUserIDHeaderKeyToContext(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, XTkpdUserIDHeaderKey, req.Header.Get("x-tkpd-userid"))
}

func XUserIDHeaderKeyFromContext(ctx context.Context) (int64, bool) {
	if ctx.Value(XTkpdUserIDHeaderKey) == nil {
		return 0, false
	}
	ua, _ := strconv.ParseInt(ctx.Value(XTkpdUserIDHeaderKey).(string), 10, 64)
	return ua, true
}

func LanguageHeaderKeyToContext(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, LanguageHeaderKey, req.Header.Get("language"))
}

func LanguageHeaderKeyFromContext(ctx context.Context) (string, bool) {
	la, _ := ctx.Value(LanguageHeaderKey).(string)
	return la, true
}

func XSourceHeaderKeyToContext(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, XSourceHeaderKey, req.Header.Get("X-Source"))
}

// JSONToContext value
func JSONToContext(ctx context.Context, req *http.Request) context.Context {

	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)

	return context.WithValue(ctx, JSONContent, body)
}

// JSONFromContext value
func JSONFromContext(ctx context.Context, i interface{}) bool {
	b, ok := ctx.Value(JSONContent).([]byte)
	if !ok {
		return false
	}

	if err := json.Unmarshal(b, i); err != nil {
		log.Println("func JSONFromContext", err)
		return false
	}
	return true
}

// Function to fetch context from *http.Request
// set timeout if QS timeout is parseable into time.Duration
// implement the save value to context funcs as optional parameter
func GetContextFromRequest(req *http.Request, opt ...func(context.Context, *http.Request) context.Context) context.Context {
	ctx := req.Context()

	for _, f := range opt {
		ctx = f(ctx, req)
	}

	return ctx
}

func GetAllContextFromRequest(req *http.Request) context.Context {
	contentType := req.Header.Get("Content-Type")
	if contentType == "application/vnd.api+json" || contentType == "application/json" {
		return GetContextFromRequest(req, UserAgentToContext, IPAddressToContext, JSONToContext, IdempotencyKeyToContext, XUserIDHeaderKeyToContext, LanguageHeaderKeyToContext, XSourceHeaderKeyToContext, OSTypeContextToContext, DateKeyToContext, XDeviceToContext, XAppVersionToContext)
	}
	return GetContextFromRequest(req, AllQSToContext, UserAgentToContext, IPAddressToContext, IdempotencyKeyToContext, XUserIDHeaderKeyToContext, LanguageHeaderKeyToContext, XSourceHeaderKeyToContext, OSTypeContextToContext, DateKeyToContext, XDeviceToContext, XAppVersionToContext)
}
