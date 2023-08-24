package main

import (
	"flag"
	"work_report/registry"
	"work_report/usecase"
)

var (
	mode           = flag.String("mode", "gen", `"gen" or "send" (default: "gen")`)
	inputFilePath  = "./data/jira.html"
	outputFilePath = "./data/report.json"
)

func main() {
	flag.Parse()

	if *mode == "send" {
		send()
	} else {
		gen()
	}
}

func mainUseCase() usecase.MainUseCase {
	r, err := registry.NewRegistry("./config/settings.json")
	if err != nil {
		panic(err)
	}

	return r.NewMainUseCase()
}

func send() {
	if err := mainUseCase().SendReport(outputFilePath); err != nil {
		panic(err)
	}
}

func gen() {
	if err := mainUseCase().GenerateReport(inputFilePath, outputFilePath); err != nil {
		panic(err)
	}
}
