## In Memory Filesystem in Go implemented as a Trie

```go
go test -v .
```

```go

fs := NewFileSystem()

if err := fs.Mkdir("/alo"); err != nil {
    log.Fatal(err)
}

if err := fs.WriteFile("/foo/bar/baz/qux.txt", "Hello"); err != nil {
    log.Fatal(err)
}

if contents, err := fs.ReadFile("/foo/bar/baz/qux.txt"); err != nil {
    log.Fatal(err)
} else {
    // Hello
    log.Println(contents)
}

if err := fs.WriteFile("/foo/bar/baz/qux.txt", " World"); err != nil {
    log.Fatal(err)
}

if contents, err := fs.ReadFile("/foo/bar/baz/qux.txt"); err != nil {
    log.Fatal(err)
} else {
    // Hello World
    log.Println(contents)
}

// {
//   "": {
//     "alo": {},
//     "foo": {
//       "bar": {
//         "baz": {
//           "qux.txt": "Hello World"
//         }
//       }
//     }
//   }
// }
fmt.Println(fs.PrettyPrint())
```
