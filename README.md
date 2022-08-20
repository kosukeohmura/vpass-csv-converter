# vpass-csv-converter

vpass-csv-converter converts csv downloaded from [Vpass](https://www.smbc-card.com/mem/index.jsp).

Converts from:

```csv
2022/7/4,ＡＭＡＺＯＮ．ＣＯ．ＪＰ,ご家族,1回払い,,'22/08,3844,3844,,,,,
2022/7/16,セブン－イレブン／ｉＤ,ご家族,1回払い,,'22/08,98,98,,,,,
```

To:

```csv
Date,Item,Amount,Purpose,Method
2022/7/4,ＡＭＡＺＯＮ．ＣＯ．ＪＰ,3844,,
2022/7/16,セブン－イレブン／ｉＤ,98,,
```

## Usage

```sh
go run main.go --src=/Users/kosukeohmura/Downloads/202208.csv
```

The `--src` opton is mandatory. It creates a file in same directory as src file by default. In this example, it creates /Users/kosukeohmura/Downloads/202208-converted.csv. You can change output path by passing `--dst` option.
