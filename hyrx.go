package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/skratchdot/open-golang/open"
)

const addJs = `
<script>
window.lock = false;
var xy = document.getElementById('pow_xs_XieyiCheck');
if (xy) {
  xy.checked = true;
  window.touzi();
} else if (jQuery) {
  console.log('no xy found');
  var cjcheck = window.setInterval(function () {
    var cjbtn = $('.Dlb_conList_right.Dlb_cj').filter(function () { return $(this).find('input').val() >= 10000; }).find('.Dlb_cjBtn').first();
    if (cjbtn.length) {
      window.clearInterval(cjcheck);
      cjbtn.click();
      var m14check = window.setInterval(function () {
        var m14 = $('#msgDiv14:visible .pow_xs_Sub');
        if (m14.length) {
          window.clearInterval(m14check);
          m14.click();
        }
      }, 200);
      var m10check = window.setInterval(function () {
        var m10 = $('#msgDiv10:visible .pow_xs_Sub');
        if (m10.length) {
          window.clearInterval(m10check);
          m10.click();
        }
      }, 200);
    }
  }, 500);
}
</script>
`

func proxy(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	body, _ := ioutil.ReadAll(r.Body)
	println(r.Method, "https://www.hengyirong.com"+r.URL.String())
	req, _ := http.NewRequest(r.Method, "https://www.hengyirong.com"+r.URL.String(), bytes.NewReader(body))
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	resp, _ := client.Do(req)
	for key, values := range resp.Header {
		for _, value := range values {
			if key == "Location" {
				url, _ := resp.Location()
				url.Scheme = r.URL.Scheme
				url.Host = r.URL.Host
				value = url.String()
			}
			if key == "Set-Cookie" {
				value = strings.Replace(value, "domain=.hengyirong.com;", "", 1)
			}
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		scanner := bufio.NewScanner(resp.Body)
		write := true
		re := regexp.MustCompile(`https?:\/\/www\.hengyirong\.com\/?`)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "soperson.com") {
				write = false
			} else if strings.Contains(line, "</body>") {
				w.Write([]byte(addJs))
				write = true
			} else if strings.Contains(line, "app_licai.png") {
				continue
			}
			if write {
				line = re.ReplaceAllString(line, "/")
				line += "\n"
				w.Write([]byte(line))
			}
		}
	} else {
		io.Copy(w, resp.Body)
	}
}

func main() {
	http.HandleFunc("/", proxy)
	go func() {
		println("HYRX 1.0 - Created By CGH")
		open.Run("http://127.0.0.1:8888/")
	}()
	http.ListenAndServe(":8888", nil)
}
