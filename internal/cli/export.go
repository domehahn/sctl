package cli

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
