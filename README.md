# sum10-solver

[SUM10 Puzzle](https://neetsdkasu.github.io/game/sum10/index.html)のソルバーを作ろうと試みたもの･･･


作ろうと試みたものの、難しくて作れなかった… orz  

なので当プログラムは、ランダムに選択していくという形で解を50万個ほど生成し、スコアのもっとも良いものを出力する、というプログラムです。  

### インストール

```bash
$ go install github.com/neetsdkasu/sum10-solver@latest
```
※Go言語のバージョンは`1.18`以上が必要です。  


### 実行方法

たとえば`SEED=5531`の解を作る場合は、  
```bash
$ sum10-solver -seed 5531
```
解が出力された`result05531.txt`というファイルが生成される。  



