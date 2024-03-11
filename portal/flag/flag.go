package flag

import (
	"strings"

	kingpin "github.com/alecthomas/kingpin/v2"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal"
)

// PushIntervalFlagName is the canonical flag name to configure the
// periodic push interval
const PushIntervalFlagName = "portal.push-interval"

// PushIntervalHelp is the help description for the portal.push-interval flag.
const PushIntervalHelp = "Periodic Push interval."

// LogIntervalFlagName is the canonical flag name to configure the
// retention period for local value buffer
const LogIntervalFlagName = "portal.log-interval"

// LogIntervalHelp is the help description for the portal.log-interval flag.
const LogIntervalHelp = "Retention period for local value buffer."

// AutoPushDelayFlagName is the canonical flag name to configure the
// minimal delay of AutoPush
const AutoPushDelayFlagName = "portal.auto-push-delay"

// AutoPushDelayHelp is the help description for the portal.auto-push-delay flag.
const AutoPushDelayHelp = "Minimal delay of AutoPush."

// TrustLevelFlagName is the canonical flag name to configure the
// trust level towards incoming requests
const TrustLevelFlagName = "portal.trust-level"

// TrustLevelHelp is the delp description for the portal.trust-level flag.
var TrustLevelHelp = "Trust level towards incoming requests for identifying information. One of: [" + strings.Join(portal.TrustLevelValues, ", ") + "]"

// AuthUsernameFlagName is the canonical flag name to configure the
// authentication principal
const AuthUsernameFlagName = "portal.auth-username"

// AuthUsernameHelp is the delp description for the portal.auth-username flag.
const AuthUsernameHelp = "Authentication principal to required from incoming requests"

// AuthPasswordFlagName is the canonical flag name to configure the
// authentication secret
const AuthPasswordFlagName = "portal.auth-password"

// AuthPasswordHelp is the delp description for the portal.auth-password flag.
const AuthPasswordHelp = "Authentication secret to required from incoming requests"

// AuthPasswordFileFlagName is the canonical flag name to configure the
// authentication secret from a file
const AuthPasswordFileFlagName = "portal.auth-password-file"

// AuthPasswordFileHelp is the delp description for the portal.auth-password-file flag.
const AuthPasswordFileHelp = "File to read the authentication secret to required from incoming requests"

// AddFlags adds the flags used by this package to the Kingpin application.
// To use the default Kingpin application, call AddFlags(kingpin.CommandLine)
func AddFlags(a *kingpin.Application, config *portal.Config) {
	a.Flag(PushIntervalFlagName, PushIntervalHelp).
		Default("10s").DurationVar(&(config.PushInterval))

	a.Flag(LogIntervalFlagName, LogIntervalHelp).
		Default("0s").DurationVar(&(config.LogInterval))

	a.Flag(AutoPushDelayFlagName, AutoPushDelayHelp).
		Default("30s").DurationVar(&(config.AutoPushDelay))

	a.Flag(AuthUsernameFlagName, AuthUsernameHelp).
		Default("").StringVar(&(config.AuthUsername))

	a.Flag(AuthPasswordFlagName, AuthPasswordHelp).
		Default("").StringVar(&(config.AuthPassword))

	a.Flag(AuthPasswordFileFlagName, AuthPasswordFileHelp).
		PlaceHolder("FILE").ExistingFileVar(&(config.AuthPasswordFile))

	a.Flag(TrustLevelFlagName, TrustLevelHelp).
		Default("client").HintOptions(portal.TrustLevelValues...).
		SetValue(&config.TrustLevel)
}
