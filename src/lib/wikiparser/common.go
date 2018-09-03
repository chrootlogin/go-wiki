package wikiparser

import (
	"bufio"
	"bytes"
	"errors"
)

const (
	EOF = rune(0)
	CONTROL_START = rune(123) // {
	CONTROL_END = rune(125) // }
	EXPR_BEGIN = "{{"
	EXPR_END = "}}"
	ARGS_START = rune(40) // )
	ARGS_STOP = rune(41) // (
	ARGS_DELIMITER = rune(44) // ,
)

type Expression struct {
	Action string
	Args []string
}

type wikiparser struct {
	r *bufio.Reader
	Expressions []Expression
}

func New(data []byte) *wikiparser {
	r := bytes.NewBuffer(data)

	return &wikiparser{
		r: bufio.NewReader(r),
	}
}

func (wp *wikiparser) Parse() {
	wp.scan()
}

func (wp *wikiparser) scan() {
	// Read the next rune.
	ch := wp.read()

	if ch == CONTROL_START {
		wp.unread()
		wp.scanBegin()
	}

	// if end of file, abort
	if ch == EOF {
		return
	}

	// scan next rune
	wp.scan()
}

func (wp *wikiparser) scanBegin() {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer

	// start reading begin
	for i := 0; i < 2; i++ {
		ch := wp.read()

		// EOF stop
		if ch == EOF {
			break
		}

		buf.WriteRune(ch)
	}

	if buf.String() == EXPR_BEGIN {
		wp.scanExpression()
	}
}

func (wp *wikiparser) scanExpression() {
	var buf bytes.Buffer

	for {
		ch := wp.read()

		// stop on eof
		if ch == EOF {
			return
		}

		buf.WriteRune(ch)

		if ch == CONTROL_END && buf.Len() > 2 {
			s := buf.String()
			end := string(s[len(s) - 2:])

			if end == EXPR_END {
				// remove expr end
				buf.Truncate(buf.Len() - 2)

				// parse expression
				expr, err := parseExpression(bytes.TrimSpace(buf.Bytes()))
				if err != nil {
					return
				}

				// finish
				wp.Expressions = append(wp.Expressions, *expr)
				return
			}
		}
	}
}

func parseExpression(expr []byte) (*Expression, error) {
	buf := bytes.NewBuffer(expr)

	var action bytes.Buffer
	var args []string

	for {
		r, _, err := buf.ReadRune()
		// stop on end
		if err != nil {
			break
		}

		if isLetter(r) {
			action.WriteRune(r)
		}

		if r == ARGS_START {
			args, err = parseArgs(buf)
			if err != nil {
				return nil, err
			}
		}
	}

	return &Expression{
		Action: action.String(),
		Args: args,
	}, nil
}

func parseArgs(buf *bytes.Buffer) ([]string, error) {
	args := make([]string, 1)

	i := 0
	for {
		r, _, err := buf.ReadRune()
		// stop on error
		if err != nil {
			break
		}

		if r == ARGS_STOP {
			return args, nil
		}

		if r == ARGS_DELIMITER {
			args = append(args, "")
			i++
			continue
		}

		args[i] += string(r)
	}

	return nil, errors.New("parsing error")
}


func (wp *wikiparser) read() rune {
	ch, _, err := wp.r.ReadRune()
	if err != nil {
		return EOF
	}

	return ch
}

// unread places the previously read rune back on the reader.
func (wp *wikiparser) unread() { _ = wp.r.UnreadRune() }

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}