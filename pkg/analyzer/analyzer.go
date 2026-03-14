package analyzer

import (
	"go/ast"
	"go/constant"
	"go/token"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var (
	cfgMu sync.RWMutex
	cfg   = DefaultConfig
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglint",
	Doc:      "checks log messages style",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func init() {
	registerFlags(&Analyzer.Flags, &cfg)
}

var loggerMethods = map[string]map[string]int{
	"slog": {
		"Info":         0,
		"Error":        0,
		"Warn":         0,
		"Debug":        0,
		"InfoContext":  1,
		"ErrorContext": 1,
		"WarnContext":  1,
		"DebugContext": 1,
	},
	"zap": {
		"Info":  0,
		"Error": 0,
		"Warn":  0,
		"Debug": 0,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	cfgMu.RLock()
	currentCfg := cfg
	cfgMu.RUnlock()

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		methodName := sel.Sel.Name
		var msgArgIdx int

		if ident, ok := sel.X.(*ast.Ident); ok {
			methods, ok := loggerMethods[ident.Name]
			if !ok {
				return
			}
			idx, ok := methods[methodName]
			if !ok {
				return
			}
			msgArgIdx = idx
		} else if innerCall, ok := sel.X.(*ast.CallExpr); ok {
			innerSel, ok := innerCall.Fun.(*ast.SelectorExpr)
			if !ok {
				return
			}
			innerIdent, ok := innerSel.X.(*ast.Ident)
			if !ok {
				return
			}
			if innerIdent.Name != "zap" {
				return
			}
			if innerSel.Sel.Name != "L" && innerSel.Sel.Name != "S" {
				return
			}
			idx, ok := loggerMethods["zap"][methodName]
			if !ok {
				return
			}
			msgArgIdx = idx
		} else {
			return
		}

		if len(call.Args) <= msgArgIdx {
			return
		}

		lit, ok := call.Args[msgArgIdx].(*ast.BasicLit)
		if !ok {
			return
		}

		if len(lit.Value) < 2 {
			return
		}
		msg := constant.StringVal(constant.MakeFromLiteral(lit.Value, token.STRING, 0))

		checkMessage(pass, lit, msg, &currentCfg)
	})

	return nil, nil
}

func SetConfig(c Config) {
	cfgMu.Lock()
	defer cfgMu.Unlock()
	cfg = c
}
