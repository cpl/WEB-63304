package main

import (
	"fmt"
	"os"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

func main() {
	_ = os.RemoveAll("./build")
	_ = os.MkdirAll("./build/static/js", 0o744)

	// language=typescript
	tsScript := `

// import {FooA} from "src/lib/foo_a"; // this is what the autocomplete injects
import { B } from "$lib"; // this is what I've written by hand and is the expected autocomplete

// neither of the imports actually resolves with the "Go to" definition or documentation

console.debug('this is the script');
B.FooB();
	`

	esBuild("$lib", "$data")
	jsScript := esBuildContent("example_script.ts", tsScript)

	someTemplate := `
<html><head><title>Example</title></head><body>
<script>` + string(jsScript) + `</script></body></html>`

	fmt.Println(someTemplate)
}

func esBuild(targets ...string) {
	result := esbuild.Build(esbuild.BuildOptions{
		// LogLevel: esbuild.LogLevelVerbose,
		// Color: esbuild.ColorAlways,

		ResolveExtensions: []string{".ts"},
		EntryPoints:       targets,
		External:          []string{"$*"},

		Bundle: true,
		Write:  true,
		Outdir: "./build/static/js/",

		Sourcemap: esbuild.SourceMapLinked,

		Format:   esbuild.FormatESModule,
		Platform: esbuild.PlatformBrowser,

		// MinifyWhitespace:  true,
		// MinifyIdentifiers: true,
		// MinifySyntax:      true,

		Metafile:    true,
		TreeShaking: esbuild.TreeShakingFalse,
		Target:      esbuild.ES2020,
		Tsconfig:    "./.config/tsconfig.json",
	})

	for _, warn := range result.Warnings {
		fmt.Println(warn)
	}

	for _, erro := range result.Errors {
		fmt.Println(erro)
	}
}

func esBuildContent(file, content string) []byte {
	result := esbuild.Build(esbuild.BuildOptions{
		// LogLevel: esbuild.LogLevelVerbose,
		// Color: esbuild.ColorAlways,

		Stdin: &esbuild.StdinOptions{
			Contents: content,
			Loader:   esbuild.LoaderTS,
		},

		ResolveExtensions: []string{".ts"},
		External:          []string{"$*"},

		Bundle: true,

		Sourcemap: esbuild.SourceMapNone,

		Format:   esbuild.FormatESModule,
		Platform: esbuild.PlatformBrowser,

		// MinifyWhitespace:  true,
		// MinifyIdentifiers: true,
		// MinifySyntax:      true,

		Metafile:    false,
		TreeShaking: esbuild.TreeShakingFalse,
		Target:      esbuild.ES2020,
		Tsconfig:    "./.config/tsconfig.json",
	})

	for _, warn := range result.Warnings {
		fmt.Println(warn)
	}

	for _, erro := range result.Errors {
		fmt.Println(erro)
	}

	for _, f := range result.OutputFiles {
		if f.Path == "<stdout>" {
			return f.Contents
		}
	}

	return nil
}
