package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	flag "github.com/ogier/pflag"
)

const urlTemplate = `https://qiita.com/api/v1/tags/%s/items?per_page=%d`

// Entry is foo
type Entry struct {
	StockCount int    `json:"stock_count"`
	Title      string `json:"title"`
	URL        string `json:"url"`
}

func constructURL(key string, num int) string {
	return fmt.Sprintf(urlTemplate, key, num)
}

func getEntries(key string, num int) ([]Entry, error) {
	url := constructURL(key, num)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	jsonByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	data := make([]Entry, num)
	if err := json.Unmarshal(jsonByte, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func main() {
	num := flag.IntP("number", "n", 10, "numbers of entry")
	peco := flag.BoolP("peco", "p", false, "title and url are joined by null chracter")
	flag.Parse()

	key := flag.Arg(0)
	if key == "" {
		fmt.Printf("Usage: qiita [options] keyword\n")
		os.Exit(1)
	}

	entries, err := getEntries(key, *num)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i, entry := range entries {
		fmt.Printf("%2d: %s [%d]", i+1, entry.Title, entry.StockCount)
		if *peco {
			fmt.Printf("\x00%s", entry.URL)
		}
		fmt.Printf("\n")
	}
}
