//line src/sql/parser.y:2
package sql

import "fmt"

//line src/sql/parser.y:9
type yySymType struct {
	yys       int
	Statement Statement
	Select    Select
	ColList   []Expr
	Expr      Expr
	Keyword   string
	Double    float64
	String    string
	Integer   int64
}

const SELECT = 57346
const DELETE = 57347
const INSERT = 57348
const UPDATE = 57349
const WHERE = 57350
const FROM = 57351
const WITH = 57352
const ROLLUP = 57353
const GROUP = 57354
const BY = 57355
const SET = 57356
const OPEN = 57357
const CLOSE = 57358
const BINOP = 57359
const UNOP = 57360
const TO = 57361
const HAVING = 57362
const COMMA = 57363
const EOF = 57364
const INTEGER = 57365
const IDENT = 57366
const DOUBLE = 57367
const STRING = 57368

var yyToknames = []string{
	"SELECT",
	"DELETE",
	"INSERT",
	"UPDATE",
	"WHERE",
	"FROM",
	"WITH",
	"ROLLUP",
	"GROUP",
	"BY",
	"SET",
	"OPEN",
	"CLOSE",
	"BINOP",
	"UNOP",
	"TO",
	"HAVING",
	"COMMA",
	"EOF",
	"INTEGER",
	"IDENT",
	"DOUBLE",
	"STRING",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line src/sql/parser.y:54
//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 11
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 18

var yyAct = []int{

	9, 11, 5, 10, 17, 12, 12, 3, 7, 2,
	6, 8, 13, 14, 15, 16, 4, 1,
}
var yyPact = []int{

	3, -1000, -1000, -15, -20, -11, -1000, -1000, -1000, -15,
	-15, -15, -15, -12, -11, -11, -11, -1000,
}
var yyPgo = []int{

	0, 17, 16, 9, 2,
}
var yyR1 = []int{

	0, 1, 3, 2, 2, 4, 4, 4, 4, 4,
	4,
}
var yyR2 = []int{

	0, 1, 2, 1, 3, 1, 1, 1, 3, 3,
	2,
}
var yyChk = []int{

	-1000, -1, -3, 4, -2, -4, 25, 23, 26, 15,
	18, 21, 17, -4, -4, -4, -4, 16,
}
var yyDef = []int{

	0, -2, 1, 0, 2, 3, 5, 6, 7, 0,
	0, 0, 0, 0, 10, 4, 9, 8,
}
var yyTok1 = []int{

	1,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c > 0 && c <= len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return fmt.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return fmt.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		fmt.Printf("lex %U %s\n", uint(char), yyTokname(c))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		fmt.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				fmt.Printf("%s", yyStatname(yystate))
				fmt.Printf("saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					fmt.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				fmt.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		fmt.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		//line src/sql/parser.y:33
		{
			l, _ := yylex.(*Lexer)
			l.Statement = yyS[yypt-0].Select
			yyVAL.Statement = yyS[yypt-0].Select
		}
	case 2:
		//line src/sql/parser.y:39
		{
			yyVAL.Select = Select{Columns: yyS[yypt-0].ColList}
		}
	case 3:
		//line src/sql/parser.y:42
		{
			yyVAL.ColList = []Expr{yyS[yypt-0].Expr}
		}
	case 4:
		//line src/sql/parser.y:43
		{
			yyVAL.ColList = append(yyS[yypt-2].ColList, yyS[yypt-0].Expr)
		}
	case 5:
		//line src/sql/parser.y:46
		{
			yyVAL.Expr = Double(yyS[yypt-0].Double)
		}
	case 6:
		//line src/sql/parser.y:47
		{
			yyVAL.Expr = Integer(yyS[yypt-0].Integer)
		}
	case 7:
		//line src/sql/parser.y:48
		{
			yyVAL.Expr = String(yyS[yypt-0].String)
		}
	case 8:
		//line src/sql/parser.y:49
		{
			yyVAL.Expr = yyS[yypt-1].Expr
		}
	case 9:
		//line src/sql/parser.y:50
		{
			yyVAL.Expr = BinaryExpr{yyS[yypt-1].Keyword, yyS[yypt-2].Expr, yyS[yypt-0].Expr}
		}
	case 10:
		//line src/sql/parser.y:51
		{
			yyVAL.Expr = UnaryExpr{yyS[yypt-1].Keyword, yyS[yypt-0].Expr}
		}
	}
	goto yystack /* stack new state and value */
}
