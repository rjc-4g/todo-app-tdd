package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []Choice     `json:"choices"`
	Error   *OpenAIError `json:"error,omitempty"`
}

type Choice struct {
	Message Message `json:"message"`
}

type OpenAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type AIClient struct {
	apiKey string
	apiURL string
	client *http.Client
}

func NewAIClient(apiKey, apiURL string) *AIClient {
	return &AIClient{
		apiKey: apiKey,
		apiURL: apiURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *AIClient) GenerateTechThemes() ([]string, error) {
	if c.apiKey == "" {
		// APIキーが設定されていない場合は、フォールバック用のテーマを返す
		fmt.Println("OpenAI API key not set, using fallback themes")
		return c.getFallbackThemes(), nil
	}

	prompt := `技術に関するトークテーマを3つ生成してください。以下の条件を満たしてください：
- プログラミング、開発、IT技術に関連する内容
- トークに適した具体的なテーマ
- 日本語で回答
- JSON配列形式で返す（例：["テーマ1", "テーマ2", "テーマ3"]）

例：
- 最近学んだ新しいプログラミング言語の特徴
- 開発環境の効率化テクニック
- コードレビューのベストプラクティス`

	request := OpenAIRequest{
		Model: "gpt-4o-mini",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   500,
		Temperature: 0.8,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Printf("API request failed, using fallback themes: %v\n", err)
		return c.getFallbackThemes(), nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read API response, using fallback themes: %v\n", err)
		return c.getFallbackThemes(), nil
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API returned error status %d, using fallback themes: %s\n", resp.StatusCode, string(body))
		return c.getFallbackThemes(), nil
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		fmt.Printf("Failed to parse API response, using fallback themes: %v\n", err)
		return c.getFallbackThemes(), nil
	}

	if openAIResp.Error != nil {
		fmt.Printf("OpenAI API error, using fallback themes: %s\n", openAIResp.Error.Message)
		return c.getFallbackThemes(), nil
	}

	if len(openAIResp.Choices) == 0 {
		fmt.Println("No choices returned from OpenAI API, using fallback themes")
		return c.getFallbackThemes(), nil
	}

	// AIの回答をパースしてテーマの配列に変換
	content := openAIResp.Choices[0].Message.Content
	var themes []string
	if err := json.Unmarshal([]byte(content), &themes); err != nil {
		// JSONパースに失敗した場合は、フォールバックテーマを返す
		return c.getFallbackThemes(), nil
	}

	// テーマが3つでない場合は調整
	if len(themes) > 3 {
		themes = themes[:3]
	} else if len(themes) < 3 {
		// 不足分をフォールバックテーマで補完
		fallback := c.getFallbackThemes()
		for len(themes) < 3 && len(fallback) > 0 {
			themes = append(themes, fallback[0])
			fallback = fallback[1:]
		}
	}

	return themes, nil
}

func (c *AIClient) getFallbackThemes() []string {
	allThemes := []string{
		"最近学んだ技術やツール",
		"開発で苦労したバグの話",
		"おすすめのVSCode拡張機能",
		"好きなプログラミング言語とその理由",
		"キーボードやマウスのこだわり",
		"リモートワークの工夫",
		"コードレビューで気をつけていること",
		"最近読んだ技術書",
		"AIツールの活用方法",
		"副業やOSSの話",
		"テストコードの書き方",
		"アーキテクチャの設計思想",
		"チーム開発のTips",
		"デバッグの極意",
		"パフォーマンス改善の経験",
	}

	// ランダムに3つ選択
	selectedThemes := c.selectRandomThemes(allThemes, 3)
	fmt.Printf("Generated fallback themes: %v\n", selectedThemes)
	return selectedThemes
}

func (c *AIClient) selectRandomThemes(themes []string, count int) []string {
	if count >= len(themes) {
		return themes
	}

	// Fisher-Yates shuffle アルゴリズムを使用
	shuffled := make([]string, len(themes))
	copy(shuffled, themes)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := time.Now().UnixNano() % int64(i+1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		// 少し待つことでランダム性を高める
		time.Sleep(1 * time.Nanosecond)
	}

	return shuffled[:count]
}
