package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var home, _ = os.UserHomeDir()
var txtfilepath = home + "/Documents/logger.txt"
var jsonfilepath = home + "/Documents/logger.json"
var bd = home + "/Documents/bd"
var bdshare = home + "/Documents/bdshare"

var sd = home + "/Documents/sd"
var logext = "txt"
var encrypttimes = "0"

func sharer(path string, version int) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	var paths []string
	var names []string
	for _, file := range files {
		nextpath := path + "/" + file.Name()
		if file.IsDir() {
			os.Mkdir(bdshare+nextpath[len(sd):], 0755)
			sharer(nextpath, version)
		} else {
			arr := strings.Split(nextpath, ".")
			verstr := arr[len(arr)-2]
			ver, _ := strconv.Atoi(verstr[2 : len(verstr)-1])
			if ver <= version {
				paths = append(paths, nextpath)
				names = append(names, arr[len(arr)-3])
				//fmt.Println(paths, names)
				if len(names) > 1 && names[len(names)-1] != names[len(names)-2] {
					copyfile(paths[len(paths)-2], bdshare+paths[len(paths)-2][len(sd):], false)
				}
			}
		}
	}
	if len(names) > 0 {
		copyfile(paths[len(paths)-1], bdshare+paths[len(paths)-1][len(sd):], false)
	}
}

func reset() {
	os.RemoveAll(bd)
	os.Remove(jsonfilepath)
	os.Remove(txtfilepath)
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func str_to_time(final string) time.Time {
	year, err := strconv.Atoi(final[0:4])
	month, err := strconv.Atoi(final[5:7])
	day, err := strconv.Atoi(final[8:10])
	hour, err := strconv.Atoi(final[11:13])
	minute, err := strconv.Atoi(final[14:16])
	second, err := strconv.Atoi(final[17:19])

	if err != nil {
		log.Fatal(err)
	}

	last := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
	return last
}

func copyfile(path string, output string, encrypt bool) {
	input, err := os.ReadFile(path)
	times, _ := strconv.Atoi(encrypttimes)
	if err != nil {
		panic(err)
	}
	if encrypt {
		if times < 0 || times > 5 {
			fmt.Println("Choose smaller encryption value")
			os.Exit(0)
		}

		for i := 0; i < times; i++ {
			input = []byte(base64.StdEncoding.EncodeToString(input))
		}
	}

	err = os.WriteFile(output, input, 0644)
	if err != nil {
		panic(err)
	}
}

func backup(ver int, last time.Time, path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}
		if file.IsDir() {
			os.Mkdir(bd+(path + "/" + file.Name())[len(path):], 0755)
			backup(ver, last, path+"/"+file.Name())
		} else if info.ModTime().After(last) {
			name := strings.Split(file.Name(), ".")
			name[len(name)-2] = name[len(name)-2] + ".(v" + strconv.Itoa(ver) + ")"

			copyfile(path+"/"+file.Name(), bd+(path + "/" + strings.Join(name, "."))[len(sd):], true)
		}

	}
}

func forjson() {
	//var last time.Time
	//var ver int
	var last time.Time
	var ver int
	if _, err := os.Stat(jsonfilepath); err == nil {
		input, err := os.ReadFile(jsonfilepath)
		var target map[string]any
		err = json.Unmarshal([]byte(input), &target)
		if err != nil {
			panic(err)
		}
		ver, err = strconv.Atoi(target["version"].(string))
		ver = ver + 1
		lastarr := target["time"].([]interface{})

		laststr := lastarr[len(lastarr)-1].(string)
		last = str_to_time(laststr)

		lastarr = append(lastarr, time.Now().Format(time.RFC3339))

		data := `{"version":"` + strconv.Itoa(ver) + `","time":[`
		for i := 0; i < len(lastarr)-1; i++ {
			data += `"` + lastarr[i].(string) + `",`
		}
		data += `"` + lastarr[len(lastarr)-1].(string) + `"]}`

		t := os.WriteFile(jsonfilepath, []byte(data), 0644)
		if t != nil {
			log.Fatal(t)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		ver = 0
		last = time.Time{}
		// var times []time.Time
		// times = append(times, time.Now())
		//times = append(times, time.Now())
		//fmt.Println(times)
		data := `{"version":"` + strconv.Itoa(ver) + `","time":["` + time.Now().Format(time.RFC3339) + `"]}`
		t := os.WriteFile(jsonfilepath, []byte(data), 0644)
		if t != nil {
			log.Fatal(t)
		}
	} else {
		log.Fatal(err)
	}
	backup(ver, last, sd)
}

func fortxt() {
	var last time.Time
	var ver int
	if _, err := os.Stat(txtfilepath); err == nil {
		content, err := os.ReadFile(txtfilepath)
		if err != nil {
			log.Fatal(err)
		}
		stringdad := strings.Split(string(content), " ")
		final := stringdad[len(stringdad)-2]
		verlast := stringdad[len(stringdad)-3][2:]
		ver, err = strconv.Atoi(verlast)

		last = str_to_time(final)
		ver = ver + 1

		f, err := os.OpenFile(txtfilepath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		if _, err = f.WriteString("v" + strconv.Itoa(ver) + " " + time.Now().Format(time.RFC3339) + " \n"); err != nil {
			panic(err)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		last = time.Time{}
		ver = 0
		f, err := os.Create(txtfilepath)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		data := []byte("v0 " + time.Now().Format(time.RFC3339) + " \n")

		_, err = f.Write(data)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
	backup(ver, last, sd)
}

func readconffile() {
	input, err := os.ReadFile("conf.json")

	var target map[string]any

	err = json.Unmarshal([]byte(input), &target)

	if err != nil {
		panic(err)
	}
	sd = target["source"].(string)
	logext = target["option"].(string)
	encrypttimes = target["encrypt"].(string)
}
func makeconffile() {
	data := `{"option":"` + logext + `","source":"` + sd + `","encrypt":"` + encrypttimes + `"}`

	os.WriteFile("conf.json", []byte(data), 0644)
}

func main() {
	var err error
	if _, err := os.Stat("conf.json"); errors.Is(err, os.ErrNotExist) {
		makeconffile()
	}
	readconffile()
	e, _ := strconv.Atoi(encrypttimes)
	enew := flag.Int("e", e, "Number of times to encrypt")
	sdnew := flag.String("sd", sd, "Source Directory")
	logextnew := flag.String("ext", logext, "Extension of logger file")
	share := flag.Int("share", -1, "Share this backup version")
	flag.Parse()
	encrypttimes = strconv.Itoa(*enew)
	sd = *sdnew
	logext = *logextnew
	makeconffile()

	os.Mkdir(sd, 0755)
	os.Mkdir(bd, 0755)

	argsWithoutProg := os.Args[1:]
	//fmt.Println(argsWithoutProg)
	if len(argsWithoutProg) > 0 && argsWithoutProg[0] == "reset" {
		c := askForConfirmation("This will delete al backup. Are you sure?")
		if c {
			reset()
			fmt.Println("Deleted backup and it's logs")
		}
		return
	}

	if *share != -1 {
		os.RemoveAll(bdshare)
		os.Mkdir(bdshare, 0755)
		sharer(bd, *share)
		return
	}

	if logext == "txt" {
		fortxt()
	} else if logext == "json" {
		forjson()
	} else {
		log.Fatal(err)
	}
}
