package ipfs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gobius/config"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// PinataIPFSClient implements the IPFSClient interface using the Pinata service.
type PinataIPFSClient struct {
	apiKey    string
	apiSecret string
	jwt       string
	baseURL   string
	client    *http.Client // HTTP client for Pinata API calls
	config    config.AppConfig
}

// PinataResponse matches the expected JSON structure from Pinata's API.
type PinataResponse struct {
	IpfsHash  string `json:"IpfsHash"`
	PinSize   int    `json:"PinSize"`
	Timestamp string `json:"Timestamp"`
}

// NewPinataIPFSClient creates a new client for interacting with Pinata.
func NewPinataIPFSClient(cfg config.AppConfig) (*PinataIPFSClient, error) {
	if cfg.IPFS.Pinata.JWT == "" {
		return nil, fmt.Errorf("pinata strategy selected but Pinata JWT is not configured")
	}

	// Configure a default HTTP client
	httpClient := &http.Client{
		Timeout: 60 * time.Second, // Default timeout, can be overridden by context
	}

	return &PinataIPFSClient{
		apiKey:    cfg.IPFS.Pinata.APIKey,
		apiSecret: cfg.IPFS.Pinata.APISecret,
		jwt:       cfg.IPFS.Pinata.JWT,
		baseURL:   cfg.IPFS.Pinata.BaseURL,
		client:    httpClient,
		config:    cfg,
	}, nil
}

// PinFileToIPFS pins a single file using Pinata.
func (pc *PinataIPFSClient) PinFileToIPFS(data []byte, filename string) (string, error) {
	cid, err := pc.pinFileViaAPI(context.Background(), pc.client, data, filename)
	if err != nil {
		return "", err
	}
	return cid, nil
}

// PinFilesToIPFS pins multiple files as a directory using Pinata.
func (pc *PinataIPFSClient) PinFilesToIPFS(ctx context.Context, taskid string, filesToAdd []IPFSFile) (string, error) {
	return pc.pinFilesViaAPI(ctx, pc.client, taskid, filesToAdd)
}

// pinFileViaAPI handles the multipart request for pinning a single file.
func (pc *PinataIPFSClient) pinFileViaAPI(ctx context.Context, client *http.Client, data []byte, filename string) (string, error) {
	// Create a new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a form file
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("error creating form file: %w", err)
	}

	// Write the file data
	_, err = part.Write(data)
	if err != nil {
		return "", fmt.Errorf("error writing file data: %w", err)
	}

	// Add pinata metadata
	metadata := map[string]string{
		"name": filename,
	}
	metadataBytes, _ := json.Marshal(metadata)
	writer.WriteField("pinataMetadata", string(metadataBytes))

	// Add pinata options
	options := map[string]int{"cidVersion": 0}
	optionsBytes, _ := json.Marshal(options)
	writer.WriteField("pinataOptions", string(optionsBytes))

	// Close the writer
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("error closing writer: %w", err)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "POST", pc.baseURL+"/pinning/pinFileToIPFS", body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+pc.jwt)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request using the client passed in (or pc.client)
	resp, err := client.Do(req) // Use the passed client
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Check status code AFTER reading body
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("pinata API returned non-OK status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Parse the response
	var pinataResp PinataResponse
	if err := json.Unmarshal(respBody, &pinataResp); err != nil {
		return "", fmt.Errorf("error parsing response body: %w, body: %s", err, string(respBody))
	}

	return pinataResp.IpfsHash, nil
}

// pinFilesViaAPI handles the multipart request for pinning multiple files.
func (pc *PinataIPFSClient) pinFilesViaAPI(ctx context.Context, client *http.Client, taskid string, filesToAdd []IPFSFile) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add each file
	for _, file := range filesToAdd {
		pinataFilePath := fmt.Sprintf("%s/%s", taskid, file.Name)
		part, err := writer.CreateFormFile("file", pinataFilePath)
		if err != nil {
			return "", fmt.Errorf("error creating form file for %s: %w", pinataFilePath, err)
		}

		// Check for nil buffer before copying
		if file.Buffer == nil {
			return "", fmt.Errorf("buffer for file %s is nil", pinataFilePath)
		}
		_, err = io.Copy(part, file.Buffer)
		if err != nil {
			return "", fmt.Errorf("error writing file data for %s: %w", pinataFilePath, err)
		}
	}
	// After loop, check if the writer has any parts added. If not (e.g., all buffers were nil), return error.
	if writer.Boundary() == "" || len(filesToAdd) == 0 /* Or check a flag set if any file was added */ {
		return "", fmt.Errorf("no valid files provided to pin to Pinata")
	}

	// Add pinata metadata (use taskid as the directory name)
	// metadata := map[string]string{
	// 	"name": taskid,
	// }
	// metadataBytes, _ := json.Marshal(metadata)
	// writer.WriteField("pinataMetadata", string(metadataBytes))

	// Add pinata options
	options := map[string]int{"cidVersion": 0}
	optionsBytes, _ := json.Marshal(options)
	writer.WriteField("pinataOptions", string(optionsBytes))

	// Close the writer
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("error closing writer: %w", err)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "POST", pc.baseURL+"/pinning/pinFileToIPFS", body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+pc.jwt)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request using the client passed in (or pc.client)
	resp, err := client.Do(req) // Use the passed client
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Check status code AFTER reading body
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("pinata API returned non-OK status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Parse the response
	var pinataResp PinataResponse
	if err := json.Unmarshal(respBody, &pinataResp); err != nil {
		return "", fmt.Errorf("error parsing response body: %w, body: %s", err, string(respBody))
	}

	return pinataResp.IpfsHash, nil
}
