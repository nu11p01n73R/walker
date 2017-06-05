# Walker

Walker lets you parallelly walks through the directories in you file system.

# Usage

```
import "github.com/nu11p01n73R/walker"

files, err := walker.Walk("/root/dir", func (files []string) []string) {
        // Filter the files accordingly
        return files
});

if err != nil {
        // Handle the error case
}

// Do something with the files
fmt.Println(len(files))
```
