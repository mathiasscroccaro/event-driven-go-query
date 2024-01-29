package use_case

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestNewGetDatasheetByUrlUseCase(t *testing.T) {
	expected := "Test PDF"

	pdfDocument := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, expected)
	}))
	defer pdfDocument.Close()

	pageContent := fmt.Sprintf(`
	<html>
		<body>
			<table>
				<tbody>
					<tr>
						<td>Datasheet</td>
						<td><a href='%s'>C90030</a></td>
					</tr>
				</tbody>
			</table>
		</body>
	</html>`, pdfDocument.URL)

	lcscPage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, pageContent)
	}))
	defer lcscPage.Close()

	datasheetReader, err := NewGetDatasheetByUrlUseCase(http.DefaultClient).Execute(lcscPage.URL)

	if err != nil {
		t.Errorf("Error while calling GetDatasheetByUrlUseCase: %v", err)
	}

	got, _ := ioutil.ReadAll(datasheetReader)

	if string(got) != expected {
		t.Errorf("Expected %s, got %s", expected, string(got))
	}
}

func TestFindDatasheetUrl(t *testing.T) {
	pageReader := strings.NewReader(`
	<html>
		<body>
			<table>
				<tbody>
					<tr>
						<td>Datasheet</td>
						<td><a href='https://example.com'>C90030</a></td>
					</tr>
				</tbody>
			</table>
		</body>
	</html>`)

	document, err := goquery.NewDocumentFromReader(pageReader)
	if err != nil {
		t.Error(err)
	}

	datasheetUrl, err := findDatasheetUrl(document)

	if err != nil {
		t.Error(err)
	}
	if datasheetUrl != "https://example.com" {
		t.Errorf("Expected https://example.com, got %s", datasheetUrl)
	}
}
