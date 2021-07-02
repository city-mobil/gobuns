package mysql

import (
	"bytes"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	_ "github.com/pingcap/parser/test_driver"
)

type stmtMethod string

const (
	stmtSelect stmtMethod = "select"
	stmtUpdate stmtMethod = "update"
	stmtInsert stmtMethod = "insert"
	stmtDelete stmtMethod = "delete"
	stmtOther  stmtMethod = "other"
)

const (
	tableOther = "_"
)

type meta struct {
	table  string
	method stmtMethod
}

type sqlParser struct {
	impl *parser.Parser
}

func newSQLParser() *sqlParser {
	return &sqlParser{
		impl: parser.New(),
	}
}

func (p *sqlParser) parse(query, charset, collation string) (*meta, error) {
	node, err := p.impl.ParseOneStmt(query, charset, collation)
	if err != nil {
		return nil, err
	}

	var m meta

	switch st := node.(type) {
	case *ast.SelectStmt:
		t, err := p.extractTable(st.From)
		if err != nil {
			return nil, err
		}
		m.table = t
		m.method = stmtSelect
	case *ast.UpdateStmt:
		t, err := p.extractTable(st.TableRefs)
		if err != nil {
			return nil, err
		}
		m.table = t
		m.method = stmtUpdate
	case *ast.DeleteStmt:
		t, err := p.extractTable(st.TableRefs)
		if err != nil {
			return nil, err
		}
		m.table = t
		m.method = stmtDelete
	case *ast.InsertStmt:
		t, err := p.extractTable(st.Table)
		if err != nil {
			return nil, err
		}
		m.table = t
		m.method = stmtInsert
	default:
		m.table = tableOther
		m.method = stmtOther
	}

	return &m, nil
}

func (p *sqlParser) extractTable(n ast.Node) (string, error) {
	var buf bytes.Buffer
	err := n.Restore(format.NewRestoreCtx(format.RestoreNameLowercase|format.RestoreKeyWordUppercase, &buf))
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
