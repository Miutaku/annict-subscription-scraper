# about

[annict](https://annict.com/)のアニメの放送情報から、各配信サービスでの作品の配信状況を返します。


# sample usage

annictにて、自分の利用している配信サービスで配信されないアニメタイトルだけを、録画サーバやレコーダーで録画したいときに有用です。
dアニメストアで配信があるかをチェックしたいときは、以下の様な感じです。
```
$ curl "http://localhost:8080/?id=10976"| jq '.services[] | select(.name == "dアニメストア")|.available
true
```


# how to use

例えば、https://annict.com/works/10976 の情報を取得する場合は以下のようにします。
```
$ curl "http://localhost:8080/?id=10976"| jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   383  100   383    0     0    290      0  0:00:01  0:00:01 --:--:--   290
{
  "services": [
    {
      "name": "バンダイチャンネル",
      "available": true
    },
    {
      "name": "ニコニコチャンネル",
      "available": false
    },
    {
      "name": "dアニメストア ニコニコ支店",
      "available": false
    },
    {
      "name": "dアニメストア",
      "available": true
    },
    {
      "name": "Amazon プライム・ビデオ",
      "available": true
    },
    {
      "name": "Netflix",
      "available": true
    },
    {
      "name": "ABEMAビデオ",
      "available": true
    }
  ]
}
```

