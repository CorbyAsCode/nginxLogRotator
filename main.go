package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

var print = fmt.Println

func check(e error) {
	if e != nil {
		print(e)
		os.Exit(1)
	}
}

func runNginx(sig chan os.Signal, done chan bool) {

	print("Starting nginxLogRotator...")
	nginx := exec.Command("nginx", "-g", "daemon off;")
	_, err := nginx.Output()
	check(err)
	signal := <-sig
	print()
	print("Stopping nginx...")

	print(signal)
	done <- true
}

func filePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func findConfigs(root string) []string {
	var files []string
	var configs []string
	files, err := filePathWalkDir(root)
	check(err)
	for _, config := range files {
		if strings.HasSuffix(config, ".conf") {
			configs = append(configs, config)
		}
	}
	return configs
}

func findLogs() {
	configs := findConfigs("/etc/nginx")

	r, err := regexp.Compile(`^\s*(error|access)_log\s+(?P<logFile>.*.log).*$`)
	check(err)
	//fmt.Printf("%#v\n", r.FindStringSubmatch("access_log /var/log/nginx/access.log;"))
	//fmt.Printf("%#v\n", r.SubexpNames())

	for _, config := range configs {
		contents, err := ioutil.ReadFile(config)
		check(err)
		print(string(contents))
		fmt.Printf("%#v\n", r.FindStringSubmatch(string(contents)))
	}
}

func main() {
	/*
		r := regexp.MustCompile(`(?P<Year>\d{4})-(?P<Month>\d{2})-(?P<Day>\d{2})`)
		fmt.Printf("%#v\n", r.FindStringSubmatch(`2015-05-27`))
		fmt.Printf("%#v\n", r.SubexpNames())
	*/
	findLogs()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)
	go runNginx(sigs, done)
	<-done

}
