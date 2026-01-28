package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// CORS対応のミドルウェア
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORSヘッダーを追加
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// プリフライトリクエスト（OPTIONS）の処理
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message: "Hello World!",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// トークテーマを取得するハンドラ
func getThemesHandler(themeService *ThemeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// GETメソッドのみ許可
		if r.Method != http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: "Method not allowed",
			})
			return
		}

		// テーマを取得
		themes, err := themeService.GetRandomThemes()
		if err != nil {
			log.Printf("Error generating themes: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: "Failed to generate themes",
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(themes)
	}
}

func main() {
	// 設定を読み込み
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// AIクライアントとサービスを初期化
	aiClient := NewAIClient(config.OpenAIAPIKey, config.OpenAIAPIURL)
	themeService := NewThemeService(aiClient)

	// ルーティング設定
	http.HandleFunc("/", withCORS(helloHandler))
	http.HandleFunc("/api/themes/randoms", withCORS(getThemesHandler(themeService)))

	log.Printf("Server starting on port %s...", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
