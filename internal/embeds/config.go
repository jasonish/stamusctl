package embeds

import (
	"embed"
	"log"
	"os"
	"runtime/debug"
	"stamus-ctl/internal/app"
	"stamus-ctl/internal/utils"
	"strings"
)

//go:embed clearndr/*
var AllConf embed.FS

// Create ClearNDR folder if it does not exist
func InitClearNDRFolder(path string) {
	clearndrConfigExist, err := utils.FolderExists(path)
	if err != nil {
		debug.PrintStack()
		panic(err)
	}
	if !clearndrConfigExist && app.Embed.IsTrue() {
		err = ExtractEmbedTo("clearndr", app.TemplatesFolder+"clearndr/embedded/")
		if err != nil {
			debug.PrintStack()
			panic(err)
		}
	}
}

func ExtractEmbedTo(embed string, outputFolder string) error {
	files := getAllFiles(embed)

	for _, file := range files {
		data, err := AllConf.ReadFile(file)
		if err != nil {
			return err
		}
		err = os.MkdirAll(outputFolder+"/"+extractPath(file), 0755)
		if err != nil {
			return err
		}
		err = os.WriteFile(outputFolder+"/"+extractPath(file)+"/"+extractFileName(file), data, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func getAllFiles(inputFolder string) []string {
	var files []string
	entries, err := AllConf.ReadDir(inputFolder)
	if err != nil {
		log.Println("Error reading dir", inputFolder)
		log.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			files = append(files, getAllFiles(inputFolder+"/"+entry.Name())...)
		} else {
			files = append(files, inputFolder+"/"+entry.Name())
		}
	}
	return files
}

func extractPath(path string) string {
	// returns everything before the last /
	array := strings.Split(path, "/")
	return strings.Join(array[1:len(array)-1], "/")
}

func extractFileName(path string) string {
	// returns everything before the last /
	array := strings.Split(path, "/")
	return strings.Join(array[len(array)-1:], "/")
}
