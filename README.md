# go-code-generate

> Commandline script to generate boilerplate user files on a project by project basis.
>
> 
## Usage
`go-code-generate [option] [arguments]`

### Generating User Code

### Usage

`go-code-generate user --domain-pkg-name [DOMAIN PKG] --domain-pkg-import-path [DOMAIN PKG PATH] --db-connection-url [DB CONN STRING] --output-path [PATH]`

#### Files Generated
* `db/db.go`
* `db/user_accessor.go`
* `domain/user.go`