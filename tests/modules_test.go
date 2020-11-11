package main

import (
	"akito/downloader"
	"fmt"
	"os"
	"testing"
)

func Test_WuxiaWorldCo_Search(t *testing.T) {
	results, err := downloader.Mod_WuxiaWorldCo.Search("divine")
	if err != nil {
		t.Errorf("[wuxiaworld.co] Error during search: %s.", err.Error())
		return
	}

	if len(results) == 0 {
		t.Errorf("[wuxiaworld.co] Expected non-zero search results.")
		return
	}
} // Tests wuxiaworld.co search

func Test_WuxiaWorldCo_Info(t *testing.T) {
	results, err := downloader.Mod_WuxiaWorldCo.Search("coiling dragon")
	if err != nil || len(results) == 0 {
		t.Errorf("[wuxiaworld.co] Error during search for novel info retrieval.")
		return
	}

	info, err := downloader.Mod_WuxiaWorldCo.NovelInfo(results[1])
	if err != nil {
		t.Errorf("[wuxiaworld.co] Error during novel info retrieval: %s.", err.Error())
		return
	}

	if info.Author != "I Eat Tomatoes" {
		t.Errorf("[wuxiaworld.co] Author incorrect, got: %s, expected: I Eat Tomatoes.", info.Author)
		return
	} else if info.Status != "Completed" {
		t.Errorf("[wuxiaworld.co] Status incorrect, got: %s, expected: Completed.", info.Status)
		return
	} else if info.Language != "Chinese" {
		t.Errorf("[wuxiaworld.co] Language incorrect, got: %s, expected: Chinese", info.Language)
		return
	} else if len(info.ChapterURLs) != 808 {
		t.Errorf("[wuxiaworld.co] Number of chapters incorrect, got: %d, expected: 808", len(info.ChapterURLs))
		return
	}
} // Tests wuxiaworld.co novel info retrieval

func Test_WuxiaWorldCo_Download(t *testing.T) {
	results, err := downloader.Mod_WuxiaWorldCo.Search("coiling dragon")
	if err != nil || len(results) == 0 {
		t.Errorf("[wuxiaworld.co] Error during search for chapter download.")
		return
	}

	info, err := downloader.Mod_WuxiaWorldCo.NovelInfo(results[1])
	if err != nil || info.Author != "I Eat Tomatoes" || info.Status != "Completed" || info.Language != "Chinese" ||
		len(info.ChapterURLs) != 808 {
		t.Errorf("[wuxiaworld.co] Error during novel info retrieval for chapter download.")
		return
	}

	// note that start and end are zero-indexed
	chapters, err := downloader.DownloadRange(downloader.Mod_WuxiaWorldCo, info, 0, len(info.ChapterURLs)-1)
	if err != nil {
		t.Errorf("[wuxiaworld.co] Error during chapter download: %s.", err.Error())
		return
	}

	fmt.Println(len(chapters))
	testFile, err := os.Create("downloaded.txt")
	if err != nil {
		panic(err)
	}
	for _, chapter := range chapters {
		_, _ = testFile.WriteString(fmt.Sprintf("== %s ==", chapter.Title))
		_, _ = testFile.WriteString("\n")
		_, _ = testFile.WriteString(chapter.Content)
		_, _ = testFile.WriteString("\n")
		_, _ = testFile.WriteString("\n=======================================\n")
		_, _ = testFile.WriteString("\n")
	}
} // Tests wuxiaworld.co singular chapter download
