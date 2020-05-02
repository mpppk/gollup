package lib

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Input は複数行からなる文字列から値をいい感じに取り出すためのメソッドを提供します.
type Input struct {
	lines [][]string
}

func (i *Input) validateColIndex(index int) error {
	if index < 0 {
		return errors.New(fmt.Sprintf("index is under zero: %d", index))
	}

	return nil
}

func (i *Input) validateRowIndex(index int) error {
	if index >= len(i.lines) {
		return errors.New(fmt.Sprintf("index(%d) is larger than lines(%d)", index, len(i.lines)))
	}

	if index < 0 {
		return errors.New(fmt.Sprintf("index is under zero: %d", index))
	}
	return nil
}

// GetLines は指定された範囲の行を返します. startRowIndexは含み、endRowIndexは含みません.存在しない行を指定した場合失敗します.
func (i *Input) GetLines(startRowIndex, endRowIndex int) ([][]string, error) {
	if err := i.validateRowIndex(startRowIndex); err != nil {
		return nil, fmt.Errorf("invalid start row index: %v", err)
	}
	if err := i.validateRowIndex(endRowIndex - 1); err != nil {
		return nil, fmt.Errorf("invalid end row index: %v", err)
	}
	return i.lines[startRowIndex:endRowIndex], nil
}

// GetStringLinesFrom は指定された行を含む、それ以降の行を返します.存在しない行を指定した場合失敗します.
func (i *Input) GetStringLinesFrom(fromIndex int) (newLines [][]string, err error) {
	for index := range i.lines {
		if index < fromIndex {
			continue
		}
		newLine, err := i.GetLine(index)
		if err != nil {
			return nil, err
		}
		newLines = append(newLines, newLine)
	}
	return
}

// GetValue は指定された行と列の値を返します.存在しない行か列を指定した場合失敗します.
func (i *Input) GetValue(rowIndex, colIndex int) (string, error) {
	line, err := i.GetLine(rowIndex)
	if err != nil {
		return "", err
	}
	if colIndex < 0 || colIndex >= len(line) {
		return "", fmt.Errorf("Invalid col index: %v ", colIndex)
	}
	return line[colIndex], nil
}

// GetFirstValue は指定した行の最初の列の値を返します. 存在しない行を指定した場合失敗します.
func (i *Input) GetFirstValue(rowIndex int) (string, error) {
	return i.GetValue(rowIndex, 0)
}

// GetColLine は指定された列を返します。値が存在しない行がある場合失敗します.
func (i *Input) GetColLine(colIndex int) (newLine []string, err error) {
	if err := i.validateColIndex(colIndex); err != nil {
		return nil, err
	}

	for i, line := range i.lines {
		if len(line) <= colIndex {
			return nil, errors.New(fmt.Sprintf("col index(%d) is larger than %dth line length(%d)", colIndex, i, len(line)))
		}
		newLine = append(newLine, line[colIndex])
	}

	return newLine, nil
}

// GetLine は指定された行を返します. 存在しない行がある場合、失敗します.
func (i *Input) GetLine(index int) ([]string, error) {
	if err := i.validateRowIndex(index); err != nil {
		return nil, err
	}
	return i.lines[index], nil
}

// ReadAsStringGridFrom は指定された行以降の行を、一文字ずつのsliceとして返します.
func (i *Input) ReadAsStringGridFrom(fromIndex int) ([][]string, error) {
	lines, err := i.GetStringLinesFrom(fromIndex)
	if err != nil {
		return nil, err
	}

	var m [][]string
	for _, line := range lines {
		if len(line) > 1 {
			return nil, fmt.Errorf("unexpected length line: %v", line)
		}

		var mLine []string
		for _, r := range line[0] {
			mLine = append(mLine, string(r))
		}
		m = append(m, mLine)
	}
	return m, nil
}

