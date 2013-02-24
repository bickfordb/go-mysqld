
package sql
import "fmt"
import "regexp"
import "strconv"
import "unicode"

type Statement interface {
  IsStatement()
}

func (s Select) IsStatement() { }
func (s Delete) IsStatement() { }

type Delete struct {}
type Update struct {}
type Insert struct {}


type Select struct {
  Columns []Expr
  From string
  Where Expr
}

type BinaryExpr struct {
  op string
  left Expr
  right Expr
}

func (f BinaryExpr) IsExpr() { }

type FuncExpr struct {
  fn string
  args []Expr
}

func (f FuncExpr) IsExpr() { }

type Expr interface {
  IsExpr()
}

type String string
func (s String) IsExpr() { }

type Integer int64
func (s Integer) IsExpr() { }

type Double float64
func (s Double) IsExpr() { }

type UnaryExpr struct {
  op string
  val Expr
}
func (s UnaryExpr) IsExpr() { }


type Lexer struct {
  buf string
  pos int
  lastPos int
  lastError error
  Statement Statement
}

var keywords = map[string]int{
  "select": SELECT,
  "insert": INSERT,
  "delete": DELETE,
  "where": WHERE,
  "from": FROM,
  "group": GROUP,
  "by": BY,
  "set": SET,
  "to": TO,
  "having": HAVING,
  "with": WITH,
  "rollup": ROLLUP,
  "[,]": COMMA,
  "-|[/+|%*&^<>]|and|or|>=|<=|=|!=|<>": BINOP,
  "[@~]": UNOP}

type Pattern struct {
  p *regexp.Regexp
  sym int
}

var intPat *regexp.Regexp = regexp.MustCompile(`^[\-+]?\d+\b`)
var doublePat *regexp.Regexp = regexp.MustCompile(`^[\-+]?\d+.\d+\b`)
var stringPat *regexp.Regexp = regexp.MustCompile(`^'([^']+)'\b`)
var patterns []*Pattern

func pat(s string, sym int) {
  s = "^(?i)" + s + `\b`
  patterns = append(patterns, &Pattern{p: regexp.MustCompile(s), sym:sym})
}

func init() {
  for s, i := range keywords {
    pat(s, i)
  }
}

func (l *Lexer) Lex(lval *yySymType) int {
  for {
    if l.pos == len(l.buf) {
      // EOF
      return 0
    }
    c := rune(l.buf[l.pos])
    if unicode.IsSpace(c) {
      l.pos += 1
      continue
    }
    remaining := l.buf[l.pos:]

    if m := intPat.FindString(remaining); m != "" {
      l.lastPos = l.pos
      l.pos += len(m)
      lval.Integer, _ = strconv.ParseInt(m, 10, 64)
      return INTEGER
    }

    if m := doublePat.FindString(remaining); m != "" {
      l.lastPos = l.pos
      l.pos += len(m)
      lval.Double, _ = strconv.ParseFloat(m, 64)
      return DOUBLE
    }

    for _, p := range patterns {
      m := p.p.FindString(remaining)
      if m != "" {
        l.lastPos = l.pos
        l.pos += len(m)
        lval.Keyword = m
        return p.sym
      }
    }
    return yyErrCode
  }
  return yyErrCode
}

type ParseError struct {
  e string
  pos int
}

func (p *ParseError) Error() string {
  return fmt.Sprintf("error: %s at: %d", p.e, p.pos)
}

func (l *Lexer) Error(e string) {
  l.lastError = &ParseError{
    pos: l.pos,
    e: e}
}

func Parse(s string) (stmt Statement, err error) {
  lexer := &Lexer{buf:s}
  ret := yyParse(lexer)
  if lexer.lastError != nil {
    err = lexer.lastError
  } else if ret != 0 {

  } else {
    stmt = lexer.Statement
  }
  return
}

