/*
Copyright © 2024 @konyu(github name)
*/
package main

import (
	"github.com/joho/godotenv"
	"github.com/konyu/StayOrGo/cmd"
)

func main() {
	// .envファイルを読み込む
	_ = godotenv.Load()

	cmd.Execute()
}
