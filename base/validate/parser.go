package validate

import (
	"bytes"
	"errors"
)

type expression struct {
	Name string
	Args []string
}

var (
	ErrParseFailed = errors.New("Parse Expression Failed")
)

func (vo expression) IsEmpty() bool {
	return vo.Name == `` && (vo.Args == nil || len(vo.Args) == 0)
}

func parseExpressions(options string) ([]expression, error) {
	var (
		exp         expression //当前处理中的Option
		bracketOpen bool       //括号模式开启
		quoteOn     bool       //引号开启
		escapeOn    bool       //转义模式

		buf bytes.Buffer
	)
	runes := []rune(options)
	retval := make([]expression, 0, 5)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '(' { //开始解析参数
			name := buf.String()
			if name == `` { //已经开始参数解析，但是名称为空
				return nil, ErrParseFailed
			}
			exp.Name = name
			buf.Reset()

			bracketOpen = true
		} else if r == ' ' || r == '\t' { //除引号模式中的空白外，一律忽略
			if quoteOn {
				buf.WriteRune(r)
			}
		} else if r == '\'' {
			if !bracketOpen {
				return nil, ErrParseFailed //只有 参数解析中才允许出现引号
			}
			if !quoteOn {
				if buf.Len() > 0 {
					return nil, ErrParseFailed //引号必须在解析开始时开启
				}
				quoteOn = true
			} else { //关闭引号
				if escapeOn {
					buf.WriteRune(r)
					escapeOn = false
				} else {
					quoteOn = false
				}
			}
		} else if r == '\\' { //转义符号处理
			if quoteOn {
				escapeOn = true
			} else {
				return nil, ErrParseFailed //如果没有开启引号模式，则转义无意义
			}
		} else if r == ',' {
			if quoteOn { //引号模式中，原样输出
				buf.WriteRune(r)
			} else {
				if !bracketOpen {
					if buf.Len() > 0 {
						exp.Name = buf.String()
						retval = append(retval, exp)
					}
					//reset
					exp.Name = ``
					exp.Args = nil
					bracketOpen = false
					quoteOn = false
					escapeOn = false
					buf.Reset()
				} else { //参数结束
					exp.Args = append(exp.Args, buf.String())
					buf.Reset()
				}
			}
		} else if r == ')' {
			if quoteOn {
				buf.WriteRune(r)
			} else {
				if !bracketOpen {
					return nil, ErrParseFailed
				}
				if buf.Len() > 0 {
					exp.Args = append(exp.Args, buf.String())
				}
				if !exp.IsEmpty() {
					retval = append(retval, exp)
				}
				//reset
				exp.Name = ``
				exp.Args = nil
				bracketOpen = false
				escapeOn = false
				quoteOn = false
				buf.Reset()
			}
		} else {
			buf.WriteRune(r)
		}
	}
	if bracketOpen || escapeOn || quoteOn { //尾部不干净
		return nil, ErrParseFailed
	}
	if buf.Len() > 0 {
		exp.Name = buf.String()
	}
	if !exp.IsEmpty() {
		retval = append(retval, exp)
	}
	if len(retval) == 0 { //统一返回值，当没有内容被解析时，统一返回nil
		retval = nil
	}
	return retval, nil
}
