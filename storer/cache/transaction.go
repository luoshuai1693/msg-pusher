/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : transaction.go
#   Created       : 2019/1/23 14:38
#   Last Modified : 2019/1/23 14:38
#   Describe      :
#
# ====================================================*/
package cache

import (
	"context"

	"uuabc.com/sendmsg/pkg/cache/redis"
	"uuabc.com/sendmsg/storer"
)

type Transaction struct {
	C *redis.Client
}

func NewTransaction() *Transaction {
	c := storer.Cache.Pipeline()
	return &Transaction{
		C: c,
	}
}

// PutBaseCache 底层缓存，跟数据库数据同步，不过期
func (t *Transaction) PutBaseCache(ctx context.Context, k string, v []byte) error {
	return put(ctx, "tx-PutBaseCache", base+k, v, 0, t.C)
}

// PutBaseTemplate 将添加的模板存入缓存中
func (t *Transaction) PutBaseTemplate(ctx context.Context, k string, v []byte) error {
	return put(ctx, "tx-PutBaseTemplate", template+k, v, 0, t.C)
}

// PutSendResult 修改发送结果，一般只有发送成功的情况才需要设置
func (t *Transaction) PutSendSuccess(ctx context.Context, k string) error {
	return put(ctx, "tx-PutSendSuccess", k+"_send", success, 0, t.C)
}

func (t *Transaction) RPushEmail(ctx context.Context, b []byte) error {
	return rPush(ctx, "tx-RPushEmail", emailDB, b, t.C)
}

func (t *Transaction) RPushWeChat(ctx context.Context, b []byte) error {
	return rPush(ctx, "tx-RPushWeChat", weChatDB, b, t.C)
}

func (t *Transaction) RPushSms(ctx context.Context, b []byte) error {
	return rPush(ctx, "tx-RPushSms", smsDB, b, t.C)
}

// LPopWeChat 从wechat队列中取一条数据
func (t *Transaction) LPopWeChat() ([]byte, error) {
	return lPop(context.Background(), weChatDB, t.C)
}

// LPopEmail 从email队列中取一条数据
func (t *Transaction) LPopEmail() ([]byte, error) {
	return lPop(context.Background(), emailDB, t.C)
}

// LPopSms 从sms队列中取一条数据
func (t *Transaction) LPopSms() ([]byte, error) {
	return lPop(context.Background(), smsDB, t.C)
}

// Commit 提交事务
func (t *Transaction) Commit() error {
	_, err := t.C.Exec()
	return err
}

// CommitParam 待参数的提交
func (t *Transaction) CommitParam() ([]interface{}, error) {
	return t.C.Exec()
}

// Rollback 回滚事务
func (t *Transaction) Rollback() error {
	return t.C.Discard()
}

func (t *Transaction) Close() error {
	return t.C.Close()
}