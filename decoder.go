package vcard

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// A Decoder parses cards.
type Decoder struct {
	s *bufio.Scanner
}

// NewDecoder creates a new Decoder reading cards from an io.Reader.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{s: bufio.NewScanner(r)}
}

// Decode parses a single card.
func (dec *Decoder) Decode() (Card, error) {
	card := make(Card)

	hasHeader := false
	for dec.s.Scan() {
		l := strings.TrimSpace(dec.s.Text())
		if l == "" {
			continue
		}

		k, v, params := parseLine(l)

		if !hasHeader {
			if k == "BEGIN" {
				if strings.ToUpper(v) != "VCARD" {
					return card, errors.New("vcard: invalid BEGIN value")
				}
				hasHeader = true
				continue
			} else {
				return card, errors.New("vcard: no BEGIN field found")
			}
		} else if k == "END" {
			if strings.ToUpper(v) != "VCARD" {
				return card, errors.New("vcard: invalid END value")
			}
			break
		}

		card[k] = append(card[k], &Field{
			Value: v,
			Params: params,
		})
	}

	return card, dec.s.Err()
}

func parseLine(l string) (k, v string, params map[string]string) {
	kv := strings.SplitN(l, ":", 2)
	if len(kv) < 2 {
		return l, "", nil
	}

	v = kv[1]
	kparams := strings.Split(kv[0], ";")
	k = strings.ToUpper(kparams[0])

	if len(kparams) > 1 {
		params = make(map[string]string)
		for i := 1; i < len(kparams); i++ {
			pk, pv := parseParam(kparams[i])
			params[pk] = pv
		}
	}

	return
}

func parseParam(s string) (k, v string) {
	kv := strings.SplitN(s, "=", 2)
	if len(kv) < 2 {
		return s, ""
	}
	return strings.ToUpper(kv[0]), kv[1]
}
