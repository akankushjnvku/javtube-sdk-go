package caribbeancom

import (
	"regexp"

	"github.com/javtube/javtube-sdk-go/provider"
	"github.com/javtube/javtube-sdk-go/provider/caribbeancom/core"
)

var _ provider.MovieProvider = (*Caribbeancom)(nil)

const (
	Name     = "Caribbeancom"
	Priority = 1000
)

const (
	baseURL  = "https://www.caribbeancom.com/"
	movieURL = "https://www.caribbeancom.com/moviepages/%s/index.html"
)

type Caribbeancom struct {
	*core.Core
}

func New() *Caribbeancom {
	return &Caribbeancom{
		Core: (&core.Core{
			BaseURL:         baseURL,
			MovieURL:        movieURL,
			DefaultName:     Name,
			DefaultPriority: Priority,
			DefaultMaker:    "カリビアンコム",
		}).Init(),
	}
}

func (carib *Caribbeancom) NormalizeID(id string) string {
	if regexp.MustCompile(`^\d{6}-\d{3}$`).MatchString(id) {
		return id
	}
	return ""
}

func init() {
	provider.RegisterMovieFactory(Name, New)
}
