![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/ucanme/gdp/master?color=%239&logoColor=red&style=for-the-badge)

## introduce
###### Gdp is a tool to generate any project layout of any language not just golang. It works only in go1.16+ (only go1.16+ support embeld).
###### Just put you template files in direcition template, then execute the command
```shell script
go run main.go new
```

follow the prompt then you project will be generated in the directory output.

##how to use

please ensure your golang version is above go1.16+
```shell script
go install github.com/ucanme/gdp
```
in any direction 
```shell script
gdp new
```
then get the generated code in output directory