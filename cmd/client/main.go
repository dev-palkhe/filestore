package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	serverAddress string
)

var rootCmd = &cobra.Command{
	Use:   "store",
	Short: "File store client",
	Long:  "A CLI client for the file store service",
}

var addCmd = &cobra.Command{
	Use:   "add [filename...]",
	Short: "Add files to the store",
	Args:  cobra.MinimumNArgs(1),
	RunE:  addFiles,
}

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List files in the store",
	RunE:  listFiles,
}

var removeCmd = &cobra.Command{
	Use:   "rm [filename]",
	Short: "Remove a file from the store",
	Args:  cobra.ExactArgs(1),
	RunE:  removeFile,
}

var updateCmd = &cobra.Command{
	Use:   "update [filename]",
	Short: "Update a file in the store",
	Args:  cobra.ExactArgs(1),
	RunE:  updateFile,
}

var wcCmd = &cobra.Command{
	Use:   "wc",
	Short: "Word count of all files",
	RunE:  wordCount,
}

var freqWordsCmd = &cobra.Command{
	Use:   "freq-words",
	Short: "Frequent words in all files",
	RunE:  freqWords,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&serverAddress, "server", "s", "http://localhost:8000", "Server address")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(wcCmd)
	rootCmd.AddCommand(freqWordsCmd) // Double, triple check this line

	freqWordsCmd.Flags().IntP("limit", "n", 10, "Number of frequent words to show")
	freqWordsCmd.Flags().String("order", "dsc", "Order of frequent words (asc or dsc)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func addFiles(cmd *cobra.Command, args []string) error {
	for _, filename := range args {
		content, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filename, err)
		}
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("filename", filename)
		writer.WriteField("content", string(content))
		writer.Close()

		contentType := writer.FormDataContentType()
		req, err := http.NewRequest("POST", serverAddress+"/add", body)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", contentType)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to add file %s: %w", filename, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to add file %s: %s", filename, string(bodyBytes))
		}

		fmt.Printf("Added file: %s\n", filename)
	}
	return nil
}

func listFiles(cmd *cobra.Command, args []string) error {
	resp, err := http.Get(serverAddress + "/list")
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}
	defer resp.Body.Close()

	var files []string
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return fmt.Errorf("failed to decode list response: %w", err)
	}

	fmt.Println("Files:")
	for _, file := range files {
		fmt.Println("-", file)
	}
	return nil
}

func removeFile(cmd *cobra.Command, args []string) error {
	filename := args[0]
	resp, err := http.Post(serverAddress+"/remove", "application/x-www-form-urlencoded", bytes.NewBufferString("filename="+filename))
	if err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to remove file %s: %s", filename, string(bodyBytes))

	}
	fmt.Printf("Removed file: %s\n", filename)
	return nil

}

func updateFile(cmd *cobra.Command, args []string) error {
	filename := args[0]
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	resp, err := http.Post(serverAddress+"/update", "application/x-www-form-urlencoded", bytes.NewBufferString("filename="+filename+"&content="+string(content)))
	if err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update file %s: %s", filename, string(bodyBytes))

	}
	fmt.Printf("Updated file: %s\n", filename)
	return nil
}

func wordCount(cmd *cobra.Command, args []string) error {
	resp, err := http.Get(serverAddress + "/wc")
	if err != nil {
		return fmt.Errorf("failed to get word count: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Word count:", string(body))
	return nil
}

func freqWords(cmd *cobra.Command, args []string) error {
	limit, _ := cmd.Flags().GetInt("limit")
	order, _ := cmd.Flags().GetString("order")

	req, err := http.NewRequest("GET", serverAddress+"/freq-words", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	if limit != 0 {
		q.Add("limit", strconv.Itoa(limit))
	}
	if order != "" {
		q.Add("order", order)
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("failed to get frequent words: %w", err)
	}
	defer resp.Body.Close()

	var freqWords map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&freqWords); err != nil {
		return fmt.Errorf("failed to decode frequent words response: %w", err)
	}

	fmt.Println("Frequent words:")
	for word, count := range freqWords {
		fmt.Printf("%s: %d\n", word, count)
	}
	return nil
}
