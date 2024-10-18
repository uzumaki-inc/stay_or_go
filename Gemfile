source 'https://rubygems.org'

# Rubyのバージョンを指定する
ruby '3.2.2'

# Railsフレームワーク
gem 'rails', '~> 7.0.0'

gem 'webmock', '~> 7.0.0', :require => false
gem "nokogiri", require: true
gem "rspec", :group => :test
gem "weakling",   :platforms => :jruby
gem "some_internal_gem", :source => "https://gems.example.com"
gem "rails_git", "2.3.8", :git => "git://github.com/rails/rails.git"
gem "rails_github", :github => "rails/rails"
gem "nokogiri_version", ">= 1.4.2"
gem "RedCloth_version", ">= 4.1.0", "< 4.2.0"

# データベースを使うためのGem（デフォルトではSQLite）
gem 'sqlite3'
gem
# Webサーバー
gem 'puma', '~> 6.0'

# Sassコンパイラ（CSSプリプロセッサ）
gem 'sassc-rails', '>= 2.1.0'

# JavaScriptコンパイラ
gem 'jsbundling-rails'

# コーヒースクリプトのサポート
gem 'coffee-rails', '~> 5.0'

# Turbolinksをサポート（ページ遷移を高速化）
gem 'turbolinks', '~> 5'

# ページレイアウト用
gem 'jbuilder', '~> 2.7'

# テスト用
group :development, :test do
  # デバッグ用
  gem 'byebug', platforms: [:mri, :mingw, :x64_mingw]

  # RSpecを使ったテスト
  gem 'rspec-rails', '~> 6.0.0'
end

# 本番環境用
group :production do
  # PostgreSQLを使う場合の設定
  gem 'pg', '~> 1.4'
end

# 開発環境用
group :development do
  # Railsコンソールでデバッグツールを使うためのgem
  gem 'web-console', '>= 4.1.0'

  # ファイル変更時に自動的にサーバーを再起動する
  gem 'listen', '~> 3.3'

  # RuboCopでコードスタイルをチェック
  gem 'rubocop', require: false

  # Seedデータの生成を簡単にする
  gem 'faker', '~> 2.20'
end

# テスト環境用
group :test do
  # ヘッドレスブラウザを使った統合テスト
  gem 'capybara', '>= 3.26'
  gem 'selenium-webdriver'
  gem 'webdrivers'
end
