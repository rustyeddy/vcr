package redeye

import "testing"

func TestWebServer(t *testing.T) {
	addr := "1.2.3.4:5678"
	tpath := "/testpath"
	w := NewWebServer(addr, tpath)
	if (w.Addr != addr) {
		t.Errorf("Address expected (%s) got (%s)", addr, w.Addr)
	}

	if (w.Basepath != tpath) {
		t.Errorf("Basepath expected (%s) got (%s)", tpath, w.Basepath)
	}

	if (w.Handlers == nil) {
		t.Errorf("Handlers expected but nil found")
	}
}
