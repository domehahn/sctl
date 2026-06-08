package registry_test

import (
	"testing"

	"github.com/domehahn/skpm/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelectVersionConstraints(t *testing.T) {
	versions := []registry.VersionInfo{
		{Version: "1.4.0"},
		{Version: "1.5.0"},
		{Version: "1.5.2"},
		{Version: "1.6.0", Deprecated: true},
		{Version: "2.0.0"},
	}

	tests := map[string]string{
		"1.5.0":          "1.5.0",
		"latest":         "2.0.0",
		"":               "2.0.0",
		"^1.5.0":         "1.5.2",
		"~1.5.0":         "1.5.2",
		">=1.4.0 <2.0.0": "1.5.2",
	}
	for constraint, want := range tests {
		got, err := registry.SelectVersion(versions, constraint, registry.ConstraintOptions{})
		require.NoError(t, err)
		assert.Equal(t, want, got)
	}
}

func TestSelectVersionPrereleaseIgnoredByDefault(t *testing.T) {
	versions := []registry.VersionInfo{{Version: "1.0.0"}, {Version: "1.1.0-alpha.1"}}
	got, err := registry.SelectVersion(versions, "latest", registry.ConstraintOptions{})
	require.NoError(t, err)
	assert.Equal(t, "1.0.0", got)

	got, err = registry.SelectVersion(versions, "latest", registry.ConstraintOptions{IncludePrerelease: true})
	require.NoError(t, err)
	assert.Equal(t, "1.1.0-alpha.1", got)
}

func TestSelectVersionYankedIgnored(t *testing.T) {
	versions := []registry.VersionInfo{{Version: "1.0.0"}, {Version: "1.1.0", Yanked: true}}
	got, err := registry.SelectVersion(versions, "latest", registry.ConstraintOptions{})
	require.NoError(t, err)
	assert.Equal(t, "1.0.0", got)
}

func TestSelectVersionDeprecatedAvoided(t *testing.T) {
	versions := []registry.VersionInfo{{Version: "1.0.0"}, {Version: "1.1.0", Deprecated: true}}
	got, err := registry.SelectVersion(versions, "latest", registry.ConstraintOptions{})
	require.NoError(t, err)
	assert.Equal(t, "1.0.0", got)
}

func TestSelectVersionNoCompatible(t *testing.T) {
	_, err := registry.SelectVersion([]registry.VersionInfo{{Version: "1.0.0"}}, "^2.0.0", registry.ConstraintOptions{})
	assert.Error(t, err)
}
