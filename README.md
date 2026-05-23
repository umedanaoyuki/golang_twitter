# golang で twiiter の API を作成する

### コンテナ作成・起動

```
# コンテナ作成＆起動
# -d をつけてバックグラウンドで起動
docker-compose up -d
```

```
# 起動中のコンテナの確認
docker-compose ps
```

### Swagger（API ドキュメント）

サーバー起動後、以下で Swagger UI を開けます。

```
http://localhost:8080/swagger/index.html
```

コメントを変更したら、プロジェクトルートでドキュメントを再生成してください。

```
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

### 動作確認

以下にアクセスできるか確認

```
http://localhost:8080/health_check
```

以下のレスポンスが返ってきたら成功

`"{status: ok}"`の json

### コンテナの停止

```
# コンテナの停止(停止するだけ)
docker-compose stop

# ストップして停止＆コンテナも削除
docker-compose down
```
