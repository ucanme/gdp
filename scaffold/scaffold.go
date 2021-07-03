package scaffold

import (
	"fmt"
	pkgErr "github.com/pkg/errors"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	GoScaffoldPath = "./out"
	ProjectPath = ""
)


type scaffold struct {
	debug bool
}

func New(debug bool) *scaffold {
	return &scaffold{debug: debug}
}

func (s *scaffold) Generate(path string) error {
	genAbsDir, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	projectName := filepath.Base(genAbsDir)
	//TODO: have to check path MUST be under the $GOPATH/src folder
	goProjectPath := path

	ProjectPath = goProjectPath
	d := data{
		AbsGenProjectPath: genAbsDir,
		ProjectPath:       goProjectPath,
		ProjectName:       projectName,
		Quit:              "<-quit",
	}

	fmt.Println("path",goProjectPath)
	if err := s.genFromTemplate(getTemplateSets(goProjectPath), d); err != nil {
		return err
	}

	if err := s.genFormStaticFle(d); err != nil {
		return err
	}


	return nil
}

type data struct {
	AbsGenProjectPath string // The Abs Gen Project Path
	ProjectPath       string //The Go import project path (eg:github.com/fooOrg/foo)
	ProjectName       string //The project name which want to generated
	Quit              string
}

type templateEngine struct {
	Templates []templateSet
	currDir   string
}

type templateSet struct {
	templateFilePath string
	templateFileName string
	genFilePath      string
}

func getTemplateSets(path string) []templateSet {
	tt := templateEngine{}
	//fmt.Printf("walk:%s\n", templatesFolder)
	path = filepath.Join(path,"./template")
	fmt.Println(path)
	err := filepath.Walk(path, tt.visit)
	fmt.Println(err)
	fmt.Println(tt.Templates)
	return tt.Templates
}

func (s *scaffold) genFromTemplate(templateSets []templateSet, d data) error {
	for _, tmpl := range templateSets {
		fmt.Println(tmpl,"tmpl")
		if err := s.tmplExec(tmpl, d); err != nil {
			return err
		}
	}
	return nil
}

func unescaped(x string) interface{} { return template.HTML(x) }

func (s *scaffold) tmplExec(tmplSet templateSet, d data) error {
	tmpl := template.New(tmplSet.templateFileName)
	tmpl = tmpl.Funcs(template.FuncMap{"unescaped": unescaped})
	tmpl, err := tmpl.ParseFiles(tmplSet.templateFilePath)
	if err != nil {
		return pkgErr.WithStack(err)
	}

	relateDir := filepath.Dir(tmplSet.genFilePath)

	distRelFilePath := filepath.Join(relateDir, filepath.Base(tmplSet.genFilePath))
	distAbsFilePath := filepath.Join(d.AbsGenProjectPath, distRelFilePath)


	if err := os.MkdirAll(filepath.Dir(distAbsFilePath), os.ModePerm); err != nil {
		panic(err)
		return pkgErr.WithStack(err)
	}

	dist, err := os.Create(distAbsFilePath)
	if err != nil {
		return pkgErr.WithStack(err)
	}
	defer dist.Close()

	//fmt.Printf("Create %s\n", distRelFilePath)
	return tmpl.Execute(dist, d)
}

func (templEngine *templateEngine) visit(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if ext := filepath.Ext(path); ext == ".tmpl" {
		templateFileName := filepath.Base(path)

		genFileBaeName := strings.TrimSuffix(templateFileName, ".tmpl") + ".go"
		genFileBasePath, err := filepath.Rel(filepath.Join(ProjectPath,"template"),path)
		if err != nil {
			return pkgErr.WithStack(err)
		}

		templ := templateSet{
			templateFilePath: path,
			templateFileName: templateFileName,
			genFilePath:      filepath.Join(GoScaffoldPath,filepath.Dir(filepath.Join(templEngine.currDir, genFileBasePath)),genFileBaeName),
		}

		templEngine.Templates = append(templEngine.Templates, templ)

	} else if mode := f.Mode(); mode.IsRegular() {
		templateFileName := filepath.Base(path)

		basepath := filepath.Join( path,GoScaffoldPath, "template")
		targpath := filepath.Join(filepath.Dir(path), templateFileName)
		genFileBasePath, err := filepath.Rel(basepath, targpath)
		if err != nil {
			return pkgErr.WithStack(err)
		}

		templ := templateSet{
			templateFilePath: path,
			templateFileName: templateFileName,
			genFilePath:      filepath.Join(templEngine.currDir, genFileBasePath),
		}

		templEngine.Templates = append(templEngine.Templates, templ)
	}

	return nil
}

func (s *scaffold) genFormStaticFle(d data) error {
	fmt.Println("aaaaaaaa",d)
	walkerFuc := func(path string, f os.FileInfo, err error) error {
		if f.Mode().IsRegular() == true {
			src, err := os.Open(path)
			if err != nil {
				return pkgErr.WithStack(err)
			}
			defer src.Close()

			basepath := filepath.Join(ProjectPath,"static")


			distRelFilePath, err := filepath.Rel(basepath, path)
			if err != nil {
				return pkgErr.WithStack(err)
			}

			distAbsFilePath := filepath.Join(d.AbsGenProjectPath,GoScaffoldPath, distRelFilePath)

			if err := os.MkdirAll(filepath.Dir(distAbsFilePath), os.ModePerm); err != nil {
				return pkgErr.WithStack(err)
			}

			dist, err := os.Create(distAbsFilePath)
			if err != nil {
				return pkgErr.WithStack(err)
			}
			defer dist.Close()

			if _, err := io.Copy(dist, src); err != nil {
				return pkgErr.WithStack(err)
			}

			fmt.Printf("Create %s \n", distRelFilePath)
		}

		return nil
	}

	walkPath := filepath.Join( ProjectPath, "static")
	return filepath.Walk(walkPath, walkerFuc)
}

func (s *scaffold) debugPrintf(format string, a ...interface{}) {
	if s.debug == true {
		fmt.Printf(format, a...)
	}
}
