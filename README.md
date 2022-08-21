# LSemiChat-Server

## 環境構築
makeコマンドが入っていれば

```
// 初期起動
make dev-up
// 停止
make dev-stop
// 再開
make dev-start
// 終了-DBにいれていたデータが吹っ飛びます
make dev-down
```

makeコマンドが入っていなければ

```
// 初期起動
docker-compose up -d
// 停止
docker-compose stop
// 再開
docker-compose start
// 終了
docker-compose down --rmi local --volumes
```

`$ make help` or `Makefile`を参考にしてください。

## memo
やりたいこと
- DBのモデルとentity分けれたら良さそう
- ormapperみたいなのを自作する。ビルダー使えばできそう？

limit_users: 0 = 無制限とか...

api/constants/secret.goの作成が必要になります。以下をconstで宣言してください。
- JWTSecret
- JWTUserIDClaimsKey
- SessionName
