package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	targetFile := "../frontend/src/types/api.ts"
	content, err := os.ReadFile(targetFile)
	if err != nil {
		panic(err)
	}
	
	sContent := string(content)

	// Use reflection to find all Values() functions in the domain package
	// Since we can't easily iterate over package functions using reflect in Go (it requires the package to be linked and symbols to be available),
	// and we are running this as a script, we will use a slightly different approach.
	// Actually, we can just manually list them OR we can use the fact that we just generated the _enumer.go files.
	// But the user wants automation.
	// 
	// A robust way in a script like this without importing "reflect" and iterating (which isn't possible for package level funcs easily)
	// is to rely on the known pattern of generated functions.
	// However, since we are importing `gcp_antigravity/backend/internal/domain`, we can try to use `reflect` on specific known types? No, that defeats the purpose.
	//
	// Let's use `go/doc` or `go/parser` again here? No, that's too heavy.
	//
	// Wait, we can just iterate over the known enums if we had a list.
	// But we want to discover them.
	//
	// Actually, since we are in the same repo, we can parse the `internal/domain` directory again to find the `Values` functions.
	// OR, we can just hardcode the list for now? No, the user explicitly asked for automation.
	//
	// Let's use the `domain` package and a map of "Type" -> "Values".
	// But we can't iterate over package exports dynamically in Go without some registry.
	//
	// Alternative: The `gen_enum_methods.go` could ALSO output a JSON file with all the enums and their values.
	// Then this script just reads that JSON and applies it.
	// This is much cleaner and decouples the Go reflection limitation.
	
	// Let's change the plan slightly:
	// 1. gen_enum_methods.go generates `enums.json` in `internal/domain` (or temp dir).
	// 2. gen_ts_enums.go reads `enums.json` and applies changes.
	
	// But I cannot change gen_enum_methods.go in this tool call (I already did).
	// Let's stick to this file.
	// I will use `go/parser` here too to find the `Values` functions in `internal/domain`.
	
	// ... Actually, I can just use the `domain` package if I register the enums.
	// But I can't modify domain code to register itself.
	
	// Let's go with parsing the `internal/domain` directory to find `func (T) Values() []string` or `func TValues() []string`.
	// Since `gen_enum_methods.go` generates `func (T) Values()`, we can look for those in `*_enumer.go` files!
	
	enums := findEnumsFromGeneratedFiles("internal/domain")
	
	for typeName, values := range enums {
		sContent = replaceEnumType(sContent, typeName, values)
	}

	err = os.WriteFile(targetFile, []byte(sContent), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully updated TypeScript enums in api.ts")
}

func findEnumsFromGeneratedFiles(dir string) map[string][]string {
	// This is a simplified parser that looks for the generated pattern in _enumer.go files
	// func (UserRole) Values() []string { return []string{ "free", "pro", "admin", } }
	
	enums := make(map[string][]string)
	
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "_enumer.go") {
			content, _ := os.ReadFile(dir + "/" + f.Name())
			s := string(content)
			
			// Extract Type Name
			// func (UserRole) Values()
			start := strings.Index(s, "func (")
			if start == -1 { continue }
			end := strings.Index(s[start:], ") Values()")
			if end == -1 { continue }
			typeName := s[start+6 : start+end]
			
			// Extract Values
			// return []string{ "free", "pro", "admin", }
			valStart := strings.Index(s, "return []string{")
			if valStart == -1 { continue }
			valEnd := strings.Index(s[valStart:], "}")
			if valEnd == -1 { continue }
			
			block := s[valStart+16 : valStart+valEnd]
			parts := strings.Split(block, ",")
			var values []string
			for _, p := range parts {
				p = strings.TrimSpace(p)
				p = strings.Trim(p, "\"")
				if p != "" {
					values = append(values, p)
				}
			}
			enums[typeName] = values
		}
	}
	return enums
}

func replaceEnumType[T ~string](content string, typeName string, values []T) string {
	// Construct the union type string: "value1" | "value2" | ...
	var quotedValues []string
	for _, v := range values {
		quotedValues = append(quotedValues, fmt.Sprintf("\"%s\"", v))
	}
	unionType := strings.Join(quotedValues, " | ")
	
	// Regex-like replacement (simple string replacement for now as we know the format from Tygo)
	// Tygo generates: export type TypeName = string;
	target := fmt.Sprintf("export type %s = string;", typeName)
	replacement := fmt.Sprintf("export type %s = %s;", typeName, unionType)
	
	return strings.ReplaceAll(content, target, replacement)
}
