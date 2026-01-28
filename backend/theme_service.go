package main

import (
	"math/rand"
	"time"
)

type ThemeService struct {
	aiClient *AIClient
}

func NewThemeService(aiClient *AIClient) *ThemeService {
	return &ThemeService{
		aiClient: aiClient,
	}
}

type ThemeResponse struct {
	Themes []string `json:"themes"`
	Count  int      `json:"count"`
}

func (s *ThemeService) GetRandomThemes() (*ThemeResponse, error) {
	// AIを使用してテーマを生成
	themes, err := s.aiClient.GenerateTechThemes()
	if err != nil {
		return nil, err
	}

	// テーマをシャッフルしてランダム性を高める
	shuffledThemes := s.shuffleThemes(themes)

	// 最大3つのテーマを返す
	if len(shuffledThemes) > 3 {
		shuffledThemes = shuffledThemes[:3]
	}

	return &ThemeResponse{
		Themes: shuffledThemes,
		Count:  len(shuffledThemes),
	}, nil
}

func (s *ThemeService) shuffleThemes(themes []string) []string {
	// 現在時刻をシードとして使用してランダム性を確保
	rand.Seed(time.Now().UnixNano())

	shuffled := make([]string, len(themes))
	copy(shuffled, themes)

	// Fisher-Yates shuffle algorithm
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled
}
