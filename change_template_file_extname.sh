find  ./template  -name "*.go" | awk -F "." '{print $2}' | xargs -I {}  mv ./{}.go ./{}.tmpl
