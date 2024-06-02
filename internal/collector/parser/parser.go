package parser

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"

	"github.com/dmitrygulevich2000/mwormix-bot/internal/utils"
)

type WallPost struct {
	PublishedAt time.Time
	Link        string
	TextContent []string
	FlatText    string
}

func CollectPosts(htmlDoc io.Reader) ([]WallPost, error) {
	doc, err := goquery.NewDocumentFromReader(htmlDoc)
	if err != nil {
		return nil, err
	}

	var posts []WallPost
	doc.Find("._post_content").Each(func(i int, s *goquery.Selection) {
		// collect text content
		content, flatText := collectText(s.FindMatcher(goquery.Single(".wall_post_text")))

		// get publish time
		publishedSelection := s.Find(".PostHeaderSubtitle__item")
		published := publishedSelection.Find(".rel_date").AttrOr("abs_time", publishedSelection.Text())
		publishedAt, err := timeFromVKString(published)
		if err != nil {
			return
		}

		// get link
		link := s.Find(".PostHeaderSubtitle__link").AttrOr("href", "")

		posts = append(posts, WallPost{
			PublishedAt: publishedAt,
			Link:        link,
			TextContent: content,
			FlatText:    flatText,
		})
	})

	return posts, nil
}

func collectText(s *goquery.Selection) ([]string, string) {
	var content []string
	var flat_content strings.Builder

	var collectFromSelection func(i int, s *goquery.Selection)
	collectFromSelection = func(i int, s *goquery.Selection) {
		if s.Get(0).Type == html.TextNode {
			text := strings.TrimSpace(s.Text())
			if len(text) > 0 {
				content = append(content, text)
				flat_content.WriteString(strings.ToLower(text))
				flat_content.WriteString("\n")
			}
			return
		}
		s.Contents().Each(collectFromSelection)
	}
	s.Each(collectFromSelection)

	return content, flat_content.String()
}

////////////////////////////////////////////////////////////////////////////////

func timeFromVKString(vkString string) (time.Time, error) {
	// TODO wrap errors, handle indices and other bad format
	vkString = strings.ReplaceAll(vkString, "\u00A0", " ")
	vkWords := strings.Split(vkString, " ")

	var date time.Time
	vkWords, err := extractDate(vkWords, &date)
	if err != nil {
		return time.Time{}, err
	}

	var hm time.Duration
	vkWords, err = extractTime(vkWords, &hm)
	if err != nil {
		return time.Time{}, err
	}

	date = date.Add(hm)
	return date, nil
}

var monthByString map[string]time.Month = map[string]time.Month{
	"янв": time.January,
	"фев": time.February,
	"мар": time.March,
	"апр": time.April,
	"май": time.May,
	"июн": time.June,
	"июл": time.July,
	"авг": time.August,
	"сен": time.September,
	"окт": time.October,
	"ноя": time.November,
	"дек": time.December,
}

func extractDate(vkWords []string, date *time.Time) ([]string, error) {
	ts := time.Now().In(utils.MSKLoc)
	ts = time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, utils.MSKLoc)

	if vkWords[0] == "сегодня" {
		*date = ts
		return vkWords[2:], nil
	}
	if vkWords[0] == "вчера" {
		*date = ts.Add(-24 * time.Hour)
		return vkWords[2:], nil
	}

	day, err := strconv.Atoi(vkWords[0])
	if err != nil {
		return vkWords, err
	}
	month, ok := monthByString[vkWords[1]]
	if !ok {
		return vkWords, fmt.Errorf("unexpected month word %s", vkWords[1])
	}
	if vkWords[2] == "в" {
		*date = time.Date(ts.Year(), month, day, 0, 0, 0, 0, utils.MSKLoc)
		return vkWords[3:], nil
	}

	year, err := strconv.Atoi(vkWords[2])
	if err != nil {
		return vkWords, err
	}
	*date = time.Date(year, month, day, 0, 0, 0, 0, utils.MSKLoc)
	return vkWords[3:], nil
}

func extractTime(vkWords []string, hm *time.Duration) ([]string, error) {
	var h, m time.Duration = 23, 59
	if len(vkWords) == 0 {
		*hm = h*time.Hour + m*time.Minute
		return vkWords, nil
	}

	cnt, err := fmt.Sscanf(vkWords[0], "%d:%d", &h, &m)
	if err != nil {
		return vkWords, err
	}
	if cnt != 2 {
		return vkWords, fmt.Errorf("word %s has wrong time format", vkWords[0])
	}

	*hm = h*time.Hour + m*time.Minute
	return vkWords[1:], nil
}
