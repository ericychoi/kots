package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"

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
	bow.Dom().Find("table.board_list tbody tr.bg1 td.subject a").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		found := false
		rawHTML, err := s.Html()
		if err != nil {
			logger.Printf("couldn't find html: %s\n", err)
			return true
		}
		if validLink.MatchString(rawHTML) {
			logger.Printf("found link: %s\n", rawHTML)
			s.Parent().Prev().Children().EachWithBreak(func(_ int, s *goquery.Selection) bool {
				scoreHTML, err := s.Html()
				if err != nil {
					logger.Printf("could not find score: %s\n", err)
					return true
				}
				score, err := strconv.Atoi(scoreHTML)
				if err != nil {
					logger.Printf("score HTML %s not an int: %s\n", scoreHTML, err)
					return true
				}
				if score < 0 {
					logger.Printf("negative score %d found\n", score)
					return true
				}

				s.Parent().Prev().ChildrenFiltered("a").EachWithBreak(func(_ int, s *goquery.Selection) bool {
					href, ok := s.Attr("href")
					if !ok {
						logger.Printf("could not find href from anchor\n")
						return true
					}

					results := validMagnetLink.FindStringSubmatch(href)
					if results == nil {
						logger.Printf("couldn't find match for %s", href)
						return true
					}

					//javascript:Mag_dn('72A3F8560D8FBB124F782035D01F961E3A8068AA')
					logger.Printf("found magnet link: %s\n", results[1])

					magnetLink := fmt.Sprintf(`magnet:?xt=urn:btih:%s`, results[1])
					fmt.Printf("%s\n", magnetLink)
					found = true
					return false
				})
				return false
			})
			return !found
		}
		return !found
	})
}
