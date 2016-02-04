package main

	
import (
	"fmt"
    "flag"
    "github.com/rtt/Go-Solr"
)

var solrHost = flag.String("solrHost", "localhost", "Solr Host")
var solrPort = flag.Int("solrPort", 8080, "Solr Port")
var solrCollection = flag.String("solrCollection", "wikipedia", "Solr Collection")
var query = flag.String("query", "*:*", "Solr Query ")
var qf = flag.String("qf", "Title^100 ContributorName^10 ContributorIP^10 Text Title_ngram^50 ContributorName_ngram^3 ContributorIP_ngram^3 Text_ngram", "Solr Query ")

func main(){

	flag.Parse()

	s, err := solr.Init(*solrHost, int(*solrPort), *solrCollection)

	if err != nil {
		fmt.Println(err)
		return
	}

	q := solr.Query{
		Params: solr.URLParamMap{
			"q":           []string{*query},
			"qf":           []string{*qf},
			"fl":          []string{"Title", "ContributorName"},
		},
		Rows: 10,
	}

	// perform the query, checking for errors
	res, err := s.Select(&q)

	if err != nil {
		fmt.Println(err)
		return
	}
	
	// grab results for ease of use later on
	results := res.Results

	// print a summary and loop over results, priting the "title" and "latlng" fields
	fmt.Println(
		fmt.Sprintf("Query: %#v\nHits: %d\nNum Results: %d\nQtime: %d\nStatus: %d\n\nResults\n-------\n",
			q,
			results.NumFound,
			results.Len(),
			res.QTime,
			res.Status))

	for i := 0; i < results.Len(); i++ {
		fmt.Println(i + 1, " - " , results.Get(i).Field("Title"))
		fmt.Println("\t", "Contributor Name/IP;", results.Get(i).Field("ContributorName"))
	}

}
