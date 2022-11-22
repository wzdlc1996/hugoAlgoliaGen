package main

import (
	"bufio"
	"errors"
	"path"
	"regexp"
	"strings"
)

// PostInfo is the information data of a single post
type PostInfo struct {
	Title      string   `json:"title"`
	Date       string   `json:"date"`
	Draft      bool     `json:"draft"`
	Tags       []string `json:"tags"`
	Categories []string `json:"categories"`
	Toc        bool     `json:"toc"`
	Summary    string   `json:"summary"`
	Content    string   `json:"content"`
	ObjectID   string   `json:"objectID"`
	Uri        string   `json:"uri"`
}

// ConfigParser returns the config object of toml/yaml/json formats.
func ConfigParser(config string) (res map[string]string, err error) {
	ext := path.Ext(config)
	if ext == "toml" {
		return tomlParser(config), nil
	}
	if ext == "yaml" {
		return yamlParser(config), nil
	}
	if ext == "json" {
		return jsonParser(config), nil
	}
	err = errors.New("not supported config file format, please check")
	return
}

var contentNeglectPattern string = `\$.*?\$|<.+?>.*?</.+?>|<.+?>`

// Some stupid parser for toml, yaml, and json. Only top level and output map[string]string

// PostParser parses the markdown file (as string), read its yaml information.
func PostParser(mdcont string) (res PostInfo) {
	mdReader := strings.NewReader(mdcont)
	scanner := bufio.NewScanner(mdReader)

	inYaml := false
	idx := -1
	for scanner.Scan() {
		idx += 1
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			inYaml = !inYaml
			continue
		}
		if inYaml {
			lineSpl := strings.SplitN(line, ":", 2)
			key, val := lineSpl[0], lineSpl[1]
			key = strings.TrimSpace(key)
			val = strings.Trim(val, "\" ")
			switch key {
			case "title":
				res.Title = val
			case "date":
				res.Date = val
			case "tags":
				res.Tags = parseList(val)
			case "categories":
				res.Categories = parseList(val)
			case "toc":
				res.Toc = val == "true"
			case "draft":
				res.Draft = val == "true"
			case "summary":
				res.Summary = val
			}

		} else {
			res.Content += line + " "
		}
	}
	res.Content = strings.Join(extractWordsFromContent(res.Content), " ")

	return
}

func tomlParser(config string) map[string]string {
	if config == "" {
		return map[string]string{}
	}
	res := map[string]string{}
	inSub := false
	for _, line := range strings.Split(config, "\n") {
		if inSub || !strings.Contains(line, "=") {
			if line != "" {
				inSub = true
			}
			continue
		}
		splitedLine := strings.SplitN(line, "=", 2)
		key, val := splitedLine[0], splitedLine[1]
		res[strings.TrimSpace(key)] = strings.Trim(val, " \"")
	}
	return res
}

func yamlParser(config string) map[string]string {
	if config == "" {
		return map[string]string{}
	}
	res := map[string]string{}
	for _, line := range strings.Split(config, "\n") {
		if !strings.Contains(line, ":") {
			continue
		}
		splitedLine := strings.SplitN(line, ":", 2)
		key, val := splitedLine[0], splitedLine[1]
		res[strings.TrimSpace(key)] = strings.Trim(val, " \"")
	}
	return res

}

func jsonParser(config string) (v map[string]string) {
	return
}

func contains(list []string, tar string) bool {
	for _, a := range list {
		if a == tar {
			return true
		}
	}
	return false
}

func parseList(input string) []string {
	// Assuming the format is "["vvv", "www", "jjj"]"
	res := []string{}
	for _, item := range strings.Split(strings.Trim(input, "[]"), ",") {
		res = append(res, strings.Trim(item, " \""))
	}
	return res
}

func extractWordsFromContent(content string) []string {
	negPatt := regexp.MustCompile(contentNeglectPattern)
	contentWithoutNeglectPatt := negPatt.ReplaceAllLiteralString(content, " ")
	wdPatt := regexp.MustCompile("[A-Za-z ]*")
	contentWordsOnly := wdPatt.FindAllString(contentWithoutNeglectPatt, -1)
	contentWordsOnly = strings.Fields(strings.Join(contentWordsOnly, " "))
	contentWords := []string{}
	for _, item := range contentWordsOnly {
		item = strings.ToLower(item)
		if !contains(Stopwords, item) {
			contentWords = append(contentWords, item)
		}
	}
	return deleteDuplicate(contentWords)
}

func deleteDuplicate(list []string) []string {
	isin := map[string]bool{}
	outlist := []string{}
	for _, item := range list {
		if _, value := isin[item]; !value {
			isin[item] = true
			outlist = append(outlist, item)
		}
	}
	return outlist
}
