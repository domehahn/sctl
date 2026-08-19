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

// Type exports for tests.
type SkillTemplate = skillTemplate
type TemplateRegistry = templateRegistry
type GraphNode = graphNode
type WorkspaceManifest = workspaceManifest

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
	ReadSkillScripts      = readSkillScripts
	LoadTemplateRegistry  = loadTemplateRegistry
	SaveTemplateRegistry  = saveTemplateRegistry
	RenameYAMLKey         = renameYAMLKey
	ReadWorkspace         = readWorkspace
	FindWorkspaceFile     = findWorkspaceFile
	WorkspaceDirs         = workspaceDirs
	FilterChangedDirs     = filterChangedDirs
	IsTemplateFile        = isTemplateFile
	SortedTemplateNames   = sortedTemplateNames
	ScriptNames           = scriptNames
	RunProjectHook        = runProjectHook
	MigrateSkillMD        = migrateSkillMD
	MigrateSkillYAMLFile  = migrateSkillYAMLFile
	ReadGraphNode         = readGraphNode
)

func TestConfig() *config.Config {
	return &config.Config{Registries: map[string]config.RegistryConfig{}}
}
