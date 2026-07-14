package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

func (u *User) UnmarshalEasyJSON(in *jlexer.Lexer) {
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if key == "Email" {
			u.Email = in.String()
		} else {
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	suffix := "." + domain

	scanner := bufio.NewScanner(r)

	var user User
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		user.Email = ""
		if err := easyjson.Unmarshal(line, &user); err != nil {
			return nil, fmt.Errorf("unmarshal error: %w", err)
		}

		if user.Email == "" {
			continue
		}

		if strings.HasSuffix(user.Email, suffix) {
			idx := strings.IndexByte(user.Email, '@')
			if idx < 0 {
				continue
			}
			domainPart := strings.ToLower(user.Email[idx+1:])
			result[domainPart]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return result, nil
}
