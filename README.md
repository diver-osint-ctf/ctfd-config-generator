# ctfd-config-generator for DIVER OSINT CTF
[English Version is here.](./README-en.md)

ctfd-config-generatorはCTFdを用いたCTFの作問を支援するためのツールです。以下の特徴があります。

- [x] CTFdをCLIベースで管理する[ctfcli](https://github.com/CTFd/ctfcli)の設定ファイルの作成
- [x] 作問に利用するディレクトリやファイルの作成

## 新しい問題の作り方
以下のコマンドで環境を作ってください。実行には、docker-compose または[Go](https://go.dev/doc/install)と Make が必要です。

```bash
git init
git submodule add https://github.com/diver-osint-ctf/ctfd-config-generator
echo "include ctfd-config-generator/Makefile" > Makefile

# 以下は好きな方を選択してください
# 1. docker-composeを使う場合
echo "version: '3'\nservices:\n  gen:\n    image: golang:1.21\n    volumes:\n      - .:/app\n    working_dir: /app" > docker-compose.yaml
docker-compose run gen bash -c "make gen"
# 2. GoとMakeを使う場合
make gen
```

実行が完了すると、以下のようなディレクトリが作成されます。

```bash
ジャンル名
└── 問題名
    ├── build           # 問題サーバで実行されるファイルを配置してください(オプショナル)
    ├── challenge.yml   # CTFdの設定を書いてください
    ├── flag.txt        # Flagを書いてください。フラグに正規表現を用いたり、複数のフラグを設定する場合はここに複数個のフラグを改行区切りで書いてください。
    ├── public          # 配布用ファイルを配置してください(オプショナル)
    ├── solver          # ソルバを配置してください(オプショナル)
    └── writeup
        └── README.md   # 作問者Writeupを配置してください
```

## コントリビュート
バグや要望などがあれば、Issueを作成するかPull Requestを作成してください。

