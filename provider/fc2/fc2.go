package fc2

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/javtube/javtube-sdk-go/common/parser"
	"github.com/javtube/javtube-sdk-go/common/random"
	"github.com/javtube/javtube-sdk-go/model"
	"github.com/javtube/javtube-sdk-go/provider"
)

var _ provider.MovieProvider = (*FC2)(nil)

const (
	Name     = "fc2"
	Priority = 10
)

const (
	baseURL   = "https://adult.contents.fc2.com/"
	movieURL  = "https://adult.contents.fc2.com/article/%s/"
	sampleURL = "https://adult.contents.fc2.com/api/v2/videos/%s/sample"
)

type FC2 struct {
	*provider.Scraper
}

func New() *FC2 {
	return &FC2{
		Scraper: provider.NewScraper(Name, Priority, colly.NewCollector(
			colly.AllowURLRevisit(),
			colly.IgnoreRobotsTxt(),
			colly.UserAgent(random.UserAgent()))),
	}
}

func (fc2 *FC2) NormalizeID(id string) string {
	if ss := regexp.MustCompile(`^(?i)(?:FC2-.*?)?(\d+)$`).FindStringSubmatch(id); len(ss) == 2 {
		return ss[1]
	}
	return ""
}

func (fc2 *FC2) GetMovieInfoByID(id string) (info *model.MovieInfo, err error) {
	return fc2.GetMovieInfoByURL(fmt.Sprintf(movieURL, id))
}

func (fc2 *FC2) GetMovieInfoByURL(u string) (info *model.MovieInfo, err error) {
	homepage, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	info = &model.MovieInfo{
		ID:            path.Base(homepage.Path),
		Number:        fmt.Sprintf("FC2-%s", path.Base(homepage.Path)),
		Provider:      fc2.Name(),
		Homepage:      homepage.String(),
		Actors:        []string{},
		PreviewImages: []string{},
		Tags:          []string{},
	}

	c := fc2.Collector()

	// Headers
	c.OnXML(`//div[@class="items_article_headerInfo"]`, func(e *colly.XMLElement) {
		info.Title = e.ChildText(`.//h3`)
		info.Tags = e.ChildTexts(`.//section[@class="items_article_TagArea"]/div/a`)
		info.Maker = e.ChildText(`.//ul/li[last()]/a`)
		{ /* score */
			class := e.ChildAttr(`.//li[@class="items_article_StarA"]/a/p/span`, "class")
			info.Score = parser.ParseScore(regexp.MustCompile(`(\d+)$`).FindString(class))
		}
		{ /* release date */
			ss := strings.Split(e.ChildText(`.//div[@class="items_article_Releasedate"]/p`), ":")
			info.ReleaseDate = parser.ParseDate(ss[len(ss)-1])
		}
	})

	// Summary
	c.OnXML(`//section[@class="items_article_Contents"]/iframe`, func(e *colly.XMLElement) {
		d := c.Clone()
		d.OnXML(`//html/body/div`, func(e *colly.XMLElement) {
			info.Summary = strings.TrimSpace(e.Text)
		})
		d.Visit(e.Request.AbsoluteURL(e.Attr("src")))
	})

	// Thumb+Cover
	c.OnXML(`//div[@class="items_article_MainitemThumb"]/span/img`, func(e *colly.XMLElement) {
		info.ThumbURL = e.Request.AbsoluteURL(e.Attr("src"))
		info.CoverURL = info.ThumbURL
	})

	// Preview Images
	c.OnXML(`//section[@class="items_article_SampleImages"]/ul/li`, func(e *colly.XMLElement) {
		info.PreviewImages = append(info.PreviewImages, e.Request.AbsoluteURL(e.ChildAttr(`.//a`, "href")))
	})

	// Preview Video
	//c.OnScraped(func(r *colly.Response) {
	//	d := c.Clone()
	//	d.OnResponse(func(r *colly.Response) {
	//		data := struct {
	//			Path string `json:"path"`
	//			Code int    `json:"code"`
	//		}{}
	//		if err := json.Unmarshal(r.Body, &data); err == nil && data.Code == http.StatusOK {
	//			info.PreviewVideoURL = data.Path
	//		}
	//	})
	//	d.Visit(fmt.Sprintf(sampleURL, info.ID))
	//})

	err = c.Visit(info.Homepage)
	return
}

func init() {
	provider.RegisterMovieFactory(Name, New)
}
