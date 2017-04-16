package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/browser"
)

const baseurl = `https://torrentwiz3.com/torrent_ent`

var logger *log.Logger

// args:
// initial search fragment "무한도전"
// regex for links `^무한도전.+151226\.HDTV\.H264\.720p-WITH$`

// output to stdout: magnet link
func main() {
	bow := surf.NewBrowser()
	logger = log.New(os.Stderr, "kots: ", log.Lshortfile)
	var fileRegex, show string

	flag.StringVar(&fileRegex, "regex", "INVALID REGEX?!!", `regexp for the file: ^무한도전.+151226\.HDTV\.H264\.720p-WITH$`)
	flag.StringVar(&show, "show", "INVALID REGEX?!!", `show name that is used for initial search: 무한도전`)
	flag.Parse()

	link := fmt.Sprintf("%s", baseurl)
	logger.Printf("link: %s", link)
	err := bow.Open(link)
	if err != nil {
		panic(err)
	}

	validLink := regexp.MustCompile(fileRegex)
	validMagnetLink := regexp.MustCompile(`magnet:\?xt=urn:btih:[A-F0-9]{40}`)

	logger.Println("found page: " + bow.Title())
	search(bow, validLink, validMagnetLink, show)
}

func search(bow *browser.Browser, validLink, validMagnetLink *regexp.Regexp, show string) {
	bow.Dom().Find("div.ranking_content div.row ul li a").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		var found bool
		link, exists := s.Attr("href")
		// logger.Printf("found link: %s\n", link)
		if !exists {
			return true
		}

		rawHTML, err := s.Children().First().Next().Html()
		// logger.Printf("found html: %s", rawHTML)
		if err != nil {
			logger.Printf("couldn't find html: %s\n", err)
			return true
		}

		if validLink.MatchString(rawHTML) {
			logger.Printf("found: %s\n", rawHTML)
			logger.Printf("opening link: %s\n", link)

			err := bow.Open(link)
			if err != nil {
				panic(err)
			}
			bow.Dom().Find("div.panel-heading div.font-12 a.view_file_download").EachWithBreak(func(_ int, ss *goquery.Selection) bool {
				href, exists := ss.Attr("href")
				if !exists {
					logger.Printf("could not find href from anchor\n")
					return true
				}

				if !validMagnetLink.MatchString(href) {
					logger.Printf("couldn't find match for %s", href)
					return true
				}
				logger.Printf("found magnet link: %s\n", href)
				fmt.Printf("%s\n", href)
				found = true
				return false
			})
		}
		return !found // not found means continue (EachWithBreak continues if it gets true)
	})
}
