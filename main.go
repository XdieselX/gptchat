package main

import (
	"fmt"
	"github.com/ian-kent/gptchat/module"
	"github.com/ian-kent/gptchat/module/memory"
	"github.com/ian-kent/gptchat/module/plugin"
	"github.com/ian-kent/gptchat/ui"
	openai "github.com/sashabaranov/go-openai"
	"os"
	"strconv"
	"strings"
)

var client *openai.Client

func init() {
	openaiAPIKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if openaiAPIKey == "" {
		ui.Warn("You haven't configured an OpenAI API key")
		fmt.Println()
		if !ui.PromptConfirm("Do you have an API key with access to the GPT-4 model?") {
			ui.Warn("You'll need an API key to use GPTChat")
			fmt.Println()
			fmt.Println("* You can get an API key at https://platform.openai.com/account/api-keys")
			fmt.Println("* You can get join the GPT-4 API waitlist at https://openai.com/waitlist/gpt-4-api")
			os.Exit(1)
		}

		openaiAPIKey = ui.PromptInput("Enter your API key:")
		if openaiAPIKey == "" {
			fmt.Println("")
			ui.Warn("You didn't enter an API key.")
			os.Exit(1)
		}
	}

	client = openai.NewClient(openaiAPIKey)

	module.Load(client, []module.Module{
		&memory.Module{},
		&plugin.Module{},
	}...)
	if err := module.LoadCompiledPlugins(); err != nil {
		fmt.Printf("error loading compiled plugins: %s", err)
		os.Exit(1)
	}
}

func main() {
	debugMode := false
	debugEnv := os.Getenv("GPT_DEBUG")
	if debugEnv != "" {
		v, err := strconv.ParseBool(debugEnv)
		if err != nil {
			ui.Warn(fmt.Sprintf("error parsing GPT_DEBUG: %s", err.Error()))
		} else {
			debugMode = v
		}
	}

	ui.Welcome(
		`Welcome to the GPT-4 client.`,
		`You can talk directly to GPT-4, or you can use /commands to interact with the client.

Use /help to see a list of available commands.`)

	chatLoop(debugMode)
}
