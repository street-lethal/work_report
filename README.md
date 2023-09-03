# Work Report

## 準備 (初回のみ)

* Linux の場合
  ```shell
  make init-linux
  ```
* Mac の場合
  ```shell
  make init-mac
  ```

## Jira

1. ブラウザーで Jira Tempo ( `/plugins/servlet/ac/io.tempo.jira/tempo-app` ) にアクセス
2. インスペクターで iframe 内にある `<html>` タグを選択し、 outer HTML をコピー
3. [data/jira.html](data/jira.html) に貼り付け
   * Mac の場合は以下コマンドで実行可能
     ```shell
     pbpaste > data/jira.html
     ``` 

## 報告用データ出力

1. [config/settings.json](config/settings.json) を編集
   * `months_ago` に対象月が何ヶ月前か入力(当月なら `0` 前月なら `1` )
   * `daily_report.starts_at` に開始時刻を入力
   * `daily_report.rest_time` に休憩時間を入力
2. データ出力
   * ```shell
     make gen
     ```
3. (結果確認) [data/report.json](data/report.json) にデータが書き込まれていれば成功 (この段階で編集したい部分があれば編集)

## 提出
### プラットフォームに自動ログインして提出する方法

1. [config/platform_id.json](config/platform_id.json) を編集
    * `email` 及び `password` にプラットフォームにログインする情報を設定
2. 提出
    * ```shell
      make send
      ```
3. (結果確認) ブラウザーでプラットフォームの作業報告書の対象月詳細画面にアクセスして、データが入力されていれば成功

### プラットフォームに自動ログインせず手動でセッション情報を設定して提出する方法

1. ブラウザーでプラットフォームにログイン
2. [data/platform_session.json](data/platform_session.json) を編集
   * `session_id` を Cookie (キーは `CAKEPHP`)からセット
   * `aws_auth` を Cookie (キーは `AWSELBAuthSessionCookie-0`) からセット
3. 提出
   * ```shell
     make send-m
     ```
4. (結果確認) ブラウザーで作業報告書の対象月詳細画面にアクセスして、データが入力されていれば成功
