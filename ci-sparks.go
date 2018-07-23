package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
)

var cinder_path string = ""
var dest_path string = "/Users/henrique/Desktop/"
var project_name string = "SparkTest"

func getDefaultCinderPath() string {
	//todo check for enviroment variable
	ciPath, _ := filepath.Abs("../../../Cinder/")

	if _, err := os.Stat(ciPath); os.IsNotExist(err) {
		fmt.Printf(".Error: \"%s does not exist\"\n", ciPath)
		return ""
	}
	return ciPath
}

func createFiles(path string) bool {

	base_project_path := cinder_path + "/blocks/__AppTemplates/BasicApp/Opengl/"
	copy.Copy(base_project_path, dest_path+"/"+project_name)
	return true
}

func buildCMakeProject() {

	cmakePath := cinder_path + "/tools/scripts/files/CMakeLists.txt"

	// write the whole body at once
	if _, err := os.Stat(dest_path + project_name); !os.IsNotExist(err) {

		createError := os.Mkdir(dest_path+project_name+"/proj", 0777)
		if createError != nil {
			panic(createError)
		}

	} else {
		panic(err)
	}

	input, err := ioutil.ReadFile(cmakePath)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(input), "\n")
	for i, _ := range lines {

		lines[i] = strings.Replace(lines[i], "__CINDER_PATH__", cinder_path, -1)
		lines[i] = strings.Replace(lines[i], "__PROJECT_NAME__", project_name, -1)
	}

	output := strings.Join(lines, "\n")

	{
		err = os.Mkdir(dest_path+project_name+"/proj/cmake/", 0777)
		if err != nil {
			panic(err)
		}
	}

	f, err := os.Create(dest_path + "/" + project_name + "/CMakeLists.txt")
	defer f.Close()

	if err != nil {
		panic(err)
	}

	f.WriteString(output)
}

func cleanUpFiles() {

	var basedir string = dest_path + project_name
	// os.Rename( basedir + "/src/_TBOX_PREFIX_App.cpp", basedir + "/src/" + project_name + "App.cpp")

	os.Remove(basedir + "/assets/_TBOX_IGNORE_")
	os.Remove(basedir + "/template.xml")

	os.RemoveAll(basedir + "/vc2015_uwp")
	os.RemoveAll(basedir + "/xcode")
	os.RemoveAll(basedir + "/xcode_ios")

	// clean up basic c++ file
	input, err := ioutil.ReadFile(basedir + "/src/_TBOX_PREFIX_App.cpp")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(input), "\n")
	for i, _ := range lines {
		lines[i] = strings.Replace(lines[i], "_TBOX_PREFIX_App", project_name+"App", -1)
	}
	output := strings.Join(lines, "\n")

	f, err := os.Create(basedir + "/src/" + project_name + "App.cpp")
	defer f.Close()

	if err != nil {
		panic(err)
	}

	f.WriteString(output)
	os.Remove(basedir + "/src/_TBOX_PREFIX_App.cpp")
}

type testStruct struct {
	CinderPath string `json:"CinderPath"`
}

func main() {
	// Open our jsonFile
	jsonFile, err := os.Open("config.json")
	defer jsonFile.Close()
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var t testStruct
	// var jsonData = []byte(jsonFile)
	err = json.Unmarshal(byteValue, &t)
	if err != nil {
		fmt.Printf("There was an error decoding the json. err = %s", err)
		return
	}
	fmt.Println(t.CinderPath)

	cinder_path = t.CinderPath

	// ---
	destStringPtr := flag.String("dest", ".", "project destination path")
	projectNamePtr := flag.String("name", "Basic", "Name of the project, default is Basic")

	flag.Parse()

	// set destination path
	dest_path = *destStringPtr
	lastChar := string(dest_path[len(dest_path)-1])
	if lastChar != "/" {
		dest_path += "/"
	}
	fmt.Println(". destination: " + dest_path)

	// set project name
	project_name = *projectNamePtr
	fmt.Println(". project name: " + project_name)

	// check if project already exisits
	if _, err := os.Stat(dest_path + project_name); !os.IsNotExist(err) {
		fmt.Println("🔥 Folder already exists, aborting!")
		return
	}

	//
	fmt.Println("cinder path: " + cinder_path)

	createFiles(dest_path)
	buildCMakeProject()
	cleanUpFiles()

	fmt.Println(" 🌋 done! ✨")
}