// GetIntLine は指定された行をIntとして返します. 存在しない行がある場合、失敗します.
func (i *Input) GetIntLine(index int) ([]int, error) {
	if err := i.validateRowIndex(index); err != nil {
		return nil, err
	}

	newLine, err := StringSliceToIntSlice(i.lines[index])
	if err != nil {
		return nil, fmt.Errorf("%dth index: %v", index, err)
	}
	return newLine, nil
}

// MustGetIntLine は指定された行をIntとして返します. 存在しない行がある場合、失敗します.
func (i *Input) MustGetIntLine(index int) []int {
	_v0, _err := i.GetIntLine(index)
	if _err != nil {
		panic(_err)
	}
	return _v0
}

// MustReadAsStringGridFrom は指定された行以降の行を、一文字ずつのsliceとして返します.
func (i *Input) MustReadAsStringGridFrom(fromIndex int) [][]string {
	_v0, _err := i.ReadAsStringGridFrom(fromIndex)
	if _err != nil {
		panic(_err)
	}
	return _v0
}

// MustNewInputFromReader はreaderから入力を読み込み、Inputを生成します.
func MustNewInputFromReader(reader *bufio.Reader) *Input {
	_v0, _err := NewInputFromReader(reader)
	if _err != nil {
		panic(_err)
	}
	return _v0
}

func NewIntGridMap(row, col int, defaultValue int) (m [][]int) {
	for i := 0; i < row; i++ {
		var newLine []int
		for j := 0; j < col; j++ {
			newLine = append(newLine, defaultValue)
		}
		m = append(m, newLine)
	}
	return
}

func StringSliceToIntSlice(line []string) (ValueLine []int, err error) {
	newLine, err := toSpecificBitIntLine(line, 64)
	if err != nil {
		return nil, err
	}
	for _, v := range newLine {
		ValueLine = append(ValueLine, int(v))
	}
	return
}

// NewInput はscannerから入力を読み込み、Inputを生成します.
func NewInput(scanner *bufio.Scanner) *Input {
	return &Input{
		lines: toLines(scanner),
	}
}

// NewInputFromReader はreaderから入力を読み込み、Inputを生成します.
func NewInputFromReader(reader *bufio.Reader) (*Input, error) {
	lines, err := toLinesFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create new Input from reader: %v", err)
	}
	return &Input{
		lines: lines,
	}, nil
}

func TrimSpaceAndNewLineCodeAndTab(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return r == ' ' || r == '\r' || r == '\n' || r == '\t'
	})
}

func toLines(scanner *bufio.Scanner) [][]string {
	var lines [][]string
	for scanner.Scan() {
		text := TrimSpaceAndNewLineCodeAndTab(scanner.Text())
		if len(text) == 0 {
			lines = append(lines, []string{})
			continue
		}
		line := strings.Split(text, " ")
		lines = append(lines, line)
	}
	return lines
}

func toLinesFromReader(reader *bufio.Reader) (lines [][]string, err error) {
	for {
		chunks, err := readLineAsChunks(reader)
		if err == io.EOF {
			return lines, nil
		}

		if err != nil {
			return nil, fmt.Errorf("failed to read line from reader: %v", err)
		}
		lineStr := TrimSpaceAndNewLineCodeAndTab(strings.Join(chunks, ""))
		line := strings.Split(lineStr, " ")
		lines = append(lines, line)
	}

}
func readLineAsChunks(reader *bufio.Reader) (chunks []string, err error) {
	for {
		chunk, isPrefix, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, string(chunk))
		if !isPrefix {
			return chunks, nil
		}
	}
}

func toSpecificBitIntLine(line []string, bitSize int) (intLine []int64, err error) {
	for j, v := range line {
		intV, err := strconv.ParseInt(v, 10, bitSize)
		if err != nil {
			return nil, fmt.Errorf(fmt.Sprintf("%dth value: %v", j, err.Error()))
		}
		intLine = append(intLine, intV)
	}
	return intLine, nil
}
