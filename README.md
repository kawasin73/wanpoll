# wanpoll

`WAN IP` をポーリングして監視し、変更を検知したら `Route53` の Aレコード を書き換えるコマンドです。

## 対象

- ソフトバンク光から支給された [`EWMTA2.1`](http://ybb.softbank.jp/support/connect/hikari/router/bbu2-setupmenu.php)

`http://172.16.255.254/GetWanIP.html` にアクセスして WanIP をぶっこぬきます。

## オプション
