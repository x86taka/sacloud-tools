# さくらクラウドのリソースをHCLに書き起こすツール

指定されたゾーンのサーバを取得し、紐付いているリソースをHCLで書き起こしている。

## Usage

envを設定する

```bash
export SAKURACLOUD_ACCESS_TOKEN=
export SAKURACLOUD_ACCESS_TOKEN_SECRET=
export SAKURACLOUD_ZONE="is1a"
```

対象のサーバリソースのtagによるフィルターがつけられる

hcl-generator {args}

{args}に`ubuntu ubuntu2004`とスペース区切りでtagを書くことができる。

この場合では、ubuntuとubuntu2004というtagが付いたサーバリソースに紐付いているリソースのHCLを`output/*.tf`に保存する
