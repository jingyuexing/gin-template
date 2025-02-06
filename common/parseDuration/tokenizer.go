package parseduration

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/jingyuexing/go-utils"
)

type TokenKind uint

const (
	NumbericKind TokenKind = iota
	HourKind
	NanosecondKind
	MicrosecondKind
	MillisecondKind
	SecondKind
	MinuteKind
	DayKind
	WeekKind
	MonthKind
	YearKind
	UnknownKind
	EOFKind
)

func (t TokenKind) String() string {
	switch t {
	case NumbericKind:
		return "numeric"
	case HourKind:
		return "hour"
	case MicrosecondKind:
		return "microsecond"
	case MillisecondKind:
		return "millisecond"
	case NanosecondKind:
		return "nanosecond"
	case MinuteKind:
		return "minute"
	case MonthKind:
		return "month"
	case WeekKind:
		return "week"
	case UnknownKind:
		return "unknown"
	case EOFKind:
		return "eof"
	}
	return ""
}

// Token 表示解析得到的词法单元
type Token struct {
	Type  TokenKind
	Value string
}

// Lexer 词法分析器
type Lexer struct {
	input string
	pos   int
}

// NewLexer 初始化一个新的 Lexer
func NewLexer(input string) *Lexer {
	return &Lexer{
		input: strings.TrimSpace(input),
		pos:   0,
	}
}

// NextToken 获取下一个词法单元
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	// 如果已经到达输入的末尾，返回EOF
	if l.pos >= len(l.input) {
		return Token{Type: EOFKind, Value: ""}
	}

	// 解析数字
	if unicode.IsDigit(rune(l.input[l.pos])) {
		start := l.pos
		for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
			l.pos++
		}
		return Token{Type: NumbericKind, Value: l.input[start:l.pos]}
	}

	// 解析单位
	if unicode.IsLetter(rune(l.input[l.pos])) {
		start := l.pos
		for l.pos < len(l.input) && unicode.IsLetter(rune(l.input[l.pos])) {
			l.pos++
		}
		value := l.input[start:l.pos]
		kind := UnknownKind
		switch value {
		case "ns", "nanosecond":
			kind = NanosecondKind
		case "us", "microsecond":
			kind = MicrosecondKind
		case "s", "second":
			kind = SecondKind
		case "m", "minute":
			kind = MinuteKind
		case "h", "hour":
			kind = HourKind
		case "D", "d", "day":
			kind = DayKind
		case "W", "w", "week":
			kind = WeekKind
		case "M", "month":
			kind = MonthKind
		case "Y", "year":
			kind = YearKind
		default:
			kind = UnknownKind

		}
		return Token{Type: kind, Value: value}
	}

	return Token{Type: EOFKind, Value: ""}
}

// skipWhitespace 跳过空白字符
func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) && unicode.IsSpace(rune(l.input[l.pos])) {
		l.pos++
	}
}

// ParseDuration 使用 Lexer 解析时间字符串
func ParseDuration(input string) (time.Duration, error) {
	lexer := NewLexer(input)
	totalDuration := time.Duration(0)

	for {
		numToken := lexer.NextToken()
		if numToken.Type == EOFKind {
			break
		}

		if numToken.Type != NumbericKind {
			return 0, fmt.Errorf("expected a numeric value, got %s", numToken.Value)
		}

		value, err := strconv.Atoi(numToken.Value)
		if err != nil {
			return 0, fmt.Errorf("invalid numeric value: %v", err)
		}

		unitToken := lexer.NextToken()
		if unitToken.Type == UnknownKind {
			return 0, fmt.Errorf("expected a unit, got %s", unitToken.Value)
		}

		duration, err := unitToDuration(value, unitToken.Type)
		if err != nil {
			return 0, err
		}

		totalDuration += duration
	}

	return totalDuration, nil
}

// unitToDuration 将单位字符串转换为time.Duration
func unitToDuration(value int, kind TokenKind) (time.Duration, error) {
	switch kind {
	case NanosecondKind:
		return time.Duration(value) * time.Nanosecond, nil
	case MicrosecondKind:
		return time.Duration(value) * time.Microsecond, nil
	case MillisecondKind:
		return time.Duration(value) * time.Millisecond, nil
	case SecondKind:
		return time.Duration(value) * time.Second, nil
	case MinuteKind:
		return time.Duration(value) * time.Minute, nil
	case HourKind:
		return time.Duration(value) * time.Hour, nil
	case DayKind:
		return time.Duration(value) * 24 * time.Hour, nil
	case WeekKind:
		return time.Duration(value) * 7 * 24 * time.Hour, nil
	case MonthKind:
		date := utils.NewDateTime()
		return time.Duration(date.Add(value, "M").Time() - date.Time()), nil
	case YearKind:
		date := utils.NewDateTime()
		return time.Duration(date.Add(value, "Y").Time() - date.Time()), nil
	default:
		return 0, fmt.Errorf("unsupported time unit")
	}
}
