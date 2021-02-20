## In Memory Filesystem in Go implemented as a Trie

go test .


```go

fs := NewFileSystem()

if err := fs.Mkdir("/alo"); err != nil {
    log.Fatal(err)
}

if err := fs.WriteFile("/foo/bar/baz/hello.txt", "World."); err != nil {
    log.Fatal(err)
}

contents, err := fs.ReadFile("/foo/bar/baz/hello.txt")
if err != nil {
    log.Fatal(err)
}

// World.
log.Println(contents)

if err := fs.WriteFile("/foo/bar/baz/hello.txt", " 💖 💖 💖 "); err != nil {
    log.Fatal(err)
}

// World. 💖 💖 💖
log.Println(contents)

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
