/*
 * gomacro - A Go intepreter with Lisp-like macros
 *
 * Copyright (C) 2017 Massimiliano Ghilardi
 *
 *     This program is free software you can redistribute it and/or modify
 *     it under the terms of the GNU General Public License as published by
 *     the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU General Public License for more details.
 *
 *     You should have received a copy of the GNU General Public License
 *     along with this program.  If not, see <http//www.gnu.org/licenses/>.
 *
 * switch.go
 *
 *  Created on May 06, 2017
 *      Author Massimiliano Ghilardi
 */

package fast

import (
	"go/ast"
	"go/token"
	r "reflect"
	"sort"
)

func (c *Comp) Switch(node *ast.SwitchStmt, labels []string) {
	initLocals := false
	var initBinds [2]int
	c, initLocals = c.pushEnvIfLocalBinds(&initBinds, node.Init)
	if node.Init != nil {
		c.Stmt(node.Init)
	}
	var ibreak int
	sort.Strings(labels)
	c.Loop = &LoopInfo{
		Break:      &ibreak,
		ThisLabels: labels,
	}

	// tag.Value (if constant) or tag.Fun() will return the tag value at runtime
	// cannot invoke e.Fun() multiple times because side effects must be applied only once!
	var e, tag *Expr
	enode := node.Tag
	if enode == nil {
		// "switch { }" without an expression means "switch true { }"
		e = exprValue(true)
		enode = &ast.Ident{NamePos: node.Pos() + 6, Name: "true"} // only for error messages
	} else {
		e = c.Expr1(enode)
	}
	if e.Const() {
		tag = e
	} else {
		tag = c.switchTag(e)
	}

	if node.Body != nil {
		list := node.Body.List
		defaulti := -1
		var defaultpos token.Pos
		seen := make(map[interface{}]token.Pos) // keep track of constant expressions in cases. error on duplicates
		n := len(list)
		for i, stmt := range list {
			switch clause := stmt.(type) {
			case *ast.CaseClause:
				canfallthrough := i < n-1 // last case cannot fallthrough
				if clause.List == nil {
					if defaulti >= 0 {
						c.Errorf("multiple defaults in switch (first at %s)", c.Fileset.Position(defaultpos))
					}
					defaulti = c.Code.Len()
					defaultpos = clause.Pos()
					c.switchDefault(clause, canfallthrough)
				} else {
					c.switchCase(clause, enode, tag, canfallthrough, seen)
				}
			default:
				c.Errorf("invalid statement inside switch: expecting case or default, found: %v <%v>", stmt, r.TypeOf(stmt))
			}
		}
		// default is executed as last, if no other case matches
		if defaulti >= 0 {
			// +1 to skip its "never matches" header
			c.Code.Append(func(env *Env) (Stmt, *Env) {
				ip := defaulti + 1
				env.IP = ip
				return env.Code[ip], env
			})
		}
	}

	// we finally know this
	ibreak = c.Code.Len()

	c = c.popEnvIfLocalBinds(initLocals, &initBinds, node.Init)
}

