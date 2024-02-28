package kpprompt

import (
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

type HistoryMessage struct {
	Type    schema.ChatMessageType
	Content string
}

type PromptTemplate struct {
	SystemTemplate string
	HistoryMessage []HistoryMessage
	UserTemplate   string
	InputKey       string
	OutPutKey      string
	Functions      []llms.FunctionDefinition
}

func (p *PromptTemplate) FormatPrompt() ([]llms.MessageContent, error) {
	var messages = make([]llms.MessageContent, 0)
	systemMessage := llms.MessageContent{
		Role:  schema.ChatMessageTypeSystem,
		Parts: []llms.ContentPart{llms.TextContent{Text: p.ParseSystemMessage(p.SystemTemplate)}},
	}
	messages = append(messages, systemMessage)

	if len(p.HistoryMessage) > 0 {
		for _, value := range p.HistoryMessage {
			historyMessage := llms.MessageContent{
				Role:  value.Type,
				Parts: []llms.ContentPart{llms.TextContent{Text: value.Content}},
			}
			messages = append(messages, historyMessage)
		}
	}

	if p.UserTemplate != "" {
		userMessage := llms.MessageContent{
			Role:  schema.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextContent{Text: p.ParseHumanMessage(p.UserTemplate)}},
		}
		messages = append(messages, userMessage)
	}
	return messages, nil
}

func (p *PromptTemplate) ParseResult(res string) any {
	return nil
}

func (p *PromptTemplate) ParseSystemMessage(sysMsg string) string {
	return sysMsg
}

func (p *PromptTemplate) ParseHumanMessage(humanMsg string) string {
	return humanMsg
}

func createAgentPrompt(system, user string, history []HistoryMessage) PromptTemplate {
	return PromptTemplate{
		SystemTemplate: system,
		UserTemplate:   user,
		HistoryMessage: history,
	}
}
