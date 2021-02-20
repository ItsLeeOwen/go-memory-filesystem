package filesystem

import (
	"fmt"
	"testing"
	"time"
)

var expectPrettyPrint string = `{
  "": {
    "a": {
      "b": {
        "c": {
          "hello.txt": "World."
        }
      }
    },
    "foo": {
      "ahoy": {},
      "bar": {
        "baz": {
          "info.txt": "Satoshi was here. 💖 💖 💖 "
        }
      },
      "zap": {}
    }
  }
}`

func TestFileSystem(t *testing.T) {
	type test = func(fs FileSystem) bool

	var tests = func() []test {
		return []test{
			func(fs FileSystem) bool {
				err := fs.Mkdir("/foo/bar/baz")
				if err != nil {
					fmt.Println("got err", err, "expected nil")
					return false
				}
				return true
			},
			func(fs FileSystem) bool {
				err := fs.Mkdir("/foo/ahoy")
				return err == nil
			},
			func(fs FileSystem) bool {
				err := fs.Mkdir("/foo/zap")
				return err == nil
			},
			func(fs FileSystem) bool {
				err := fs.WriteFile("/a/b/c/hello.txt", "World.")
				if err != nil {
					fmt.Println("got err", err, "expected nil")
					return false
				}
				return true
			},
			func(fs FileSystem) bool {
				err := fs.WriteFile("/foo/bar/baz/info.txt", "Satoshi was here.")
				return err == nil
			},
			func(fs FileSystem) bool {
				data, err := fs.ReadFile("/foo/bar/baz/info.txt")
				if err != nil {
					return false
				}
				return data == "Satoshi was here."
			},
			func(fs FileSystem) bool {
				err := fs.WriteFile("/foo/bar/baz/info.txt", " 💖 💖 💖 ")
				return err == nil
			},
			func(fs FileSystem) bool {
				data, err := fs.ReadFile("/foo/bar/baz/info.txt")
				if err != nil {
					return false
				}
				return data == `Satoshi was here. 💖 💖 💖 `
			},
			func(fs FileSystem) bool {
				content, err := fs.PrettyPrint()
				if err != nil {
					panic(err)
				}

				return content == expectPrettyPrint
			},
		}
	}

	start := time.Now()

	fs := NewFileSystem()

	for i, tt := range tests() {
		testname := fmt.Sprintf("%d", i)

		t.Run(testname, func(t *testing.T) {
			ok := tt(fs)
			if !ok {
				t.Errorf("🛑 test %d failed", i)
			}
		})
	}

	fmt.Println(time.Since(start))

	fmt.Println(fs.PrettyPrint())

	// @todo convert to test
	// simulate concurrent read/write requests
	// nRequest := 1000
	// var wg sync.WaitGroup
	// wg.Add(nRequest * 2)

	// go func(fs FileSystem) {
	// 	for i := 0; i < nRequest; i++ {
	// 		line := fmt.Sprintf("\n+%d", i+1)
	// 		go func(line string) {
	// 			defer wg.Done()

	// 			delay := rand.Intn(100-1) + 1
	// 			time.Sleep(time.Millisecond * time.Duration(delay))

	// 			if err := fs.WriteFile(infoFilePath, line); err != nil {
	// 				panic(err)
	// 			}
	// 		}(line)
	// 	}
	// }(fs)

	// go func(fs FileSystem) {
	// 	for i := 0; i < nRequest; i++ {
	// 		go func() {
	// 			defer wg.Done()

	// 			delay := rand.Intn(100-1) + 1
	// 			time.Sleep(time.Millisecond * time.Duration(delay))

	// 			data, err := fs.ReadFile(infoFilePath)
	// 			if err != nil {
	// 				panic(err)
	// 			}

	// 			lines := strings.Split(data, "\n")

	// 			fmt.Println("------\n", len(lines), "lines read\n", lines, "\n------\n")
	// 		}()
	// 	}
	// }(fs)

	// wg.Wait()

	// // finale read
	// data, err = fs.ReadFile(infoFilePath)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(data)

	fmt.Println("👋")
}
