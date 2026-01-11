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

### 動作確認

```
http://localhost:8080/health_check
```

にアクセスできるか確認

以下のレスポンスが返ってきたら成功

`"{status: ok}"`の json

### コンテナの停止

```
# コンテナの停止(停止するだけ)
docker-compose stop

# ストップして停止＆コンテナも削除
docker-compose down
```
