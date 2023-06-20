package trust

import (
	"context"
	"strings"

	"github.com/streamingfast/dauth"
)

func Register() {
	dauth.Register("trust", func(configURL string) (dauth.Authenticator, error) {
		return &trustPlugin{}, nil
	})

}

type trustPlugin struct {
}

func (t *trustPlugin) Authenticate(ctx context.Context, path string, headers map[string][]string, ipAddress string) (context.Context, error) {
	out := make(dauth.TrustedHeaders)
	for key, values := range headers {
		out.Set(key, strings.ToLower(values[0]))
	}
	return out.ToContext(ctx), nil
}
