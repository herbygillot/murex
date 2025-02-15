package parameters

import (
	"errors"
	"strconv"
	"strings"

	"github.com/lmorg/murex/lang/types"
)

var errTooFew = errors.New("too few parameters")

// Byte gets a single parameter as a []byte slice
func (p *Parameters) Byte(pos int) ([]byte, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.Len() <= pos {
		return []byte{}, errTooFew
	}
	return []byte(p.params[pos]), nil
}

// ByteAll returns all parameters as one space-delimited []byte slice
func (p *Parameters) ByteAll() []byte {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return []byte(strings.Join(p.params, " "))
}

// ByteAllRange returns all parameters within range as one space-delimited []byte slice
// `start` is first point in array. `end` is last. Set `end` to `-1` if you want `[n:]`.
func (p *Parameters) ByteAllRange(start, end int) []byte {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if end == -1 {
		return []byte(strings.Join(p.params[start:], " "))
	}
	return []byte(strings.Join(p.params[start:end], " "))
}

// String gets a single parameter as string
func (p *Parameters) String(pos int) (string, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.Len() <= pos {
		return "", errTooFew
	}
	return p.params[pos], nil
}

// StringArray returns all parameters as a slice of strings
func (p *Parameters) StringArray() []string {
	p.mutex.RLock()

	params := make([]string, len(p.params))
	copy(params, p.params)

	p.mutex.RUnlock()

	return params
}

// StringAll returns all parameters as one space-delimited string
func (p *Parameters) StringAll() string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return strings.Join(p.params, " ")
}

// StringAllRange returns all parameters within range as one space-delimited string.
// `start` is first point in array. `end` is last. Set `end` to `-1` if you want `[n:]`.
func (p *Parameters) StringAllRange(start, end int) string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	switch {
	case len(p.params) == 0:
		return ""
	case end == -1:
		return strings.Join(p.params[start:], " ")
	default:
		return strings.Join(p.params[start:end], " ")
	}
}

// Int gets parameter as integer
func (p *Parameters) Int(pos int) (int, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.Len() <= pos {
		return 0, errTooFew
	}
	return strconv.Atoi(p.params[pos])
}

// Uint32 gets parameter as Uint32
func (p *Parameters) Uint32(pos int) (uint32, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.Len() <= pos {
		return 0, errTooFew
	}
	i, err := strconv.ParseUint(p.params[pos], 10, 32)
	return uint32(i), err
}

// Bool gets parameter as boolean
func (p *Parameters) Bool(pos int) (bool, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.Len() <= pos {
		return false, errTooFew
	}
	return types.IsTrue([]byte(p.params[pos]), 0), nil
}

// Block get parameter as a code block or JSON block
func (p *Parameters) Block(pos int) ([]rune, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	switch {
	case p.Len() <= pos:
		return []rune{}, errTooFew

	case len(p.params[pos]) < 2:
		return []rune{}, errors.New("not a valid code block. Too few characters. Code blocks should be surrounded by curly brace, eg `{ code }`")

	case p.params[pos][0] != '{':
		return []rune{}, errors.New("not a valid code block. Missing opening curly brace. Found `" + string([]byte{p.params[pos][0]}) + "` instead.")

	case p.params[pos][len(p.params[pos])-1] != '}':
		return []rune{}, errors.New("not a valid code block. Missing closing curly brace. Found `" + string([]byte{p.params[pos][len(p.params[pos])-1]}) + "` instead.")

	}

	return []rune(p.params[pos][1 : len(p.params[pos])-1]), nil
}

// Len returns the number of parameters
func (p *Parameters) Len() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return len(p.params)
}

// TokenLen returns the number of parameters from ParamTokens. This method is
// not recommended for casual use - instead use Len() - however if you are
// monitoring other processes and those processes might be in a state prior to
// execution then you might want to check the parameter length before it has
// been parsed. This is this methods _only_ use case.
func (p *Parameters) TokenLen() int {
	var count int
	for _, tokens := range p.Tokens {
		for i := range tokens {
			if tokens[i].Type != TokenTypeNil && tokens[i].Type != TokenTypeNamedPipe {
				count++
				break
			}
		}
	}

	return count
}
