package dotenv

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"template/common/mapping"
	"github.com/jingyuexing/go-utils"
)

type Token int

const (
	EQUAL Token = iota
	LBRACKETS
	RBRACKETS
	DOUBLE_QUOTE
	LBRACES
	RBRACES
	LGROUP
	RAW
	RGROUP
	SINGLE_QUOTE
	COLON
	COMMAS
	HASH
)

const (
	ON      = "on"
	Off     = "off"
	True    = "true"
	False   = "false"
	Disable = "disable"
	Enable  = "enable"
	No      = "no"
	Yes     = "yes"
	N       = "n"
	Y       = "y"
	Allow   = "allow"
)

var TokenValues = map[Token]rune{
	EQUAL:        '=',
	LBRACKETS:    '[',
	RBRACKETS:    ']',
	DOUBLE_QUOTE: '"',
	LBRACES:      '{',
	RBRACES:      '}',
	LGROUP:       '(',
	RAW:          '`',
	RGROUP:       ')',
	SINGLE_QUOTE: '\'',
	COLON:        ':',
	COMMAS:       ',',
	HASH:         '#',
}

type stack struct {
	stack int
}

func (s *stack) Push(ch rune) {
	if isExtend(ch, '{', '[', '(', '<') {
		s.stack++
	}
}

func (s *stack) Pop(ch rune) {
	if isExtend(ch, '}', ']', ')', '>') {
		s.stack--
	}
}

func (s *stack) Balance() bool {
	return s.stack == 0
}

type TokenKind int

const (
	Text TokenKind = iota
	Number
	ENV
	Variable
	Placeholder
	Identifier
	Equal
	Comment
	Commas
	Array
	JSON
	Tuple
	EXPORT
	Brackets
	EOF
	Hash
)

type EnvToken struct {
	Kind  TokenKind
	Value string
}

func isNumeric(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch rune) bool {
	return InRange(ch, 'a', 'z') || InRange(ch, 'A', 'Z') || isExtend(ch, '+', '-', '*', '/', '!', '%', '~') || ch >= '\xff'
}

func InRange(ch rune, begin, end rune) bool {
	return ch >= begin && ch <= end
}

func isExtend(ch rune, char ...rune) bool {
	for _, ch_ := range char {
		if ch == ch_ {
			return true
		}
	}
	return false
}

func isWhitespace(ch rune) bool {
	return ch == '\n' || ch == ' ' || ch == '\r' || ch == '\t'
}

type DotENV struct {
	content   string
	current   int
	maxLength int
	env       map[string]any
	exports   map[string]any
	vars      map[string]any
	delimiter string
	tokens    []EnvToken
}

func New(text, delimiter string) *DotENV {
	d := &DotENV{
		delimiter: delimiter,
		env:       make(map[string]any),
	}
	d.SetContent(text)
	return d
}

func (d *DotENV) SetContent(text string) {
	d.content = text
}

func (d *DotENV) Parse() *DotENV {
	d.tokens = d.tokenize(d.content)
	d.parse(d.tokens)
	return d
}

func (d *DotENV) Inject(variables ...any) {
	for i := 0; i < len(variables); i += 2 {
		key, ok1 := variables[i].(string)
		value, ok2 := variables[i+1].(any)
		if ok1 && ok2 {
			d.vars[key] = value
		}
	}
}

