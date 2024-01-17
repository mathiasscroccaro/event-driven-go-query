package use_case

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type GetDatasheetByUrlUseCase struct {
	httpClient *http.Client
}

func NewGetDatasheetByUrlUseCase(client *http.Client) *GetDatasheetByUrlUseCase {
	return &GetDatasheetByUrlUseCase{
		httpClient: client,
	}
}

func fetchPage(client *http.Client, url string) (io.Reader, error) {
	response, error := client.Get(url)
	if error != nil {
		return nil, error
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(content), nil
}

func findDatasheetUrl(doc *goquery.Document) (string, error) {
	selector := fmt.Sprintf("td:contains('%s')", "Datasheet")
	datasheetTableRow := doc.Find(selector)

	if datasheetTableRow.Length() > 0 {
		linkToDatasheetTableRow := datasheetTableRow.Next()
		targetDownloadLink, exists := linkToDatasheetTableRow.Find("a").Attr("href")

		if exists {
			return targetDownloadLink, nil
		}
	}

	return "", errors.New("target download link not found")
}

func downloadDatasheet(client *http.Client, url string) (io.Reader, error) {
	return fetchPage(client, url)
}

func (use_case *GetDatasheetByUrlUseCase) Execute(url string) (io.Reader, error) {
	pageReader, err := fetchPage(use_case.httpClient, url)
	if err != nil {
		return nil, err
	}

	document, err := goquery.NewDocumentFromReader(pageReader)
	if err != nil {
		return nil, err
	}
	datasheetUrl, err := findDatasheetUrl(document)

	datasheetReader, err := downloadDatasheet(use_case.httpClient, datasheetUrl)
	if err != nil {
		return nil, err
	}

	return datasheetReader, nil
}

type WriteToFileUseCase struct {
	reader io.Reader
}

func NewWriteToFileUseCase(reader io.Reader) *WriteToFileUseCase {
	return &WriteToFileUseCase{
		reader: reader,
	}
}

func (use_case *WriteToFileUseCase) Execute(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Error while creating the file: %v", err)
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, use_case.reader)
	if err != nil {
		log.Printf("Error while writing to the file: %v", err)
		return err
	}

	return nil
}
