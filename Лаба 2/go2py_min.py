from __future__ import annotations

# - type T struct { ... }
# - func (r T) Method(...) ...
# - :=, multi-assign, multi-return
# - for init;cond;post and for cond
# - if/else

from dataclasses import dataclass
import argparse
import sys


@dataclass(frozen=True)
class Tok:
    k: str
    v: str
    p: int


class Lex:
    KW = {"type", "struct", "func", "return", "if", "else", "for", "var", "true", "false", "package", "import"}
    TWO = {":=", "==", "!=", "<=", ">=", "&&", "||"}
    ONE = set("+-*/%<>=!.,;:(){}[]")

    def __init__(self, s: str) -> None:
        self.s, self.i = s, 0

    def _p(self, k: int = 0) -> str:
        j = self.i + k
        return self.s[j] if 0 <= j < len(self.s) else ""

    def _a(self, k: int = 1) -> None:
        self.i += k

    def _skip(self) -> None:
        while self.i < len(self.s):
            c = self._p()
            if c in " \t\r\n":
                self._a();
                continue
            if c == "/" and self._p(1) == "/":
                while self.i < len(self.s) and self._p() != "\n":
                    self._a()
                continue
            break

    def _read(self, pred) -> str:
        st = self.i
        while self.i < len(self.s) and pred(self._p()):
            self._a()
        return self.s[st:self.i]

    def _str(self) -> Tok:
        st = self.i
        self._a()  # "
        out: list[str] = []
        while self.i < len(self.s):
            c = self._p()
            if c == '"':
                self._a();
                break
            if c == "\\":
                self._a()
                esc = self._p()
                out.append({"n": "\n", "t": "\t", '"': '"', "\\": "\\"}.get(esc, esc))
                self._a()
            else:
                out.append(c)
                self._a()
        return Tok("STRING", "".join(out), st)

    def toks(self) -> list[Tok]:
        t: list[Tok] = []
        while True:
            self._skip()
            if self.i >= len(self.s):
                t.append(Tok("EOF", "", self.i))
                return t

            st = self.i
            c = self._p()
            if c.isalpha() or c == "_":
                w = self._read(lambda ch: ch.isalnum() or ch == "_")
                t.append(Tok("KW" if w in self.KW else "IDENT", w, st))
                continue
            if c.isdigit():
                n = self._read(lambda ch: ch.isdigit())
                t.append(Tok("INT", n, st))
                continue
            if c == '"':
                t.append(self._str())
                continue

            two = c + self._p(1)
            if two in self.TWO:
                t.append(Tok(two, two, st))
                self._a(2)
                continue
            if c in self.ONE:
                t.append(Tok(c, c, st))
                self._a()
                continue
            raise SyntaxError(f"Bad char {c!r} at {self.i}")


