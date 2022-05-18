package aventertainments

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAVE_GetMovieInfoByID(t *testing.T) {
	provider := New()
	for _, item := range []string{
		"9865",
		"10161",
		"140930",
		"115855",
	} {
		info, err := provider.GetMovieInfoByID(item)
		data, _ := json.MarshalIndent(info, "", "\t")
		assert.True(t, assert.NoError(t, err) && assert.True(t, info.Valid()))
		t.Logf("%s", data)
	}
}

func TestAVE_SearchMovie(t *testing.T) {
	provider := New()
	for _, item := range []string{
		"lldv-12",
		"mcbd-25",
	} {
		results, err := provider.SearchMovie(provider.TidyKeyword(item))
		data, _ := json.MarshalIndent(results, "", "\t")
		if assert.NoError(t, err) {
			for _, result := range results {
				assert.True(t, result.Valid())
			}
		}
		t.Logf("%s", data)
	}
}