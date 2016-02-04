package wiki

import (
    "regexp"
    "strings"
    //"fmt"
)

var rRef, _ = regexp.Compile("(?smi:<ref.+?</ref>)")
var rCite, _ = regexp.Compile("(?smi:{+.+?}+)")
var rLink, _ = regexp.Compile("(?smi:\\[\\[(?:.+?\\|)?(.+?)\\]\\])")
var rMarkup, _ = regexp.Compile("(?smi:=+|'+|<br\\s*/?>|<div.+?>|</div>|<math>.+?</math>|<.+?>|\\[http.+?\\]|[[File.+?]]|<!--+.+?-->|#REDIRECT)")
var rSpaces, _ = regexp.Compile("\\s+")
var rNotAlphaNumeric, _ = regexp.Compile("[^a-zA-Z0-9]")

func NormalizeText(text string) string {
	text = RemoveCite(text)
	text = RemoveRef(text)
	
	text = RemoveLink(text)
	text = RemoveMarkup(text)
	text = RemoveDuplication(text)
	return text
}

func RemoveCite(text string) string {
	text = rCite.ReplaceAllString(text, " ")
	return TrimSpaces(text)
}

func RemoveRef(text string) string {
	text = rRef.ReplaceAllString(text, " ")
	return TrimSpaces(text)
}

func RemoveLink(text string) string {
	if rLink.MatchString(text) == true {
        result:= rLink.FindAllStringSubmatch(text,-1)
        for _, v := range result {
    		text = strings.Replace(text, v[0], v[1], 5000)
		}	
        return TrimSpaces(text)
    } else {
        return text
    }
}

func RemoveMarkup(text string) string {
	text = rMarkup.ReplaceAllString(text, " ")
	return TrimSpaces(text)
}

func TrimSpaces(text string) string{
	text = rSpaces.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)
	return text
}

func RemoveDuplication(text string) string{
	text = strings.ToLower(text)
	text = rNotAlphaNumeric.ReplaceAllString(text, " ")
	terms := strings.Split(text, " ")
	mapTerms := make(map[string]int)
	for _, value := range terms {
		value = strings.TrimSpace(value)
		if( "" != value){
			mapTerms[value] = 1	
		}
	}
	newText := ""
	for key, _ := range mapTerms {
    	newText = newText + " " + key
	}
	newText = strings.TrimSpace(newText)
	//fmt.Println("SAIDA " + newText)
	return newText
}