func (d *DotENV) tokenize(text string) []EnvToken {
	var tokens []EnvToken
	_text := []rune(text)
	d.content = text
	d.maxLength = len(_text)
	current := 0
	length := len(_text)
	for current < length {
		ch := _text[current]
		if isNumeric(ch) {
			tokens = append(tokens, d.parseNumber(current))
			current = d.current
		} else if isLetter(ch) {
			tokens = append(tokens, d.parseKey(current))
			current = d.current
		} else if ch == TokenValues[HASH] {
			tokens = append(tokens, EnvToken{Kind: Hash, Value: string(ch)})
			current++
			tokens = append(tokens, d.parseComment(current))
			current = d.current
		} else if ch == TokenValues[DOUBLE_QUOTE] || ch == TokenValues[SINGLE_QUOTE] {
			tokens = append(tokens, d.parseString(ch, current)...)
			current = d.current
		} else if ch == TokenValues[RAW] {
			tokens = append(tokens, d.parseRawString(current))
			current = d.current
		} else if ch == TokenValues[LBRACKETS] || ch == TokenValues[RBRACKETS] {
			tokens = append(tokens, EnvToken{Kind: Brackets, Value: string(ch)})
			current++
		} else if ch == TokenValues[LBRACES] {
			tokens = append(tokens, d.parseJSON(current))
			current = d.current
		} else if ch == TokenValues[COMMAS] {
			tokens = append(tokens, EnvToken{Kind: Commas, Value: string(ch)})
			current++
		} else if ch == TokenValues[EQUAL] {
			tokens = append(tokens, EnvToken{Kind: Equal, Value: string(ch)})
			current++
		} else if isWhitespace(ch) {
			current++
		} else {
			current++
		}
	}
	d.tokens = tokens
	return d.tokens
}

// parseString 解析字符串并返回Token列表
func (d *DotENV) parseString(entry rune, current int) []EnvToken {
	d.current = current
	_content := []rune(d.content)
	var tokens []EnvToken
	value := ""
	inTemplate := false
	templateName := ""

	for d.current < d.maxLength {
		ch := _content[d.current]

		// 检查是否进入或退出模板变量
		if ch == '$' && d.peekNext() == '{' {
			if value != "" {
				tokens = append(tokens, EnvToken{Kind: Text, Value: value})
				value = ""
			}
			inTemplate = true
			d.current += 2 // 跳过 "${"
			templateName = ""
			continue
		} else if ch == '{' && d.peekNext() != '{' {
			if value != "" {
				tokens = append(tokens, EnvToken{Kind: Text, Value: value})
				value = ""
			}
			inTemplate = true
			d.current++ // 跳过 "{"
			templateName = ""
			continue
		}

		// 退出模板解析
		if inTemplate && ch == '}' {
			if val, ok := d.vars[templateName]; ok {
				tokens = append(tokens, EnvToken{Kind: Variable, Value: utils.ToString(val)})
			} else {
				tokens = append(tokens, EnvToken{Kind: Placeholder, Value: templateName})
			}
			inTemplate = false
			d.current++
			continue
		}

		// 解析模板变量名
		if inTemplate {
			if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
				templateName += string(ch)
			} else {
				// 如果出现未预期的字符，返回错误标记
				tokens = append(tokens, EnvToken{Kind: Hash, Value: templateName})
				inTemplate = false
			}
		} else {
			// 普通字符直接加入
			value += string(ch)
		}

		d.current++
	}

	// 处理剩余的普通文本
	if value != "" {
		tokens = append(tokens, EnvToken{Kind: Text, Value: value})
	}
	return tokens

}

// peekNext 查看下一个字符而不前进 pos
func (d *DotENV) peekNext() rune {
	if d.current+1 < d.maxLength {
		return rune(d.content[d.current+1])
	}
	return 0
}

func (d *DotENV) isNumberic(token EnvToken) bool {
	return token.Kind == Number
}

func (d *DotENV) isString(token EnvToken) bool {
	return token.Kind == Text
}

func (d *DotENV) isJSON(token EnvToken) bool {
	return token.Kind == JSON

}

func ToBoolean(token EnvToken) int {
	val := strings.ToLower(token.Value)
	switch val {
	case ON, True, Y, Yes, Allow, Enable:
		return 1
	case Off, False, N, No, Disable:
		return 0
	default:
		return 3
	}
}

func (d *DotENV) isValue(token EnvToken) bool {
	switch token.Kind {
	case Text, Number, JSON, Array, Tuple:
		return true
	default:
		return false
	}
}

