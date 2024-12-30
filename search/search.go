package search

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"encoding/json"
)

//TODO: potentially add another flag to no use the ignore patterns / or maybe to ignore a particular pattern / or a flag to add and remove patterns from the config from the cli using the jd command. potentially also add in a way to parse the json for any list so users can organize patterns as they see fit.

//TODO: possibly add a -c or something flag for continued search that will not stop on the first found search but instead return all matching dirs found.

//TODO: maybe incorporate a markdown generator representing the files system. could be cool idea from boot.dev https://github.com/aymaneallaoui/dirscanner could be under jd -G starting dir found using the jd command
type Config struct {
	IgnorePatterns []string `json:"ignore_patterns"`
	RestrictedDirs []string `json:"restrictedDirs"`
	Loaded         bool
	Patterns       []string
}
	
var config = &Config{}

// LoadPatterns loads configuration from a JSON file specified by the CONFIG_PATH environment variable,
// merging ignore patterns and restricted directories into the Patterns slice. It ensures patterns are
// only loaded once and sets the Loaded flag to true upon success.
func (c *Config) LoadPatterns()error{
	if c.Loaded{
		return nil
	}

	configpath := os.Getenv("CONFIG_PATH")
	file, err := os.Open(configpath)

	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	
	if err := decoder.Decode(&config); err != nil {
		return err
	}

	c.Patterns = append(config.IgnorePatterns, config.RestrictedDirs...)
	c.Loaded = true

	return nil
}

// jumpDirectory takes in the name (case insensitive) of the directory you want to find and a starting directory to begin the search from. This will start a depth-first search from the starting directory to locate the directory passed in in the first argument. the return will be the full path to the directory. Pass . as the second argument if you want to start the search form your current directory.
func JumpDirectory(name string, currDir string) string {
	var err error
	if !config.Loaded {
		err = config.LoadPatterns()
	}
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	os.Chdir(currDir)
	cleanPaths, path := getDirs(currDir, name)

	if path != ""{
		return path
	}

	for _, absPath := range cleanPaths{
		path = JumpDirectory(name, absPath)
		if path != "" {
			break
		}
	}

	return path
}


func searchVisDirs(dirs []fs.DirEntry, name string)([]string, string) {
	var visDirs []string

	for _, dir := range dirs{
		if dir.Name()[0] == '.' || !dir.IsDir() || slices.Contains(config.Patterns, dir.Name()) {
			continue
		}else if checkMatch(name, dir){
			path, _ := filepath.Abs(dir.Name())
			return nil, path
		}

		if len(dirs) != 0{
			path, err := filepath.Abs(dir.Name())
			if err != nil {
				fmt.Println(err)
			}
			visDirs = append(visDirs, path)
		}
		
	}
	return visDirs, ""
}

func checkMatch(name string, dir fs.DirEntry)bool{	
	return strings.EqualFold(dir.Name(), name ) 
}

func getDirs(dir string, name string)([]string, string){
	dirs, err := os.ReadDir(dir)
	
	if err != nil {
		fmt.Println(err)
	}

	return searchVisDirs(dirs, name)
}