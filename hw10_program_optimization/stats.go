package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Email struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat, 1000)
	var elmail Email
	reader := bufio.NewReader(r)
	for i := 0; ; i++ {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return result, err
		}
		if err := json.Unmarshal(line, &elmail); err != nil {
			return nil, fmt.Errorf("line processing error: %w", err)
		}
		if strings.Contains(elmail.Email, "@") && strings.HasSuffix(elmail.Email, "."+domain) {
			r := strings.ToLower(strings.Split(elmail.Email, "@")[1])
			count := strings.Count(r, ".")
			if count > 1 {
				s := strings.Split(r, ".")
				n := len(s)
				r = s[n-2] + "." + s[n-1]
			}
			result[r]++
		}
	}
	return result, nil
}
