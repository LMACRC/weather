package http

type Config struct {
	// Path is the name of the HTTP path for incoming Ecowitt requests
	Path string

	Dev struct {
		// ForwardTo specifies a HTTP address to forward incoming requests.
		ForwardTo string `toml:"forward_to" mapstructure:"forward_to"`
	}
}
