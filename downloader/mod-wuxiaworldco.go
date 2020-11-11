package downloader

import (
	"akito/utilities"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cockroachdb/errors"
	"github.com/valyala/fasthttp"
	"net/url"
	"strings"
)

// WuxiaWorld.co (Chinese) (Aggregate) - wuxiaworld.co
// Last Updated: 11/5/2020
// ===============================
// [X] Search       [ ] Download
// [X] Novel Info   [ ] Test cases
// ===============================
// Remarks:
// -

// Main downloader struct
type modWuxiaWorldCo struct {
	BasicDownloader
}

// Exported module variable
var WuxiaWorldCo *modWuxiaWorldCo

func init() {
	WuxiaWorldCo = &modWuxiaWorldCo{
		BasicDownloader: BasicDownloader{
			WebsiteName: "WuxiaWorld.co",
			WebsiteURL:  "wuxiaworld.co",
			Languages:   []string{"Chinese"},
			WebsiteType: "Aggregate",
			LastUpdated: "11/5/2020",
		},
	}
	Modules = append(Modules, WuxiaWorldCo)
}

func (m modWuxiaWorldCo) Search(searchTerm string) ([]NovelBasic, error) {
	httpClient := fasthttp.Client{} // eventually replace with "global" usage

	const searchURLFmt = "https://www.wuxiaworld.co/search/%s/%d" // formatting for search URL
	const resultsQuery = ".list-item .item-info"                  // info for each search result
	escapedTerm := url.PathEscape(searchTerm)                     // URL-encode escape for search

	var totalResults int
	var pagesResult [][]NovelBasic  // slice of slices of results from each page
	for pageNum := 1; ; pageNum++ { // forever, iterate while search results exist
		// perform search request
		searchURL := fmt.Sprintf(searchURLFmt, escapedTerm, pageNum)
		response, err := utilities.RequestGET(&httpClient, searchURL)
		if err != nil { // return if request error
			return nil, ErrRequest
		}
		respReader := bytes.NewReader(response.Body()) // performance issue?
		fasthttp.ReleaseResponse(response)

		// parse HTML for querying purposes
		document, err := goquery.NewDocumentFromReader(respReader)
		if err != nil { // return if request error
			return nil, ErrParseHTML
		}

		// query search results from document
		searchResults := document.Find(resultsQuery)
		if len(searchResults.Nodes) == 0 { // no more results - none or last page
			break
		}
		pageResults := make([]NovelBasic, len(searchResults.Nodes))

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
				foundErr = errors.Wrapf(ErrGoQuery, "get novel URL")
				return
			}

			pageResults[index] = NovelBasic{
				Name:     novelName,
				NameTrim: trimmed,
				NovelURL: fmt.Sprintf("https://www.%s%s", m.WebsiteURL, novelURL),
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
	arrangedResults := make([]NovelBasic, totalResults)
	for _, pageResults := range pagesResult {
		for _, result := range pageResults {
			arrangedResults[currentIndex] = result
			currentIndex++
		}
	}

	return arrangedResults, nil
}

func (m modWuxiaWorldCo) NovelInfo(basic NovelBasic) (*NovelInfo, error) {
	httpClient := fasthttp.Client{} // eventually replace with "global" usage

	// retrieve synopsis page information
	response, err := utilities.RequestGET(&httpClient, basic.NovelURL)
	if err != nil { // return if request error
		return nil, ErrRequest
	}
	respReader := bytes.NewReader(response.Body()) // performance issue?
	fasthttp.ReleaseResponse(response)

	// parse HTML for querying purposes
	document, err := goquery.NewDocumentFromReader(respReader)
	if err != nil { // return if request error
		return nil, ErrParseHTML
	}

	// query author name from document
	authorQuery := ".name"
	searchResults := document.Find(authorQuery)
	if len(searchResults.Nodes) == 0 { // author name not found
		return nil, errors.Wrapf(ErrGoQuery, "get novel name")
	}
	author := searchResults.Text() // should capture only author name

	// query status from document
	statusQuery := ".book-state .txt"
	searchResults = document.Find(statusQuery)
	if len(searchResults.Nodes) == 0 { // status name not found
		return nil, errors.Wrapf(ErrGoQuery, "get novel status")
	}
	status := searchResults.Text() // should capture only author name

	// language is always Chinese
	language := "Chinese"

	// get chapters list
	chapterQuery := ".chapter-item" // all chapter a-href
	searchResults = document.Find(chapterQuery)
	if len(searchResults.Nodes) == 0 { // chapters not found
		return nil, errors.Wrapf(ErrGoQuery, "get chapter elements")
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
			foundErr = errors.Wrapf(ErrGoQuery, "get chapter URLS")
			return
		}

		chapterURLs[index] = fmt.Sprintf("https://www.%s%s", m.WebsiteURL, chapterURL)
	})

	info := NovelInfo{
		NovelBasic:  basic,
		Author:      author,
		Status:      status,
		Language:    language,
		ChapterURLs: chapterURLs,
	}

	return &info, nil
}

