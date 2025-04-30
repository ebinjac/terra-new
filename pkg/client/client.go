package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Client handles HTTP operations
type Client struct {
	httpClient *http.Client
	token      string
}

// NewClient creates a new HTTP client with the given token
func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{},
		token:      token,
	}
}

// DownloadFile downloads a file from the given URL using Bearer token authentication
func (c *Client) DownloadFile(url, destPath string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error downloading file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}

// UploadFiles uploads the required files to the specified endpoint
func (c *Client) UploadFiles(endpoint string, req *UploadRequest) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add files to the multipart form
	files := map[string]string{
		"tf-plan":      req.TFPlan,
		"mapping-file": req.MappingFile,
		"tfgraph-file": req.TFGraphFile,
	}

	for field, path := range files {
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("error opening file %s: %v", path, err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile(field, filepath.Base(path))
		if err != nil {
			return fmt.Errorf("error creating form file: %v", err)
		}

		if _, err = io.Copy(part, file); err != nil {
			return fmt.Errorf("error copying file content: %v", err)
		}
	}

	// Add other form fields
	_ = writer.WriteField("product-id", req.ProductID)
	_ = writer.WriteField("name", req.Name)

	if err := writer.Close(); err != nil {
		return fmt.Errorf("error closing writer: %v", err)
	}

	req2, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req2.Header.Set("Content-Type", writer.FormDataContentType())
	req2.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req2)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
