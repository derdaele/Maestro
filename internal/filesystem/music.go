package filesystem

import (
	"os"
	"path"
	"path/filepath"
)

var musicFileExtensions = map[string]bool{
	".wma": true,
	".mp3": true,
}

func isMusicFile(fileInfo os.FileInfo) bool {
	enabled, ok := musicFileExtensions[path.Ext(fileInfo.Name())]
	return ok && enabled
}

// ListMusicFiles returns a channel enumerating files in a directory
func ListMusicFiles(directory string) chan string {
	ch := make(chan string)

	go func() {
		filepath.Walk(directory, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && isMusicFile(info) {
				ch <- filePath
			}

			return nil
		})
		close(ch)
	}()

	return ch
}