func (d *DotENV) parse(tokens []EnvToken) {
	if len(tokens) == 0 {
		return
	}
	current := 0
	maxTokenLength := len(tokens)
	var tokenCache []EnvToken
	name := ""
	for current < maxTokenLength {
		currentToken := tokens[current]
		switch currentToken.Kind {
		case ENV:
			if name == "" {
				tokenCache = append(tokenCache, currentToken)
			} else {
				d.env[name] = currentToken.Value
				currentToken.Kind = Text
				name = ""
			}
		case Brackets:
			current++
			if current >= maxTokenLength {
				break
			}
			currentToken = tokens[current]
			var value []any
			for current < maxTokenLength && currentToken.Kind != Brackets {
				if currentToken.Kind == Number {
					value = append(value, d.convertNumber(currentToken))
				} else if currentToken.Kind == Commas {
					current++
					if current >= maxTokenLength {
						break
					}
					currentToken = tokens[current]
					continue
				} else {
					value = append(value, currentToken.Value)
				}
				current++
				if current < maxTokenLength {
					currentToken = tokens[current]
				}
			}
			if name != "" {
				d.env[name] = value
				name = ""
			}
		case Equal:
			if current+1 < maxTokenLength && tokens[current+1].Kind == ENV {
				current++ // 跳到 = 后的 ENV token
				currentToken = tokens[current]
				if name == "" && len(tokenCache) > 0 {
					result := ToBoolean(currentToken)
					name = tokenCache[len(tokenCache)-1].Value
					if result == 3 {
						d.env[name] = currentToken.Value // 将后一个 key token 作为值
					} else if result == 1 {
						d.env[name] = true
					} else {
						d.env[name] = false
					}
					name = ""
				}
			} else if len(tokenCache) > 0 {
				name = tokenCache[len(tokenCache)-1].Value
				tokenCache = tokenCache[:len(tokenCache)-1]
			}
		case Number:
			if name != "" {
				if strings.Contains(currentToken.Value, ".") {
					d.env[name] = d.convertNumber(currentToken)
				}
				d.env[name] = d.convertInt(currentToken)
				name = ""
			}
		case EXPORT:
			current++
			key_count := 0
			for current < maxTokenLength {
				currentToken = tokens[current]
				if currentToken.Kind == ENV {
					name = currentToken.Value
					current++
					if current < maxTokenLength {
						currentToken = tokens[current]
					}
					key_count++
					if key_count > 1 {
						break
					}
				} else if currentToken.Kind == EXPORT || currentToken.Kind == Equal {
					current++
				} else if d.isValue(currentToken) {
					if name != "" {
						switch currentToken.Kind {
						case JSON:
							d.exports[name] = d.convertJSON(currentToken)
						case Text:
							d.exports[name] = currentToken.Value
						case Number:
							if strings.Contains(currentToken.Value, ".") {
								d.exports[name] = d.convertNumber(currentToken)
							}

							d.exports[name] = d.convertInt(currentToken)
						}
					}
				}
				current++
			}
		case Text:
			if name != "" {
				d.env[name] = currentToken.Value
				name = ""
			}
		case JSON:
			if name != "" {
				d.env[name] = currentToken.Value // You may need to implement JSON parsing here
				name = ""
			}
		}
		current++
	}
}

func (d *DotENV) parseKey(current int) EnvToken {
	d.current = current
	value := ""
	_content := []rune(d.content)
	ch := _content[d.current]
	for d.current < d.maxLength && (isLetter(ch) || isNumeric(ch) || isExtend(ch, []rune(d.delimiter)...)) {
		value += string(ch)
		d.current++
		if d.current < d.maxLength {
			ch = _content[d.current]
		}
	}
	return EnvToken{Kind: ENV, Value: value}
}

