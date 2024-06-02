package collector

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"

	"github.com/dmitrygulevich2000/mwormix-bot/internal/collector/parser"
	"github.com/dmitrygulevich2000/mwormix-bot/internal/storage"
)

type BonusCollector struct {
	searches storage.Storage
	pageURL  string
}

func New(storage storage.Storage, pageUrl string) *BonusCollector {
	return &BonusCollector{
		searches: storage,
		pageURL:  pageUrl,
	}
}

type Bonus struct {
	PublishedAt time.Time
	Link        string
	Code        string
	ValidUntil  string
}

func (c *BonusCollector) CollectNewBonuses() ([]Bonus, error) {
	now := time.Now()
	lastSearch, err := c.searches.LastSearchTime()
	if err != nil {
		return nil, err
	}

	bonuses, err := c.collectBonusesAfter(lastSearch)
	if err == nil {
		c.searches.StoreSearch(now)
	}
	return bonuses, err
}

var codeRegexp regexp.Regexp = *regexp.MustCompile("Промокод: ([a-zA-Z0-9]+)")
var validUntilRegexp regexp.Regexp = *regexp.MustCompile("действител(ен|ьна|ьно) (.+)$")

// дата публикации:
// 25 дек 2023 - тут неразрывный пробел
// 1 янв в 16:38
// вчера в 12:00
// сегодня в 12:00

// ---- ниже есть атрибут abs_time
// час назад
// 4 минуты назад
// две/три минуты назад
// минуту назад
// 12 секунд назад
// только что

// действителен до 04:00 - следующий день
// действителен до 02.01.24 23:59 - дата в 2 слова

// слово "промокод" не в подходящем посте

func (c *BonusCollector) collectBonusesAfter(ts time.Time) ([]Bonus, error) {
	ts = ts.Truncate(time.Minute)

	req, err := http.NewRequest("GET", c.pageURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content := charmap.Windows1251.NewDecoder().Reader(resp.Body)

	posts, err := parser.CollectPosts(content)
	if err != nil {
		return nil, err
	}

	result := []Bonus{}
	for _, post := range posts {
		if post.PublishedAt.Compare(ts) < 0 {
			continue
		}
		if !strings.Contains(post.FlatText, "промокод") {
			continue
		}

		var code, valid string
		for _, line := range post.TextContent {
			codeMatch := codeRegexp.FindStringSubmatch(line)
			if codeMatch != nil {
				code = codeMatch[1]
			}

			validMatch := validUntilRegexp.FindStringSubmatch(line)
			if validMatch != nil {
				valid = validMatch[2]
			}
		}

		result = append(result, Bonus{
			PublishedAt: post.PublishedAt,
			Link:        "https://vk.com" + post.Link,
			Code:        code,
			ValidUntil:  valid,
		})
	}

	return result, nil
}
