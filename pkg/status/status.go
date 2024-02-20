package status

type StatusCode int

const (
	// Informational Response
	Continue              StatusCode = 100
	SwitchingProtocols    StatusCode = 101
	Processing            StatusCode = 102
	EarlyHints            StatusCode = 103
	ResponseIsStale       StatusCode = 110
	RevalidationFailed    StatusCode = 111
	DisconnectedOperation StatusCode = 112
	HeuristicExpiration   StatusCode = 113
	// Success
	OK                          StatusCode = 200
	Created                     StatusCode = 201
	Accepted                    StatusCode = 202
	NonAuthoritativeInformation StatusCode = 203
	NoContent                   StatusCode = 204
	ResetContent                StatusCode = 205
	PartialContent              StatusCode = 206
	MultiStatus                 StatusCode = 207
	AlreadyReported             StatusCode = 208
	ThisIsFine                  StatusCode = 218 // unofficial
	ImUsed                      StatusCode = 226
	// Redirection
	MultipleChoices   StatusCode = 300
	MovedPermanently  StatusCode = 301
	Found             StatusCode = 302
	SeeOther          StatusCode = 303
	NotModified       StatusCode = 304
	UseProxy          StatusCode = 305
	SwitchProxy       StatusCode = 306
	TemporaryRedirect StatusCode = 307
	PermanentRedirect StatusCode = 308
	// Client Errors
	BadRequest                       StatusCode = 400
	Unauthorized                     StatusCode = 401
	PaymentRequired                  StatusCode = 402
	Forbidden                        StatusCode = 403
	NotFound                         StatusCode = 404
	MethodNotAllowed                 StatusCode = 405
	NotAcceptable                    StatusCode = 406
	ProxyAuthenticationRequired      StatusCode = 407
	RequestTimeout                   StatusCode = 408
	Conflict                         StatusCode = 409
	Gone                             StatusCode = 410
	LengthRequired                   StatusCode = 411
	PreconditionFailed               StatusCode = 412
	PayloadTooLarge                  StatusCode = 413
	URITooLong                       StatusCode = 414
	UnsupportedMediaType             StatusCode = 415
	RangeNotSatisfiable              StatusCode = 416
	ExpectationFailed                StatusCode = 417
	ImATeapot                        StatusCode = 418
	PageExpired                      StatusCode = 419 // Unofficial
	MethodFailure                    StatusCode = 420 // Unofficial
	MisdirectedRequest               StatusCode = 421
	UnprocessableContent             StatusCode = 422
	Locked                           StatusCode = 423
	FailedDependency                 StatusCode = 424
	TooEarly                         StatusCode = 425
	UpgradeRequired                  StatusCode = 426
	PreconditionRequired             StatusCode = 428
	TooManyRequests                  StatusCode = 429
	RequestHeaderFieldsTooLarge      StatusCode = 431
	NoResponse                       StatusCode = 444
	BlockedByWindowsParentalControls StatusCode = 450 // Unofficial
	UnavailableForLegalReasons       StatusCode = 451
	InvalidToken                     StatusCode = 498 // Unofficial
	TokenRequired                    StatusCode = 499 // Unofficial
	// Server Errors
	InternalServerError           StatusCode = 500
	NotImplemented                StatusCode = 501
	BadGateway                    StatusCode = 502
	ServiceUnavailable            StatusCode = 503
	GatewayTimeout                StatusCode = 504
	HTTPVersionNotSupported       StatusCode = 505
	VariantAlsoNegotiates         StatusCode = 506
	InsufficientStorage           StatusCode = 507
	LoopDetected                  StatusCode = 508
	BandwidthLimit                StatusCode = 509 // Unofficial
	NotExtended                   StatusCode = 510
	NetworkAuthenticationRequired StatusCode = 511
	WebServerReturnedUnknownError StatusCode = 520
	WebServerIsDown               StatusCode = 521 // Cloudflare
	ConnectionTimedOut            StatusCode = 522 // Cloudflare
	OriginIsUnreachable           StatusCode = 523 // Cloudflare
	TimeoutOccurred               StatusCode = 524 // Cloudflare
	RailgunError                  StatusCode = 527 // Cloudflare
	SiteOverloaded                StatusCode = 529 // Unofficial
	SiteFrozen                    StatusCode = 530 // Unofficial
	NetworkReadTimeout            StatusCode = 598 // Unofficial
	NetworkConnectTimeout         StatusCode = 599 // Unofficial
)
