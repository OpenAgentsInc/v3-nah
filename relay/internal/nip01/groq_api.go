package nip01

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

const GroqAPIURL = "https://api.groq.com/openai/v1/audio/transcriptions"

func TranscribeAudio(audioData, format string) (string, error) {
	// Decode base64 audio data
	decodedAudio, err := base64.StdEncoding.DecodeString(audioData)
	if err != nil {
		return "", fmt.Errorf("failed to decode audio data: %v", err)
	}

	// Create a buffer to write our multipart form
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add the audio file
	part, err := writer.CreateFormFile("file", "audio."+format)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = io.Copy(part, bytes.NewReader(decodedAudio))
	if err != nil {
		return "", fmt.Errorf("failed to copy audio data: %v", err)
	}

	// Add other form fields
	writer.WriteField("model", "distil-whisper-large-v3-en")
	writer.WriteField("temperature", "0")
	writer.WriteField("response_format", "json")
	writer.WriteField("language", "en")

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %v", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", GroqAPIURL, &body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GROQ_API_KEY"))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// Parse the JSON response
	var result struct {
		Text string `json:"text"`
	}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	return result.Text, nil
}