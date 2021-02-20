## In Memory Filesystem in Go implemented as a Trie

```go
go test -v .
```

```go

fs := NewFileSystem()

if err := fs.Mkdir("/alo"); err != nil {
    log.Fatal(err)
}

if err := fs.WriteFile("/foo/bar/baz/hello.txt", "World."); err != nil {
    log.Fatal(err)
}

if contents, err := fs.ReadFile("/foo/bar/baz/hello.txt"); err != nil {
    log.Fatal(err)
} else {
    // World.
    log.Println(contents)
}

if err := fs.WriteFile("/foo/bar/baz/hello.txt", " 💖 💖 💖 "); err != nil {
    log.Fatal(err)
}

if contents, err := fs.ReadFile("/foo/bar/baz/hello.txt"); err != nil {
    log.Fatal(err)
} else {
    // World. 💖 💖 💖
    log.Println(contents)
}

// {
//   "": {
//     "alo": {},
//     "foo": {
//       "bar": {
//         "baz": {
//           "hello.txt": "World. 💖 💖 💖 "
//         }
//       }
//     }
//   }
// }
fmt.Println(fs.PrettyPrint())
```
