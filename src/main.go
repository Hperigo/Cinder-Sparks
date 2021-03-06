package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
)

var cinder_path string = ""
var dest_path string = ""
var app_name string = ""

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func createFiles(path string) bool {

	fmt.Println("copying files to...:  " + dest_path + "/" + app_name)

	base_project_path := "./templates/BasicApp/"
	copy.Copy(base_project_path, dest_path+app_name)
	return true
}

func buildCMakeProject() {

	cmakePath := "./templates/CMakeLists.txt"

	// write the whole body at once
	if _, err := os.Stat(dest_path + app_name); !os.IsNotExist(err) {

		createError := os.Mkdir(dest_path+app_name+"/proj", 0777)
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
		lines[i] = strings.Replace(lines[i], "__PROJECT_NAME__", app_name, -1)
		lines[i] = strings.Replace(lines[i], "__app_name__", app_name, -1)
	}

	output := strings.Join(lines, "\n")

	{
		err = os.Mkdir(dest_path+app_name+"/proj/cmake/", 0777)
		if err != nil {
			panic(err)
		}
	}

	f, err := os.Create(dest_path + "/" + app_name + "/CMakeLists.txt")
	defer f.Close()

	if err != nil {
		panic(err)
	}

	f.WriteString(output)
}

func cleanUpFiles() {

	var basedir string = dest_path + app_name
	// os.Rename( basedir + "/src/_TBOX_PREFIX_App.cpp", basedir + "/src/" + app_name + "App.cpp")

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
		lines[i] = strings.Replace(lines[i], "_TBOX_PREFIX_App", app_name+"App", -1)
	}
	output := strings.Join(lines, "\n")

	f, err := os.Create(basedir + "/src/" + app_name + "App.cpp")
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

func parseJson() {

	jsonFile, err := os.Open("./templates/config.json")
	defer jsonFile.Close()
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	// {
	// 	fileUrl := "https://libcinder.org/static/releases/cinder_0.9.1_mac.zip"
	// 	err := DownloadFile("cinder.zip", fileUrl)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var t testStruct
	// var jsonData = []byte(jsonFile)
	err = json.Unmarshal(byteValue, &t)
	if err != nil {
		fmt.Printf("There was an error decoding the json. err = %s", err)
		return
	}

	cinder_path, _ = filepath.Abs(t.CinderPath)
	fmt.Println(". CinderPath: " + cinder_path)

}

func main() {

	parseJson()

	// ---
	destStringPtr := flag.String("dest", ".", "project destination path")
	appNamePtr := flag.String("name", "Basic", "Name of the project, default is Basic")
	isStandAlonePtr := flag.Bool("standalone", false, "creates a santadalone folder with cinder and a app folder")
	flag.Parse()

	fmt.Println(*isStandAlonePtr)

	// set destination path
	dest_path = *destStringPtr
	lastChar := string(dest_path[len(dest_path)-1])
	if lastChar != "/" {
		dest_path += "/"
	}

	fmt.Println(". Destination: " + dest_path)

	// set project name
	app_name = *appNamePtr
	fmt.Println(". Project name: " + app_name)

	if *isStandAlonePtr {
		dest_path = dest_path + app_name + "Poject/"

		fmt.Println("standalone path is: " + dest_path)

		createError := os.Mkdir(dest_path, 0777)
		if createError != nil {
			panic(createError)
		}
		cinder_path = dest_path + "lib/"
		createError = os.Mkdir(cinder_path, 0777)
		if createError != nil {
			panic(createError)
		}

		// fmt.Println(runtime.GOOS)
		// if runtime.GOOS == "windows" {
		// 	fmt.Println("Hello from Windows")
		// } else if runtime.GOOS == "darwin" {
		// 	fmt.Println("Hello from dawrwin")
		// }

		// git clone git@github.com:whatever folder-na
		// DownloadFile( cinder_path,
	}

	// check if project already exisits
	if _, err := os.Stat(dest_path + app_name); !os.IsNotExist(err) {
		fmt.Println("🔥 Folder already exists, aborting!")
		return
	}

	createError := os.Mkdir(dest_path+"/"+app_name, 0777)
	if createError != nil {
		panic(createError)
	}

	createFiles(dest_path)
	buildCMakeProject()
	cleanUpFiles()

	fmt.Println("Created at: " + dest_path)
	fmt.Println(" 🌋 done! ✨")
}