func (m modWuxiaWorldCo) Download(dlInfo DownloadInfo) {
	// assume that index is valid btw
	// check for existence of previous error
	if *dlInfo.FoundErr != nil {
		return
	}

	httpClient := fasthttp.Client{} // eventually replace with "global" usage

	// retrieve chapter page
	response, err := utilities.RequestGET(&httpClient, dlInfo.NovelInfo.ChapterURLs[dlInfo.Index])
	if err != nil { // return if request error
		if *dlInfo.FoundErr != nil { // check if previous error
			*dlInfo.FoundErr = errors.Wrapf(ErrRequest, "chapter %d", dlInfo.Index)
		}
		return
	}
	respReader := bytes.NewReader(response.Body()) // performance issue?
	fasthttp.ReleaseResponse(response)

	dlInfo.SWG.Done() // let parsing not be throttled

	// parse HTML for querying purposes
	document, err := goquery.NewDocumentFromReader(respReader)
	if err != nil { // return if request error
		if *dlInfo.FoundErr != nil { // check if previous error
			*dlInfo.FoundErr = errors.Wrapf(ErrParseHTML, "chapter %d", dlInfo.Index)
		}
		return
	}

	// query chapter title from document
	titleQuery := ".chapter-title"
	searchResults := document.Find(titleQuery)
	if len(searchResults.Nodes) == 0 { // title not found
		if *dlInfo.FoundErr != nil { // check if previous error
			*dlInfo.FoundErr = errors.Wrapf(
				ErrParseHTML,
				"chapter %d: get chapter title",
				dlInfo.Index,
			)
		}
		return
	}
	chapterTitle := searchResults.Text()

	// query chapter content from document
	// note: contains ad content, not sure how parsing works
	contentQuery := ".chapter-entity"
	searchResults = document.Find(contentQuery)
	if len(searchResults.Nodes) == 0 { // title not found
		if *dlInfo.FoundErr != nil { // check if previous error
			*dlInfo.FoundErr = errors.Wrapf(
				ErrParseHTML,
				"chapter %d: get chapter content",
				dlInfo.Index,
			)
		}
		return
	}

	// iterate over div content and splice out br and script
	// ensure that only one newline between lines using fancy boolean
	flipNewline := true // prevents unneeded newlines
	var chapterContent string
	searchResults.Contents().Each(func(_ int, sel *goquery.Selection) {
		// replace br with newline, add trimmed content
		if sel.Is("br") && !flipNewline {
			flipNewline = true
			chapterContent += "\n"
		} else if !sel.Is("script") { // skip script
			flipNewline = false
			content := sel.Text()

			// skip annoying "please go to wuxiaworld.co" thing and empty lines
			if !strings.HasPrefix(content, "Please go tohttps") && content != "" {
				chapterContent += strings.TrimSpace(sel.Text())
			}
		}
	})

	dlInfo.Chapters[dlInfo.Index-dlInfo.Start] = &NovelChapter{
		Title:   chapterTitle,
		Content: chapterContent,
	}
}
