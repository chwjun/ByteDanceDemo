// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"bytedancedemo/model"
)

func newRelation(db *gorm.DB, opts ...gen.DOOption) relation {
	_relation := relation{}

	_relation.relationDo.UseDB(db, opts...)
	_relation.relationDo.UseModel(&model.Relation{})

	tableName := _relation.relationDo.TableName()
	_relation.ALL = field.NewAsterisk(tableName)
	_relation.ID = field.NewInt64(tableName, "id")
	_relation.CreatedAt = field.NewTime(tableName, "created_at")
	_relation.UpdatedAt = field.NewTime(tableName, "updated_at")
	_relation.DeletedAt = field.NewField(tableName, "deleted_at")
	_relation.UserID = field.NewInt64(tableName, "user_id")
	_relation.FollowingID = field.NewInt64(tableName, "following_id")
	_relation.Followed = field.NewInt64(tableName, "followed")

	_relation.fillFieldMap()

	return _relation
}

type relation struct {
	relationDo

	ALL         field.Asterisk
	ID          field.Int64 // 主键
	CreatedAt   field.Time  // 记录创建时间
	UpdatedAt   field.Time  // 记录更新时间
	DeletedAt   field.Field // 软删除时间
	UserID      field.Int64 // 用户id
	FollowingID field.Int64 // user id关注的用户id
	Followed    field.Int64 // 默认0表示未关注，1表示已关注

	fieldMap map[string]field.Expr
}

func (r relation) Table(newTableName string) *relation {
	r.relationDo.UseTable(newTableName)
	return r.updateTableName(newTableName)
}

func (r relation) As(alias string) *relation {
	r.relationDo.DO = *(r.relationDo.As(alias).(*gen.DO))
	return r.updateTableName(alias)
}

func (r *relation) updateTableName(table string) *relation {
	r.ALL = field.NewAsterisk(table)
	r.ID = field.NewInt64(table, "id")
	r.CreatedAt = field.NewTime(table, "created_at")
	r.UpdatedAt = field.NewTime(table, "updated_at")
	r.DeletedAt = field.NewField(table, "deleted_at")
	r.UserID = field.NewInt64(table, "user_id")
	r.FollowingID = field.NewInt64(table, "following_id")
	r.Followed = field.NewInt64(table, "followed")

	r.fillFieldMap()

	return r
}

func (r *relation) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := r.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (r *relation) fillFieldMap() {
	r.fieldMap = make(map[string]field.Expr, 7)
	r.fieldMap["id"] = r.ID
	r.fieldMap["created_at"] = r.CreatedAt
	r.fieldMap["updated_at"] = r.UpdatedAt
	r.fieldMap["deleted_at"] = r.DeletedAt
	r.fieldMap["user_id"] = r.UserID
	r.fieldMap["following_id"] = r.FollowingID
	r.fieldMap["followed"] = r.Followed
}

func (r relation) clone(db *gorm.DB) relation {
	r.relationDo.ReplaceConnPool(db.Statement.ConnPool)
	return r
}

func (r relation) replaceDB(db *gorm.DB) relation {
	r.relationDo.ReplaceDB(db)
	return r
}

type relationDo struct{ gen.DO }

