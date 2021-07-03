package gdp

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	GeneratedProjectsPath = "output"
	ProjectPath = "app"
)


type GdpApp struct {
	debug bool
	Files embed.FS
	ProjectPath string
	Templates []templateSet
	ReplaceRules map[string]interface{}
}



type templateSet struct {
	EmbedPath string
	templateFilePath string
	templateFileDirEntry fs.DirEntry
	genFilePath      string
}


func New(debug bool,Files embed.FS) (*GdpApp,error) {
	dir,err := os.Getwd();
	if err!=nil{
		return nil,err
	}
	outputPath := path.Join(dir,GeneratedProjectsPath)
	st, err :=os.Stat(outputPath)
	if err == nil &&  !st.IsDir(){
		return nil,errors.New("output file exit,please remove or move it")
	}
	if err!=nil{
		err := os.Mkdir(outputPath,os.ModePerm)
		if err!=nil{
			return nil,errors.New("make output dir fail"+err.Error())
		}
	}

	return &GdpApp{debug: debug,Files: Files,ProjectPath:path.Join(outputPath,ProjectPath)},nil
}




//generate template files
func StepOne()  {

}

//template replace
func StepTwo()  {

}


func (g *GdpApp) Generate(files embed.FS) error {
	g.AddTempToSetFromStatics()
	g.AddTempToSetFromTmpl()
	g.genFromTemplate()
	return nil
}




func (g *GdpApp)AddTempToSetFromTmpl() error {
	f,err:= fs.Sub(g.Files,"template")
	if err!=nil{
		return errors.New("read template files fail")
	}
	fs.WalkDir(f,".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir(){
			return nil
		}
		g.visit(d,path)
		return nil;
	})
	return nil
}

func (g *GdpApp)AddTempToSetFromStatics() error {
	f,err:= fs.Sub(g.Files,"statics")
	if err!=nil{
		return errors.New("read template files fail")
	}
	fs.WalkDir(f,".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir(){
			return nil
		}
		g.visit(d,path)
		return nil;
	})
	return nil
}


func (g *GdpApp) genFromTemplate() error {
	for _, tmpl := range g.Templates {
		if err := g.tmplExec(tmpl, g.ReplaceRules); err != nil {
			return err
		}
	}
	return nil
}
//
func unescaped(x string) interface{} { return template.HTML(x) }

func (g *GdpApp) tmplExec(tmplSet templateSet, d map[string]interface{}) error {
	tmpl := template.New(tmplSet.templateFilePath)
	tmpl = tmpl.Funcs(template.FuncMap{"unescaped": unescaped})
	data,err := g.Files.ReadFile(tmplSet.EmbedPath)
	tmpl,err  = tmpl.Parse(string(data))
	if err := os.MkdirAll(filepath.Dir(tmplSet.genFilePath), os.ModePerm); err != nil {
		panic(err)
	}

	dist, err := os.Create(tmplSet.genFilePath)
	if err != nil {
		panic(err)
	}
	defer dist.Close()

	//fmt.Printf("Create %s\n", distRelFilePath)
	return tmpl.Execute(dist, d)
}

func (g *GdpApp) visit(d fs.DirEntry,path string) error {
	templ := templateSet{
		templateFilePath: path,
		templateFileDirEntry: d,
	}
	if ext := filepath.Ext(path); ext == ".tmpl" {
		genFilePath := filepath.Join(g.ProjectPath,strings.TrimSuffix(path, ".tmpl") + ".go")
		templ.genFilePath = genFilePath
		templ.EmbedPath = filepath.Join("./template",path)
	}else {
		genFilePath := filepath.Join(g.ProjectPath,path)
		templ.EmbedPath = filepath.Join("./statics",path)
		templ.genFilePath = genFilePath
	}

	g.Templates = append(g.Templates, templ)
	return nil
}

func (s *GdpApp) debugPrintf(format string, a ...interface{}) {
	if s.debug == true {
		fmt.Printf(format, a...)
	}
}