package analyzer

import (
	"fmt"
	"go/ast"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

var sensitiveKeywords = []string{
	"password",
	"passwd",
	"secret",
	"token",
	"api_key",
	"apikey",
	"auth",
	"credential",
	"private_key",
	"credit_card",
}

func checkMessage(pass *analysis.Pass, lit *ast.BasicLit, msg string, cfg *Config) {
	if len(msg) == 0 {
		return
	}

	if cfg.CheckEnglish && !isEnglish(msg) {
		pass.Reportf(lit.Pos(), "log message must be in English: %q", msg)
		return
	}

	if cfg.CheckLowercase {
		checkLowercase(pass, lit, msg)
	}
	if cfg.CheckSpecialChars {
		checkSpecialChars(pass, lit, msg)
	}
	if cfg.CheckSensitive {
		checkSensitiveData(pass, lit, msg, cfg)
	}
}

func isEnglish(msg string) bool {
	for _, r := range msg {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func checkLowercase(pass *analysis.Pass, lit *ast.BasicLit, msg string) {
	firstRune, _ := utf8.DecodeRuneInString(msg)
	if unicode.IsUpper(firstRune) {
		rawContent := lit.Value[1 : len(lit.Value)-1]
		rawFirst, rawSize := utf8.DecodeRuneInString(rawContent)
		delimiter := string(lit.Value[0])
		fixed := string(unicode.ToLower(rawFirst)) + rawContent[rawSize:]
		fixedWithQuotes := delimiter + fixed + delimiter

		pass.Report(analysis.Diagnostic{
			Pos:     lit.Pos(),
			Message: fmt.Sprintf("log message must start with lowercase letter: %q", msg),
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "convert first letter to lowercase",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     lit.Pos(),
							End:     lit.End(),
							NewText: []byte(fixedWithQuotes),
						},
					},
				},
			},
		})
	}
}

func checkSpecialChars(pass *analysis.Pass, lit *ast.BasicLit, msg string) {
	for _, r := range msg {
		if r > 0x7E {
			reportSpecialCharsViolation(pass, lit, msg, "special characters or emoji")
			return
		}
		switch r {
		case '!', '?', '@', '#', '$', '%', '^', '&', '*', '~', '`', '|', '\\':
			reportSpecialCharsViolation(pass, lit, msg, "special characters")
			return
		}
	}
}

func reportSpecialCharsViolation(pass *analysis.Pass, lit *ast.BasicLit, msg, kind string) {
	fixed := removeSpecialChars(msg)
	delimiter := string(lit.Value[0])
	fixedWithQuotes := delimiter + fixed + delimiter

	pass.Report(analysis.Diagnostic{
		Pos:     lit.Pos(),
		Message: fmt.Sprintf("log message must not contain %s: %q", kind, msg),
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: fmt.Sprintf("remove %s", kind),
				TextEdits: []analysis.TextEdit{
					{
						Pos:     lit.Pos(),
						End:     lit.End(),
						NewText: []byte(fixedWithQuotes),
					},
				},
			},
		},
	})
}

func removeSpecialChars(msg string) string {
	var result strings.Builder
	for _, r := range msg {
		if r > 0x7E {
			continue
		}
		switch r {
		case '!', '?', '@', '#', '$', '%', '^', '&', '*', '~', '`', '|', '\\':
			continue
		default:
			result.WriteRune(r)
		}
	}
	return strings.TrimRight(result.String(), " ")
}

func checkSensitiveData(pass *analysis.Pass, lit *ast.BasicLit, msg string, cfg *Config) {
	msgLower := strings.ToLower(msg)

	words := strings.FieldsFunc(msgLower, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_'
	})

	for _, word := range words {
		for _, keyword := range sensitiveKeywords {
			if word == keyword {
				pass.Reportf(lit.Pos(), "log message must not contain sensitive data (found %q): %q", keyword, msg)
				return
			}
		}
		for _, keyword := range cfg.ExtraPatterns {
			if word == strings.ToLower(keyword) {
				pass.Reportf(lit.Pos(), "log message must not contain sensitive data (found %q): %q", keyword, msg)
				return
			}
		}
	}
}
