package dao

import (
	"context"

	"bbs-back/base/common"
	"bbs-back/models/entity"

	"github.com/beego/beego/v2/client/orm"
)

type Message struct {
	entity.Message
	Comment *Comment `json:"comment" orm:"-"`
}

func (m *Message) TableName() string {
	return "message"
}

func (m *Message) Read(ctx context.Context) (*Message, error) {
	res := new(Message)
	res.Id = m.Id
	err := ORM.ReadWithCtx(ctx, res)
	if err != nil {
		return nil, err
	}
	comment := new(Comment)
	comment.Id = res.CommentId
	res.Comment, _ = comment.Read()
	return res, err
}

func (m *Message) Find(ctx context.Context, page *common.Page, cols ...string) ([]*Message, error) {
	qs := m.createQsByParam(ctx)
	qs = qs.OrderBy("type", "-create_time")
	qs = qs.Limit(page.GetPageSize(), page.GetOffset())
	var res []*Message
	if len(cols) > 0 {
		_, err := qs.All(&res, cols...)
		return res, err
	}
	_, err := qs.All(&res)
	if err != nil {
		return nil, err
	}
	comment := new(Comment)
	for _, message := range res {
		comment.Id = message.CommentId
		message.Comment, _ = comment.Read()
	}
	return res, err
}

func (m *Message) Count(ctx context.Context) (int64, error) {
	return m.createQsByParam(ctx).Count()
}

func (m *Message) createQsByParam(ctx context.Context) orm.QuerySeter {
	qs := ORM.QueryTableWithCtx(ctx, m)
	if m.Id != 0 {
		qs = qs.Filter("id", m.Id)
	}
	if m.Status != 0 {
		qs = qs.Filter("status", m.Status)
	}
	if m.Type != 0 {
		qs = qs.Filter("type", m.Type)
	}
	if m.Kind != 0 {
		qs = qs.Filter("kind", m.Kind)
	}
	if m.CreatorId != 0 {
		qs = qs.Filter("creator_id")
	}
	if m.TargetId != 0 {
		qs = qs.Filter("target_id", m.TargetId)
	}
	return m.TimeFilter(qs)
}

func (m *Message) Insert(ctx context.Context) error {
	if m.Status == 0 {
		m.Status = 1
	}
	if m.Type == 0 {
		m.Type = 1
	}
	_, err := ORM.InsertWithCtx(ctx, m)
	return err
}

// 根据主健更新
func (m *Message) Update(ctx context.Context) error {

	var cols []string
	if m.Status != 0 {
		cols = append(cols, "status")
	}
	if m.TargetId != 0 {
		cols = append(cols, "target_id")
	}
	if m.CreatorId != 0 {
		cols = append(cols, "creator_id")
	}
	if m.CommentId != 0 {
		cols = append(cols, "comment_id")
	}
	if m.Kind != 0 {
		cols = append(cols, "kind")
	}
	if m.Type != 0 {
		cols = append(cols, "type")
	}
	if len(cols) == 0 {
		return nil
	}
	_, err := ORM.UpdateWithCtx(ctx, m, cols...)
	return err
}

// 根据查询条件更新
func (m *Message) Updates(ctx context.Context, target *Message) (int64, error) {
	params := orm.Params{}
	if m.Status != 0 {
		params["status"] = m.Status
	}
	if m.TargetId != 0 {
		params["target_id"] = m.TargetId
	}
	if m.CreatorId != 0 {
		params["creator_id"] = m.CreatorId
	}
	if m.CommentId != 0 {
		params["comment_id"] = m.CommentId
	}
	if m.Kind != 0 {
		params["kind"] = m.Kind
	}
	return m.createQsByParam(ctx).Update(params)
}

func (m *Message) Delete(ctx context.Context) error {
	m.Status = -1
	_, err := ORM.UpdateWithCtx(ctx, m, "status")
	return err
}
