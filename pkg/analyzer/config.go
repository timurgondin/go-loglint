package analyzer

import (
	"flag"
	"strings"
)

type Config struct {
	CheckLowercase    bool
	CheckEnglish      bool
	CheckSpecialChars bool
	CheckSensitive    bool
	ExtraPatterns     []string
}

var DefaultConfig = Config{
	CheckLowercase:    true,
	CheckEnglish:      true,
	CheckSpecialChars: true,
	CheckSensitive:    true,
	ExtraPatterns:     nil,
}

type extraPatternsFlag []string

func (f *extraPatternsFlag) String() string {
	return strings.Join(*f, ",")
}

func (f *extraPatternsFlag) Set(s string) error {
	patterns := strings.Split(s, ",")
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p != "" {
			*f = append(*f, p)
		}
	}
	return nil
}

func registerFlags(fs *flag.FlagSet, cfg *Config) {
	fs.BoolVar(&cfg.CheckLowercase, "check-lowercase", true, "check that log messages start with lowercase letter")
	fs.BoolVar(&cfg.CheckEnglish, "check-english", true, "check that log messages are in English")
	fs.BoolVar(&cfg.CheckSpecialChars, "check-special-chars", true, "check that log messages have no special characters")
	fs.BoolVar(&cfg.CheckSensitive, "check-sensitive", true, "check that log messages have no sensitive data")

	fs.Var((*extraPatternsFlag)(&cfg.ExtraPatterns), "extra-patterns", "comma-separated list of additional sensitive patterns")
}
