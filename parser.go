package css

import "fmt"

func ParseInlineCss(input string) (map[string]string, error) {
	l := lex(input)
	res := make(map[string]string)
	for {
		item := l.nextItem()
		if item.typ == itemEOF {
			break
		}

		if item.typ == itemError {
			return nil, fmt.Errorf("unexpected error at position %d: %v", l.pos, item)
		}
		prop := item.val

		item = l.nextItem()
		if item.typ == itemError {
			return nil, fmt.Errorf("unexpected error at position %d: %v", l.pos, item)
		}
		res[prop] = item.val
	}
	return res, nil
}
