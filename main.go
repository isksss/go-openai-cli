package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"github.com/joho/godotenv"
)

var writer *bufio.Writer

// .envファイルからAPI_KEYを取得する
func getApiKey() string {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	API_KEY := os.Getenv("OPENAI_API_KEY")

	if API_KEY == "" {
		panic("OPENAI_API_KEY is NULL.")
	}

	// fmt.Println(API_KEY)
	return API_KEY
}

func printw(s string) {
	fmt.Println(s)
	fmt.Fprintln(writer, s)
	writer.Flush()
}

func getTime() string {
	now := time.Now()
	nowStr := now.Format("2006-01-02-15:04:05")

	return nowStr
}

func getFileName() string {

	// Get user's home directory
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	appDir := filepath.Join(usr.HomeDir, ".config", "go-openai-cli", "dist")
	err = os.MkdirAll(appDir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	filename := appDir + "/text-" + getTime() + ".txt"

	return filename
}

func main() {
	filename := getFileName()
	file, err := os.Create(filename)

	if err != nil {
		panic(err)
	}

	defer file.Close()
	writer = bufio.NewWriter(file)

	printw("####################################")
	printw(getTime())

	API_KEY := getApiKey()

	for {
		printw("#--------------------")
		printw(getTime())

		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("You: ")
		fmt.Fprintln(writer, "You: ")
		scanner.Scan()
		input := scanner.Text()
		fmt.Fprintln(writer, input)
		writer.Flush()

		if input == "exit" {
			break
		}

		client := openai.NewClient(API_KEY)

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: input,
					},
				},
			},
		)

		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			fmt.Fprintf(writer, "ChatCompletion error: %v\n", err)
			writer.Flush()
			return
		}

		fmt.Printf("AI: %v\n", resp.Choices[0].Message.Content)
		fmt.Fprintf(writer, "AI: %v\n", resp.Choices[0].Message.Content)
		writer.Flush()

	}
}
