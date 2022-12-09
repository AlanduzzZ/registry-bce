package instance

import (
	"context"
	"fmt"

	beego_orm "github.com/beego/beego/orm"

	"github.com/goharbor/harbor/src/lib/orm"
	"github.com/goharbor/harbor/src/lib/q"
	"github.com/goharbor/harbor/src/pkg/p2p/preheat/models/provider"
)

// DAO for instance
type DAO interface {
	Create(ctx context.Context, instance *provider.Instance) (int64, error)
	Get(ctx context.Context, id int64) (*provider.Instance, error)
	GetByName(ctx context.Context, name string) (*provider.Instance, error)
	Update(ctx context.Context, instance *provider.Instance, props ...string) error
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context, query *q.Query) (total int64, err error)
	List(ctx context.Context, query *q.Query) (ins []*provider.Instance, err error)
}

// New instance dao
func New() DAO {
	return &dao{}
}

// ListInstanceQuery defines the query params of the instance record.
type ListInstanceQuery struct {
	Page     uint
	PageSize uint
	Keyword  string
}

type dao struct{}

var _ DAO = (*dao)(nil)

// Create adds a new distribution instance.
func (d *dao) Create(ctx context.Context, instance *provider.Instance) (id int64, err error) {
	var o beego_orm.Ormer
	o, err = orm.FromContext(ctx)
	if err != nil {
		return
	}

	id, err = o.Insert(instance)
	if err != nil {
		if e := orm.AsConflictError(err, "instance %s already exists", instance.Name); e != nil {
			err = e
		}
		return
	}
	return
}

// Get gets instance from db by id.
func (d *dao) Get(ctx context.Context, id int64) (*provider.Instance, error) {
	var o, err = orm.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	di := provider.Instance{ID: id}
	if err = o.Read(&di, "ID"); err != nil {
		if e := orm.AsNotFoundError(err, "instance %d not found", id); e != nil {
			err = e
		}
		return nil, err
	}

	return &di, err
}

// Get gets instance from db by name.
func (d *dao) GetByName(ctx context.Context, name string) (instance *provider.Instance, err error) {
	o, err := orm.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	instance = &provider.Instance{Name: name}
	if err = o.Read(instance, "Name"); err != nil {
		if e := orm.AsNotFoundError(err, "instance %s not found", name); e != nil {
			err = e
		}
		return nil, err
	}
	return
}

// Update updates distribution instance.
func (d *dao) Update(ctx context.Context, instance *provider.Instance, props ...string) error {
	var trans = func(ctx context.Context) (err error) {
		o, err := orm.FromContext(ctx)
		if err != nil {
			return
		}

		// check default instances first
		if instance.Default {
			_, err = o.Raw(fmt.Sprintf("UPDATE %s SET is_default = false WHERE id != ?", instance.TableName()), instance.ID).Exec()
			if err != nil {
				return
			}
		}

		_, err = o.Update(instance, props...)
		return
	}
	return orm.WithTransaction(trans)(orm.SetTransactionOpNameToContext(ctx, "tx-prehead-update"))
}

// Delete deletes one distribution instance by id.
func (d *dao) Delete(ctx context.Context, id int64) error {
	var o, err = orm.FromContext(ctx)
	if err != nil {
		return err
	}

	_, err = o.Delete(&provider.Instance{ID: id})
	return err
}

// List count instances by query params.
func (d *dao) Count(ctx context.Context, query *q.Query) (total int64, err error) {
	qs, err := orm.QuerySetterForCount(ctx, &provider.Instance{}, query)
	if err != nil {
		return 0, err
	}
	return qs.Count()
}

// List lists instances by query params.
func (d *dao) List(ctx context.Context, query *q.Query) (ins []*provider.Instance, err error) {
	ins = []*provider.Instance{}
	qs, err := orm.QuerySetter(ctx, &provider.Instance{}, query)
	if err != nil {
		return nil, err
	}
	if _, err = qs.All(&ins); err != nil {
		return nil, err
	}
	return ins, nil
}
