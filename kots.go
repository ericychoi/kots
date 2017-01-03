package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/browser"
)

const baseurl = `https://torrentkim5.net/bbs`

var logger *log.Logger

// args:
// initial search fragment "무한도전"
// regex for links `^무한도전.+151226\.HDTV\.H264\.720p-WITH$`

// output to stdout: magnet link
// http://www.tosarang2.net/bbs/magnet:?xt=urn:btih:61A35E62EFED0544ACACC2F1E3FF4DB3DC6B9A36&dn=무한도전.E456.151128.HDTV.H264.720p-WITH
func main() {
	bow := surf.NewBrowser()
	logger = log.New(os.Stderr, "kots: ", log.Lshortfile)
	var fileRegex, show string

	flag.StringVar(&fileRegex, "regex", "INVALID REGEX?!!", `regexp for the file: ^무한도전.+151226\.HDTV\.H264\.720p-WITH$`)
	flag.StringVar(&show, "show", "INVALID REGEX?!!", `show name that is used for initial search: 무한도전`)
	flag.Parse()

	encodedKeyword := url.QueryEscape(show)
	err := bow.Open(baseurl + `/s.php?k=` + encodedKeyword)
	if err != nil {
		panic(err)
	}

	validLink := regexp.MustCompile(fileRegex)
	validMagnetLink := regexp.MustCompile(`Mag_dn\('([A-F0-9]{40})'\)`)

	logger.Println("found page: " + bow.Title())

	search(bow, validLink, validMagnetLink, show)
}

func search(bow *browser.Browser, validLink, validMagnetLink *regexp.Regexp, show string) {
	bow.Dom().Find("table.board_list tbody tr.bg1 td.subject a").Each(func(_ int, s *goquery.Selection) {
		rawHTML, err := s.Html()
		if err != nil {
			panic(err)
		}

		if validLink.MatchString(rawHTML) {
			logger.Printf("found link: %s\n", rawHTML)

			// get back to a couple of siblings
			s.Parent().Prev().Prev().Children().Each(func(_ int, s *goquery.Selection) {
				href, ok := s.Attr("href")
				if !ok {
					logger.Fatalf("could not find href from anchor\n")
				}

				results := validMagnetLink.FindStringSubmatch(href)
				if results == nil {
					logger.Fatalf("couldn't find match for %s", href)
				}

				//javascript:Mag_dn('72A3F8560D8FBB124F782035D01F961E3A8068AA')
				logger.Printf("found magnet link: %s\n", results[1])

				magnetLink := fmt.Sprintf(`magnet:?xt=urn:btih:%s`, results[1])
				fmt.Printf("%s\n", magnetLink)
			})
		}
	})
}
