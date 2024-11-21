/*
Copyright © 2024 @konyu
*/
package main

import (
	"github.com/joho/godotenv"
	"github.com/uzumaki-inc/stay_or_go/cmd"
)

func main() {
	// .envファイルを読み込む
	_ = godotenv.Load()

	cmd.Execute()
}
