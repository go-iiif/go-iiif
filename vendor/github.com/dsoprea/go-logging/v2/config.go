package log

import (
	"fmt"
	"os"
	"strings"
)

// Config keys.
const (
	ckFormat                 = "LogFormat"
	ckDefaultAdapterName     = "LogDefaultAdapterName"
	ckLevelName              = "LogLevelName"
	ckIncludeNouns           = "LogIncludeNouns"
	ckExcludeNouns           = "LogExcludeNouns"
	ckExcludeBypassLevelName = "LogExcludeBypassLevelName"
)

// Other constants
const (
	defaultFormat    = "{{.Noun}}: [{{.Level}}] {{if eq .ExcludeBypass true}} [BYPASS]{{end}} {{.Message}}"
	defaultLevelName = levelNameInfo
)

// Config
var (
	// Alternative format.
	format = defaultFormat

	// Alternative adapter.
	defaultAdapterName = ""

	// Alternative level at which to display log-items
	levelName = LogLevelName(strings.ToLower(string(defaultLevelName)))

	// Configuration-driven comma-separated list of nouns to include.
	includeNouns = ""

	// Configuration-driven comma-separated list of nouns to exclude.
	excludeNouns = ""

	// excludeBypassLevelName is the level at which to disregard exclusion (if
	// the severity of a message meets or exceed this, always display).
	excludeBypassLevelName LogLevelName
)

// Other
var (
	configurationLoaded = false
)

// GetDefaultAdapterName returns the default adapter name. May be empty.
func GetDefaultAdapterName() string {
	return defaultAdapterName
}

// SetDefaultAdapterName sets the default adapter. If not set, the first one
// registered will be used.
func SetDefaultAdapterName(name string) {
	defaultAdapterName = name
}

// LoadConfiguration loads the effective configuration.
func LoadConfiguration(cp ConfigurationProvider) {
	configuredDefaultAdapterName := cp.DefaultAdapterName()

	if configuredDefaultAdapterName != "" {
		defaultAdapterName = configuredDefaultAdapterName
	}

	includeNouns = cp.IncludeNouns()
	excludeNouns = cp.ExcludeNouns()
	excludeBypassLevelName = cp.ExcludeBypassLevelName()

	f := cp.Format()
	if f != "" {
		format = f
	}

	ln := cp.LevelName()
	if ln != "" {
		levelName = LogLevelName(strings.ToLower(string(ln)))
	}

	configurationLoaded = true
}

func getConfigState() map[string]interface{} {
	return map[string]interface{}{
		"format":                 format,
		"defaultAdapterName":     defaultAdapterName,
		"levelName":              levelName,
		"includeNouns":           includeNouns,
		"excludeNouns":           excludeNouns,
		"excludeBypassLevelName": excludeBypassLevelName,
	}
}

func setConfigState(config map[string]interface{}) {
	format = config["format"].(string)

	defaultAdapterName = config["defaultAdapterName"].(string)

	levelName = config["levelName"].(LogLevelName)
	levelName = LogLevelName(strings.ToLower(string(levelName)))

	includeNouns = config["includeNouns"].(string)
	excludeNouns = config["excludeNouns"].(string)
	excludeBypassLevelName = config["excludeBypassLevelName"].(LogLevelName)
}

func getConfigDump() string {
	return fmt.Sprintf(
		"Current configuration:\n"+
			"  FORMAT=[%s]\n"+
			"  DEFAULT-ADAPTER-NAME=[%s]\n"+
			"  LEVEL-NAME=[%s]\n"+
			"  INCLUDE-NOUNS=[%s]\n"+
			"  EXCLUDE-NOUNS=[%s]\n"+
			"  EXCLUDE-BYPASS-LEVEL-NAME=[%s]",
		format, defaultAdapterName, levelName, includeNouns, excludeNouns, excludeBypassLevelName)
}

// IsConfigurationLoaded indicates whether a config has been loaded.
func IsConfigurationLoaded() bool {
	return configurationLoaded
}

// ConfigurationProvider describes minimal configuration implementation.
type ConfigurationProvider interface {
	// Alternative format (defaults to .
	Format() string

	// Alternative adapter (defaults to "appengine").
	DefaultAdapterName() string

	// Alternative level at which to display log-items (defaults to
	// "info").
	LevelName() LogLevelName

	// Configuration-driven comma-separated list of nouns to include. Defaults
	// to empty.
	IncludeNouns() string

	// Configuration-driven comma-separated list of nouns to exclude. Defaults
	// to empty.
	ExcludeNouns() string

	// Level at which to disregard exclusion (if the severity of a message
	// meets or exceed this, always display). Defaults to empty.
	ExcludeBypassLevelName() LogLevelName
}

