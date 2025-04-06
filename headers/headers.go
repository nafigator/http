package headers

import "net/textproto"

// HTTP headers RFC2616 https://datatracker.ietf.org/doc/html/rfc2616#section-14.
const (
	Accept                        = "Accept"
	AcceptCharset                 = "Accept-Charset"
	AcceptEncoding                = "Accept-Encoding"
	AcceptLanguage                = "Accept-Language"
	AcceptRanges                  = "Accept-Ranges"
	Age                           = "Age"
	Allow                         = "Allow"
	Authorization                 = "Authorization"
	CacheControl                  = "Cache-Control"
	Connection                    = "Connection"
	ContentEncoding               = "Content-Encoding"
	ContentLanguage               = "Content-Language"
	ContentLength                 = "Content-Length"
	ContentLocation               = "Content-Location"
	ContentMD5                    = "Content-MD5"
	ContentRange                  = "Content-Range"
	ContentType                   = "Content-Type"
	Date                          = "Date"
	ETag                          = "ETag"
	Expect                        = "Expect"
	Expires                       = "Expires"
	From                          = "From"
	DoNotTrack                    = "DNT"
	IfMatch                       = "If-Match"
	IfModifiedSince               = "If-Modified-Since"
	IfNoneMatch                   = "If-None-Match"
	IfRange                       = "If-Range"
	IfUnmodifiedSince             = "If-Unmodified-Since"
	MaxForwards                   = "Max-Forwards"
	ProxyAuthorization            = "Proxy-Authorization"
	Pragma                        = "Pragma"
	Range                         = "Range"
	Referer                       = "Referer"
	UserAgent                     = "User-Agent"
	TE                            = "TE"
	Via                           = "Via"
	Warning                       = "Warning"
	Cookie                        = "Cookie"
	Origin                        = "Origin"
	AcceptDatetime                = "Accept-Datetime"
	XRequestedWith                = "X-Requested-With"
	AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	AccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	AccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	AccessControlMaxAge           = "Access-Control-Max-Age"
	AccessControlRequestMethod    = "Access-Control-Request-Method"
	AccessControlRequestHeaders   = "Access-Control-Request-Headers"
	AcceptPatch                   = "Accept-Patch"
	ContentDisposition            = "Content-Disposition"
	LastModified                  = "Last-Modified"
	Link                          = "Link"
	Location                      = "Location"
	P3P                           = "P3P"
	ProxyAuthenticate             = "Proxy-Authenticate"
	Refresh                       = "Refresh"
	RetryAfter                    = "Retry-After"
	Server                        = "Server"
	SetCookie                     = "Set-Cookie"
	StrictTransportSecurity       = "Strict-Transport-Security"
	TransferEncoding              = "Transfer-Encoding"
	Upgrade                       = "Upgrade"
	Vary                          = "Vary"
	WWWAuthenticate               = "WWW-Authenticate"

	// Deprecated X-headers https://datatracker.ietf.org/doc/html/rfc6648.

	XFrameOptions          = "X-Frame-Options"
	XXSSProtection         = "X-XSS-Protection"
	ContentSecurityPolicy  = "Content-Security-Policy"
	XContentSecurityPolicy = "X-Content-Security-Policy"
	XWebKitCSP             = "X-WebKit-CSP"
	XContentTypeOptions    = "X-Content-Type-Options"
	XPoweredBy             = "X-Powered-By"
	XUACompatible          = "X-UA-Compatible"
	XForwardedProto        = "X-Forwarded-Proto"
	XHTTPMethodOverride    = "X-HTTP-Method-Override"
	XForwardedFor          = "X-Forwarded-For"
	XRealIP                = "X-Real-IP"
	XCSRFToken             = "X-CSRF-Token" //nolint: gosec // False positive
	XRatelimitLimit        = "X-Ratelimit-Limit"
	XRatelimitRemaining    = "X-Ratelimit-Remaining"
	XRatelimitReset        = "X-Ratelimit-Reset"
)

// Normalize formats the input header to the formation of "Xxx-Xxx".
func Normalize(header string) string {
	return textproto.CanonicalMIMEHeaderKey(header)
}
