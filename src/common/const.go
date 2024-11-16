package common

const (
	// CtxKeyTraceID is the key for the TraceId value in contexts
	CtxKeyTraceID = "TraceId"

	// BackgroundTraceID default, locally generated traceId
	BackgroundTraceID = "background-trace-id"

	// HeaderTraceID is the trace ID header.
	HeaderTraceID = "X-Trace-Id"

	// HeaderContentType is content type header.
	HeaderContentType = "Content-Type"

	// HeaderAuthorization is the value of the authorization header.
	HeaderAuthorization = "Authorization"

	// HeaderValueContentTypeJson represents the value for this content type
	HeaderValueContentTypeJson = "application/json"
)
