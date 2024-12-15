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

// regex searches for: {{ extends "yada yada yada" }}
var extendsRegex = regexp.MustCompile(`\{\{\s*extends\s*\"(.*?)\"\}\}`)

func parseTemplate(baseTemplate, templateName string, templates map[string]string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"url_for": UrlFor,
	}
	tmpl := template.New(templateName).Funcs(funcMap)

	var processTemplate func(name string) (*template.Template, error)
	processTemplate = func(name string) (*template.Template, error) {
		content, ok := templates[name]
		if !ok {
			return nil, fmt.Errorf("template %s not found", name)
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

			var err error
			tmpl, err = tmpl.Parse(parentContent)
			if err != nil {
				return nil, fmt.Errorf("error parsing parent template %s: %v", parentName, err)
			}
		} else {
			var err error
			tmpl, err = tmpl.Parse(content)
			if err != nil {
				return nil, fmt.Errorf("error parsing template %s: %v", name, err)
			}
		}

		return tmpl, nil
	}

	return processTemplate(templateName)
}

func replaceBlocks(parentContent, childContent string) string {
	childBlocks := extractBlocks(childContent)

	for blockName, blockContent := range childBlocks {
		parentContent = replaceBlockContent(parentContent, blockName, blockContent)
	}

	return parentContent
}

func extractBlocks(content string) map[string]string {
	re := regexp.MustCompile(`\{\{\s*block\s+\"(\w+)\"\s*\}\}(.*?)\{\{\s*end\s*\}\}`)
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

	finalTmpl, err := parseTemplate("base.html", name, templates)
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