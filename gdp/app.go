package gdp

import (
	"embed"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"html/template"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

var (
	GeneratedProjectsPath = "output"
	ProjectPath = "app"
)


type GdpApp struct {
	Name string
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


func New(name string,debug bool,Files embed.FS) (*GdpApp,error) {
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

	if st,err := os.Stat(path.Join(outputPath,name));err == nil && st.IsDir(){
		return nil,errors.New("directory output/"+name+" exist please remove it")
	}


	replaceRules,err := loadRules()
	if err!=nil{
		return nil,err
	}

	replaceRules["app_name"] = name
	return &GdpApp{Name:name,debug: debug,Files: Files,ProjectPath:path.Join(outputPath,name),ReplaceRules: replaceRules},nil
}

func loadRules() (map[string]interface{},error) {
	var err error
	rules := map[string]interface{}{}
	if ok,_ := PathExists("./replace_rules.toml");ok{
		_,err = toml.DecodeFile("./replace_rules.toml",&rules)
		if err!=nil{
			return rules,errors.New("load replace rules fail"+err.Error())
		}
	}
	return rules,nil
}

func PathExists(path string) (bool, error) {

	st, err := os.Stat(path)
	if err == nil && !st.IsDir()  {
		return true, nil
	}
	if err== nil && st.IsDir(){
		return false,errors.New(path +" is a dir")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (g *GdpApp) Generate(files embed.FS) error {
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

	genFilePath := filepath.Join(g.ProjectPath,path)
	templ.genFilePath = genFilePath
	templ.EmbedPath = filepath.Join("./template",path)
	g.Templates = append(g.Templates, templ)
	return nil
}

func (s *GdpApp) debugPrintf(format string, a ...interface{}) {
	if s.debug == true {
		fmt.Printf(format, a...)
	}
}
