# gno

This a **very cheap** knock-off of [nobuild](https://github.com/tsoding/nobuild)

I liked the idea of having your build system in the same language as your codebase.

I have a couple of go projects using makefiles so I thought it was a good idea to
build my own implementation

To use it I recommend to fork the repo or copy [**gno.go**](./gno.go) and plug it in your codebase.
This ensures nothing breaks incase I update something.

## Usage
<details>
<summary>
Dir structure
</summary>

```console
.
├── cmd
│   └── gno.go
│
└── Rest of your codebase..
```
</details>

```go
// gno.go
import "where/you/have/the/code"

func main(){
	g := gno.New()
	g.BootstrapBuild("output_dir", "binary_name", "source")
	g.CopyResources("assets")
	g.CopyResources("more_assets")
	g.CopyResources("even_single_files")
	g.AddCommand("echo", "hello", "world")
	g.AddCommand("./output_dir/binary_name")
	g.Build() // Builds the binary
	g.RunCommandsSync()
}
```
