package downloader

import "github.com/remeh/sizedwaitgroup"

// BasicDownloader holds a basic info and utility functions for downloader modules.
type BasicDownloader struct {
	Downloader
	WebsiteURL string   // website URL
	Type       string   // [Unused] type of website (translator, aggregate, etc.)
	Languages  []string // website languages (original or translated)
}

// WebDownloader represents an individual website downloader "module".
// Contains functions for retrieval, search, etc. and metadata.
type Downloader interface {
	Search(string) ([]NovelBasic, error) // searches website using the given query string
	NovelInfo(NovelBasic) (*NovelInfo, error)
	Download(info DownloadInfo) // retrieves information from a novel page
}

// NovelBasic represents a novel's basic metadata.
// Includes formatted novel name for keeping track (sometimes punctuation differs between websites)
type NovelBasic struct {
	Name     string // Novel name
	NameTrim string // Trimmed novel name (lowercase, no space, no punctuation)
	NovelURL string // URL for novel synopsis for future retrieval
}

// NovelInfo represents an individual novel (from a website).
// Contains metadata and slice of chapter URLs for easy access.
type NovelInfo struct {
	NovelBasic           // basic metadata
	Author      string   // Author name (hopefully english?)
	Status      string   // Ongoing, Completed, Dropped (?)
	Language    string   // Novel (translated from) language / English
	ChapterURLs []string // Novel chapter URLs (how to best number?)
}

// NovelChapter represents the content of a novel chapter.
type NovelChapter struct {
	Title   string // chapter title
	Content string // actual content
}

// DownloadInfo contains info for downloading a given chapter
type DownloadInfo struct {
	SWG       *sizedwaitgroup.SizedWaitGroup
	Chapters  []*NovelChapter // downloaded chapters
	NovelInfo *NovelInfo      // novel stuff
	Start     int             // starting index
	Index     int             // index of chapters
	FoundErr  *error          // global error thing?
}
