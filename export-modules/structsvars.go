package exporter

import downloader "akito/download-modules"

// BasicExporter holds basic info and utility functions for exporter modules.
type BasicExporter struct {
	Exporter
	ExportName string // type of file being exported
	Extension  string // output file type extension
}

// Exporter represents an individual exporter "module".
// Contains functions for parsing, formatting, and outputting.
// Note that out represents an arbitrary output file thing to write to.
type Exporter interface {
	InitExport() (interface{}, error)                            // initializes output file
	WriteInfo(interface{}, *downloader.NovelInfo) error          // writes metadata and stuff
	WriteChapters(interface{}, []*downloader.NovelChapter) error // writes table of contents and chapters
	SaveFile(interface{}, string) error                          // saves file using given filename
}
