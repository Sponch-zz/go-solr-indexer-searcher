package main

	
import (
  //  "bufio"
    "fmt"
    "flag"
  //  "io"
  //  "io/ioutil"
    "strings"
    "os"
    "encoding/xml"
    "encoding/json"
	"github.com/wkp/solr"
    "github.com/wkp/wiki"
)

var inputFile = flag.String("infile", "enwiki-latest-pages-articles.xml", "Input file path")
var solrHost = flag.String("solrHost", "localhost", "Solr Host")
var solrPort = flag.String("solrPort", "8984", "Solr Port")
var solrCollection = flag.String("solrCollection", "wikipedia", "Solr Collection")
var skip = flag.Int("skip", 0, "number to skip")


var listPage []*Page

type Page struct { 
	Id string
    Title string `xml:"title"` 
    ContributorName string `xml:"revision>contributor>username"` 
    ContributorIP string `xml:"revision>contributor>ip"` 
    Text string `xml:"revision>text"` 
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {

	flag.Parse()
	url := "http://" +  *solrHost + ":" + *solrPort + "/solr/"+ *solrCollection + "/update"

	xmlFile, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

    decoder := xml.NewDecoder(xmlFile) 

	initial := -1
	var inElement string
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			// If we just read a StartElement token
			inElement = se.Name.Local
			// ...and its name is "page"
			if inElement == "page" {
				var p Page
				// decode a whole chunk of following XML into the
				// variable p which is a Page (se above)
				if(initial < *skip){
					initial=initial+ 1
					fmt.Printf("Skipping : %v\n",  initial)
					continue
				}
				decoder.DecodeElement(&p, &se)
			    accumulePage(&p, url)
			}
		default:
		}
	}
	
}

func accumulePage(page *Page, url string){
	page.generateId()
	page.Text = wiki.NormalizeText(page.Text)
	fmt.Println("Title : " + page.Title)
	listPage = append(listPage, page)
	numPages := 100
	if len(listPage) == numPages {
		fmt.Printf("indexing %v pages", numPages)
		jsonPage, _ := json.Marshal(listPage)
    	json := string(jsonPage)
    	//fmt.Println(json)
    	go solr.IndexPages(json, url)
    	listPage = listPage[len(listPage):]
	}
}

func (p *Page) generateId() {
	p.Id = strings.ToLower(p.Title)
	p.Id = strings.Replace(p.Id, " ", "_", -1)
}
