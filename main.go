package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"
)

type Post struct {
	ID        string     `json:"id"`
	Kind      int        `json:"kind"`
	Pubkey    string     `json:"pubkey"`
	CreatedAt int        `json:"created_at"`
	Content   string     `json:"content"`
	Tags      [][]string `json:"tags"`
	Sig       string     `json:"sig"`
}

var htmlBaseHaed = `
<html>
<head>
<meta charset="UTF-8">
<title>ryo_grid Nostr posts</title>
<link rel="stylesheet" href="./posts.css">
</head>
<body>
`

var htmlBaseTail = `
</body>
</html>
`

func main() {
	f, err := os.Open("./ryo_grid_posts_dump2.json")
	if err != nil {
		fmt.Println("error")
	}

	allPosts := make([]Post, 0)
	yearMonthList := make([]string, 0)
	curMonth := ""
	curYear := ""

	r := bufio.NewReader(f)
	for {
		b, err_ := r.ReadBytes('\n')
		if err_ == io.EOF {
			break
		}
		var tmpPost Post
		json.Unmarshal(b, &tmpPost)
		//if !strings.Contains(tmpPost.Content, "\"e\"") && !strings.Contains(tmpPost.Content, "\"created_at\"") {
		if tmpPost.Kind == 1 {
			allPosts = append(allPosts, tmpPost)
		}
	}

	sort.Slice(allPosts, func(ii, jj int) bool { return allPosts[ii].CreatedAt < allPosts[jj].CreatedAt })

	monthPosts := make([]Post, 0)
	for _, post_ := range allPosts {
		dtFromUnix := time.Unix(int64(post_.CreatedAt), 0)
		tmpMonth := dtFromUnix.Month().String()
		if curMonth != tmpMonth {
			curYearMonth := curYear + curMonth

			outputHtml := htmlBaseHaed
			for _, monthPost := range monthPosts {
				outputHtml += "<p class='datetime'><strong>"
				dt := time.Unix(int64(monthPost.CreatedAt), 0)
				outputHtml += dt.Format("2006/01/02 15:04:05")
				outputHtml += "</strong></p><br/>"
				outputHtml += "<p class='content'>"
				outputHtml += monthPost.Content
				outputHtml += "</p><br/><br/>"
			}
			outputHtml += htmlBaseTail

			yearMonthList = append(yearMonthList, curYearMonth)
			if curYearMonth != "" {
				ioutil.WriteFile("./"+curYearMonth+".html", []byte(outputHtml), 0444)
			}

			curMonth = tmpMonth
			curYear = strconv.Itoa(dtFromUnix.Year())
			monthPosts = make([]Post, 0)
		}
		monthPosts = append(monthPosts, post_)
	}

	indexHtml := htmlBaseHaed
	for _, yearMonth := range yearMonthList {
		indexHtml += "<p><a href='./" + yearMonth + ".html'>" + yearMonth + "</a></p><br/>"
	}
	indexHtml += htmlBaseTail
	ioutil.WriteFile("./index.html", []byte(indexHtml), 0444)
}
