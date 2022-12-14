package honcrawler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Youngkingman/gentlemanSpider/settings"
)

var patternWinFile = regexp.MustCompile(`[/\:*?"<>|]`)

func genDirNameAndFilter(hd *HonDetail) (s string) {
	s = hd.Title + "["
	for _, v := range hd.Tags {
		s = s + "_" + v
	}
	s = s + "]"
	s = patternWinFile.ReplaceAllString(s, "")
	return
}

// `/photos-index-aid-169728.html`` =>`/photos-index-page-i-aid-169728.html`
func pageUrlTrans(u string, i int) string {
	strs := strings.Split(u, "-")
	ret := fmt.Sprintf("%s-%s-page-%d-%s-%s", strs[0], strs[1], i, strs[2], strs[3])
	return ret
}

// filter the wanted tags
func parseTages(tags []string) bool {
	if !settings.CrawlerSetting.EnableFilter {
		return true
	}
	for _, v := range tags {
		if settings.WantedTagsSet[v] {
			return true
		}
	}
	return false
}

// implementation of set
type set map[string]struct{}

func (s set) has(key string) bool {
	_, ok := s[key]
	return ok
}

func (s set) insert(key string) {
	s[key] = struct{}{}
}
