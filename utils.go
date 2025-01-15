package goblet

import (
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

func UrlFor(endpoint string) string {
	return fmt.Sprintf("/%s", endpoint)
}

var extendsRegex = regexp.MustCompile(`\{\{\s*extends\s*\"(.*?)\"\}\}`)

func parseTemplate(baseTemplate *template.Template, templateName string, templates map[string]string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"url_for": UrlFor,
		"extends": func(string, ...interface{}) string { return "" },
		"block":   func(string, interface{}) string { return "" },
		"define":  func(string) string { return "" },
		"if":      func(interface{}) string { return "" },
		"else":    func() string { return "" },
		"end":     func() string { return "" },
		"range":   func(interface{}) string { return "" },
		"with":    func(interface{}) string { return "" },
		"template": func(string, ...interface{}) string { return "" },
	}

	tmpl, err := baseTemplate.Clone()
	if err != nil {
		return nil, fmt.Errorf("error cloning base template: %v", err)
	}
	tmpl = tmpl.Funcs(funcMap)

	content, ok := templates[templateName]
	if !ok {
		return nil, fmt.Errorf("template %s not found", templateName)
	}

	match := extendsRegex.FindStringSubmatch(content)
	if match != nil {
		parentName := match[1]
		parentContent, ok := templates[parentName]
		if !ok {
			return nil, fmt.Errorf("parent template %s not found", parentName)
		}

		content = extendsRegex.ReplaceAllString(content, "")
		parentContent = replaceBlocks(parentContent, content)

		_, err = tmpl.Parse(parentContent)
		if err != nil {
			return nil, fmt.Errorf("error parsing parent template %s: %v", parentName, err)
		}
	} else {
		_, err = tmpl.Parse(content)
		if err != nil {
			return nil, fmt.Errorf("error parsing template %s: %v", templateName, err)
		}
	}

	return tmpl, nil
}

func replaceBlocks(parentContent, childContent string) string {
	childBlocks := extractBlocks(childContent)

	for blockName, blockContent := range childBlocks {
		parentContent = replaceBlockContent(parentContent, blockName, blockContent)
	}

	return parentContent
}

func extractBlocks(content string) map[string]string {
	re := regexp.MustCompile(`\{\{\s*block\s+\"(\w+)\"\s*\.\s*\}\}(.*?)\{\{\s*end\s*\}\}`)
	matches := re.FindAllStringSubmatch(content, -1)

	blocks := make(map[string]string)
	for _, match := range matches {
		blockName := match[1]
		blockContent := match[2]
		blocks[blockName] = blockContent
	}

	return blocks
}

func replaceBlockContent(content, blockName, blockContent string) string {
	re := regexp.MustCompile(fmt.Sprintf(`\{\{\s*block\s+\"%s\"\s*\.\s*\}\}(.*?)\{\{\s*end\s*\}\}`, blockName))
	return re.ReplaceAllString(content, fmt.Sprintf("{{block \"%s\" .}}%s{{end}}", blockName, blockContent))
}

func Extends(tpl *template.Template, name string, data interface{}) (*template.Template, error) {
	templates, err := loadTemplateFiles("templates")
	if err != nil {
		return nil, err
	}

	finalTmpl, err := parseTemplate(tpl, name, templates)
	if err != nil {
		return nil, err
	}

	return finalTmpl, nil
}

func loadTemplateFiles(dir string) (map[string]string, error) {
	templates := make(map[string]string)
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".html") {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			relPath, _ := filepath.Rel(dir, path)
			templates[relPath] = string(content)
		}
		return nil
	})
	return templates, err
}
