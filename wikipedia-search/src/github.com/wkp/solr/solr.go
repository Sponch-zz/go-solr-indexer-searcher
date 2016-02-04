package solr

import (
  //  "bufio"
    "fmt"
    "bytes"
    "io/ioutil"
    "net/http"
    "time"
)

func Show(){
	fmt.Println("Spooooonch on the water!!!!")
}

func IndexPages(jsonPages string, url string){
	time.Sleep(1000 * time.Millisecond)

	var jsonStr = []byte(jsonPages)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    client.Timeout = 0
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
}