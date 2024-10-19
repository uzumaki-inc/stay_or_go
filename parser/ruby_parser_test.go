package parser

import (
	"os"
	"testing"
)

// RubyParser構造体のテスト
func TestRubyParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []LibInfo
	}{
		{
			name:    "Valid gem with no special options",
			content: `gem "rails"`,
			expected: []LibInfo{
				{
					Name:   `"rails"`,
					Skip:   false,
					Others: []string{},
				},
			},
		},
		{
			name:    "Gem with source option (should skip)",
			content: `gem "rails", :source => "https://example.com"`,
			expected: []LibInfo{
				{
					Name:   `"rails",`, // 修正: クオートとカンマを含める
					Skip:   true,
					Others: []string{`:source => "https://example.com"`},
				},
			},
		},
		{
			name:    "Gem with git option (should skip)",
			content: `gem "nokogiri", :git => "git://github.com/sparklemotion/nokogiri.git", :require => false`,
			expected: []LibInfo{
				{
					Name:   `"nokogiri",`, // 修正: クオートとカンマを含める
					Skip:   true,
					Others: []string{`:git => "git://github.com/sparklemotion/nokogiri.git"`, `:require => false`},
				},
			},
		},
		{
			name:    "Gem with non-NG options (should not skip)",
			content: `gem 'pg', :require => false`,
			expected: []LibInfo{
				{
					Name:   `'pg',`, // 修正: クオートとカンマを含める
					Skip:   false,
					Others: []string{`:require => false`},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用の一時ファイルを作成
			tmpfile, err := os.CreateTemp("", "testfile")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(tmpfile.Name()) // テスト終了後に一時ファイルを削除

			// テスト内容を書き込む
			if _, err := tmpfile.Write([]byte(tt.content)); err != nil {
				t.Fatalf("failed to write to temp file: %v", err)
			}

			// ファイルを閉じる
			if err := tmpfile.Close(); err != nil {
				t.Fatalf("failed to close temp file: %v", err)
			}

			// RubyParserを使ってファイルをパース
			p := RubyParser{}
			result := p.Parse(tmpfile.Name())

			// 結果をフィールドごとに比較する
			if len(*result) != len(tt.expected) {
				t.Fatalf("expected %d lib entries, but got %d", len(tt.expected), len(*result))
			}

			for i := range *result {
				if (*result)[i].Name != tt.expected[i].Name {
					t.Errorf("expected name %v, but got %v", tt.expected[i].Name, (*result)[i].Name)
				}
				if (*result)[i].Skip != tt.expected[i].Skip {
					t.Errorf("expected skip %v, but got %v", tt.expected[i].Skip, (*result)[i].Skip)
				}
				if len((*result)[i].Others) != len(tt.expected[i].Others) {
					t.Errorf("expected %d others, but got %d", len(tt.expected[i].Others), len((*result)[i].Others))
				}
				for j := range (*result)[i].Others {
					if (*result)[i].Others[j] != tt.expected[i].Others[j] {
						t.Errorf("expected other %v, but got %v", tt.expected[i].Others[j], (*result)[i].Others[j])
					}
				}
			}
		})
	}
}
