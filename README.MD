running application from command line, skipping test files:
    go run $(ls -t | grep -v _test | grep .go) 8080

GoLang project Directory structure:
    https://medium.com/golang-learn/go-project-layout-e5213cdcfaa2
    https://vsupalov.com/go-folder-structure/
    https://github.com/golang-standards/project-layout
    
Building a web application:
    https://www.sohamkamani.com/blog/2017/09/13/how-to-build-a-web-application-in-golang/
    
JSON Parsing:
    https://www.sohamkamani.com/blog/2017/10/18/parsing-json-in-golang/
    https://blog.golang.org/json-and-go
    
error handling:
    https://golangbot.com/panic-and-recover/
    
git pre-commit hook is in file ".git/hooks/pre-commit". Modified to run "go fmt" automatically before every commit
    https://go.googlesource.com/go/+/dev.tls/misc/git/pre-commit