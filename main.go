package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func getDBDir() string {
	home, _ := os.UserHomeDir()

	return home + "/.godir"
}

func getDBFilePath() string {
	return getDBDir() + "/godir.db"
}

func ensureDB() {
	if _, err := os.Stat(getDBDir()); os.IsNotExist(err) {
		// fmt.Println("Creating DB Dir: " + getDBDir())

		if err := os.Mkdir(getDBDir(), os.ModePerm); err != nil {
			log.Fatalln("Error creating db dir: ", err)
		}

		// fmt.Println("Creating DB File: " + getDBFilePath())

		f, e := os.Create(getDBFilePath())
		if e != nil {
			log.Fatalln("Error creating db file: ", e)
		}

		defer f.Close()
	}
}

// func dirExists(dir string) bool {
// 	if _, err := os.Stat(dir); os.IsNotExist(err) {
// 		return true
// 	}

// 	return false
// }

func readDB() []string {
	ensureDB()

	dbBytes, err := os.ReadFile(getDBFilePath()) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	lines := strings.Split(string(dbBytes), "\n")

	return lines
}

func writeDB(db []string) {

	// fmt.Println("getDBFilePath(): ", getDBFilePath())

	flushDB()

	file, err := os.OpenFile(getDBFilePath(), os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("Error opening db file: ", err)
	}

	defer file.Close()

	dbStr := strings.Join(db, "\n")

	// fmt.Println("dbStr: ", dbStr)

	_, err = file.WriteString(dbStr)

	if err != nil {
		log.Fatalln("Error writing db file: ", err)
	}
}

func addDir(dir string) {
	db := readDB()

	// fmt.Println("db[0]: ", string(db[0]))

	// if !dirExists(dir) {
	// 	log.Fatalln("Dir Not Exists")
	// }

	if db[len(db)-1] != dir {
		db = append(db, dir)
	}

	writeDB(db)

	fmt.Print(dir)
}

func flushDB() {
	f, e := os.Create(getDBFilePath())
	if e != nil {
		log.Fatalln("Error flushing db file: ", e)
	}

	defer f.Close()
}

func popLastDir() string {
	db := readDB()

	lastDir := pop(&db)

	writeDB(db)

	// fmt.Println("DB: ", db)

	return lastDir
}

func pop(alist *[]string) string {
	f := len(*alist)
	rv := (*alist)[f-1]
	*alist = (*alist)[:f-1]
	return rv
}

// https://stackoverflow.com/questions/52435908/how-to-change-the-shells-current-working-directory-in-go
// ? https://github.com/JohnStarich/goenable?tab=readme-ov-file#write-a-plugin

func main() {

	args := os.Args[1:]

	if len(args) <= 0 {
		panic("Missing directory to chdir to")
	}

	if args[0] == "--prev" {
		currentDir := popLastDir()

		file, err := os.Create(getDBDir() + "/.toggle")
		if err != nil {
			log.Fatalln("Error creating toggle dir file: ", err)
		}

		defer file.Close()

		file.WriteString(currentDir)

		fmt.Print(popLastDir())
	} else if args[0] == "--flush" {
		fmt.Println("FLUSHING DB")

		flushDB()
	} else if args[0] == "-" {
		bytes, err := os.ReadFile(getDBDir() + "/.toggle")
		if err != nil {
			log.Fatalln("Error reading toggle dir file: ", err)
		}

		toggleDir := string(bytes)

		currentDir := popLastDir()

		file, err := os.Create(getDBDir() + "/.toggle")
		if err != nil {
			log.Fatalln("Error creating toggle dir file: ", err)
		}

		defer file.Close()

		file.WriteString(currentDir)

		addDir(toggleDir)

	} else {
		addDir(args[0])
	}
}
