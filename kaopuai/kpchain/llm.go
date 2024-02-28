package kpchain

import (
	"context"
	"errors"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/kaopuai/kpprompt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

type LLMChain struct {
	Prompt           kpprompt.PromptTemplate
	LLM              llms.Model
	CallbacksHandler callbacks.Handler
	Options          *chainCallOption

	OutputKey string
}

func NewLLMChain(llm llms.Model, prompt kpprompt.PromptTemplate, opts ...ChainCallOption) *LLMChain {
	opt := &chainCallOption{}
	for _, o := range opts {
		o(opt)
	}
	chain := &LLMChain{
		Prompt:  prompt,
		LLM:     llm,
		Options: opt,
	}

	return chain
}

func (l *LLMChain) Call(ctx context.Context, options ...ChainCallOption) (any, error) {
	promptValue, err := l.Prompt.FormatPrompt()
	if err != nil {
		return "", err
	}
	if len(l.Prompt.Functions) > 0 {
		options = append(options, WithFunctions(l.Prompt.Functions))
	}

	res, err := GenerateFromSinglePrompt(ctx, l.LLM, promptValue, options...)
	if err != nil {
		return "", err
	}

	value := l.Prompt.ParseResult(res)
	return value, nil
}

func (l *LLMChain) GetMemory() schema.Memory {
	return nil
}

func (l *LLMChain) GetInputKeys() []string {
	return nil
}

func (l *LLMChain) GetOutputKeys() []string {
	return nil
}

func GenerateFromSinglePrompt(ctx context.Context, llm llms.Model, message []llms.MessageContent, options ...ChainCallOption) (string, error) {
	resp, err := llm.GenerateContent(ctx, message, getLLMCallOptions(options...)...)
	if err != nil {
		return "", err
	}

	choices := resp.Choices
	if len(choices) < 1 {
		return "", errors.New("empty response from model")
	}
	c1 := choices[0]
	return c1.Content, nil
}
