package modules

import (
	"akito/downloader"
	"akito/utilities"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/valyala/fasthttp"
	"net/url"
)

// WuxiaWorld.co (Chinese) (Aggregate) - wuxiaworld.co
// Last Updated: 11/5/2020
// ===============================
// [ ] Search       [ ] Download
// [ ] Novel Info   [ ] Test cases
// ===============================
// Remarks:
// -

// Main downloader struct
type WuxiaWorldCo struct {
	downloader.BasicDownloader
	downloader.WebDownloader
}

// Exported module variable
var Mod_WuxiaWorldCo = func() *WuxiaWorldCo {
	return &WuxiaWorldCo{
		BasicDownloader: downloader.BasicDownloader{
			WebsiteName: "WuxiaWorld.co",
			WebsiteURL:  "wuxiaworld.co",
			Languages:   []string{"Chinese"},
			WebsiteType: "Aggregate",
			LastUpdated: "11/5/2020",
		},
	}
}()

func (m WuxiaWorldCo) Search(searchTerm string) ([]downloader.NovelBasic, error) {
	httpClient := fasthttp.Client{} // eventually replace with "global" usage

	const searchURLFmt = "https://www.wuxiaworld.co/search/%s/%d" // formatting for search URL
	const resultsQuery = ".list-item .item-info"                  // info for each search result
	escapedTerm := url.PathEscape(searchTerm)                     // URL-encode escape for search

	var totalResults int
	var pagesResult [][]downloader.NovelBasic // slice of slices of results from each page
	for pageNum := 1; ; pageNum++ {           // forever, iterate while search results exist
		// perform search request
		searchURL := fmt.Sprintf(searchURLFmt, escapedTerm, pageNum)
		response, err := utilities.RequestGET(&httpClient, searchURL)
		if err != nil { // return if request error
			return nil, downloader.RequestError
		}
		respReader := bytes.NewReader(response.Body()) // performance issue?
		fasthttp.ReleaseResponse(response)

		// parse HTML for querying purposes
		document, err := goquery.NewDocumentFromReader(respReader)
		if err != nil { // return if request error
			return nil, downloader.ParseHTMLError
		}

		// query search results from document
		searchResults := document.Find(resultsQuery)
		if len(searchResults.Nodes) == 0 { // no more results - none or last page
			break
		}
		pageResults := make([]downloader.NovelBasic, len(searchResults.Nodes))

		// iterate over query result and get info
		var foundErr error // for errors caught during iteration
		searchResults.Each(func(index int, sel *goquery.Selection) {
			// quick break if error encountered
			if foundErr != nil {
				return
			}

			aElement := sel.Find("a")
			novelName := aElement.Text()               // get full novel name
			trimmed := utilities.TrimString(novelName) // trim uppercase and special characters
			novelURL, exists := aElement.Attr("href")
			if !exists { // return if attribute not found
				foundErr = downloader.GoQueryError
				return
			}

			pageResults[index] = downloader.NovelBasic{
				Name:     novelName,
				NameTrim: trimmed,
				NovelURL: fmt.Sprintf("https://%s%s", m.WebsiteURL, novelURL),
			}
		})
		if foundErr != nil { // return if error caught
			return nil, foundErr
		}

		totalResults += len(pageResults) // increment number of total results
		pagesResult = append(pagesResult, pageResults)
	}

	// transform results from 2-dim to 1-dim slice
	// I know it's disgusting, I'd love to find an easier way
	var currentIndex int
	arrangedResults := make([]downloader.NovelBasic, totalResults)
	for _, pageResults := range pagesResult {
		for _, result := range pageResults {
			arrangedResults[currentIndex] = result
			currentIndex++
		}
	}

	return arrangedResults, nil
}

func (m WuxiaWorldCo) NovelInfo(basic downloader.NovelBasic) (*downloader.NovelInfo, error) {
	httpClient := fasthttp.Client{} // eventually replace with "global" usage

	// retrieve synopsis page information
	response, err := utilities.RequestGET(&httpClient, basic.NovelURL)
	if err != nil { // return if request error
		fmt.Println(err.Error())
		return nil, downloader.RequestError
	}
	respReader := bytes.NewReader(response.Body()) // performance issue?
	fasthttp.ReleaseResponse(response)

	// parse HTML for querying purposes
	document, err := goquery.NewDocumentFromReader(respReader)
	if err != nil { // return if request error
		return nil, downloader.ParseHTMLError
	}

	// query author name from document
	authorQuery := ".name"
	searchResults := document.Find(authorQuery)
	if len(searchResults.Nodes) == 0 { // author name not found
		return nil, downloader.GoQueryError
	}
	author := searchResults.Text() // should capture only author name

	// query status from document
	statusQuery := ".book-state .txt"
	searchResults = document.Find(statusQuery)
	if len(searchResults.Nodes) == 0 { // status name not found
		return nil, downloader.GoQueryError
	}
	status := searchResults.Text() // should capture only author name

	// language is always Chinese
	language := "Chinese"

	// get chapters list
	chapterQuery := ".chapter-item" // all chapter a-href
	searchResults = document.Find(chapterQuery)
	if len(searchResults.Nodes) == 0 { // chapters not found
		return nil, downloader.GoQueryError
	}

	var foundErr error // for errors caught during iteration
	chapterURLs := make([]string, len(searchResults.Nodes))
	searchResults.Each(func(index int, sel *goquery.Selection) {
		// quick break if error encountered
		if foundErr != nil {
			return
		}

		chapterURL, exists := sel.Attr("href")
		if !exists {
			foundErr = downloader.GoQueryError
			return
		}

		chapterURLs[index] = chapterURL
	})

	info := downloader.NovelInfo{
		NovelBasic: basic,
		Author:     author,
		Status:     status,
		Language:   language,
		Chapters:   chapterURLs,
	}

	return &info, nil
}
