# wanpoll

ルーターの管理画面に表示される `WAN IP` をポーリングして監視し、変更を検知したら `Route53` の `Aレコード` を書き換えるコマンドです。

一定間隔でIPアドレスを確認し、変更されている時のみ Route53 にAPI経由でレコードを書き換えます。

ルーターに一時的に接続できなくなった場合や Route53 のリクエストに失敗した場合は、標準エラー出力にエラーを出力しポーリングを継続します。

自宅のルーターが DDNS (Dynamic DNS) に対応していなかったため、ルーターの設定ページをハックしてこのツールを作りました。ルーターから `パブリックIP` をぶっこ抜いて動的なIPアドレスを `DNS` に登録します。

## 対象

自宅ルーターの管理画面に表示された **IPv4** のアドレスを正規表現で取得します。IPアドレスが1つしか表示されていないページの URL をあらかじめルーターの管理画面から探して設定する必要があります。

以下のルーターでは、正しく動くことを確認しています。

- ソフトバンク光から支給された [`EWMTA2.1`](http://ybb.softbank.jp/support/connect/hikari/router/bbu2-setupmenu.php)
  - `http://172.16.255.254/GetWanIP.html` にアクセスして WanIP をぶっこぬきます。

## Install

[リリースページ](https://github.com/kawasin73/wanpoll/releases) にビルド済みの実行バイナリがアップロードされているため、OS と アーキテクチャに応じて適切なバイナリをダウンロードしてください。

```bash
curl -o wanpoll https://github.com/kawasin73/wanpoll/releases/download/v0.1.0/wanpoll-darwin-amd64
chmod 755 wanpoll
```

## Build

実行バイナリをビルドする場合は以下のコマンドで行います。`wanpoll` は、パッケージマネージャーに `dep` を利用しています。

```bash
git clone https://github.com/kawasin73/wanpoll.git && cd wanpoll
dep ensure
go build
```

## Example

```bash
$ ./wanpoll -hz xxxxxxxxxx -name=home.kawasin73.com,usa-reminder.kawasin73.com -user=xxxxxxxx -password=xxxxxxx -interval=10 -ttl=10
2018/07/11 23:00:23 detect new ip address "" -> "60.xxx.xxx.xxx"
2018/07/11 23:00:24 update record "home.kawasin73.com" == "60.xxx.xxx.xxx"
2018/07/11 23:00:24 update record "usa-reminder.kawasin73.com" == "60.xxx.xxx.xxx"
```

## オプション

```bash
$ go build && ./wanpoll -h
Usage of ./wanpoll:
  -hz string
    	hosted zone Id in route53
  -interval int
    	interval time polling to routerIpPage (second) (default 1)
  -name string
    	recode name separated with ',' comma
  -password string
    	basic auth password for routerIpPage
  -region string
    	aws region where route53 is placed (default "ap-northeast-1")
  -ttl int
    	dns record ttl (second) (default 60)
  -url string
    	routerIpPage url which shows global wan ip address (default "http://172.16.255.254/GetWanIP.html")
  -user string
    	basic auth username for routerIpPage
```

## AWS の Route53

`wanpoll` は AWS の Route53 にのみ対応しています。
Route53 の Hosted Zone のID を調べておく必要があります。

レコードの更新は、`ChangeResourceRecordSets` API を利用します。

- https://docs.aws.amazon.com/Route53/latest/APIReference/API_ChangeResourceRecordSets.html

### 認証

Route53 のリージョンは、オプションで設定します。

Route53 にアクセスする認証情報は、以下の3つの方法で設定できます。

- EC2 の IAM role (自宅のソフトバンク光の環境から使うことを想定しているので、この方法は利用しないはずです)
- shared credential file
- 環境変数
  - `AWS_ACCESS_KEY_ID`
  - `AWS_SECRET_ACCESS_KEY`
  - `AWS_SESSION_TOKEN` (optional) (一時トークンは使うのは難しいので使わないと思います)

詳しくは以下の URL を参照してください。

- https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials

### 権限

IAM には、インラインポリシーとして、以下のポリシーを追加します。

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "<Sid>",
            "Effect": "Allow",
            "Action": [
                "route53:ChangeResourceRecordSets"
            ],
            "Resource": [
                "arn:aws:route53:::hostedzone/<zoneid>"
            ]
        }
    ]
}
```

参考URLは以下の通りです。

- https://docs.aws.amazon.com/ja_jp/Route53/latest/DeveloperGuide/r53-api-permissions-ref.html#required-permissions-resource-record-sets
- https://docs.aws.amazon.com/ja_jp/general/latest/gr/aws-arns-and-namespaces.html#arn-syntax-route53

## LICENSE

MIT
