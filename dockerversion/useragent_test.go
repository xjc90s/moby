package dockerversion

import (
	"context"
	"testing"

	"github.com/moby/moby/v2/pkg/useragent"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestDockerUserAgent(t *testing.T) {
	tests := []struct {
		doc      string
		ctx      context.Context
		metadata []useragent.VersionInfo
		expected string
	}{
		{
			doc:      "daemon user-agent",
			ctx:      t.Context(),
			expected: getDaemonUserAgent(),
		},
		{
			doc: "daemon user-agent custom metadata",
			ctx: t.Context(),
			metadata: []useragent.VersionInfo{
				{Name: "hello", Version: "world"},
				{Name: "foo", Version: "bar"},
			},
			expected: getDaemonUserAgent() + ` hello/world foo/bar`,
		},
		{
			doc:      "daemon user-agent with upstream",
			ctx:      context.WithValue(t.Context(), UAStringKey{}, "Magic-Client/1.2.3 (linux)"),
			expected: getDaemonUserAgent() + ` UpstreamClient(Magic-Client/1.2.3 \(linux\))`,
		},
		{
			doc: "daemon user-agent with upstream and custom metadata",
			ctx: context.WithValue(t.Context(), UAStringKey{}, "Magic-Client/1.2.3 (linux)"),
			metadata: []useragent.VersionInfo{
				{Name: "hello", Version: "world"},
				{Name: "foo", Version: "bar"},
			},
			expected: getDaemonUserAgent() + ` hello/world foo/bar UpstreamClient(Magic-Client/1.2.3 \(linux\))`,
		},
		{
			doc:      "daemon user-agent with upstream special chars",
			ctx:      context.WithValue(t.Context(), UAStringKey{}, `Magic-Client/1.2.3 (linux); \ test`),
			expected: getDaemonUserAgent() + ` UpstreamClient(Magic-Client/1.2.3 \(linux\)\; \\ test)`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.doc, func(t *testing.T) {
			assert.Check(t, is.Equal(DockerUserAgent(tc.ctx, tc.metadata...), tc.expected))
		})
	}
}
