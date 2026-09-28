# 画像生成 API 多言語デモ：Python・Node.js・C++・Go・Rust・PHP でテキスト生成から画像編集まで

このガイドでは、開発者向けの多言語画像生成サンプルプロジェクトを詳しく紹介します。6つの実装は同じ実用的な流れをたどります。まず OpenAI-compatible Images API の `/images/generations` にテキストプロンプトを送り、返された画像をローカルに保存します。次に、その保存済み画像ファイルを `/images/edits` の `image` フィールドに添付し、別の指示で編集します。プロトコルを追いながら動作を確かめたり、各言語の設計を比較したりするための教材です。

これは学習と統合検証のためのサンプルであり、OpenAI 公式 SDK でも、そのまま一般公開できる本番サービスでもありません。既定では独立したサービス [AI-ROUTER](https://ai-router.dev/ja) を利用しますが、API のベース URL は設定可能で、互換形式を提供する別のゲートウェイにも向けられます。モデルの利用可否、権限、料金、クォータ、画像サイズ、レスポンス形式はアカウントや接続先の設定に依存します。サンプルにモデル名が記載されていても、すべてのアカウントで利用できるという意味ではありません。

## このプロジェクトで学べること

- **画像生成から編集までの一連の処理：**プロンプトから初稿を生成し、ファイルとして保存した後、そのファイルを入力にして画像編集を行います。
- **Images API の実装ポイント：**Bearer 認証、生成用 JSON、編集用 multipart アップロード、レスポンス解析、ローカルファイル保存を確認できます。
- **6つの言語・エコシステムの違い：**業務上の順序と API はそろえたうえで、設定、HTTP クライアント、型の扱い、CLI の構成を比較できます。
- **アプリのプロトタイプ作成：**API サーバー、worker、キュー処理、バッチタスク、社内ツールの基礎として再利用できます。一般公開する際は認証、制限、コンテンツ管理などを別途実装してください。

通常のワークフローでは、生成と編集のために少なくとも2回 API を呼び出します。実行前にアカウントの料金、課金単位、データ保持条件を確認してください。

## 6言語のアーキテクチャ

各実装では、設定、HTTP クライアント、ワークフロー、CLI の責務を分けています。CLI にネットワーク処理を詰め込まないことで、クライアントや処理手順を別のアプリケーションから利用しやすくしています。一方、実際のファイル構成や実行方法は各言語に合わせています。

| 言語 | モジュールと役割 | 実際の実行方法 | サンプルでわかること |
| --- | --- | --- | --- |
| **Python** | `image_demo/config.py` が設定を検証し、`client.py` が標準ライブラリ `urllib` で JSON・multipart・画像レスポンスを処理します。`workflow.py` は生成物を保存してから編集に渡し、`cli.py` が引数と表示を担当します。 | `python/` ディレクトリで `python main.py --prompt ... --edit-prompt ...`。1コマンドで2段階を実行します。 | 追加の実行時依存がなく、スクリプト、Notebook、バッチ処理に組み込みやすい構成です。 |
| **Node.js** | `src/config.mjs` に共通設定、`src/client.mjs` に組み込みの `fetch` を使う通信処理、`src/workflow.mjs` に2段階の制御、`src/cli.mjs` に CLI を分けています。テストには標準の `node:test` を利用します。 | `nodejs/` ディレクトリで `npm run demo -- --prompt ... --edit-prompt ...`。 | npm の追加ランタイム依存なしで、JavaScript サーバーやキュー、serverless に応用できます。 |
| **C++** | `include/image_client.hpp` が再利用可能な API、`src/image_client.cpp` が libcurl による HTTP・multipart・デコード処理、`src/main.cpp` が CLI とワークフローを担当します。CMake でライブラリ、実行ファイル、オフラインテストを構成します。 | `cpp/` で CMake ビルド後、`imagegen --prompt ... --edit-prompt ...` を実行します。 | C++17、libcurl、CMake を使ったネイティブ向けの責務分離を確認できます。 |
| **Go** | `internal/config` が環境設定、`internal/client` が `net/http`・`encoding/json`・multipart、`cmd/imagegen` が CLI とワークフローを担当します。`httptest` を使ったクライアントテストがあります。 | `go/` ディレクトリで `go run ./cmd/imagegen --prompt ... --edit-prompt ...`。 | 標準ライブラリ中心の実装で、小さなサービスや常駐 worker に適しています。 |
| **Rust** | `config.rs` が設定、`client.rs` が `reqwest` と `serde` による通信、`types.rs` がレスポンス型、`workflow.rs` が処理順、`artifacts.rs` が保存を担当します。`main.rs` は `clap` で `generate`・`edit`・`workflow` を提供します。 | `rust/` で `cargo run --release -- workflow ...`。個別に generate と edit を実行することもできます。 | 型を明示した非同期クライアントと、通信・ファイルシステム・アプリケーション処理の分離を示します。 |
| **PHP** | `Config` が環境設定、`ImageClient` が cURL 通信、`ImageData` がレスポンスの正規化、`ArtifactStore` が保存、`Workflow` が処理の接続を担当します。`bin/image-demo.php` が CLI です。 | `php/` で `php bin/image-demo.php workflow ...`。`generate` と `edit` も利用できます。 | Composer の実行時パッケージは必須ではありません。PHP の cURL・JSON 拡張が必要です。 |

どの実装でも、編集リクエストの前に生成画像を保存します。2段階目で失敗しても最初の画像は残るため、確認や個別リトライができます。各モジュールの責務は[アーキテクチャガイド](../ARCHITECTURE.md)、言語を追加する際の方針は[コントリビューションガイド](../../CONTRIBUTING.md)を参照してください。

## OpenAI-compatible API のリクエスト形式

既定の API ベース URL は `https://api.ai-router.dev/v1` です。`base_url` は `/v1` を含むルート URL として設定し、クライアントがその後に endpoint を追加します。生成・編集の各 API リクエストにはアカウントの API key を Bearer token として付与します。

### 1. テキストから画像を生成：`POST /images/generations`

リクエストは `application/json` です。次は互換 API でよく使われる形の例です。実際に使用できるオプション項目はサービスとアカウントによって異なります。

```http
POST https://api.ai-router.dev/v1/images/generations
Authorization: Bearer $AI_ROUTER_API_KEY
Content-Type: application/json
```

```json
{
  "model": "gpt-image-2",
  "prompt": "雨上がりの屋上にある小さなガラス温室、柔らかな朝日、エディトリアル写真",
  "size": "1024x1024",
  "n": 1,
  "response_format": "b64_json"
}
```

クライアントはレスポンスの `data` 配列に含まれる先頭の画像を処理します。各実装は `b64_json` をデコードでき、サービスが画像 URL を返す場合は URL から画像を取得する処理も備えています。base64 形式なら、別の一時 URL へのアクセスを挟まず編集ステップへつなげられます。最終的なレスポンス形式は接続先サービスに従ってください。

### 2. 画像を編集：`POST /images/edits`

編集リクエストは `multipart/form-data` です。テキスト指示だけでなく、1つ目の処理で保存した画像を `image` フィールドに添付します。

| フィールド | 役割 |
| --- | --- |
| `model` | アカウントで有効になっており、画像編集に対応するモデル ID。 |
| `prompt` | 変更内容の指示。維持したい被写体や構図も明記すると意図が伝わりやすくなります。 |
| `size` | 任意の出力サイズ。`1024x1024` などが使えるかはモデルとアカウント次第です。 |
| `response_format` | サンプルでは移植しやすい `b64_json` を指定します。サービス設定によって URL が返る場合もあります。 |
| `image` | 必須の画像ファイル。文生図ステップで保存した画像をアップロードします。 |

編集結果も `data` を含む画像レスポンスです。共有フィールド一覧は[API 契約ドキュメント](../API_CONTRACT.md)を参照してください。

### モデル選択と互換性に関する注意

`--model` や `AI_ROUTER_MODEL` の値は、ゲートウェイへ渡すモデル ID です。`gpt-image-2`、`gpt-image-2.5`、`gpt-image-2.5-flare`、`gpt-image-2.5-burst` などを試す場合は、アカウントのモデル一覧に表示され、かつ必要な生成・編集操作で使用可能であることを確認してください。サンプルが権限を拡張するわけではありません。モデル名、料金、対応サイズが将来も同じであることを保証しません。既定の `1024x1024` はサンプル値であり、すべてのモデルで使えるという保証ではありません。

「OpenAI-compatible」は、endpoint やデータ形式の一部に互換性があり、既存の統合パターンを応用しやすいという意味です。AI-ROUTER が OpenAI の製品であることや、任意パラメータ、モデル、挙動が完全一致することを意味しません。必ず利用するサービスの最新ドキュメントを確認してください。

## 環境変数と設定

リポジトリ直下の `.env.example` は設定例です。各 CLI は `.env` を自動的に読み込みません。シェルで環境変数を設定するか、信頼できる dotenv ローダーを利用してください。実際の API key を Git にコミットしないでください。

共通の基本設定例：

```bash
export AI_ROUTER_API_KEY="自分のAPIキー"
export AI_ROUTER_BASE_URL="https://api.ai-router.dev/v1"
export AI_ROUTER_MODEL="gpt-image-2"
export AI_ROUTER_SIZE="1024x1024"
```

環境変数の読み込み方には言語ごとの差があります。以下は現在のコードが参照する値です。

| 実装 | API key | Base URL | モデル・サイズ・タイムアウト |
| --- | --- | --- | --- |
| Python | `IMAGE_API_KEY`、`OPENAI_API_KEY`、`AI_ROUTER_API_KEY` の順に確認 | `IMAGE_API_BASE_URL`、`AI_ROUTER_BASE_URL` | `IMAGE_MODEL` / `AI_ROUTER_MODEL`、`IMAGE_SIZE` / `AI_ROUTER_SIZE`、`IMAGE_QUALITY`、`IMAGE_TIMEOUT_SECONDS`（秒） |
| Node.js | `IMAGE_API_KEY`、`OPENAI_API_KEY`、`AI_ROUTER_API_KEY` | `IMAGE_API_BASE_URL`、`AI_ROUTER_BASE_URL` | `IMAGE_MODEL` / `AI_ROUTER_MODEL`、`IMAGE_SIZE` / `AI_ROUTER_SIZE`、`IMAGE_QUALITY`、`IMAGE_TIMEOUT_MS`（ミリ秒） |
| C++ | `OPENAI_API_KEY`、`AI_ROUTER_API_KEY` | `OPENAI_BASE_URL`、`AI_ROUTER_BASE_URL` | `IMAGE_MODEL` / `AI_ROUTER_MODEL`、`IMAGE_SIZE` / `AI_ROUTER_SIZE`、`IMAGE_RESPONSE_FORMAT`、`IMAGE_API_TIMEOUT` / `AI_ROUTER_IMAGE_TIMEOUT`（秒） |
| Go | `OPENAI_API_KEY`、`AI_ROUTER_API_KEY` | `OPENAI_BASE_URL`、`AI_ROUTER_BASE_URL` | `IMAGE_MODEL` / `AI_ROUTER_MODEL`、`IMAGE_SIZE` / `AI_ROUTER_SIZE`、`IMAGE_RESPONSE_FORMAT`、`IMAGE_API_TIMEOUT` / `AI_ROUTER_IMAGE_TIMEOUT`（秒） |
| Rust | `AI_ROUTER_API_KEY`、`OPENAI_API_KEY` | `AI_ROUTER_BASE_URL` | `AI_ROUTER_TIMEOUT_SECONDS`。モデル・サイズはサブコマンド引数で指定でき、CLI は `AI_ROUTER_MODEL`・`AI_ROUTER_SIZE` も読み取ります。 |
| PHP | `AI_ROUTER_API_KEY`、`OPENAI_API_KEY` | `AI_ROUTER_BASE_URL` | `AI_ROUTER_TIMEOUT_SECONDS`。モデル・サイズは CLI 引数またはコマンド側の既定値を使います。 |

多くの CLI オプションは環境変数の既定値を上書きします。PHP には `--api-key` オプションがないため、環境変数経由で渡してください。Rust の `--base-url` はトップレベルオプションで、`--model` と `--size` は各サブコマンドのオプションです。

## インストールと実行手順

以下のパスはリポジトリのルートを基準とします。完全な処理は API を2回呼び出します。実際の API を利用する前に、key 不要のローカルテストを実行できます。

### Python（3.10 以上）

このデモは Python 標準ライブラリのみを使います。まず仮想環境を作成します。

```bash
cd marketing/image-generation-multilang/python
python -m venv .venv
```

使用するシェルに合わせて有効化してください。

```bash
# macOS / Linux
source .venv/bin/activate

# Windows PowerShell
.\.venv\Scripts\Activate.ps1
```

key を設定して、生成と編集を連続実行します。

```bash
export IMAGE_API_KEY="自分のAPIキー"
python main.py \
  --prompt "雨上がりの屋上にあるガラス温室、映画のような編集写真" \
  --edit-prompt "温室の構造は維持して、光を暖かい夕日の色にする" \
  --output-dir ./outputs
```

PowerShell では `$env:IMAGE_API_KEY = "自分のAPIキー"` を設定し、そのシェルの改行ルールに従って実行してください。`outputs/` に生成画像と編集画像が保存され、拡張子はレスポンスに応じます。`python -m pip install -e .` を実行すれば開発用の `image-demo` エントリーポイントもインストールできます。オフラインテストは `python -m unittest discover -s tests -v` です。

### Node.js（18.17 以上）

Node.js の組み込み `fetch` を使用し、実行時の追加 npm パッケージはありません。まずテストします。

```bash
cd marketing/image-generation-multilang/nodejs
npm test
```

生成と編集を実行します。

```bash
export AI_ROUTER_API_KEY="自分のAPIキー"
npm run demo -- \
  --prompt "雨上がりの屋上にあるガラス温室、映画のような編集写真" \
  --edit-prompt "温室の構造は維持して、光を暖かい夕日の色にする" \
  --model gpt-image-2
```

PowerShell では `$env:AI_ROUTER_API_KEY` を使います。既定の出力先は `nodejs/outputs/generated.png` と `nodejs/outputs/edited.png` です。`--size`、`--quality`、`--base-url`、ミリ秒単位の `--timeout`、`--output-dir` も指定できます。`npm test` は実 API を呼ばないテストです。

### C++17（CMake 3.20 以上・libcurl）

C++17 対応コンパイラ、CMake、libcurl の開発用ファイルを準備してください。Windows で vcpkg を使う場合の依存インストール例：

```powershell
vcpkg install curl:x64-windows
```

vcpkg の場合は CMake に `vcpkg.cmake` toolchain ファイルを指定します。一般的な構成・ビルド・オフラインテストは次の通りです。

```bash
cd marketing/image-generation-multilang/cpp
cmake -S . -B build -DBUILD_TESTING=ON
cmake --build build --config Release
ctest --test-dir build --output-on-failure
```

key を設定して実行します。

```bash
export AI_ROUTER_API_KEY="自分のAPIキー"
./build/imagegen \
  --base-url https://api.ai-router.dev/v1 \
  --model gpt-image-2 \
  --prompt "雨上がりの屋上にあるガラス温室" \
  --edit-prompt "建物の輪郭は保ったまま、柔らかな夕日の光を加える" \
  --output generated.png --edited-output edited.png
```

Visual Studio のマルチ構成ジェネレーターでは通常 `build/Release/imagegen.exe` に実行ファイルが出力されます。CTest は API key 不要のローカルテストです。CLI では `--size`、`--response-format`、`--timeout`、出力パスも指定できます。

### Go（1.22 以上）

標準ライブラリを使用するため、実行用の追加 Go パッケージは不要です。

```bash
cd marketing/image-generation-multilang/go
go test ./...
go vet ./...
```

二段階の処理を実行します。

```bash
export AI_ROUTER_API_KEY="自分のAPIキー"
go run ./cmd/imagegen \
  --base-url https://api.ai-router.dev/v1 \
  --model gpt-image-2 \
  --prompt "雨上がりの屋上にあるガラス温室" \
  --edit-prompt "建物の輪郭は保ったまま、柔らかな夕日の光を加える" \
  --output generated.png --edited-output edited.png
```

既定では実行時のカレントディレクトリに `generated.png` と `edited.png` が作られます。`go test` はローカルの `httptest.Server` を使い、実際の画像 API や料金を発生させません。

### Rust（stable Rust toolchain）

依存ライブラリは Cargo が `Cargo.toml` から取得します。初回ビルド時には crates.io へのネットワーク接続が必要になることがあります。

```bash
cd marketing/image-generation-multilang/rust
export AI_ROUTER_API_KEY="自分のAPIキー"
```

生成画像を保存してから編集する workflow コマンド：

```bash
cargo run --release -- workflow \
  --prompt "雨上がりの屋上にあるガラス温室" \
  --edit-prompt "構図は保ったまま、暖かな吊り下げ照明を追加する" \
  --model gpt-image-2 --size 1024x1024 \
  --generated output/generated.png \
  --edited output/edited.png
```

出力先の親ディレクトリは自動作成されます。`cargo run -- generate --prompt ... --output ...` と `cargo run -- edit --image ... --prompt ... --output ...` で個別に実行することもできます。`cargo test` でテストターゲットを確認できます。このクライアントは `b64_json` を指定し、URL 形式の画像レスポンスも処理します。

### PHP（8.1 以上、cURL・JSON 拡張）

プロジェクト内に簡易 autoloader があるため、デモを実行するだけなら Composer のパッケージインストールは不要です。PHP CLI で拡張が有効か確認してください。

```bash
cd marketing/image-generation-multilang/php
php -m
export AI_ROUTER_API_KEY="自分のAPIキー"
```

生成から編集までを実行します。

```bash
php bin/image-demo.php workflow \
  --prompt "雨上がりの屋上にあるガラス温室" \
  --edit-prompt "構図は保ったまま、暖かな吊り下げ照明を追加する" \
  --model gpt-image-2 --size 1024x1024 \
  --generated output/php-generated.png \
  --edited output/php-edited.png
```

`generate` と `edit` を個別に実行することもできます。既定の保存先は `output/php-generated.png` と `output/php-edited.png` で、API URL は `AI_ROUTER_BASE_URL` から読み込みます。API key は `AI_ROUTER_API_KEY` または `OPENAI_API_KEY` 環境変数から読み込み、CLI フラグでは指定しません。Composer アプリに組み込む場合は `composer.json` の PSR-4 設定を使い、必要に応じて `composer dump-autoload` を実行してください。

## 生成ファイルと実行コスト

出力はテキストから作成した初稿と、その画像を入力にして作成した編集版の2つです。保存先は言語ごとの既定値か CLI オプションで決まります。通常の出力フォルダーは `.gitignore` の対象ですが、コミット前に `git status` で確認してください。顧客データや個人情報を含む画像を公開リポジトリへ誤って追加しないようにしましょう。

サンプルは自動リトライを行いません。画像 API のリトライは二重リクエストや追加課金につながる場合があります。編集に失敗しても生成済みの画像は残るため、まず初稿を確認し、必要な処理だけを意図的に再実行してください。

## トラブルシューティング

| 症状 | 確認すること |
| --- | --- |
| `401`・`403`・API key がない | CLI を起動した同じシェルに key を設定したか、環境変数名がその言語に合っているか確認します。key の有効性とアカウント権限も確認してください。ログやソースに秘密を出さないでください。 |
| モデルが見つからない、権限エラー | アカウントのモデル一覧と正確な ID を確認します。生成と編集の両方をサポートするモデルかも確認してください。 |
| 画像サイズや quality が拒否される | 選択したモデルとアカウントに対応する値を使います。`1024x1024` は例の既定値であり、上流の対応確認なしに保証される値ではありません。 |
| タイムアウト、接続失敗 | Base URL、TLS、プロキシ、ネットワーク、タイムアウト値の単位を確認します。Node.js はミリ秒、他の多くの実装は秒です。課金状態が不明なときに無制限に再試行しないでください。 |
| レスポンスに画像がない、base64 エラー | 上流が返した JSON 形式と互換性を確認します。URL を使う実装では、一時 URL の期限切れやダウンロード制限も確認します。 |
| 生成は成功したが編集に失敗する | 生成ファイルの存在、読み取り権限、ファイルサイズ、編集用 prompt、選択モデルの編集対応を確認します。生成画像を残し、編集段階だけを調査してください。 |
| CMake が CURL を見つけない | libcurl 開発用ファイルと CMake の toolchain/config を確認します。vcpkg 使用時は構成時に `vcpkg.cmake` を指定してください。 |
| Rust の初回ビルドが失敗する | stable Rust/Cargo と依存レジストリへの接続を確認します。`cargo check` で依存・コンパイルエラーを切り分けてください。 |
| PHP の cURL 関数が使えない | CLI が参照する `php.ini` で `curl` 拡張を有効にします。Web サーバー用 PHP と CLI の設定は異なることがあります。 |

## セキュリティ、プライバシー、責任ある利用

- API key は環境変数かシークレット管理サービスに保管してください。Git、shell history、issue、スクリーンショットへ実キーを残さないでください。公開 Web アプリではブラウザーにサービスキーを渡してはいけません。
- HTTPS を使い、base URL が意図したサービスであることを確認してください。実際の画像処理、ログ、保管期間、利用データの扱いは接続先サービスの規約・プライバシー文書で確認します。
- プロンプトや入力画像・生成画像には個人情報や機密情報が含まれる場合があります。利用権限のある画像だけを送信し、保存先のアクセス制御、保持期間、削除手順を決めてください。
- API レスポンス内の画像 URL は外部入力として扱います。公開サービスにする場合は、接続先 allow-list、リダイレクト、最大サイズ、タイムアウト、MIME と実ファイル形式の検証を行い、SSRF や過剰なリソース消費に備えてください。
- 一般ユーザーに提供する前に、利用者認証、利用上限、rate limit、コンテンツチェック、キャンセル、監視、ログの秘匿化を追加してください。このデモはそれらを完成済み機能として提供していません。
- サービス利用規約、入力画像の権利、適用される法令を守り、公開前に人の目で品質と内容を確認してください。

追加の注意点は[セキュリティガイド](../SECURITY.md)にまとめています。

## 自分のサービスへ拡張するには

1. **ローカル保存をオブジェクトストレージに置き換える：**アクセス制御、暗号化、ライフサイクルを設定しつつ、`generate → save → edit → save` の処理境界を保ちます。
2. **入力とコストを制御する：**プロンプトの長さ、アップロードの MIME・サイズ、ユーザー別・モデル別のクォータを設定します。キャンセルとタイムアウトも全層に伝えます。
3. **リトライは課金仕様を確認してから追加する：**再試行可能な一時エラーに限り、回数を制限した backoff を使います。冪等性と請求の扱いを確認せず機械的に再送しないでください。
4. **必要な範囲で観測する：**request ID、処理段階、HTTP 状態、所要時間を記録し、プロンプト、個人情報、画像 URL は最小化またはマスキングします。API key は記録しません。
5. **API 変更をテストで固定する：**フィールドを追加したら共通契約、対象となる各言語実装、オフラインテストを更新します。multipart は mock server で検証すれば、毎回実 API のクォータを消費する必要はありません。
6. **新しい言語を段階的に追加する：**設定、独立したクライアント、workflow、CLI を分け、境界ごとにテストします。CLI に通信処理を重複実装せず、provider 固有フィールドがすべてのサービスで共通だと仮定しないでください。

拡張方針は[アーキテクチャ](../ARCHITECTURE.md)と[共通 API 契約](../API_CONTRACT.md)、比較に使える入力例は[プロンプト集](../../examples/prompts.md)をご覧ください。

## 関連リンク

- [AI-ROUTER 日本語ホーム](https://ai-router.dev/ja)：サービス、アカウント、API 接続を確認するための入口です。利用可能なモデルや機能はアカウントごとに確認してください。
- [API ベース URL](https://api.ai-router.dev/v1)：デモの既定値です。このリンクだけで認証や API 仕様の説明を代替するものではありません。
- [API 契約](../API_CONTRACT.md) · [プロジェクト構成](../ARCHITECTURE.md) · [セキュリティ](../SECURITY.md)
- [プロンプト例](../../examples/prompts.md) · [コントリビューションガイド](../../CONTRIBUTING.md)

この多言語サンプルが役立った場合は、再現可能なテスト結果や具体的な改善案をリポジトリで共有してください。サイトへのリンクは AI-ROUTER の紹介を目的としています。本ドキュメントの公開によって検索順位、流入数、特定の生成品質が保証されるものではありません。
