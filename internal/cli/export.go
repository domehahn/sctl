package cli

import "github.com/domehahn/skpm/v2/internal/config"

var (
	IsLocalPath            = isLocalPath
	ParseSkillArg          = parseSkillArg
	IsValidSkillName       = isValidSkillName
	BuildConfigYAML        = buildConfigYAML
	FileExists             = fileExists
	WriteAtomic            = writeAtomic
	UpdateGitignore        = updateGitignore
	ValidateConfig         = validateConfig
	ValidateRegistryConfig = validateRegistryConfig
	MaskTokens             = maskTokens
	CopyDir                = copyDir
	ConfigFilePath         = configFilePath
	FormatGitTag           = formatGitTag
	GitIgnoreBlock         = gitignoreBlock
)

// ConfigErrorSlice is the exported type for []configError used in tests.
type ConfigError = configError

// Exports for new commands — used by tests/unit/cli/new_commands_test.go.
var (
	HumanBytes            = humanBytes
	ParseSkillAtVersion   = parseSkillAtVersion
	ParseSkillRef         = parseSkillRef
	ScaffoldFiles         = scaffoldFiles
	ParseChangelogEntries = parseChangelogEntries
	AddChangelogMessage   = addChangelogMessage
	SnapshotPath          = snapshotPath
	UpsertLink            = upsertLink
	IsLinkedSkill         = isLinkedSkill
	DefaultSkillRoots     = defaultSkillRoots
	ConfigGetScalar       = configGetScalar
	ConfigSetScalar       = configSetScalar
	SemverDelta           = semverDelta
	IsWithinBase          = isWithinBase
)

func TestConfig() *config.Config {
	return &config.Config{Registries: map[string]config.RegistryConfig{}}
}
