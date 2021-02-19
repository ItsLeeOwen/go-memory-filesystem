package filesystem

import (
	"fmt"
	"testing"
	"time"
)

func TestFileSystem(t *testing.T) {
	type test = func(fs FileSystemer) bool

	var tests = func() []test {
		return []test{
			func(fs FileSystemer) bool {
				err := fs.Mkdir("/foo/bar/baz")
				return err != nil
			},
			func(fs FileSystemer) bool {
				err := fs.Mkdir("/foo")
				return err == nil
			},
			func(fs FileSystemer) bool {
				err := fs.Mkdir("/foo/bar/baz")
				return err != nil
			},
			func(fs FileSystemer) bool {
				err := fs.Mkdir("/foo/bar")
				return err == nil
			},
			func(fs FileSystemer) bool {
				err := fs.Mkdir("/foo/bar/baz")
				return err == nil
			},
			func(fs FileSystemer) bool {
				err := fs.WriteFile("/doesnotexist/info.txt", "Hello.")
				return err != nil
			},
			func(fs FileSystemer) bool {
				err := fs.WriteFile("/foo/bar/baz/info.txt", "Satoshi was here.")
				return err == nil
			},
			func(fs FileSystemer) bool {
				data, err := fs.ReadFile("/foo/bar/baz/info.txt")
				if err != nil {
					return false
				}
				return data == "Satoshi was here."
			},
			func(fs FileSystemer) bool {
				err := fs.WriteFile("/foo/bar/baz/info.txt", " 💖 💖 💖 ")
				return err == nil
			},
			func(fs FileSystemer) bool {
				data, err := fs.ReadFile("/foo/bar/baz/info.txt")
				if err != nil {
					return false
				}
				return data == `Satoshi was here. 💖 💖 💖 `
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

	// @todo convert to test
	// simulate concurrent read/write requests
	// nRequest := 1000
	// var wg sync.WaitGroup
	// wg.Add(nRequest * 2)

	// go func(fs FileSystemer) {
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

	// go func(fs FileSystemer) {
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
