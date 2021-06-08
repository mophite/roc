package namespace

type Service = string

const (
	DefaultVersion      Service = "v1.0.0"
	DefaultSchema               = "goroc"
	DefaultConfigSchema         = "configroc"
)

type Header = string

const (
	DefaultHeaderTrace   Header = "X-Idempotency-Key"
	DefaultHeaderVersion        = "X-Api-Version"
	DefaultHeaderToken          = "X-Api-Token"
	DefaultHeaderAddress        = "X-Api-Address"
)

type Schema = string

type Scope = string

type RequestChannel = string
