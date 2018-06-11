package css

import (
	"testing"
	"fmt"
)

func TestParseInline(t *testing.T) {
	cases := []struct {
		input  string
		result map[string]string
		hasErr bool
	}{
		{"", map[string]string{}, false},
		{"width: 5px", map[string]string{"width": "5px"}, false},
		{"width: 5px; height:		20px", map[string]string{"width": "5px", "height": "20px"}, false},
		{"width: 5px; height:				20px;", map[string]string{"width": "5px", "height": "20px"}, false},
		{`width: 5px;height:20px;`, map[string]string{"width": "5px", "height": "20px"}, false},
		{`background-image: url("http://image.com?w=1231&h=88s9")`, map[string]string{"background-image": `url("http://image.com?w=1231&h=88s9")`}, false},
		{`width: 5px;
	height:20px;`, map[string]string{"width": "5px", "height": "20px"}, false},
		{"width: ", nil, true},
		{":", nil, true},
	}

	for i, c := range cases {
		tn, c1 := fmt.Sprintf("%d", i), c
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			res, err := ParseInlineCss(c1.input)

			if c1.hasErr && err == nil {
				t.Fatal("expected error but finished without any")
			}

			if !c1.hasErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(res) != len(c1.result) {
				t.Fatalf("expected and actual results dont match each other by length. Expected len is %d. Actual len is %d", len(c1.result), len(res))
			}

			for k, v := range c1.result {
				if v != res[k] {
					t.Fatalf("expected value of property %q doesnt mactch with actual one. Expected %q. Got %q", k, v, res[k])
				}
			}
		})
	}
}
