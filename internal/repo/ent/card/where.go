// Code generated by entc, DO NOT EDIT.

package card

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/realHoangHai/awesome/internal/repo/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int64) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int64) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int64) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int64) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int64) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int64) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int64) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int64) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int64) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// CardNo applies equality check predicate on the "card_no" field. It's identical to CardNoEQ.
func CardNo(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCardNo), v))
	})
}

// Ccv applies equality check predicate on the "ccv" field. It's identical to CcvEQ.
func Ccv(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCcv), v))
	})
}

// Expires applies equality check predicate on the "expires" field. It's identical to ExpiresEQ.
func Expires(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldExpires), v))
	})
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Card {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldName), v...))
	})
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Card {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldName), v...))
	})
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// CardNoEQ applies the EQ predicate on the "card_no" field.
func CardNoEQ(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCardNo), v))
	})
}

// CardNoNEQ applies the NEQ predicate on the "card_no" field.
func CardNoNEQ(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCardNo), v))
	})
}

// CardNoIn applies the In predicate on the "card_no" field.
func CardNoIn(vs ...string) predicate.Card {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCardNo), v...))
	})
}

// CardNoNotIn applies the NotIn predicate on the "card_no" field.
func CardNoNotIn(vs ...string) predicate.Card {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCardNo), v...))
	})
}

// CardNoGT applies the GT predicate on the "card_no" field.
func CardNoGT(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCardNo), v))
	})
}

// CardNoGTE applies the GTE predicate on the "card_no" field.
func CardNoGTE(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCardNo), v))
	})
}

// CardNoLT applies the LT predicate on the "card_no" field.
func CardNoLT(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCardNo), v))
	})
}

// CardNoLTE applies the LTE predicate on the "card_no" field.
func CardNoLTE(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCardNo), v))
	})
}

// CardNoContains applies the Contains predicate on the "card_no" field.
func CardNoContains(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldCardNo), v))
	})
}

// CardNoHasPrefix applies the HasPrefix predicate on the "card_no" field.
func CardNoHasPrefix(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldCardNo), v))
	})
}

// CardNoHasSuffix applies the HasSuffix predicate on the "card_no" field.
func CardNoHasSuffix(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldCardNo), v))
	})
}

// CardNoEqualFold applies the EqualFold predicate on the "card_no" field.
func CardNoEqualFold(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldCardNo), v))
	})
}

// CardNoContainsFold applies the ContainsFold predicate on the "card_no" field.
func CardNoContainsFold(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldCardNo), v))
	})
}

// CcvEQ applies the EQ predicate on the "ccv" field.
func CcvEQ(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCcv), v))
	})
}

// CcvNEQ applies the NEQ predicate on the "ccv" field.
func CcvNEQ(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCcv), v))
	})
}

// CcvIn applies the In predicate on the "ccv" field.
func CcvIn(vs ...string) predicate.Card {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCcv), v...))
	})
}

// CcvNotIn applies the NotIn predicate on the "ccv" field.
func CcvNotIn(vs ...string) predicate.Card {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCcv), v...))
	})
}

// CcvGT applies the GT predicate on the "ccv" field.
func CcvGT(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCcv), v))
	})
}

// CcvGTE applies the GTE predicate on the "ccv" field.
func CcvGTE(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCcv), v))
	})
}

// CcvLT applies the LT predicate on the "ccv" field.
func CcvLT(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCcv), v))
	})
}

// CcvLTE applies the LTE predicate on the "ccv" field.
func CcvLTE(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCcv), v))
	})
}

// CcvContains applies the Contains predicate on the "ccv" field.
func CcvContains(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldCcv), v))
	})
}

// CcvHasPrefix applies the HasPrefix predicate on the "ccv" field.
func CcvHasPrefix(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldCcv), v))
	})
}

// CcvHasSuffix applies the HasSuffix predicate on the "ccv" field.
func CcvHasSuffix(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldCcv), v))
	})
}

// CcvEqualFold applies the EqualFold predicate on the "ccv" field.
func CcvEqualFold(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldCcv), v))
	})
}

// CcvContainsFold applies the ContainsFold predicate on the "ccv" field.
func CcvContainsFold(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldCcv), v))
	})
}

// ExpiresEQ applies the EQ predicate on the "expires" field.
func ExpiresEQ(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldExpires), v))
	})
}

// ExpiresNEQ applies the NEQ predicate on the "expires" field.
func ExpiresNEQ(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldExpires), v))
	})
}

// ExpiresIn applies the In predicate on the "expires" field.
func ExpiresIn(vs ...string) predicate.Card {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldExpires), v...))
	})
}

// ExpiresNotIn applies the NotIn predicate on the "expires" field.
func ExpiresNotIn(vs ...string) predicate.Card {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldExpires), v...))
	})
}

// ExpiresGT applies the GT predicate on the "expires" field.
func ExpiresGT(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldExpires), v))
	})
}

// ExpiresGTE applies the GTE predicate on the "expires" field.
func ExpiresGTE(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldExpires), v))
	})
}

// ExpiresLT applies the LT predicate on the "expires" field.
func ExpiresLT(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldExpires), v))
	})
}

// ExpiresLTE applies the LTE predicate on the "expires" field.
func ExpiresLTE(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldExpires), v))
	})
}

// ExpiresContains applies the Contains predicate on the "expires" field.
func ExpiresContains(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldExpires), v))
	})
}

// ExpiresHasPrefix applies the HasPrefix predicate on the "expires" field.
func ExpiresHasPrefix(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldExpires), v))
	})
}

// ExpiresHasSuffix applies the HasSuffix predicate on the "expires" field.
func ExpiresHasSuffix(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldExpires), v))
	})
}

// ExpiresEqualFold applies the EqualFold predicate on the "expires" field.
func ExpiresEqualFold(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldExpires), v))
	})
}

// ExpiresContainsFold applies the ContainsFold predicate on the "expires" field.
func ExpiresContainsFold(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldExpires), v))
	})
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Card {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Card {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreatedAt), v))
	})
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Card {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldUpdatedAt), v...))
	})
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Card {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldUpdatedAt), v...))
	})
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdatedAt), v))
	})
}

// HasUser applies the HasEdge predicate on the "user" edge.
func HasUser() predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(UserTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserWith applies the HasEdge predicate on the "user" edge with a given conditions (other predicates).
func HasUserWith(preds ...predicate.User) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(UserInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Card) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Card) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Card) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		p(s.Not())
	})
}
