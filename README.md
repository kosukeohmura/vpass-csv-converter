# vpass-csv-converter

vpass-csv-converter converts csv downloaded from [Vpass](https://www.smbc-card.com/mem/index.jsp). Vpass has two transaction list format: Fixed or non-fixed.

Fixed format:

```csv
大村　幸佑　様,0000-0000-0000-0***,ＶＩＳＡ
2022/08/05,ヨドバシカメラ　通信販売,4853,１,１,4853,
2022/08/15,ＡＭＡＺＯＮ．ＣＯ．ＪＰ,3400,１,１,3400,
大村　幸佑　様,0000-0000-0000-0***,ＡｐｐｌｅＰａｙ／ｉＤ
2022/08/18,ファミリーマート／ｉＤ,340,１,１,340,ﾌｱﾐﾘ-ﾏ-ﾄ/ID
,,,,,123456
,,,,,123456,
```

Non-fixed format:

```csv
2022/7/4,ＡＭＡＺＯＮ．ＣＯ．ＪＰ,ご家族,1回払い,,'22/08,3844,3844,,,,,
2022/7/16,セブン－イレブン／ｉＤ,ご家族,1回払い,,'22/08,98,98,,,,,
```

This tool convert to:

```csv
Date,Item,Amount,Purpose,Method
2022/7/4,ＡＭＡＺＯＮ．ＣＯ．ＪＰ,3844,,
2022/7/16,セブン－イレブン／ｉＤ,98,,
```

## Usage

```sh
go run main.go --src=/Users/kosukeohmura/Downloads/202208.csv
```

### Options

- `--src`: (Mandatory) You can change output path by passing this option.
- `--dst`: (Optional) You can change output path by passing this option. The default output path is the same directory as the src file.
- `--srcfixed`: (Optional) Specify source CSV is fixed or not. Vpass has two types of transaction list: fixed and non-fixed. The default value is `true`.
