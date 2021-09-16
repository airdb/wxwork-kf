package schema

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type ModelMeta struct {
	ID        ModelPk `gorm:"type:varchar(36);primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// ModelPk 主键
// 数据中存 bit 还是 string
type ModelPk string

func (pk *ModelPk) renew() {
	*pk = ModelPk(uuid.New().String())
}

func (ModelPk) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{ModelPkCreateClause{Field: f}}
}

type ModelPkCreateClause struct {
	Field *schema.Field
}

func (v ModelPkCreateClause) Name() string {
	return ""
}

func (v ModelPkCreateClause) Build(clause.Builder) {
}

func (v ModelPkCreateClause) MergeClause(*clause.Clause) {
}

func (v ModelPkCreateClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.String() == "" {
		// create new value if empty
		if cv, zero := v.Field.ValueOf(stmt.ReflectValue); !zero {
			if cvv, ok := cv.(ModelPk); ok {
				log.Println(cvv)
				return
			}
		}

		nv := uuid.New().String()
		stmt.AddClause(clause.Set{{Column: clause.Column{Name: v.Field.DBName}, Value: nv}})
		stmt.SetColumn(v.Field.DBName, nv, true)
	}
}
