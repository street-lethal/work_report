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
2. `...` > `生データをダウンロード` > `CSV` で対象月データを CSV 形式でダウンロード
3. ダウンロードした CSV ファイルを `jira.csv` にリネームし [data/](data) ディレクトリー下へ配置

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

### 注意

実行すると [data/jira.csv](data/jira.csv) が `data/jira_****_**_**_******.csv` にリネームされるので、再実行する場合は [data/jira.csv](data/jira.csv) にリネームし直してから実行

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
