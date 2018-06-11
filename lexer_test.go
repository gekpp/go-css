package css

import (
	"testing"
	"fmt"
)

func TestLexProp(t *testing.T) {
	cases := []struct {
		input string
		item  item
	}{
		{"1234", item{itemError, ""}},
		{"", item{itemEOF, ""}},
		{":", item{itemError, ""}},
		{"abc", item{itemProp, "abc"}},
		{"abc:", item{itemProp, "abc"}},
		{"width : 50px;", item{itemProp, "width"}},
	}

	for i, c := range cases {
		j, c1 := fmt.Sprintf("%d", i), c
		t.Run(j, func(t *testing.T) {
			t.Parallel()
			l := lex(c1.input)
			l.state = lexProp
			nextItem := l.nextItem()
			if nextItem.typ != c1.item.typ {
				t.Fatalf("unexpected type of next item. Expected %v but got %v", c1.item.typ, nextItem.typ)
			}
			if nextItem.typ == itemError {
				return
			}

			if nextItem.val != c1.item.val {
				t.Fatalf("unexpected value of next item. Expected %v but got %v", c1.item.val, nextItem.val)
			}
		})
	}
}

func TestLexDelimiter(t *testing.T) {
	cases := []struct {
		input string
		item  item
	}{
		{"", item{itemError, ""}},
		{" ", item{itemError, ""}},
		{"  ", item{itemError, ""}},
		{"\t", item{itemError, "abc"}},
		{" \t", item{itemError, "abc"}},
		{"\t ", item{itemError, "abc"}},
		{"\n", item{itemError, "abc"}},
		{"\n ", item{itemError, "abc"}},
		{"   5px", item{itemError, ""}},
		{"   :\t\t 5px", item{itemValue, "5px"}},
		{"   :\t\t \"5px\"", item{itemValue, `"5px"`}},
		{"   :\t\t '5px'", item{itemValue, `'5px'`}},
	}

	for i, c := range cases {
		j, c1 := fmt.Sprintf("%d", i), c
		t.Run(j, func(t *testing.T) {
			t.Parallel()
			l := lex(c1.input)
			l.state = lexDelimiter
			nextItem := l.nextItem()
			if nextItem.typ != c1.item.typ {
				t.Fatalf("unexpected type of next item. Expected %v but got %v", c1.item.typ, nextItem.typ)
			}

			if nextItem.typ == itemError {
				return
			}

			if nextItem.val != c1.item.val {
				t.Fatalf("unexpected value of next item. Expected %v but got %v", c1.item.val, nextItem.val)
			}
		})
	}
}

func TestLexValue(t *testing.T) {
	cases := []struct {
		input string
		item  item
	}{
		{"", item{itemError, ""}},
		{" ", item{itemError, ""}},
		{"  ", item{itemError, ""}},
		{"\t", item{itemError, ""}},
		{" \t", item{itemError, ""}},
		{"\t ", item{itemError, ""}},
		{"\n", item{itemError, ""}},
		{"\n ", item{itemError, ""}},
		{"   5px", item{itemValue, "5px"}},
		{"   \t\t 5px", item{itemValue, "5px"}},
		{"   \t\t \"5px\";", item{itemValue, `"5px"`}},
		{`   'url(\'http://img.com\')`, item{itemError, ""}},
		{`   "url('http://img.com')`, item{itemError, ""}},
		{`   'url(\'http://img.com\')'`, item{itemValue, `'url(\'http://img.com\')'`}},
		{`		"url(\"http://img.com\")"`, item{itemValue, `"url(\"http://img.com\")"`}},
	}

	for i, c := range cases {
		j, c1 := fmt.Sprintf("%d", i), c
		t.Run(j, func(t *testing.T) {
			t.Parallel()
			l := lex(c1.input)
			l.state = lexValue
			nextItem := l.nextItem()
			if nextItem.typ != c1.item.typ {
				t.Fatalf("unexpected type of next item. Expected %v but got %v", c1.item.typ, nextItem.typ)
			}

			if nextItem.typ == itemError {
				return
			}

			if nextItem.val != c1.item.val {
				t.Fatalf("unexpected value of next item. Expected %v but got %v", c1.item.val, nextItem.val)
			}
		})
	}
}
