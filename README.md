# bottlemail
https://bottlemail.appspot.com/

## これは何?
ボトルメールを簡単に送りあえるアプリをつくりました!
App Engine(Golang) + CloudSQL(MySQL) の構成です

## コマンド一覧

### 本番環境
- デプロイ
$ gcloud app deploy

### ローカル環境
- CloudSQLへの接続用プロキシ
$ ./cloud_sql_proxy -instances=bottlemail:asia-east2:bottlemaildb=tcp:3306
これでローカルホストを通じてCloudSQL側のデータベースにアクセス可能

- デバック
$ dev_appserver.py app.yaml
