//go:generate go run ./var_gen.go

package main

import (
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"text/template"

	"github.com/gobuffalo/flect"
)

type Var struct {
	Name       string
	OmitGetter bool
	OmitPrefix bool
	OmitDocs   bool
	Computed   string
	DefaultTo  string
	Categories []string
	Desc       string
}

type Vars []Var

func (v Vars) Len() int {
	return len(v)
}

func (v Vars) Less(i, j int) bool {
	return v[i].Name < v[j].Name
}

func (v Vars) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

type varOp func(*Var)

func omitGetter(v *Var) {
	v.OmitGetter = true
}
func omitPrefix(v *Var) {
	v.OmitPrefix = true
}
func unused(v *Var) {
	v.OmitDocs = true
}
func build(v *Var) {
	v.Categories = append(v.Categories, "Build")
}
func deploy(v *Var) {
	v.Categories = append(v.Categories, "Deploy")
}
func computed(v *Var) {
	v.Computed = "yes"
}
func defaultTo(s string) varOp {
	return func(v *Var) {
		v.DefaultTo = s
	}
}

func nVar(name, desc string, ops ...varOp) Var {
	v := Var{Name: name, Desc: desc}
	for _, o := range ops {
		o(&v)
	}
	return v
}

var vars = Vars{
	nVar("Application", "If provided, name of the (sub)application to compile", unused, computed),
	nVar("ComposeFile", "Compose file(s) for deploying", omitPrefix, deploy),
	nVar("ComposeTemplate", "A template docker-compose file that may contain mutations for the compose file", deploy),
	nVar("Context", "Docker build context", build, defaultTo("<project root>")),
	nVar("DefaultBranch", "Target branch to compare for BELUGA_ENVIROMENT.", omitGetter, defaultTo("`master`")),
	nVar("Dockerfile", "Dockerfile to build", build, defaultTo("`Dockerfile` in context")),
	nVar("DockerHost", "Docker instance for building/deployin", build, deploy, omitGetter, omitPrefix),
	nVar("Domain", "Domain name of the stack", deploy, computed),
	nVar("Environment", "Environment name", deploy, computed, defaultTo("`review`, `staging`, or `production`")),
	nVar("Image", "First image listed in BELUGA_IMAGES; doesn't affect pushng", computed, deploy),
	nVar("Images", "Docker images to push after build", computed, build, deploy),
	nVar("ImagesTemplate", "Go template for a space-separated list of Docker images to push after build", defaultTo("yes")),
	nVar("Overrides", "YAML document with environment names or patterns as keys and variables to override as values"),
	nVar("Registry", "Docker registry for pushing", build),
	nVar("RegistryPassword", "Password for Docker registry", build),
	nVar("RegistryUsername", "Username for Docker registry", build),
	nVar("StackName", "Name of the compose/swamrm/etc. stack", deploy),
	nVar("Version", "Version of the application being built/deployed", computed, deploy),
}

var varsByCategory map[string]Vars
var categories []string

func init() {
	m := make(map[string]Vars)
	uncategorized := Vars{}
	for _, v := range vars {
		for _, c := range v.Categories {
			m[c] = append(m[c], v)
		}
		if len(v.Categories) == 0 {
			uncategorized = append(uncategorized, v)
		}
	}
	m["Uncategorized"] = uncategorized

	categories = make([]string, 0, len(m))
	for c := range m {
		categories = append(categories, c)
	}
	sort.Strings(categories)
	varsByCategory = m
}

var t = template.Must(template.New("").Parse(`
// DO NOT EDIT: This file is generated by var_gen.go

package beluga

const (
	{{- range .Variables}}
	{{ .VarName }} = "{{ .EnvName }}"
	{{- end }}
)

var knownVarNames = []string{
	{{range .Variables}}
	{{- .VarName }},
{{end}} }

{{ range .Variables}}
{{- if not .OmitGetter }}
{{.Comment}}
func (e Environment) {{ .GoName }}() string {
	return e[{{ .VarName }}]
}
{{ end -}}
{{ end }}
`))

var mt = template.Must(template.New("").Parse(`
{{- define "table" -}}
| Name | Description | Default | Computed |
| ---- | ----------- | ------- | -------- |
{{- range . -}}
	{{- if not .OmitDocs }}
| {{ .EnvName }} | {{ .Desc }} | {{ .DefaultTo }} | {{ .Computed }} |
	{{- end -}}
{{- end -}}
{{- end -}}

<!-- DO NOT EDIT: This file is generated by var_gen.go -->

# Beluga Environment Variables

{{ range .Categories -}}

## {{ . }}

{{ template "table" (index $.VariablesByCategory .) }}

{{ end -}}
`))

type Ctx struct {
	Variables           Vars
	Categories          []string
	VariablesByCategory map[string]Vars
}

func templateToFile(t *template.Template, filename string) error {
	output, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = t.Execute(output, Ctx{
		Variables:           vars,
		Categories:          categories,
		VariablesByCategory: varsByCategory,
	})
	if err != nil {
		return err
	}
	err = output.Close()
	if err != nil {
		return err
	}
	return nil
}

const docsDir = "../../docs"

func main() {
	sort.Sort(vars)
	const filename = "../../var.go"

	err := templateToFile(t, filename)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(docsDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = templateToFile(mt, path.Join(docsDir, "variables.md"))
	if err != nil {
		panic(err)
	}

	err = exec.Command("go", "fmt", filename).Run()
	if err != nil {
		panic(err)
	}
}

func (v Var) GoName() string {
	return flect.Pascalize(v.Name)
}

func (v Var) VarName() string {
	return "var" + flect.Pascalize(v.Name)
}

func (v Var) EnvName() string {
	n := strings.ToUpper(flect.Underscore(v.Name))
	if !v.OmitPrefix {
		n = "BELUGA_" + n
	}
	return n
}

func (v Var) Comment() string {
	if v.Desc == "" {
		return ""
	}
	const maxLength = 80 - 3 // 3 chars for comment
	comment := "//"
	lineLength := 0
	for _, word := range strings.Split(v.Desc, " ") {
		lineLength += len(word) + 1
		if lineLength >= maxLength {
			comment += "\n// " + word
			lineLength = len(word) + 1
		} else {
			comment += " " + word
		}
	}
	return comment
}
