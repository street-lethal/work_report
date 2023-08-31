# Work Report

## 準備 (初回のみ)

1. 設定ファイル等生成
   * ```shell
     make init
     ```
2. 実行ファイル(バイナリー)のリンク
   * Linux の場合
     ```shell
     make init-linux
     ```
   * Mac の場合
     ```shell
     make init-mac
     ```

## Jira

1. Jira Dashboard ( `/jira/dashboard/***` ) にブラウザーでアクセス
2. インスペクターで `Tempo User Timesheet` の部分の要素を選択し、 inner HTML をコピー
   * iframe になっているので iframe 内にある `<html>` タグを丸ごとコピーでOK
3. [data/jira.html](data/jira.html) に貼り付け
   * Mac の場合は以下コマンドで実行可能
     ```shell
     pbpaste > data/jira.html
     ``` 

## 報告用データ出力

1. [config/settings.json](config/settings.json) を編集
   * `months_ago` に対象月が何ヶ月前か入力(当月なら `0` 前月なら `1` )
   * `daily_report.starts_at` に開始時刻を入力
   * `daily_report.ends_at` に終了時刻を入力
   * `daily_report.rest_time` に休憩時間を入力
   * `holidays` に祝日の数字を入力 (土日は入力不要)
2. データ出力
   * ```shell
     make gen
     ```
     (ソースコードを直接実行する場合は以下)
     ```shell
     make gen-s
     ```
3. [data/report.json](data/report.json) にデータが書き込まれていれば成功 (この段階で編集したい部分があれば編集)

## 提出
### プラットフォームに自動ログインして提出する方法

1. ブラウザーでプラットフォームにログイン
2. 作業報告書の対象月詳細画面へ遷移
3. [config/settings.json](config/settings.json) を編集
    * `report_id` を URL からセット( `/p/workreport/***/` の `***` の部分)
4. [config/platform_id.json](config/platform_id.json) を編集
    * `email` 及び `password` にプラットフォームにログインする情報を設定
5. 提出
    * ```shell
      make send
      ```
      (ソースコードを直接実行する場合は以下)
      ```shell
      make send-s
      ```
6. ブラウザーの画面をリロードして、データが入力されていれば成功

### プラットフォームに自動ログインせず手動でセッション情報を設定して提出する方法

1. ブラウザーでプラットフォームにログイン
2. 作業報告書の対象月詳細画面へ遷移
3. [config/settings.json](config/settings.json) を編集
   * `report_id` を URL からセット( `/p/workreport/***/` の `***` の部分)
4. [data/platform_session.json](data/platform_session.json) を編集
   * `session_id` を Cookie (キーは `CAKEPHP`)からセット
   * `aws_auth` を Cookie (キーは `AWSELBAuthSessionCookie-0`) からセット
5. 提出
   * ```shell
     make send-m
     ```
     (ソースコードを直接実行する場合は以下)
     ```shell
     make send-m-s
     ```
6. ブラウザーの画面をリロードして、データが入力されていれば成功
