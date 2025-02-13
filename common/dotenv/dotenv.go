package dotenv

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"template/common/mapping"
	"unicode"

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
	// 修改字符判断逻辑，支持UTF-8
	return unicode.IsLetter(ch) || ch == '_' || ch == '-' || ch >= '\u0080'
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
	nested    bool // 是否启用嵌套结构解析
}

func New(text, delimiter string) *DotENV {
	d := &DotENV{
		delimiter: delimiter,
		env:       make(map[string]any),
		exports:   make(map[string]any),
		vars:      make(map[string]any),
		nested:    true, // 默认启用嵌套结构
	}
	d.SetContent(text)
	return d
}

// SetNested 设置是否启用嵌套结构解析
func (d *DotENV) SetNested(nested bool) {
	d.nested = nested
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

		// 处理空行和纯注释行
		if ch == '\n' {
			current++
			continue
		}

		// 跳过行首空白字符
		if isWhitespace(ch) {
			current++
			continue
		}

		// 优先处理注释
		if ch == TokenValues[HASH] {
			commentToken := d.parseComment(current)
			if len(commentToken.Value) > 0 {
				tokens = append(tokens, EnvToken{Kind: Hash, Value: string(ch)})
				tokens = append(tokens, commentToken)
			}
			current = d.current
			continue
		}

		// 其他token处理
		switch {
		case isNumeric(ch):
			tokens = append(tokens, d.parseNumber(current))
			current = d.current
		case isLetter(ch):
			tokens = append(tokens, d.parseKey(current))
			current = d.current
		case ch == TokenValues[DOUBLE_QUOTE] || ch == TokenValues[SINGLE_QUOTE]:
			tokens = append(tokens, d.parseString(ch, current)...)
			current = d.current
		case ch == TokenValues[RAW]:
			tokens = append(tokens, d.parseRawString(current))
			current = d.current
		case ch == TokenValues[LBRACKETS] || ch == TokenValues[RBRACKETS]:
			tokens = append(tokens, EnvToken{Kind: Brackets, Value: string(ch)})
			current++
		case ch == TokenValues[LBRACES]:
			tokens = append(tokens, d.parseJSON(current))
			current = d.current
		case ch == TokenValues[COMMAS]:
			tokens = append(tokens, EnvToken{Kind: Commas, Value: string(ch)})
			current++
		case ch == TokenValues[EQUAL]:
			tokens = append(tokens, EnvToken{Kind: Equal, Value: string(ch)})
			current++
		default:
			current++
		}
	}
	return tokens
}

// parseString 解析字符串并返回Token列表
func (d *DotENV) parseString(entry rune, current int) []EnvToken {
	d.current = current + 1 // 跳过开始的引号
	_content := []rune(d.content)
	var tokens []EnvToken
	value := ""
	escaped := false

	for d.current < d.maxLength {
		ch := _content[d.current]

		// 处理转义字符
		if escaped {
			switch ch {
			case 'n':
				value += "\n"
			case 't':
				value += "\t"
			case 'r':
				value += "\r"
			default:
				value += string(ch)
			}
			escaped = false
			d.current++
			continue
		}

		if ch == '\\' {
			escaped = true
			d.current++
			continue
		}

		// 处理字符串结束
		if ch == entry {
			d.current++ // 跳过结束的引号
			break
		}

		// 处理变量替换
		if ch == '$' && d.peekNext() == '{' {
			if value != "" {
				tokens = append(tokens, EnvToken{Kind: Text, Value: value})
				value = ""
			}
			varToken := d.parseVariable()
			tokens = append(tokens, varToken)
			continue
		}

		value += string(ch)
		d.current++
	}

	if value != "" {
		tokens = append(tokens, EnvToken{Kind: Text, Value: value})
	}

	return tokens
}

func (d *DotENV) parseVariable() EnvToken {
	d.current += 2 // 跳过 "${"
	varName := ""

	for d.current < d.maxLength {
		ch := rune(d.content[d.current])
		if ch == '}' {
			d.current++
			break
		}
		varName += string(ch)
		d.current++
	}

	// 查找变量值
	if val, ok := d.vars[varName]; ok {
		return EnvToken{Kind: Variable, Value: utils.ToString(val)}
	}

	// 如果找不到变量，返回占位符
	return EnvToken{Kind: Placeholder, Value: varName}
}

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
				d.setNestedValue(name, currentToken.Value)
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
				d.setNestedValue(name, value)
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
						d.setNestedValue(name, currentToken.Value) // 将后一个 key token 作为值
					} else if result == 1 {
						d.setNestedValue(name, true)
					} else {
						d.setNestedValue(name, false)
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
					d.setNestedValue(name, d.convertNumber(currentToken))
				} else {
					d.setNestedValue(name, d.convertInt(currentToken))
				}
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
				d.setNestedValue(name, currentToken.Value)
				name = ""
			}
		case JSON:
			if name != "" {
				d.setNestedValue(name, d.convertJSON(currentToken))
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
	d.current = current + 1 // 跳过 # 符号
	value := ""

	for d.current < d.maxLength {
		ch := rune(d.content[d.current])
		if ch == '\n' || ch == '\r' {
			break
		}
		value += string(ch)
		d.current++
	}

	// 处理行末
	if d.current < d.maxLength && rune(d.content[d.current]) == '\r' {
		d.current++
		if d.current < d.maxLength && rune(d.content[d.current]) == '\n' {
			d.current++
		}
	} else if d.current < d.maxLength && rune(d.content[d.current]) == '\n' {
		d.current++
	}

	return EnvToken{Kind: Comment, Value: strings.TrimSpace(value)}
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
	// 先检查导出变量
	if val, ok := d.exports[field]; ok {
		return val
	}

	// 检查嵌套结构
	if d.nested && strings.Contains(field, d.delimiter) {
		parts := strings.Split(field, d.delimiter)
		current := d.env

		for _, part := range parts {
			if next, ok := current[part]; ok {
				if mapVal, isMap := next.(map[string]any); isMap {
					current = mapVal
				} else {
					return next
				}
			} else {
				return nil
			}
		}
		return current
	}

	// 常规查找
	if val, ok := d.env[field]; ok {
		return val
	}

	return nil
}

func (d *DotENV) Load(text string) {
	d.parse(d.tokenize(text))
}

func (d *DotENV) String() string {
	var text []string

	var flatten func(prefix string, value any) []string
	flatten = func(prefix string, value any) []string {
		var result []string

		switch v := value.(type) {
		case map[string]any:
			for k, val := range v {
				newPrefix := k
				if prefix != "" {
					newPrefix = prefix + d.delimiter + k
				}
				result = append(result, flatten(newPrefix, val)...)
			}
		default:
			result = append(result, prefix+"="+fmt.Sprintf("%v", value))
		}

		return result
	}

	// 扁平化所有键值对
	text = flatten("", d.env)

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

// setNestedValue 设置嵌套值
func (d *DotENV) setNestedValue(key string, value any) {
	if !d.nested || !strings.Contains(key, d.delimiter) {
		d.env[key] = value
		return
	}

	parts := strings.Split(key, d.delimiter)
	current := d.env

	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		if _, exists := current[part]; !exists {
			current[part] = make(map[string]any)
		}

		if next, ok := current[part].(map[string]any); ok {
			current = next
		} else {
			// 如果当前节点不是map，创建新的map
			newMap := make(map[string]any)
			current[part] = newMap
			current = newMap
		}
	}

	current[parts[len(parts)-1]] = value
}
