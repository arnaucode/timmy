package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

//dataEntry is the map used to create the array of maps, where the templatejson data is stored
type dataEntry map[string]string

func readFile(path string) string {
	dat, err := ioutil.ReadFile(path)
	check(err)
	return string(dat)
}

func getDataFromJson(path string) []dataEntry {
	var entries []dataEntry
	file, err := ioutil.ReadFile(path)
	check(err)
	content := string(file)
	var rawEntries []*json.RawMessage
	json.Unmarshal([]byte(content), &rawEntries)
	for i := 0; i < len(rawEntries); i++ {
		rawEntryMarshaled, err := json.Marshal(rawEntries[i])
		check(err)
		var newDataEntry map[string]string
		json.Unmarshal(rawEntryMarshaled, &newDataEntry)
		entries = append(entries, newDataEntry)
	}

	return entries
}

func generateFromTemplateAndData(templateContent string, entries []dataEntry) string {

	var newContent string

	for i := 0; i < len(entries); i++ {
		var entryContent string
		entryContent = templateContent
		//first, get the map keys
		var keys []string
		for key, _ := range entries[i] {
			keys = append(keys, key)
		}
		//now, replace the keys with the values
		for j := 0; j < len(keys); j++ {
			entryContent = strings.Replace(entryContent, "{{"+keys[j]+"}}", entries[i][keys[j]], -1)
			entryContent = strings.Replace(entryContent, "[[i]]", strconv.Itoa(i), -1)
		}

		newContent = newContent + entryContent
	}
	return newContent
}
func getTemplateParameters(line string) (string, string) {
	var templatePath string
	var data string
	line = strings.Replace(line, "<konstrui-template ", "", -1)
	line = strings.Replace(line, "></konstrui-template>", "", -1)
	attributes := strings.Split(line, " ")
	//fmt.Println(attributes)
	for i := 0; i < len(attributes); i++ {
		attSplitted := strings.Split(attributes[i], "=")
		if attSplitted[0] == "html" {
			templatePath = strings.Replace(attSplitted[1], `"`, "", -1)
		}
		if attSplitted[0] == "data" {
			data = strings.Replace(attSplitted[1], `"`, "", -1)
		}
	}
	return templatePath, data
}

func useTemplate(templatePath string, dataPath string) string {
	filepath := rawFolderPath + "/" + templatePath
	templateContent := readFile(filepath)
	entries := getDataFromJson(rawFolderPath + "/" + dataPath)
	generated := generateFromTemplateAndData(templateContent, entries)
	return generated
}

func putTemplates(folderPath string, filename string) string {
	var fileContent string
	f, err := os.Open(folderPath + "/" + filename)
	check(err)
	scanner := bufio.NewScanner(f)
	lineCount := 1
	for scanner.Scan() {
		currentLine := scanner.Text()
		if strings.Contains(currentLine, "<konstrui-template") && strings.Contains(currentLine, "</konstrui-template>") {
			templatePath, data := getTemplateParameters(currentLine)
			fileContent = fileContent + useTemplate(templatePath, data)
		} else {
			fileContent = fileContent + currentLine
		}
		lineCount++
	}
	return fileContent
}

func writeFile(path string, newContent string) {
	err := ioutil.WriteFile(path, []byte(newContent), 0644)
	check(err)
}

func generatePageFromTemplateAndData(templateContent string, entry dataEntry) string {
	var entryContent string
	entryContent = templateContent
	//first, get the map keys
	var keys []string
	for key, _ := range entry {
		keys = append(keys, key)
	}
	//now, replace the keys with the values
	for j := 0; j < len(keys); j++ {
		entryContent = strings.Replace(entryContent, "{{"+keys[j]+"}}", entry[keys[j]], -1)
	}
	return entryContent
}
func getHtmlAndDataFromRepeatPages(page RepeatPages) (string, []dataEntry) {
	templateContent := readFile(rawFolderPath + "/" + page.HtmlPage)
	data := getDataFromJson(rawFolderPath + "/" + page.Data)
	return templateContent, data
}