// EnvironmentConfigurationProvider configuration-provider.
type EnvironmentConfigurationProvider struct {
}

// NewEnvironmentConfigurationProvider returns a new
// EnvironmentConfigurationProvider.
func NewEnvironmentConfigurationProvider() *EnvironmentConfigurationProvider {
	return new(EnvironmentConfigurationProvider)
}

// Format returns the format string.
func (ecp *EnvironmentConfigurationProvider) Format() string {
	return os.Getenv(ckFormat)
}

// DefaultAdapterName returns the name of the default-adapter.
func (ecp *EnvironmentConfigurationProvider) DefaultAdapterName() string {
	return os.Getenv(ckDefaultAdapterName)
}

// LevelName returns the current level-name.
func (ecp *EnvironmentConfigurationProvider) LevelName() LogLevelName {
	return LogLevelName(os.Getenv(ckLevelName))
}

// IncludeNouns returns inlined set of effective include nouns.
func (ecp *EnvironmentConfigurationProvider) IncludeNouns() string {
	return os.Getenv(ckIncludeNouns)
}

// ExcludeNouns returns inlined set of effective exclude nouns.
func (ecp *EnvironmentConfigurationProvider) ExcludeNouns() string {
	return os.Getenv(ckExcludeNouns)
}

// ExcludeBypassLevelName returns the level, if any, of the current bypass level
// for the excluded nouns.
func (ecp *EnvironmentConfigurationProvider) ExcludeBypassLevelName() LogLevelName {
	return LogLevelName(os.Getenv(ckExcludeBypassLevelName))
}

// StaticConfigurationProvider configuration-provider.
type StaticConfigurationProvider struct {
	format                 string
	defaultAdapterName     string
	levelName              LogLevelName
	includeNouns           string
	excludeNouns           string
	excludeBypassLevelName LogLevelName
}

// NewStaticConfigurationProvider returns a new StaticConfigurationProvider
// struct.
func NewStaticConfigurationProvider() *StaticConfigurationProvider {
	return new(StaticConfigurationProvider)
}

// SetFormat sets the message format/layout.
func (scp *StaticConfigurationProvider) SetFormat(format string) {
	scp.format = format
}

// SetDefaultAdapterName sets the default adapter name.
func (scp *StaticConfigurationProvider) SetDefaultAdapterName(adapterName string) {
	scp.defaultAdapterName = adapterName
}

// SetLevelName sets the effective level (using the name).
func (scp *StaticConfigurationProvider) SetLevelName(levelName LogLevelName) {
	scp.levelName = LogLevelName(strings.ToLower(string(levelName)))
}

// SetLevel sets the effective level (using the constant).
func (scp *StaticConfigurationProvider) SetLevel(level LogLevel) {
	scp.levelName = levelNameMapR[level]
}

// SetIncludeNouns sets an inlined set of nouns to include.
func (scp *StaticConfigurationProvider) SetIncludeNouns(includeNouns string) {
	scp.includeNouns = includeNouns
}

// SetExcludeNouns sets an inlined set of nouns to exclude.
func (scp *StaticConfigurationProvider) SetExcludeNouns(excludeNouns string) {
	scp.excludeNouns = excludeNouns
}

// SetExcludeBypassLevelName sets a specific level to for the noun exclusions
// (e.g. hide them at/below INFO but show ERROR logging for everything).
func (scp *StaticConfigurationProvider) SetExcludeBypassLevelName(excludeBypassLevelName LogLevelName) {
	scp.excludeBypassLevelName = excludeBypassLevelName
}

// Format returns the format string.
func (scp *StaticConfigurationProvider) Format() string {
	return scp.format
}

// DefaultAdapterName returns the name of the default-adapter.
func (scp *StaticConfigurationProvider) DefaultAdapterName() string {
	return scp.defaultAdapterName
}

// LevelName returns the current level-name.
func (scp *StaticConfigurationProvider) LevelName() LogLevelName {
	return scp.levelName
}

// IncludeNouns returns inlined set of effective include nouns.
func (scp *StaticConfigurationProvider) IncludeNouns() string {
	return scp.includeNouns
}

// ExcludeNouns returns inlined set of effective exclude nouns.
func (scp *StaticConfigurationProvider) ExcludeNouns() string {
	return scp.excludeNouns
}

// ExcludeBypassLevelName returns the level, if any, of the current bypass level
// for the excluded nouns.
func (scp *StaticConfigurationProvider) ExcludeBypassLevelName() LogLevelName {
	return scp.excludeBypassLevelName
}

func init() {
	// Do the initial configuration-load from the environment. We gotta seed it
	// with something for simplicity's sake.
	ecp := NewEnvironmentConfigurationProvider()
	LoadConfiguration(ecp)
}
