## プラットフォームに自動ログインせず手動でセッション情報を設定して提出する方法

1. ブラウザーでプラットフォームにログイン
2. [data/platform_session.json](data/platform_session.json) を編集
   * `session_id` を Cookie (キーは `CAKEPHP`)からセット
   * `aws_auth` を Cookie (キーは `AWSELBAuthSessionCookie-0`) からセット
3. 提出
   * ```shell
     make send-m
     ```
4. (結果確認) ブラウザーで作業報告書の対象月詳細画面にアクセスして、データが入力されていれば成功
