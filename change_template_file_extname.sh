find  ./template  -name "*.tmpl" | awk -F "." '{print $2}' | xargs -I {}  mv ./{}.tmpl ./{}.go
