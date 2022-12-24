package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go/format"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
)

type UserCodeGenerationFields struct {
	DomainPkgName       string
	DomainPkgImportPath string
	DBConnectionUrl     string
}

func NewUserCommand() *cobra.Command {
	var domainPkgName string
	var domainPkgImportPath string
	var dbConnUrl string
	var outputPath string
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Generate user go code",
		Long:  ``,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			data := UserCodeGenerationFields{
				DomainPkgName:       domainPkgName,
				DBConnectionUrl:     dbConnUrl,
				DomainPkgImportPath: domainPkgImportPath,
			}

			generateUserCode(data, outputPath)
		},
	}

	cmd.Flags().StringVarP(&domainPkgName, "domain-pkg-name", "d", "", "--domain-pkg-name")
	cmd.Flags().StringVarP(&domainPkgImportPath, "domain-pkg-import-path", "p", "", "--domain-pkg-import-path")
	cmd.Flags().StringVarP(&dbConnUrl, "db-connection-url", "u", "", "--db-connection-url")
	cmd.Flags().StringVarP(&outputPath, "output-path", "o", ".", "--output-path")

	return cmd
}

func generateUserCode(data UserCodeGenerationFields, outputDir string) {
	t, err := loadTemplates()
	if err != nil {
		fmt.Println(err)
		return
	}

	processTemplate(t, "db.tmpl", path.Join(outputDir, "db", "db.go"), data)
	processTemplate(t, "user_accessor.tmpl", path.Join(outputDir, "db", "user_accessor.go"), data)
	processTemplate(t, "user.tmpl", path.Join(outputDir, "domain", "user.go"), data)
}

func loadTemplates() (*template.Template, error) {
	box, err := rice.FindBox("./templates")
	if err != nil {
		return nil, err
	}

	var filenames []string
	_ = box.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		filenames = append(filenames, path)
		return nil
	})

	var tmpl *template.Template
	for _, filename := range filenames {
		name := path.Base(filename)
		s, err := box.String(name)
		if err != nil {
			return nil, errors.Wrap(err, "box.String")
		}
		if tmpl == nil {
			tmpl = template.Must(template.New(name).Parse(s))
		} else {
			tmpl = template.Must(tmpl.New(name).Parse(s))
		}
	}

	return tmpl, nil
}

func processTemplate(templates *template.Template, templateName string, outputFilePath string, data UserCodeGenerationFields) {
	var processed bytes.Buffer
	err := templates.ExecuteTemplate(&processed, templateName, data)
	if err != nil {
		log.Fatalf("Unable to parse data into template: %v\n", err)
	}
	formatted, err := format.Source(processed.Bytes())
	if err != nil {
		log.Fatalf("Could not format processed template: %v\n", err)
	}

	fmt.Println("Writing file: ", outputFilePath)
	f, err := create(outputFilePath)

	if err != nil {
		fmt.Println(err.Error())
	}
	w := bufio.NewWriter(f)
	w.WriteString(string(formatted))
	w.Flush()
}

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}
