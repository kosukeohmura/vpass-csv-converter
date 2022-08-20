package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var (
	errLogger = log.New(os.Stderr, "[error] ", log.LstdFlags)
)

const (
	// exitCodeOK ok
	exitCodeOK int = 0
	// exitCodeError error
	exitCodeError int = 1
)

type (
	srcRecord struct {
		Date     string
		ShopName string
		Amount   int
	}

	outputRecord struct {
		Date   string
		Item   string
		Amount int
	}
)

var (
	srcFilepath string
	dstFilepath string
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	flag.StringVar(
		&srcFilepath,
		"src",
		"",
		"Path of usage list CSV file downloaded from Vpass.")

	flag.StringVar(
		&dstFilepath,
		"dst",
		"",
		"Output CSV file path. By default, it outputs in same directory as src file path, with -output prefix.")

	os.Args = args
	flag.Parse()
	if srcFilepath == "" {
		errLogger.Println(fmt.Errorf("specify src file path with --src option"))
		return exitCodeError
	}
	if dstFilepath == "" {
		dstFilepath = fmt.Sprintf("%s-converted.csv", strings.TrimSuffix(srcFilepath, filepath.Ext(srcFilepath)))
	}

	srcFile, err := os.Open(srcFilepath)
	if err != nil {
		errLogger.Println(fmt.Errorf("failed to open src file: %w", err))
		return exitCodeError
	}
	srcRecords, err := loadSrcRecords(srcFile)
	if err != nil {
		errLogger.Println(fmt.Errorf("failed to load src records: %w", err))
		return exitCodeError
	}
	srcFile.Close()

	parsedRecords, err := parseSrcRecords(srcRecords)
	if err != nil {
		errLogger.Println(fmt.Errorf("failed to parse src records: %w", err))
		return exitCodeError
	}

	outputRecords := make([]*outputRecord, 0, len(parsedRecords))
	for i := range parsedRecords {
		outputRecords = append(outputRecords, newOutputRecord(parsedRecords[i]))
	}

	dstFile, err := os.Create(dstFilepath)
	if err != nil {
		errLogger.Println(fmt.Errorf("failed to create dst file: %w", err))
		return exitCodeError
	}
	if err := writeOutputRecords(outputRecords, dstFile); err != nil {
		errLogger.Println(fmt.Errorf("failed to write output to file: %w", err))
		return exitCodeError
	}

	return exitCodeOK
}

func loadSrcRecords(srcReader io.Reader) ([][]string, error) {
	records, err := csv.NewReader(transform.NewReader(srcReader, japanese.ShiftJIS.NewDecoder())).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read src file: %w", err)
	}
	return records, nil
}

func parseSrcRecords(records [][]string) ([]*srcRecord, error) {
	srcRecords := make([]*srcRecord, 0, len(records))
	for i := range records {
		line := i + 1
		amount, err := strconv.Atoi(records[i][6])
		if err != nil {
			return nil, fmt.Errorf("failed to convert amount text \"%s\": line %d", records[i][6], line)
		}
		srcRecords = append(srcRecords, newSrcRecord(records[i][0], records[i][1], amount))
	}

	return srcRecords, nil
}

func writeOutputRecords(outputRecords []*outputRecord, writer io.Writer) error {
	header := []string{"Date", "Item", "Amount", "Purpose", "Method"}
	records := [][]string{header}
	for i := range outputRecords {
		records = append(records, []string{outputRecords[i].Date, outputRecords[i].Item, strconv.Itoa(outputRecords[i].Amount), "", ""})
	}
	return csv.NewWriter(writer).WriteAll(records)
}

func newSrcRecord(date string, shopName string, amount int) *srcRecord {
	return &srcRecord{Date: date, ShopName: shopName, Amount: amount}
}

func newOutputRecord(srcRecord *srcRecord) *outputRecord {
	return &outputRecord{Date: srcRecord.Date, Item: srcRecord.ShopName, Amount: srcRecord.Amount}
}
