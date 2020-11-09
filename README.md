# novel-downloader
Automated crawler for searching, downloading, and exporting novels from various online sources. 

## Installation
`go get -u github.com/albert-sun/novel-downloader`

## Website Support
The below table lists current support for websites along with development remarks, if specified. 

| Symbol | Meaning |
|:------:|:--------|
| ✅ | Full website support |
| 🕒 | Development in-progress |
| ⚠️ | Current issues being resolved (check remarks) |

| Name | Type | Language(s) | Support | Remarks |
|:-----|:-----|:------------|:-------:|:--------|
| WuxiaWorld.co | Aggregate | Chinese | 🕒 ||


## Quickstart Guide
Website support is separated into individual modules which can be accessed via the Mod_WebsiteName exported variables. Each module represents an (implemented) interface containing functions for searching, retrieving info for, and downloading novels.

```
// search via wuxiaworld.co, returns slice of results
results, err := modules.Mod_WuxiaWorldCo.Search("I Shall Seal The Heavens") 
if err != nil { 
    // handle error
} else if len(results) == 0 {
    // handle lack of search results
}

// retrieve info from URL with basic info and chapter URLs
novelInfo, err := modules.Mod_WuxiaWorldCo.NovelInfo(results[0])
if err != nil {
    // handle error
}

// download all chapters, returns a slice of chapter contents
allChapters, err := modules.Mod_WuxiaWorldCo.DownloadAll(novelInfo)
if err != nil {
    // handle error
}

// to be implemented: exporting to different filetypes
```

## Contributions
Contributions are always appreciated, especially in the form of new modules for websites and export formats. Thank you for your help!

Released under the GNU General Public License
