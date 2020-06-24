# じゃんけん Bot

## 概要

`@zyankenbot ぐー`のようにリプライをすると、じゃんけんが行われる。
「ぐー」「ちょき」「ぱー」以外でリプライすると`ぐー、ちょき、ぱーのいずれかで` と返ってくる。
引き分けになることはない。

### あなたが勝ちの場合

`今日は負けを認めます。ただ、勝ち逃げは許しませんよ` と本田圭佑っぽくリプライが返ってくる。

### あなたが負けの場合

`俺の勝ち！何で負けたか、明日まで考えといてください` と本田圭佑っぽくリプライが返ってくる。

## セットアップ

`.env`を設定する

```
CONSUMER_KEY=<your twitter app's consumer key>
CONSUMER_SECRET=<your twitter app's consumer secret>
ACCESS_TOKEN_KEY=<your account's access token>
ACCESS_TOKEN_SECRET=<your account's access secret>
```