func (d *DotENV) convertJSON(token EnvToken) map[string]any {
	m := make(map[string]any, 0)
	return m
}
func (d *DotENV) parseNumber(current int) EnvToken {
	d.current = current
	_content := []rune(d.content)
	ch := _content[d.current]
	value := ""
	dot := 0
	for d.current < d.maxLength && (isNumeric(ch) || ch == '.') {
		value += string(ch)
		if ch == '.' {
			dot++
		}
		d.current++
		if d.current < d.maxLength {
			ch = _content[d.current]
		}
	}
	if dot > 1 {
		return EnvToken{Kind: Text, Value: value}
	}
	return EnvToken{Kind: Number, Value: value}
}

func (d *DotENV) parseRawString(current int) EnvToken {
	d.current = current + 1
	_content := []rune(d.content)
	value := ""
	ch := _content[d.current]
	for d.current < d.maxLength && ch != TokenValues[RAW] {
		value += string(ch)
		d.current++
		if d.current < d.maxLength {
			ch = _content[d.current]
		}
	}
	d.current++
	return EnvToken{Kind: Text, Value: value}
}

func (d *DotENV) parseComment(current int) EnvToken {
	d.current = current
	value := ""
	ch := d.content[d.current]
	for d.current < d.maxLength && ch != '\n' {
		value += string(ch)
		d.current++
		if d.current < d.maxLength {
			ch = d.content[d.current]
		}
	}
	return EnvToken{Kind: Comment, Value: value}
}
func (d *DotENV) Environment(env map[string]string) {
	allENV := d.flattenNestedDict(d.env)
	for name := range allENV {
		os.Setenv(name, string(env[name]))
	}
}
func (d *DotENV) convertInt(token EnvToken) int64 {
	value, _ := strconv.ParseInt(token.Value, 10, 64)
	return int64(value)
}

func (d *DotENV) convertNumber(token EnvToken) float64 {
	if strings.Contains(token.Value, ".") {
		value, _ := strconv.ParseFloat(token.Value, 64)
		return value
	}
	return -1
}

func (d *DotENV) parseJSON(current int) EnvToken {
	d.current = current
	_content := []rune(d.content)
	stack := &stack{}
	ch := _content[d.current]
	value := ""
	for d.current < d.maxLength {
		value += string(ch)
		if ch == TokenValues[LBRACES] {
			stack.Push(ch)
		}
		if ch == TokenValues[RBRACES] {
			stack.Pop(ch)
		}
		if stack.Balance() {
			break
		}
		d.current++
		if d.current < d.maxLength {
			ch = _content[d.current]
		}
	}
	d.current++
	return EnvToken{Kind: JSON, Value: value}
}

func (d *DotENV) Bind(config any) error {
	return mapping.BindMapToStruct(d.env, config, func(sf reflect.StructField) string {
		return sf.Tag.Get("env")
	})
}
func (d *DotENV) Get(field string) any {
	val, ok := d.exports[field]
	if !ok {
		val, ok = d.env[field]
	}
	if !ok {
		return ""
	}
	return val
}

func (d *DotENV) Load(text string) {
	d.parse(d.tokenize(text))
}

func (d *DotENV) String() string {
	var text []string
	for key, value := range d.env {
		if mapValue, ok := value.(map[string]any); ok {
			temp := map[string]any{key: mapValue}
			flattened := d.flattenNestedDict(temp)
			for flattenKey, flattenValue := range flattened {
				text = append(text, flattenKey+"="+fmt.Sprintf("%v", flattenValue))
			}
		} else {
			text = append(text, key+"="+fmt.Sprintf("%v", value))
		}
	}
	return strings.Join(text, "\n")
}

func (d *DotENV) flattenNestedDict(nestedDict map[string]any) map[string]any {
	flatDict := make(map[string]any)
	for k, v := range nestedDict {
		switch child := v.(type) {
		case map[string]any:
			flattenedChild := d.flattenNestedDict(child)
			for childKey, childValue := range flattenedChild {
				flatDict[k+d.delimiter+childKey] = childValue
			}
		default:
			flatDict[k] = child
		}
	}
	return flatDict
}
