package ipfs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type PinataClient struct {
	apiKey    string
	apiSecret string
	jwt       string
	baseURL   string
}

type PinataResponse struct {
	IpfsHash  string `json:"IpfsHash"`
	PinSize   int    `json:"PinSize"`
	Timestamp string `json:"Timestamp"`
}

func NewPinataClient(apiKey, apiSecret, jwt string) *PinataClient {
	return &PinataClient{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		jwt:       jwt,
		baseURL:   "https://arbius.mypinata.cloud",
	}
}

func (pc *PinataClient) PinFileToIPFS(data []byte, filename string) string {
	ctx := context.Background()
	
	// Create a new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// Create a form file
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		fmt.Printf("Error creating form file: %v\n", err)
		return ""
	}
	
	// Write the file data
	_, err = part.Write(data)
	if err != nil {
		fmt.Printf("Error writing file data: %v\n", err)
		return ""
	}
	
	// Add pinata metadata
	metadata := map[string]string{
		"name": filename,
	}
	metadataBytes, _ := json.Marshal(metadata)
	writer.WriteField("pinataMetadata", string(metadataBytes))
	
	// Close the writer
	err = writer.Close()
	if err != nil {
		fmt.Printf("Error closing writer: %v\n", err)
		return ""
	}
	
	// Create the request
	req, err := http.NewRequestWithContext(ctx, "POST", pc.baseURL+"/pinning/pinFileToIPFS", body)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return ""
	}
	
	// Set headers
	req.Header.Set("Authorization", "Bearer "+pc.jwt)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return ""
	}
	defer resp.Body.Close()
	
	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return ""
	}
	
	// Parse the response
	var pinataResp PinataResponse
	err = json.Unmarshal(respBody, &pinataResp)
	if err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return ""
	}
	
	return pinataResp.IpfsHash
}

func (pc *PinataClient) PinFilesToIPFS(ctx context.Context, taskid string, filesToAdd []IPFSFile) (string, error) {
	// Create a new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	// Add each file
	for _, file := range filesToAdd {
		part, err := writer.CreateFormFile("file", file.Name)
		if err != nil {
			return "", fmt.Errorf("error creating form file: %v", err)
		}
		
		_, err = io.Copy(part, file.Buffer)
		if err != nil {
			return "", fmt.Errorf("error writing file data: %v", err)
		}
	}
	
	// Add pinata metadata
	metadata := map[string]string{
		"name": taskid,
	}
	metadataBytes, _ := json.Marshal(metadata)
	writer.WriteField("pinataMetadata", string(metadataBytes))
	
	// Close the writer
	err := writer.Close()
	if err != nil {
		return "", fmt.Errorf("error closing writer: %v", err)
	}
	
	// Create the request
	req, err := http.NewRequestWithContext(ctx, "POST", pc.baseURL+"/pinning/pinFileToIPFS", body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	
	// Set headers
	req.Header.Set("Authorization", "Bearer "+pc.jwt)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()
	
	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}
	
	// Parse the response
	var pinataResp PinataResponse
	err = json.Unmarshal(respBody, &pinataResp)
	if err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}
	
	return pinataResp.IpfsHash, nil
} 