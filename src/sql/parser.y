%{

package sql

import "fmt"

%}

%union {
  Statement Statement
  Select Select
  ColList []Expr
  Expr Expr
  Keyword string
  Double float64
  String string
  Integer int64
}

%token <Keyword> SELECT DELETE INSERT UPDATE WHERE FROM WITH ROLLUP GROUP BY SET OPEN CLOSE BINOP UNOP TO HAVING COMMA EOF
%token <Integer> INTEGER
%token <Ident> IDENT
%token <Double> DOUBLE
%token <String> STRING

%type <Statement> statement
%type <ColList> cols
%type <Select> select
%type <Expr> expr

%%
statement: select
         {
           l, _ := yylex.(*Lexer)
           l.Statement = $1
           $$ = $1
         }

select: SELECT cols { $$ = Select{Columns: $2} }
;

cols: expr { $$ = []Expr{$1} }
  | cols COMMA expr { $$ = append($1, $3) }
;

expr: DOUBLE { $$ = Double($1) }
 | INTEGER { $$ = Integer($1) }
 | STRING { $$ = String($1) }
 | OPEN expr CLOSE { $$ = $2 }
 | expr BINOP expr { $$ = BinaryExpr{$2, $1, $3}}
 | UNOP expr { $$ = UnaryExpr{$1, $2} }
;

%%