type IRelationDo interface {
	gen.SubQuery
	Debug() IRelationDo
	WithContext(ctx context.Context) IRelationDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IRelationDo
	WriteDB() IRelationDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IRelationDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IRelationDo
	Not(conds ...gen.Condition) IRelationDo
	Or(conds ...gen.Condition) IRelationDo
	Select(conds ...field.Expr) IRelationDo
	Where(conds ...gen.Condition) IRelationDo
	Order(conds ...field.Expr) IRelationDo
	Distinct(cols ...field.Expr) IRelationDo
	Omit(cols ...field.Expr) IRelationDo
	Join(table schema.Tabler, on ...field.Expr) IRelationDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IRelationDo
	RightJoin(table schema.Tabler, on ...field.Expr) IRelationDo
	Group(cols ...field.Expr) IRelationDo
	Having(conds ...gen.Condition) IRelationDo
	Limit(limit int) IRelationDo
	Offset(offset int) IRelationDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IRelationDo
	Unscoped() IRelationDo
	Create(values ...*model.Relation) error
	CreateInBatches(values []*model.Relation, batchSize int) error
	Save(values ...*model.Relation) error
	First() (*model.Relation, error)
	Take() (*model.Relation, error)
	Last() (*model.Relation, error)
	Find() ([]*model.Relation, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Relation, err error)
	FindInBatches(result *[]*model.Relation, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Relation) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IRelationDo
	Assign(attrs ...field.AssignExpr) IRelationDo
	Joins(fields ...field.RelationField) IRelationDo
	Preload(fields ...field.RelationField) IRelationDo
	FirstOrInit() (*model.Relation, error)
	FirstOrCreate() (*model.Relation, error)
	FindByPage(offset int, limit int) (result []*model.Relation, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IRelationDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (r relationDo) Debug() IRelationDo {
	return r.withDO(r.DO.Debug())
}

func (r relationDo) WithContext(ctx context.Context) IRelationDo {
	return r.withDO(r.DO.WithContext(ctx))
}

func (r relationDo) ReadDB() IRelationDo {
	return r.Clauses(dbresolver.Read)
}

func (r relationDo) WriteDB() IRelationDo {
	return r.Clauses(dbresolver.Write)
}

func (r relationDo) Session(config *gorm.Session) IRelationDo {
	return r.withDO(r.DO.Session(config))
}

func (r relationDo) Clauses(conds ...clause.Expression) IRelationDo {
	return r.withDO(r.DO.Clauses(conds...))
}

func (r relationDo) Returning(value interface{}, columns ...string) IRelationDo {
	return r.withDO(r.DO.Returning(value, columns...))
}

func (r relationDo) Not(conds ...gen.Condition) IRelationDo {
	return r.withDO(r.DO.Not(conds...))
}

func (r relationDo) Or(conds ...gen.Condition) IRelationDo {
	return r.withDO(r.DO.Or(conds...))
}

func (r relationDo) Select(conds ...field.Expr) IRelationDo {
	return r.withDO(r.DO.Select(conds...))
}

func (r relationDo) Where(conds ...gen.Condition) IRelationDo {
	return r.withDO(r.DO.Where(conds...))
}

func (r relationDo) Order(conds ...field.Expr) IRelationDo {
	return r.withDO(r.DO.Order(conds...))
}

func (r relationDo) Distinct(cols ...field.Expr) IRelationDo {
	return r.withDO(r.DO.Distinct(cols...))
}

func (r relationDo) Omit(cols ...field.Expr) IRelationDo {
	return r.withDO(r.DO.Omit(cols...))
}

func (r relationDo) Join(table schema.Tabler, on ...field.Expr) IRelationDo {
	return r.withDO(r.DO.Join(table, on...))
}

func (r relationDo) LeftJoin(table schema.Tabler, on ...field.Expr) IRelationDo {
	return r.withDO(r.DO.LeftJoin(table, on...))
}

func (r relationDo) RightJoin(table schema.Tabler, on ...field.Expr) IRelationDo {
	return r.withDO(r.DO.RightJoin(table, on...))
}

func (r relationDo) Group(cols ...field.Expr) IRelationDo {
	return r.withDO(r.DO.Group(cols...))
}

func (r relationDo) Having(conds ...gen.Condition) IRelationDo {
	return r.withDO(r.DO.Having(conds...))
}

func (r relationDo) Limit(limit int) IRelationDo {
	return r.withDO(r.DO.Limit(limit))
}

func (r relationDo) Offset(offset int) IRelationDo {
	return r.withDO(r.DO.Offset(offset))
}

func (r relationDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IRelationDo {
	return r.withDO(r.DO.Scopes(funcs...))
}

func (r relationDo) Unscoped() IRelationDo {
	return r.withDO(r.DO.Unscoped())
}

func (r relationDo) Create(values ...*model.Relation) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Create(values)
}

func (r relationDo) CreateInBatches(values []*model.Relation, batchSize int) error {
	return r.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (r relationDo) Save(values ...*model.Relation) error {
	if len(values) == 0 {
		return nil
	}
	return r.DO.Save(values)
}

func (r relationDo) First() (*model.Relation, error) {
	if result, err := r.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Relation), nil
	}
}

func (r relationDo) Take() (*model.Relation, error) {
	if result, err := r.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Relation), nil
	}
}

func (r relationDo) Last() (*model.Relation, error) {
	if result, err := r.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Relation), nil
	}
}

func (r relationDo) Find() ([]*model.Relation, error) {
	result, err := r.DO.Find()
	return result.([]*model.Relation), err
}

func (r relationDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Relation, err error) {
	buf := make([]*model.Relation, 0, batchSize)
	err = r.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (r relationDo) FindInBatches(result *[]*model.Relation, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return r.DO.FindInBatches(result, batchSize, fc)
}

func (r relationDo) Attrs(attrs ...field.AssignExpr) IRelationDo {
	return r.withDO(r.DO.Attrs(attrs...))
}

func (r relationDo) Assign(attrs ...field.AssignExpr) IRelationDo {
	return r.withDO(r.DO.Assign(attrs...))
}

func (r relationDo) Joins(fields ...field.RelationField) IRelationDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Joins(_f))
	}
	return &r
}

func (r relationDo) Preload(fields ...field.RelationField) IRelationDo {
	for _, _f := range fields {
		r = *r.withDO(r.DO.Preload(_f))
	}
	return &r
}

func (r relationDo) FirstOrInit() (*model.Relation, error) {
	if result, err := r.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Relation), nil
	}
}

func (r relationDo) FirstOrCreate() (*model.Relation, error) {
	if result, err := r.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Relation), nil
	}
}

func (r relationDo) FindByPage(offset int, limit int) (result []*model.Relation, count int64, err error) {
	result, err = r.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = r.Offset(-1).Limit(-1).Count()
	return
}

func (r relationDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = r.Count()
	if err != nil {
		return
	}

	err = r.Offset(offset).Limit(limit).Scan(result)
	return
}

func (r relationDo) Scan(result interface{}) (err error) {
	return r.DO.Scan(result)
}

func (r relationDo) Delete(models ...*model.Relation) (result gen.ResultInfo, err error) {
	return r.DO.Delete(models)
}

func (r *relationDo) withDO(do gen.Dao) *relationDo {
	r.DO = *do.(*gen.DO)
	return r
}
