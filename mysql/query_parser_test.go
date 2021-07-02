package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseQuery(t *testing.T) {
	type args struct {
		query     string
		charset   string
		collation string
	}
	tests := []struct {
		name    string
		args    args
		want    *meta
		wantErr bool
	}{
		{
			name: "SelectStatement",
			args: args{
				query: "SELECT * FROM meta WHERE id = 1",
			},
			want: &meta{
				table:  "meta",
				method: stmtSelect,
			},
			wantErr: false,
		},
		{
			name: "SelectStatementWithJoin",
			args: args{
				query: "SELECT * FROM meta AS m, data AS d WHERE m.id = d.id",
			},
			want: &meta{
				table:  "(meta AS m) JOIN data AS d",
				method: stmtSelect,
			},
			wantErr: false,
		},
		{
			name: "SelectStatementWithSubSelect",
			args: args{
				query: "SELECT * FROM meta AS m WHERE m.id IN (SELECT id FROM data WHERE id < 100)",
			},
			want: &meta{
				table:  "meta AS m",
				method: stmtSelect,
			},
			wantErr: false,
		},
		{
			name: "SelectStatementWithLeftJoin",
			args: args{
				query: "SELECT * FROM meta m LEFT JOIN data d ON m.id = d.id",
			},
			want: &meta{
				table:  "meta AS m LEFT JOIN data AS d ON m.id=d.id",
				method: stmtSelect,
			},
			wantErr: false,
		},
		{
			name: "UpdateStatement",
			args: args{
				query: "UPDATE meta SET name = 'bob', version = 2 WHERE version = 1",
			},
			want: &meta{
				table:  "meta",
				method: stmtUpdate,
			},
			wantErr: false,
		},
		{
			name: "InsertStatement",
			args: args{
				query: "INSERT INTO meta (name, version) VALUES ('bob', 2)",
			},
			want: &meta{
				table:  "meta",
				method: stmtInsert,
			},
			wantErr: false,
		},
		{
			name: "DeleteStatement",
			args: args{
				query: "DELETE FROM meta WHERE id = 2",
			},
			want: &meta{
				table:  "meta",
				method: stmtDelete,
			},
			wantErr: false,
		},
		{
			name: "OtherStatement",
			args: args{
				query: "EXPLAIN DELETE FROM meta WHERE id = 2",
			},
			want: &meta{
				table:  tableOther,
				method: stmtOther,
			},
			wantErr: false,
		},
		{
			name: "BadStatement",
			args: args{
				query: "DELETE BRO meta WHERE id = 2",
			},
			want:    nil,
			wantErr: true,
		},
	}

	p := newSQLParser()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.parse(tt.args.query, tt.args.charset, tt.args.collation)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func BenchmarkParseQuery(b *testing.B) {
	query := `
SELECT id, name, last_name, middle_name, email, outer_uuid
FROM driver_devices dd
LEFT JOIN drivers dr ON dd.id_driver = dr.id
WHERE dd.token IN (?, ?)
	AND dd.active = 1 
	AND (dd.token_expire IS NULL OR dd.token_expire <= NOW())
`
	p := newSQLParser()

	var m *meta
	var err error

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m, err = p.parse(query, "", "")
		require.NoError(b, err)
	}

	require.NotNil(b, m)
}
