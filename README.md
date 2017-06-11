# Walker

Walker lets you parallely walks through the directories in you file system.
For each directory at the root level, walker starts a go routine. Each of the
go routines will recursively get all the files from the directories and
sub directories.

The walker also accepts a format/filter function as the second argument.
The format function can be used to filter and/or format the files 
returned by the walk.

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
