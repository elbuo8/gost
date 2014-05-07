# Gost

Simple package to interact with [GitHub's Gist API](developer.github.com/v3/gists/).

### Installation

```bash
go get github.com/elbuo8/gost
```

### Example

```go
gostClient := gost.New("TOKEN")

gist := &gost.GistFile{Filename: "Test.go", Content: "package go"}

_, err := g.Create("SampleFile", false, gist)
if err != nil {
  t.Errorf("%v", err)
}
```

### Contributing

Feel free to submit pull requests after running:
```bash
go test
```

### MIT License