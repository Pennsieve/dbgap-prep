package analytetype

import "testing"

func TestRoundTrip(t *testing.T) {
	for _, g := range []Type{DNA, RNA, DNARNA} {
		if got, err := FromString(g.String()); err != nil || got != g {
			t.Errorf("%v: round-trip got %v, err %v", g, got, err)
		}
	}
}
