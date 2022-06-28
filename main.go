package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func main() {
	replace := "Something"
	replaceLow := strings.ToLower(replace[0:1]) + replace[1:]
	replaceUp := replace
	replaceTitle := matchFirstCap.ReplaceAllString(replace, "${1}_${2}")
	replaceTitle = matchAllCap.ReplaceAllString(replaceTitle, "${1}_${2}")
	replaceTitle = strings.ToLower(replaceTitle)

	repoFileName := "template_repo.gotmp"
	serviceFileName := "template_service.gotmp"
	err := generateFile(replaceLow, replaceUp, replaceTitle, repoFileName, serviceFileName)
	if err != nil {
		fmt.Println(err)
	}
}

func generateFile(replaceLow, replaceUp, title, repo, service string) error {
	fileRepo, err := os.Open(repo)
	if err != nil {
		fmt.Println(err)
	}
	defer fileRepo.Close()
	_ = os.Remove("result")
	os.Mkdir("result", 0777)

	resFileRepo, err := os.Create("result/" + title + "_repository.")
	if err != nil {
		fmt.Println(errors.New("failed creating file, " + err.Error()))
	}
	defer resFileRepo.Close()

	var lineText string
	scanner := bufio.NewScanner(fileRepo)
	for scanner.Scan() {
		lineText = scanner.Text()
		lineText = strings.Replace(lineText, "xxx", replaceUp, -1)
		_, err = resFileRepo.WriteString(lineText + "\n")
		if err != nil {
			fmt.Println(errors.New("failed writing file, " + err.Error()))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	fileService, err := os.Open(service)
	if err != nil {
		fmt.Println(err)
	}
	defer fileService.Close()
	resFileService, err := os.Create("result/" + title + "_services.")
	if err != nil {
		fmt.Println(errors.New("failed creating file, " + err.Error()))
	}
	defer resFileService.Close()
	scannerService := bufio.NewScanner(fileService)
	for scannerService.Scan() {
		lineText = scannerService.Text()
		lineText = strings.Replace(lineText, "xxx", replaceLow, -1)
		lineText = strings.Replace(lineText, "yyy", replaceUp, -1)
		_, err = resFileService.WriteString(lineText + "\n")
		if err != nil {
			fmt.Println(errors.New("failed writing file, " + err.Error()))
		}
	}

	if err := scannerService.Err(); err != nil {
		fmt.Println(err)
	}
	return nil
}
