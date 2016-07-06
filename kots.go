package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
)

const baseurl = `http://www.tosarang2.net`

// args:
// initial search fragment "무한도전"
// regex for links `^무한도전.+151226\.HDTV\.H264\.720p-WITH$`

// output to stdout: magnet link
// http://www.tosarang2.net/bbs/magnet:?xt=urn:btih:61A35E62EFED0544ACACC2F1E3FF4DB3DC6B9A36&dn=무한도전.E456.151128.HDTV.H264.720p-WITH

func main() {
	bow := surf.NewBrowser()
	logger := log.New(os.Stderr, "kots: ", log.Lshortfile)
	err := bow.Open(baseurl)
	if err != nil {
		panic(err)
	}

	var fileRegex, show string

	flag.StringVar(&fileRegex, "regex", "INVALID REGEX?!!", `regexp for the file: ^무한도전.+151226\.HDTV\.H264\.720p-WITH$`)
	flag.StringVar(&show, "show", "INVALID REGEX?!!", `show name that is used for initial search: 무한도전`)
	flag.Parse()

	validLink := regexp.MustCompile(fileRegex)
	validMagnetLink := regexp.MustCompile(`magnet:\?xt=urn:btih:\w+`)

	// <a href="http://www.tosarang2.net/bbs/board.php?bo_table=torrent_kortv_ent" title="한국TV > 예능/오락">예능</a>
	logger.Println("found page: " + bow.Title())

	err = bow.Click("a[title='한국TV > 예능/오락']")
	if err != nil {
		panic(err)
	}

	searchFm, err := bow.Form("form#fsearch")
	if err != nil {
		panic(err)
	}

	searchFm.Input("stx", show)
	err = searchFm.Submit()
	if err != nil {
		panic(err)
	}

	bow.Dom().Find("div#bo_l_list table tbody tr td.td_subject a").Each(func(_ int, s *goquery.Selection) {
		rawHTML, err := s.Html()
		if err != nil {
			panic(err)
		}

		if validLink.MatchString(rawHTML) {
			href, ok := s.Attr("href")
			if !ok {
				logger.Fatalf("could not find href from anchor\n")
			}

			logger.Printf("found link: %s => %s\n", rawHTML, href)

			err = bow.Open(href)
			if err != nil {
				panic(err)
			}

			bow.Dom().Find("div.bo_v_file a").Each(func(_ int, s *goquery.Selection) {
				rawHTML, err := s.Html()
				if err != nil {
					panic(err)
				}

				if validMagnetLink.MatchString(rawHTML) {
					href, ok := s.Attr("href")
					if !ok {
						logger.Fatalf("could not find href from anchor\n")
					}

					logger.Printf("found magnet link: %s => %s\n", rawHTML, href)
					fmt.Printf("%s\n", href)
				}
			})
		}
	})
}
