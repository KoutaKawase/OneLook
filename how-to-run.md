# tmux でアプリを常駐させる方法

Go のシングルバイナリを tmux のセッション内で動かし、ターミナルを閉じても動き続けるようにする手順です。

## 前提

- tmux がインストールされていること（未インストールの場合は `brew install tmux`）
- バイナリのパス: `/Users/yourname/bin/myapp`（環境に合わせて変更してください）
- アプリは `localhost:8080` で待ち受けること

## 起動

```bash
tmux new -d -s myapp '/Users/yourname/bin/myapp'
```

- `-d`：デタッチ状態（バックグラウンド）で起動します。
- `-s myapp`：セッション名を `myapp` にします。

動作確認：

```bash
curl http://localhost:8080
```

## セッションの確認

```bash
tmux ls
```

`myapp` が一覧に出ていれば起動中です。

## ログ（出力）を見る

セッションにアタッチすると、アプリの出力をそのまま見られます。

```bash
tmux attach -t myapp
```

アタッチ中にデタッチ（セッションは残したまま抜ける）するには：

1. `Ctrl+b` を押す
2. 手を離してから `d` を押す

## 停止

```bash
tmux kill-session -t myapp
```

セッションごと終了するため、アプリも停止します。

## 再起動（バイナリを更新したとき）

```bash
tmux kill-session -t myapp
tmux new -d -s myapp '/Users/yourname/bin/myapp'
```

## よく使うコマンド一覧

| 操作 | コマンド |
|---|---|
| 起動 | `tmux new -d -s myapp '/Users/yourname/bin/myapp'` |
| 一覧表示 | `tmux ls` |
| 中身を見る（アタッチ） | `tmux attach -t myapp` |
| 抜ける（デタッチ） | `Ctrl+b` → `d` |
| 停止 | `tmux kill-session -t myapp` |

## 注意点

- **Mac の再起動後は自動で立ち上がりません。** 再起動後に自動起動させたい場合は、launchd（LaunchAgent）の利用を検討してください。
- アプリが異常終了すると、そのペインが閉じてセッションも消えます。`tmux ls` で見当たらない場合は再度起動してください。
- 出力をファイルに残したい場合は、起動コマンドをリダイレクトします。

```bash
tmux new -d -s myapp '/Users/yourname/bin/myapp >> /Users/yourname/myapp.log 2>&1'
```
