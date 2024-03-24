package octopus

type statusCode int
type codeMessage string

const (
	StatusContinue                      statusCode = 100
	StatusSwitchingProtocols            statusCode = 101
	StatusProcessing                    statusCode = 102
	StatusEarlyHints                    statusCode = 103
	StatusOK                            statusCode = 200
	StatusCreated                       statusCode = 201
	StatusAccepted                      statusCode = 202
	StatusNonAuthoritative              statusCode = 203
	StatusNoContent                     statusCode = 204
	StatusResetContent                  statusCode = 205
	StatusPartialContent                statusCode = 206
	StatusMultiStatus                   statusCode = 207
	StatusAlreadyReported               statusCode = 208
	StatusIMUsed                        statusCode = 226
	StatusMultipleChoices               statusCode = 300
	StatusMovedPermanently              statusCode = 301
	StatusFound                         statusCode = 302
	StatusSeeOther                      statusCode = 303
	StatusNotModified                   statusCode = 304
	StatusUseProxy                      statusCode = 305
	StatusTemporaryRedirect             statusCode = 307
	StatusPermanentRedirect             statusCode = 308
	StatusBadRequest                    statusCode = 400
	StatusUnauthorized                  statusCode = 401
	StatusPaymentRequired               statusCode = 402
	StatusForbidden                     statusCode = 403
	StatusNotFound                      statusCode = 404
	StatusMethodNotAllowed              statusCode = 405
	StatusNotAcceptable                 statusCode = 406
	StatusProxyAuthRequired             statusCode = 407
	StatusRequestTimeout                statusCode = 408
	StatusConflict                      statusCode = 409
	StatusGone                          statusCode = 410
	StatusLengthRequired                statusCode = 411
	StatusPreconditionFailed            statusCode = 412
	StatusPayloadTooLarge               statusCode = 413
	StatusURITooLong                    statusCode = 414
	StatusUnsupportedMedia              statusCode = 415
	StatusRangeNotSatisfiable           statusCode = 416
	StatusExpectationFailed             statusCode = 417
	StatusTeapot                        statusCode = 418
	StatusMisdirectedRequest            statusCode = 421
	StatusUnprocessableEntity           statusCode = 422
	StatusLocked                        statusCode = 423
	StatusFailedDependency              statusCode = 424
	StatusTooEarly                      statusCode = 425
	StatusUpgradeRequired               statusCode = 426
	StatusPreconditionRequired          statusCode = 428
	StatusTooManyRequests               statusCode = 429
	StatusRequestHeaderFieldsTooLarge   statusCode = 431
	StatusUnavailableForLegalReasons    statusCode = 451
	StatusInternalServerError           statusCode = 500
	StatusNotImplemented                statusCode = 501
	StatusBadGateway                    statusCode = 502
	StatusServiceUnavailable            statusCode = 503
	StatusGatewayTimeout                statusCode = 504
	StatusHTTPVersionNotSupported       statusCode = 505
	StatusVariantAlsoNegotiates         statusCode = 506
	StatusInsufficientStorage           statusCode = 507
	StatusLoopDetected                  statusCode = 508
	StatusNotExtended                   statusCode = 510
	StatusNetworkAuthenticationRequired statusCode = 511
)

var statusMessages = map[statusCode]codeMessage{
	StatusContinue:                      "Continue",
	StatusSwitchingProtocols:            "Switching Protocols",
	StatusProcessing:                    "Processing",
	StatusEarlyHints:                    "Early Hints",
	StatusOK:                            "OK",
	StatusCreated:                       "Created",
	StatusAccepted:                      "Accepted",
	StatusNonAuthoritative:              "Non-Authoritative Information",
	StatusNoContent:                     "No Content",
	StatusResetContent:                  "Reset Content",
	StatusPartialContent:                "Partial Content",
	StatusMultiStatus:                   "Multi-Status",
	StatusAlreadyReported:               "Already Reported",
	StatusIMUsed:                        "IM Used",
	StatusMultipleChoices:               "Multiple Choices",
	StatusMovedPermanently:              "Moved Permanently",
	StatusFound:                         "Found",
	StatusSeeOther:                      "See Other",
	StatusNotModified:                   "Not Modified",
	StatusUseProxy:                      "Use Proxy",
	StatusTemporaryRedirect:             "Temporary Redirect",
	StatusPermanentRedirect:             "Permanent Redirect",
	StatusBadRequest:                    "Bad Request",
	StatusUnauthorized:                  "Unauthorized",
	StatusPaymentRequired:               "Payment Required",
	StatusForbidden:                     "Forbidden",
	StatusNotFound:                      "Not Found",
	StatusMethodNotAllowed:              "Method Not Allowed",
	StatusNotAcceptable:                 "Not Acceptable",
	StatusProxyAuthRequired:             "Proxy Authentication Required",
	StatusRequestTimeout:                "Request Timeout",
	StatusConflict:                      "Conflict",
	StatusGone:                          "Gone",
	StatusLengthRequired:                "Length Required",
	StatusPreconditionFailed:            "Precondition Failed",
	StatusPayloadTooLarge:               "Payload Too Large",
	StatusURITooLong:                    "URI Too Long",
	StatusUnsupportedMedia:              "Unsupported Media Type",
	StatusRangeNotSatisfiable:           "Range Not Satisfiable",
	StatusExpectationFailed:             "Expectation Failed",
	StatusTeapot:                        "I'm a teapot",
	StatusMisdirectedRequest:            "Misdirected Request",
	StatusUnprocessableEntity:           "Unprocessable Entity",
	StatusLocked:                        "Locked",
	StatusFailedDependency:              "Failed Dependency",
	StatusTooEarly:                      "Too Early",
	StatusUpgradeRequired:               "Upgrade Required",
	StatusPreconditionRequired:          "Precondition Required",
	StatusTooManyRequests:               "Too Many Requests",
	StatusRequestHeaderFieldsTooLarge:   "Request Header Fields Too Large",
	StatusUnavailableForLegalReasons:    "Unavailable For Legal Reasons",
	StatusInternalServerError:           "Internal Server Error",
	StatusNotImplemented:                "Not Implemented",
	StatusBadGateway:                    "Bad Gateway",
	StatusServiceUnavailable:            "Service Unavailable",
	StatusGatewayTimeout:                "Gateway Timeout",
	StatusHTTPVersionNotSupported:       "HTTP Version Not Supported",
	StatusVariantAlsoNegotiates:         "Variant Also Negotiates",
	StatusInsufficientStorage:           "Insufficient Storage",
	StatusLoopDetected:                  "Loop Detected",
	StatusNotExtended:                   "Not Extended",
	StatusNetworkAuthenticationRequired: "Network Authentication Required",
}
