package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Mldlr/url-shortener/internal/analyzers"
	critic "github.com/go-critic/go-critic/checkers/analyzer"
	useStdLib "github.com/sashamelentyev/usestdlibvars/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

// Config is the name of the configuration file for this program.
const Config = `config.json`

// ConfigData represents the data structure of the configuration file.
type ConfigData struct {
	Staticcheck []string
	Simple      []string
	Stylecheck  []string
}

func main() {
	// Get the path to the executable file.
	appFile, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// Read the configuration file.
	data, err := os.ReadFile(filepath.Join(filepath.Dir(appFile), Config))
	if err != nil {
		panic(err)
	}
	var cfg ConfigData
	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	// Initialize the list of analyzers to run.
	mychecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		critic.Analyzer,
		analyzers.Analyzer,
		useStdLib.New(),
	}
	// Create a map of the enabled analyzers.
	checks := make(map[string]bool)
	for _, v := range cfg.Staticcheck {
		checks[v] = true
	}
	// Add enabled analyzers to the list of analyzers to run.
	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	for _, v := range simple.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	for _, v := range stylecheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	// Run the analyzers.
	multichecker.Main(
		mychecks...,
	)
}