// compile a statement that evaluates e.Fun() once,
// and return an expression that returns cached e.Fun() result
func (c *Comp) switchTag(e *Expr) *Expr {
	efun := e.Fun
	var cachefun I
	var stmt Stmt
	switch e.Type.Kind() {
	case r.Bool:
		var val bool
		cachefun = func(*Env) bool {
			return val
		}
		efun := efun.(func(*Env) bool)
		stmt = func(env *Env) (Stmt, *Env) {
			val = efun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case r.Int:
		var val int
		cachefun = func(*Env) int {
			return val
		}
		efun := efun.(func(*Env) int)
		stmt = func(env *Env) (Stmt, *Env) {
			val = efun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case r.Int8:
		var val int8
		cachefun = func(*Env) int8 {
			return val
		}
		efun := efun.(func(*Env) int8)
		stmt = func(env *Env) (Stmt, *Env) {
			val = efun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case r.Int16:
		var val int16
		cachefun = func(*Env) int16 {
			return val
		}
		efun := efun.(func(*Env) int16)
		stmt = func(env *Env) (Stmt, *Env) {
			val = efun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case r.Int32:
		var val int32
		cachefun = func(*Env) int32 {
			return val
		}
		efun := efun.(func(*Env) int32)
		stmt = func(env *Env) (Stmt, *Env) {
			val = efun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case r.Int64:
		var val int64
		cachefun = func(*Env) int64 {
			return val
		}
		efun := efun.(func(*Env) int64)
		stmt = func(env *Env) (Stmt, *Env) {
			val = efun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case r.Uint:
		var val uint
		cachefun = func(*Env) uint {
			return val
		}
		efun := efun.(func(*Env) uint)
		stmt = func(env *Env) (Stmt, *Env) {
			val = efun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case r.Uint8:
		var val uint8
		cachefun = func(*Env) uint8 {
			return val
		}
		efun := efun.(func(*Env) uint8)
		stmt = func(env *Env) (Stmt, *Env) {
			val = efun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case r.Uint16:
		var val uint16
		cachefun = func(*Env) uint16 {
			return val
		}
		efun := efun.(func(*Env) uint16)
		stmt = func(env *Env) (Stmt, *Env) {
			val = efun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case r.Uint32:
		var val uint32
		cachefun = func(*Env) uint32 {
			return val
		}
		efun := efun.(func(*Env) uint32)
		stmt = func(env *Env) (Stmt, *Env) {
			val = efun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case r.Uint64:
		var val uint64
		cachefun = func(*Env) uint64 {
			return val
		}
		efun := efun.(func(*Env) uint64)
		stmt = func(env *Env) (Stmt, *Env) {
			val = efun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case r.Uintptr:
		var val uintptr
		cachefun = func(*Env) uintptr {
			return val
		}
		efun := efun.(func(*Env) uintptr)
		stmt = func(env *Env) (Stmt, *Env) {
			val = efun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case r.Float32:
		var val float32
		cachefun = func(*Env) float32 {
			return val
		}
		efun := efun.(func(*Env) float32)
		stmt = func(env *Env) (Stmt, *Env) {
			val = efun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case r.Float64:
		var val float64
		cachefun = func(*Env) float64 {
			return val
		}
		efun := efun.(func(*Env) float64)
		stmt = func(env *Env) (Stmt, *Env) {
			val = efun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case r.Complex64:
		var val complex64
		cachefun = func(*Env) complex64 {
			return val
		}
		efun := efun.(func(*Env) complex64)
		stmt = func(env *Env) (Stmt, *Env) {
			val = efun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case r.Complex128:
		var val complex128
		cachefun = func(*Env) complex128 {
			return val
		}
		efun := efun.(func(*Env) complex128)
		stmt = func(env *Env) (Stmt, *Env) {
			val = efun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	case r.String:
		var val string
		cachefun = func(*Env) string {
			return val
		}
		efun := efun.(func(*Env) string)
		stmt = func(env *Env) (Stmt, *Env) {
			val = efun(env)
			env.IP++
			return env.Code[env.IP], env
		}
	default:
		val := r.Zero(e.Type)
		cachefun = func(*Env) r.Value {
			return val
		}
		if efun, ok := efun.(func(*Env) (r.Value, []r.Value)); ok {
			stmt = func(env *Env) (Stmt, *Env) {
				val, _ = efun(env)
				env.IP++
				return env.Code[env.IP], env
			}
		} else {
			efun := e.AsX1()
			stmt = func(env *Env) (Stmt, *Env) {
				val = efun(env)
				env.IP++
				return env.Code[env.IP], env
			}
		}
	}
	c.Code.Append(stmt)
	return exprFun(e.Type, cachefun)
}

// switchDefault compiles a case in a switch.
func (c *Comp) switchCase(node *ast.CaseClause, tagnode ast.Expr, tag *Expr, canfallthrough bool, seen map[interface{}]token.Pos) {
	cmpfuns := make([]func(*Env) bool, 0)
	cmpnode := &ast.BinaryExpr{Op: token.EQL, X: tagnode} // only for error messages

	// compile a comparison of tag against each expression
	sometrue := false
	for _, enode := range node.List {
		e := c.Expr1(enode)
		cmpnode.OpPos = enode.Pos()
		cmpnode.Y = enode
		cmp := c.Eql(cmpnode, tag, e)
		if e.Const() {
			c.updateSeen(seen, e.Value, enode.Pos())
			if tag.Const() {
				// constant propagation
				flag := cmp.EvalConst(CompileDefaults)
				if r.ValueOf(flag).Bool() {
					sometrue = true
					break // always matches, no need to check further expressions
				} else {
					// can never match, skip this expression
					continue
				}
			}
		}
		// constants are handled above. only add non-constant comparisons to cmpfuns
		cmpfuns = append(cmpfuns, cmp.Fun.(func(*Env) bool))
	}
	// compile like "if tag == e1 || tag == e2 ... { }"
	// and keep track ourselves of where to jump if no expression matches
	//
	// always occupy a Code slot for cmpfuns, even if nothing to do.
	// reason: fallthrough from previous case skips such slot and jumps to current body
	var iend int
	var stmt Stmt
	switch len(cmpfuns) {
	case 0:
		if sometrue {
			stmt = stmtNop
		} else {
			// compile anyway, a fallthrough from previous case may still reach the current body
			stmt = func(env *Env) (Stmt, *Env) {
				ip := iend
				env.IP = ip
				return env.Code[ip], env
			}
		}
	case 1:
		cmpfun := cmpfuns[0]
		if sometrue {
			stmt = func(env *Env) (Stmt, *Env) {
				// keep side effects
				cmpfun(env)
				env.IP++
				return env.Code[env.IP], env
			}
		} else {
			stmt = func(env *Env) (Stmt, *Env) {
				var ip int
				if cmpfun(env) {
					ip = env.IP + 1
				} else {
					ip = iend
				}
				env.IP = ip
				return env.Code[ip], env
			}
		}
	case 2:
		cmpfuns := [...]func(*Env) bool{
			cmpfuns[0],
			cmpfuns[1],
		}
		if sometrue {
			stmt = func(env *Env) (Stmt, *Env) {
				// keep side effects
				_ = cmpfuns[0](env) || cmpfuns[1](env)
				env.IP++
				return env.Code[env.IP], env
			}
		} else {
			stmt = func(env *Env) (Stmt, *Env) {
				var ip int
				if cmpfuns[0](env) || cmpfuns[1](env) {
					ip = env.IP + 1
				} else {
					ip = iend
				}
				env.IP = ip
				return env.Code[ip], env
			}
		}
	default:
		if sometrue {
			stmt = func(env *Env) (Stmt, *Env) {
				for _, cmpfun := range cmpfuns {
					// keep side effects
					if cmpfun(env) {
						break
					}
				}
				env.IP++
				return env.Code[env.IP], env
			}
		} else {
			stmt = func(env *Env) (Stmt, *Env) {
				ip := iend
				for _, cmpfun := range cmpfuns {
					if cmpfun(env) {
						ip = env.IP + 1
						break
					}
				}
				env.IP = ip
				return env.Code[ip], env
			}
		}
	}
	c.Code.Append(stmt)
	c.switchCaseBody(node.Body, canfallthrough)
	// we finally know where to jump if match fails
	iend = c.Code.Len()
}

// switchDefault compiles the default case in a switch
func (c *Comp) switchDefault(node *ast.CaseClause, canfallthrough bool) {
	var iend int
	c.Code.Append(func(env *Env) (Stmt, *Env) {
		// jump to the next case. we must always add this statement for three reasons:
		// 1) if default is entered normally, it always fails to match and jumps to the next case
		// 2) if the previous case ends with fallthrough, it will skip this statement and jump to default's body
		// 3) if the switch ends without any matching case, it will manually jump to default's body (skipping this statement)
		ip := iend
		env.IP = ip
		return env.Code[ip], env
	})
	c.switchCaseBody(node.Body, canfallthrough)
	// we finally know where to jump if match fails
	iend = c.Code.Len()
}

// switchCaseBody compiles the body of a case in a switch
func (c *Comp) switchCaseBody(list []ast.Stmt, canfallthrough bool) {
	isfallthrough := false
	n := len(list)
	if n != 0 {
		isfallthrough = isFallthrough(list[n-1])
		if isfallthrough {
			if canfallthrough {
				n--
				list = list[:n]
			} else {
				c.Errorf("cannot fallthrough final case in switch")
				return
			}
		}

		// c.List creates a new scope... not accurate, compiled Go doesn't.
		// but at least isolates per-case variables, as compiled Go does
		if n != 0 {
			c.List(list)
		}
	}
	// after executing the case body, either break or fallthrough
	if isfallthrough {
		c.Code.Append(stmtFallthrough)
	} else {
		c.jumpOut(0, c.Loop.Break)
	}
}

// stmtFallThrough executes a fallthrough statement - only works inside a switch,
// and cannot be used in the last switch of a case
func stmtFallthrough(env *Env) (Stmt, *Env) {
	env.IP += 2 // +2 to skip the comparisons in next case, and jump directly to its body
	return env.Code[env.IP], env
}

func isFallthrough(node ast.Stmt) bool {
	switch node := node.(type) {
	case *ast.BranchStmt:
		return node.Tok == token.FALLTHROUGH
	default:
		return false
	}
}

// keep track of constant expressions in cases. error on duplicates
func (c *Comp) updateSeen(seen map[interface{}]token.Pos, val interface{}, pos token.Pos) {
	prev, found := seen[val]
	if found {
		c.Errorf("duplicate case %v <%v> in switch\n\tprevious case at %s", val, r.TypeOf(val), c.Fileset.Position(prev))
		return
	}
	seen[val] = pos
}