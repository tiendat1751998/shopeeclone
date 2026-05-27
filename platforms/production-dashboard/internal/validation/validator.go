package validation

import (
	"regexp"
	"strings"
)

var serviceNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,253}[a-zA-Z0-9]$`)
var versionRegex = regexp.MustCompile(`^v?\d+\.\d+\.\d+.*$`)
var environmentRegex = regexp.MustCompile(`^(development|staging|production)$`)
var severityRegex = regexp.MustCompile(`^(critical|high|medium|low)$`)

func IsValidServiceName(name string) bool {
	return serviceNameRegex.MatchString(name) && !strings.Contains(name, "..")
}

func IsValidVersion(version string) bool {
	return versionRegex.MatchString(version)
}

func IsValidEnvironment(env string) bool {
	return environmentRegex.MatchString(env)
}

func IsValidSeverity(severity string) bool {
	return severityRegex.MatchString(severity)
}

func SanitizeString(input string, maxLen int) string {
	s := strings.TrimSpace(input)
	if len(s) > maxLen {
		s = s[:maxLen]
	}
	return s
}
