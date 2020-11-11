package downloader

import (
	"github.com/remeh/sizedwaitgroup"
)

const maxDownloads = 100 // number of "threads"

// DownloadRange downloads a range of chapters for a novel given the starting and ending index.
// Note that both the starting and ending index are zero-indexed instead of one-indexed from the user.
// Returns a slice of downloaded chapter contents and error, if any.
func DownloadRange(dl Downloader, info *NovelInfo, start int, end int) ([]*NovelChapter, error) {
	// primitive validation for starting and ending indexes
	if start > end || start < 0 || end < 0 || start >= len(info.ChapterURLs) || end >= len(info.ChapterURLs) {
		return nil, ErrInvalidRange
	}

	// for processing downloaded chapters
	downloaded := make([]*NovelChapter, end-start+1)
	swg := sizedwaitgroup.New(maxDownloads) // last swg for channel

	// concurrently download chapters using SWG
	// note: goroutine writes directly to array
	var foundErr error
	for index := start; index <= end; index++ {
		swg.Add() // foundErr check inside Download
		go dl.Download(DownloadInfo{
			SWG:       &swg,
			Chapters:  downloaded,
			NovelInfo: info,
			Start:     start,
			Index:     index,
			FoundErr:  &foundErr,
		}) // concurrently download
	}

	swg.Wait() // wait for downloads to finish

	if foundErr == nil {
		return downloaded, nil
	}

	return nil, foundErr
}

// DownloadAll provides a simple wrapper for downloading all chapters of a novel.
func DownloadAll(dl Downloader, info *NovelInfo) ([]*NovelChapter, error) {
	return DownloadRange(dl, info, 0, len(info.ChapterURLs)-1)
}
