package converter

import (
	"strings"
	"testing"
)

func TestProcessHTML(t *testing.T) {
	input := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<header><h1>Header Content</h1></header>
	<main><p>Main Content</p></main>
	<footer><p>Footer Content</p></footer>
</body>
</html>`

	result, err := ProcessHTML([]byte(input), StripConfig{ElementsToStrip: []string{}})
	if err != nil {
		t.Fatalf("processHTML failed: %v", err)
	}

	resultStr := string(result)

	// Check that header and footer are removed
	if strings.Contains(resultStr, "<header>") {
		t.Error("Result still contains <header> tag")
	}
	if strings.Contains(resultStr, "<footer>") {
		t.Error("Result still contains <footer> tag")
	}

	// Check that main content is preserved
	if !strings.Contains(resultStr, "<main>") {
		t.Error("Result missing <main> tag")
	}
	if !strings.Contains(resultStr, "Main Content") {
		t.Error("Result missing main content")
	}
}

func TestProcessHTMLWithElementsToStrip(t *testing.T) {
	input := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<div><p>Content to keep</p></div>
	<span><p>Content to strip</p></span>
</body>
</html>`

	result, err := ProcessHTML([]byte(input), StripConfig{ElementsToStrip: []string{"span"}})
	if err != nil {
		t.Fatalf("processHTML failed: %v", err)
	}

	resultStr := string(result)

	// Check that span is removed
	if strings.Contains(resultStr, "<span>") {
		t.Error("Result still contains <span> tag")
	}
}

func TestProcessHTMLWithDataLLMKeep(t *testing.T) {
	input := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<header data-llm="keep"><h1>Header Content</h1></header>
	<footer><p>Footer Content</p></footer>
</body>
</html>`

	result, err := ProcessHTML([]byte(input), StripConfig{ElementsToStrip: []string{}})
	if err != nil {
		t.Fatalf("processHTML failed: %v", err)
	}

	resultStr := string(result)
	// Check that header with data-llm="keep" is preserved
	if !strings.Contains(resultStr, "<header data-llm=\"keep\">") {
		t.Error("Result missing <header> tag with data-llm=\"keep\"")
	}
	// Check that footer is removed
	if strings.Contains(resultStr, "<footer>") {
		t.Error("Result still contains <footer> tag")
	}
}

func TestProcessHTMLWithDataLLMDrop(t *testing.T) {
	input := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<div data-llm="drop"><h1>drop this</h1></header>
	<footer><p>Footer Content</p></footer>
</body>
</html>`
	result, err := ProcessHTML([]byte(input), StripConfig{ElementsToStrip: []string{}})
	if err != nil {
		t.Fatalf("processHTML failed: %v", err)
	}

	resultStr := string(result)

	if strings.Contains(resultStr, "<div data-llm=\"drop\">") {
		t.Error("Result still contains <div> tag with data-llm=\"drop\"")
	}

	if strings.Contains(resultStr, "drop this") {
		t.Error("Result still contains <div> tag with data-llm=\"drop\"")
	}
}
