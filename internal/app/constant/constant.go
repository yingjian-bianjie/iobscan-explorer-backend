package constant

const (
	Success = "success"
	Fail    = "fail"

	EnvNameZkServices   = "ZK_SERVICES"
	EnvNameZkUsername   = "ZK_USERNAME"
	EnvNameZkPasswd     = "ZK_PASSWD"
	EnvNameZkConfigPath = "ZK_CONFIG_PATH"

	DefaultZkConfigPath = "/iobscan-explorer-backend/config"

	DefaultTimezone   = "UTC"
	DefaultTimeFormat = "2006-01-02 15:04:05"

	HeaderAuthorization      = "Authorization"
	HeaderPagination         = "x-pagination"
	HeaderContentDisposition = "Content-Disposition"
	HeaderXForwardedFor      = "X-Forwarded-For"
	HeaderTimestamp          = "X-Timestamp"
	HeaderSignature          = "X-Signature"

	DefaultRetryTimes = 6
	Guest             = "1"
	NetworkDelay      = 15

	DefaultGinCacheTtl = 10
)