class Go2Py:
    PREC = {"||": 1, "&&": 2, "==": 3, "!=": 3, "<": 4, "<=": 4, ">": 4, ">=": 4, "+": 5, "-": 5, "*": 6, "/": 6, "%": 6}

    def __init__(self, toks: list[Tok]) -> None:
        self.t, self.i = toks, 0
        self.out: list[str] = []
        self.ind = 0
        self.structs: dict[str, list[str]] = {}
        self.methods: dict[str, list[list[str]]] = {}
        self.funcs: list[list[str]] = []
        self._recv: str | None = None

    def cur(self) -> Tok:
        return self.t[self.i]

    def is_(self, k: str, v: str | None = None) -> bool:
        tok = self.cur()
        return tok.k == k and (v is None or tok.v == v)

    def eat(self, k: str, v: str | None = None) -> Tok:
        tok = self.cur()
        if tok.k != k or (v is not None and tok.v != v):
            need = f"{k}:{v}" if v is not None else k
            raise SyntaxError(f"Need {need}, got {tok.k}:{tok.v!r} at {tok.p}")
        self.i += 1
        return tok

    def maybe(self, k: str, v: str | None = None) -> bool:
        if self.is_(k, v):
            self.i += 1
            return True
        return False

    def emit(self, s: str = "") -> None:
        self.out.append("    " * self.ind + s)

    def gen(self) -> str:
        while not self.is_("EOF"):
            if self.is_("KW", "package") or self.is_("KW", "import"):
                self._skip_decl()
                continue
            if self.is_("KW", "type"):
                self._struct_decl()
                continue
            if self.is_("KW", "func"):
                self._func_decl()
                continue
            raise SyntaxError(f"Unexpected {self.cur().k}:{self.cur().v!r} at {self.cur().p}")

        self.emit("from __future__ import annotations")
        self.emit("from dataclasses import dataclass")
        self.emit()

        # structs first (Go style)
        for name, lines in self.structs.items():
            self.out += lines
            ms = self.methods.get(name)
            if ms:
                for m in ms:
                    self.out += m
                    self.emit()
            self.emit()

        for f in self.funcs:
            self.out += f
            self.emit()

        if any(line.startswith("def main") for f in self.funcs for line in f):
            self.emit('if __name__ == "__main__":')
            self.ind += 1
            self.emit("main()")
            self.ind -= 1

        return "\n".join(self.out).rstrip() + "\n"

    def _skip_decl(self) -> None:
        self.eat("KW")
        if self.is_("IDENT") or self.is_("STRING"):
            self.i += 1
        if self.maybe("("):
            while not self.maybe(")") and not self.is_("EOF"):
                self.i += 1

    def _type_name(self) -> str:
        tok = self.cur()
        if tok.k in {"IDENT", "KW"}:
            return self.eat(tok.k).v
        raise SyntaxError(f"Need type at {tok.p}")

    def _struct_decl(self) -> None:
        self.eat("KW", "type")
        name = self.eat("IDENT").v
        self.eat("KW", "struct")
        self.eat("{")
        fields: list[tuple[str, str]] = []
        while not self.maybe("}"):
            fields.append((self.eat("IDENT").v, self._type_name()))
            self.maybe(";")

        lines: list[str] = ["@dataclass", f"class {name}:"]
        if not fields:
            lines.append("    pass")
        else:
            for fn, ft in fields:
                py_t = {"int": "int", "string": "str", "bool": "bool", "float32": "float", "float64": "float"}.get(ft, ft)
                lines.append(f"    {fn}: {py_t}")
        self.structs[name] = lines

    def _func_decl(self) -> None:
        self.eat("KW", "func")
        recv_name = recv_type = None
        if self.maybe("("):
            recv_name = self.eat("IDENT").v
            recv_type = self._type_name()
            self.eat(")")
        name = self.eat("IDENT").v

        self.eat("(")
        params: list[str] = []
        if not self.is_(")"):
            while True:
                pn = self.eat("IDENT").v
                _pt = self._type_name()
                params.append(pn)
                if not self.maybe(","):
                    break
        self.eat(")")

        if self.maybe("("):
            self._type_name()
            while self.maybe(","):
                self._type_name()
            self.eat(")")
        else:
            # single return type
            if self.is_("IDENT") or (self.is_("KW") and self.cur().v not in {"return", "if", "else", "for", "var", "func", "type", "struct"}):
                self._type_name()

        prev_recv = self._recv
        self._recv = recv_name

        header = [f"def {name}({', '.join((['self'] if recv_name else []) + params)}):"]
        body = self._block_lines()
        if len(body) == 1 and body[0].strip() == "":
            body = ["    pass"]

        block = header + body

        self._recv = prev_recv

        if recv_type:
            block = [("    " + ln) if ln else "" for ln in block]
            self.methods.setdefault(recv_type, []).append(block)
        else:
            self.funcs.append(block)

    def _block_lines(self) -> list[str]:
        self.eat("{")
        lines: list[str] = []
        while not self.maybe("}"):
            lines += self._stmt_lines(base_indent=1)
            self.maybe(";")
        if not lines:
            return ["    pass"]
        return lines

    def _stmt_lines(self, base_indent: int) -> list[str]:
        if self.is_("KW", "return"):
            self.eat("KW", "return")
            if self.is_("}") or self.is_(";"):
                return ["    " * base_indent + "return"]
            vals = [self._expr(0)[0]]
            while self.maybe(","):
                vals.append(self._expr(0)[0])
            return ["    " * base_indent + ("return " + ", ".join(vals))]

        if self.is_("KW", "if"):
            self.eat("KW", "if")
            cond, _ = self._expr(0)
            self.eat("{")
            then_lines: list[str] = []
            while not self.maybe("}"):
                then_lines += self._stmt_lines(base_indent + 1)
                self.maybe(";")
            if not then_lines:
                then_lines = ["    " * (base_indent + 1) + "pass"]
            out = ["    " * base_indent + f"if {cond}:"] + then_lines
            if self.is_("KW", "else"):
                self.eat("KW", "else")
                self.eat("{")
                else_lines: list[str] = []
                while not self.maybe("}"):
                    else_lines += self._stmt_lines(base_indent + 1)
                    self.maybe(";")
                if not else_lines:
                    else_lines = ["    " * (base_indent + 1) + "pass"]
                out += ["    " * base_indent + "else:"] + else_lines
            return out

        if self.is_("KW", "for"):
            self.eat("KW", "for")
            # for { ... }
            if self.is_("{"):
                self.eat("{")
                body: list[str] = []
                while not self.maybe("}"):
                    body += self._stmt_lines(base_indent + 1)
                    self.maybe(";")
                if not body:
                    body = ["    " * (base_indent + 1) + "pass"]
                return ["    " * base_indent + "while True:"] + body

            save = self.i
            first = self._for_piece_lines(base_indent)
            if self.maybe(";"):
                init_lines = first
                cond = "True"
                if not self.is_(";"):
                    cond, _ = self._expr(0)
                self.eat(";")
                post_lines: list[str] = []
                if not self.is_("{"):
                    post_lines = self._for_piece_lines(base_indent + 1)
                self.eat("{")
                body: list[str] = []
                while not self.maybe("}"):
                    body += self._stmt_lines(base_indent + 1)
                    self.maybe(";")
                if not body and not post_lines:
                    body = ["    " * (base_indent + 1) + "pass"]
                out = init_lines + ["    " * base_indent + f"while {cond}:"] + body
                if post_lines:
                    out += post_lines
                return out

            self.i = save
            cond, _ = self._expr(0)
            self.eat("{")
            body: list[str] = []
            while not self.maybe("}"):
                body += self._stmt_lines(base_indent + 1)
                self.maybe(";")
            if not body:
                body = ["    " * (base_indent + 1) + "pass"]
            return ["    " * base_indent + f"while {cond}:"] + body

        if self.is_("KW", "var"):
            self.eat("KW", "var")
            names = [self.eat("IDENT").v]
            while self.maybe(","):
                names.append(self.eat("IDENT").v)
            # optional type
            if not self.is_("=") and (self.is_("IDENT") or self.is_("KW")):
                self._type_name()
            rhs = "None"
            if self.maybe("="):
                rhs, _ = self._expr(0)
            return ["    " * base_indent + f"{', '.join(names)} = {rhs}"]

        if self.is_("IDENT"):
            save = self.i
            lvs = [self._lvalue()]
            while self.maybe(","):
                lvs.append(self._lvalue())
            if self.is_("=") or self.is_(":="):
                self.i += 1
                rhs, _ = self._expr(0)
                return ["    " * base_indent + f"{', '.join(lvs)} = {rhs}"]
            self.i = save

        expr, _ = self._expr(0)
        return ["    " * base_indent + expr]

    def _for_piece_lines(self, base_indent: int) -> list[str]:
        if self.is_("KW", "var"):
            return self._stmt_lines(base_indent=base_indent)
        if self.is_("IDENT"):
            save = self.i
            lvs = [self._lvalue()]
            while self.maybe(","):
                lvs.append(self._lvalue())
            if self.is_("=") or self.is_(":="):
                self.i += 1
                rhs, _ = self._expr(0)
                return ["    " * base_indent + f"{', '.join(lvs)} = {rhs}"]
            self.i = save
        expr, _ = self._expr(0)
        return ["    " * base_indent + expr]

    def _lvalue(self) -> str:
        base = self.eat("IDENT").v
        if self._recv and base == self._recv:
            base = "self"
        while self.maybe("."):
            base = base + "." + self.eat("IDENT").v
        return base

    def _expr(self, minp: int) -> tuple[str, int]:
        left, lp = self._unary()
        while True:
            op = self.cur().k
            prec = self.PREC.get(op)
            if prec is None or prec < minp:
                break
            self.i += 1
            right, rp = self._expr(prec + 1)
            pyop = "and" if op == "&&" else "or" if op == "||" else "//" if op == "/" else op
            expr = f"{left} {pyop} {right}"
            left, lp = self._paren(expr, prec, lp), prec
        return left, lp

    def _paren(self, s: str, prec: int, child_prec: int) -> str:
        return f"({s})" if child_prec and child_prec < prec else s

    def _unary(self) -> tuple[str, int]:
        if self.is_("!") or self.is_("-"):
            op = self.cur().k
            self.i += 1
            x, xp = self._unary()
            py = "not " if op == "!" else "-"
            return f"{py}{x}", 10
        return self._primary()

    def _primary(self) -> tuple[str, int]:
        tok = self.cur()
        if tok.k == "IDENT":
            name = self.eat("IDENT").v
            if self._recv and name == self._recv:
                name = "self"
            return self._postfix(name)
        if tok.k == "INT":
            return self.eat("INT").v, 0
        if tok.k == "STRING":
            return repr(self.eat("STRING").v), 0
        if tok.k == "KW" and tok.v in {"true", "false"}:
            return ("True" if self.eat("KW").v == "true" else "False"), 0
        if self.maybe("("):
            e, p = self._expr(0)
            self.eat(")")
            return self._postfix(e)
        raise SyntaxError(f"Bad expr {tok.k}:{tok.v!r} at {tok.p}")

    def _postfix(self, base: str) -> tuple[str, int]:
        while True:
            if self.maybe("."):
                base = base + "." + self.eat("IDENT").v
                continue
            if self.maybe("("):
                args: list[str] = []
                if not self.is_(")"):
                    a, _ = self._expr(0)
                    args.append(a)
                    while self.maybe(","):
                        a, _ = self._expr(0)
                        args.append(a)
                self.eat(")")

                # fmt.Println/Print/Printf -> print
                if base.startswith("fmt.") and base.split(".", 1)[1] in {"Println", "Print", "Printf"}:
                    base = f"print({', '.join(args)})"
                else:
                    base = f"{base}({', '.join(args)})"
                continue
            if self.maybe("{"):
                # Struct literal: T{A: 1, B: 2} -> T(A=1, B=2)
                fields: list[str] = []
                if not self.is_("}"):
                    while True:
                        k = self.eat("IDENT").v
                        self.eat(":")
                        v, _ = self._expr(0)
                        fields.append(f"{k}={v}")
                        if not self.maybe(","):
                            break
                self.eat("}")
                base = f"{base}({', '.join(fields)})"
                continue
            break
        return base, 0


def translate(src: str) -> str:
    return Go2Py(Lex(src).toks()).gen()


def main(argv: list[str] | None = None) -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("input")
    ap.add_argument("-o", "--output")
    ns = ap.parse_args(argv)
    text = open(ns.input, "r", encoding="utf-8").read()
    out = translate(text)
    if ns.output:
        open(ns.output, "w", encoding="utf-8").write(out)
    else:
        sys.stdout.write(out)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
