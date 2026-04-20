package hw03frequencyanalysis

import (
	"regexp"
	"slices"
	"sort"
	"strings"
)

func Top10(text string) []string {
	re := regexp.MustCompile("^[\"«']*[^\"'.,!?;:]*[.,!?;:\"»']+$")
	sre := regexp.MustCompile("[^\"«»'.,!?;:]+")
	sre.Longest()

	re2 := regexp.MustCompile("^.*[-]+$")
	sre2 := regexp.MustCompile("[^\"«»'.,!?;:]*[^\"»'.,!?;:-]+")
	sre2.Longest()

	var res []string

	m := make(map[string]int)

	rvals := []struct {
		num int
		val string
	}{}

	for part := range strings.FieldsSeq(text) {
		part = strings.ToLower(part)

		switch {
		case re.MatchString(part):
			part = sre.FindString(part)
			if part != "" {
				m[part]++
			}
		case re2.MatchString(part):
			part = sre2.FindString(part)
			if part != "" {
				m[part]++
			}
		default:
			m[part]++
		}
	}

	for key, num := range m {
		rvals = append(rvals,
			struct {
				num int
				val string
			}{num: num, val: key})

		sort.Slice(rvals, func(i, j int) bool {
			return (rvals[i].num > rvals[j].num) || ((rvals[i].num == rvals[j].num) && (rvals[i].val < rvals[j].val))
		})
		if len(rvals) > 10 {
			rvals = slices.Delete(rvals, 10, 11)
		}
	}

	for i := 0; i < len(rvals); i++ {
		res = append(res, rvals[i].val)
	}

	return res
}
