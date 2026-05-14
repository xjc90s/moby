package dockerversion

import (
	"context"
	"runtime"
	"strings"
	"sync"

	"github.com/moby/moby/v2/pkg/parsers/kernel"
	"github.com/moby/moby/v2/pkg/useragent"
)

// UAStringKey is used as key type for user-agent string in net/context struct
type UAStringKey struct{}

// DockerUserAgent is the User-Agent the Docker client uses to identify itself.
// In accordance with RFC 7231 (5.5.3) is of the form:
//
//	[docker client's UA] UpstreamClient([upstream client's UA])
func DockerUserAgent(ctx context.Context, extraVersions ...useragent.VersionInfo) string {
	ua := useragent.AppendVersions(getDaemonUserAgent(), extraVersions...)
	if upstreamUA := getUpstreamUserAgent(ctx); upstreamUA != "" {
		ua += " " + upstreamUA
	}
	return ua
}

var (
	daemonUAOnce sync.Once
	daemonUA     string
)

// getDaemonUserAgent returns the user-agent to use for requests made by
// the daemon.
//
// It includes:
//
// - the docker version
// - go version
// - git-commit
// - kernel version
// - os
// - architecture
func getDaemonUserAgent() string {
	daemonUAOnce.Do(func() {
		httpVersion := make([]useragent.VersionInfo, 0, 6)
		httpVersion = append(httpVersion, useragent.VersionInfo{Name: "docker", Version: Version})
		httpVersion = append(httpVersion, useragent.VersionInfo{Name: "go", Version: runtime.Version()})
		httpVersion = append(httpVersion, useragent.VersionInfo{Name: "git-commit", Version: GitCommit})
		if kernelVersion, err := kernel.GetKernelVersion(); err == nil {
			httpVersion = append(httpVersion, useragent.VersionInfo{Name: "kernel", Version: kernelVersion.String()})
		}
		httpVersion = append(httpVersion, useragent.VersionInfo{Name: "os", Version: runtime.GOOS})
		httpVersion = append(httpVersion, useragent.VersionInfo{Name: "arch", Version: runtime.GOARCH})
		daemonUA = useragent.AppendVersions("", httpVersion...)
	})
	return daemonUA
}

// getUpstreamUserAgent returns the previously saved user-agent context stored
// in ctx, if one exists, and formats it as:
//
//	UpstreamClient(<upstream user agent string>)
//
// It returns an empty string if no user-agent is present in the context.
func getUpstreamUserAgent(ctx context.Context) string {
	upstreamUA, ok := ctx.Value(UAStringKey{}).(string)
	if !ok || upstreamUA == "" {
		return ""
	}

	return "UpstreamClient(" + escapeStr(upstreamUA) + ")"
}

// escapeStr escapes and sanitizes s for use in a User-Agent comment.
func escapeStr(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	for i := range len(s) {
		switch c := s[i]; c {
		case '(', ')', ';', '\\':
			b.WriteByte('\\')
			b.WriteByte(c)
		case '\t':
			b.WriteByte(c)
		default:
			if c >= 0x20 && c != 0x7f {
				b.WriteByte(c)
			}
		}
	}

	return b.String()
}
