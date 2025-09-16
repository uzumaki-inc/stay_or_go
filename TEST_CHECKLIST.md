# Go テスト作成チェックリスト

対象: テスト未作成の `.go` ファイルに対する、実装前のチェック項目一覧。

## 対象ファイル

- analyzer/parameter_weights.go
- cmd/root.go
- main.go
- parser/go_parser.go
- parser/parser.go
- presenter/csv_presenter.go
- presenter/markdown_presenter.go
- presenter/tsv_presenter.go

## 共通タスク

- [ ] テーブル駆動テストで正常系/境界/エラー系を網羅
- [ ] 外部呼び出しは `jarcoal/httpmock` でモック
- [ ] 一時ファイル/ディレクトリを使う場合は `t.Cleanup` で片付け
- [ ] 出力は `bytes.Buffer` または `os.Pipe` で捕捉しアサート
- [ ] `go test ./...` がローカルで成功することを確認

## analyzer/parameter_weights.go

- [x] `NewParameterWeights` が定数どおりのデフォルト値を返す（`analyzer/parameter_weights_test.go`）
- [x] `NewParameterWeightsFromConfiFile` が一時 YAML の値を正しく読み込む
- [x] 異常系（存在しないパス）は `os.Exit(1)` となることをサブプロセスで検証
  - 例: `exec.Command(os.Args[0], "-test.run=TestHelperProcess_WeightsExit")`
  - 実行メモ: サンドボックス環境では Go モジュール取得/ビルドキャッシュ制約により `go test` 実行が不可。ローカルでは `GOCACHE` を任意ディレクトリに設定して実行:
    - 例: `GOCACHE=.gocache go test ./analyzer -v`

## cmd/root.go

- [x] 未サポート言語指定でエラーメッセージと終了コード 1 を出す（サブプロセス）
- [x] `-f` 未対応フォーマットでエラーを出す（サブプロセス）
- [x] `-g` 未指定時に `GITHUB_TOKEN` 未設定ならエラー（サブプロセス）
- [x] 言語ごとのデフォルト入力パスが設定（`go`→`go.mod` / `ruby`→`Gemfile`）
- [x] `--verbose` でデバッグ出力（`Selected Language: go` 等）が stderr に出る
  - 追加テスト: `cmd/root_test.go`
  - 実行例: `GOCACHE=.gocache go test ./cmd -v`

## main.go

- [ ] スモークテスト（実行してパニックせず起動する）
  - 必要に応じて `cmd.Execute` をスタブ/サブプロセスで検証

## parser/go_parser.go

- [x] `require (...)` ブロックから `// indirect` を除外して抽出できる
- [x] `replace (...)` の対象は `Skip` と `SkipReason` が設定される
- [x] `GetRepositoryURL` が `proxy.golang.org` のレスポンスから GitHub URL を設定する
- [x] GitHub 以外の場合は `Skip` と理由が設定される（HTTP モックで検証）
  - 追加テスト: `parser/go_parser_test.go`
  - 実行例: `GOCACHE=.gocache go test ./parser -v`

## parser/parser.go

- [x] `SelectParser("go")`/`("ruby")` が対応型を返す（`parser/parser_core_test.go`）
- [x] 未対応言語でエラーを返す（`errors.Is(err, parser.ErrUnsupportedLanguage)`）
- [x] `NewLibInfo` と各 `With...` オプションの反映を確認
  - 実行例: `GOCACHE=.gocache go test ./parser -run TestSelectParser -v`

## presenter/\*\_presenter.go

- [x] CSV/TSV/Markdown の `Display` が想定ヘッダ・ボディを出力する
  - 既存: `presenter/presenter_test.go`（通常ケース）
  - 追加: `presenter/presenters_missing_test.go`（Repo情報欠損 + Skip理由）
- [x] `makeBody` が区切り文字ごとに正しい整形を行い、欠損値を `N/A` とする（上記テストで検証）
- [x] `SelectPresenter` がフォーマットに応じて型を返す（`presenter/select_presenter_test.go`）
- 実行メモ: サンドボックスでは `go test` 実行が制限される場合があります。ローカルでは次を推奨:
  - `GOCACHE=.gocache go test ./presenter -v`

## parser/ruby_parser.go

- [x] `source`/`platforms`/`install_if` ブロック内の gem は Skip される
- [x] `GetRepositoryURL` で `homepage_uri`/`source_code_uri` が GitHub でない場合は Skip
  - 追加テスト: `parser/ruby_parser_edge_test.go`
  - 実行例: `GOCACHE=.gocache go test ./parser -v`

## 実行コマンド例

- `go test ./...`
- リント: `make lint`（自動修正は `make lintFix`）
