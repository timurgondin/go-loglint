package plugin

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/timurgondin/go-loglint/pkg/analyzer"
)

func init() {
	register.Plugin("loglint", New)
}

type Settings struct {
	CheckLowercase    bool     `mapstructure:"check-lowercase"`
	CheckEnglish      bool     `mapstructure:"check-english"`
	CheckSpecialChars bool     `mapstructure:"check-special-chars"`
	CheckSensitive    bool     `mapstructure:"check-sensitive"`
	ExtraPatterns     []string `mapstructure:"extra-patterns"`
}

type plugin struct {
	settings Settings
}

var _ register.LinterPlugin = new(plugin)

func New(settings any) (register.LinterPlugin, error) {
	s := Settings{
		CheckLowercase:    analyzer.DefaultConfig.CheckLowercase,
		CheckEnglish:      analyzer.DefaultConfig.CheckEnglish,
		CheckSpecialChars: analyzer.DefaultConfig.CheckSpecialChars,
		CheckSensitive:    analyzer.DefaultConfig.CheckSensitive,
		ExtraPatterns:     analyzer.DefaultConfig.ExtraPatterns,
	}

	if settings != nil {
		if err := mapstructure.Decode(settings, &s); err != nil {
			return nil, err
		}
	}

	return &plugin{settings: s}, nil
}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	analyzer.SetConfig(analyzer.Config{
		CheckLowercase:    p.settings.CheckLowercase,
		CheckEnglish:      p.settings.CheckEnglish,
		CheckSpecialChars: p.settings.CheckSpecialChars,
		CheckSensitive:    p.settings.CheckSensitive,
		ExtraPatterns:     p.settings.ExtraPatterns,
	})

	return []*analysis.Analyzer{
		analyzer.Analyzer,
	}, nil
}

func (p *plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
