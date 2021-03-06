/**
 * Copyright 2015 @ z3q.net.
 * name : payment_service.go
 * author : jarryliu
 * date : 2016-07-03 13:24
 * description :
 * history :
 */
package rsi

import (
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/service/thrift/idl/gen-go/define"
	"go2o/core/service/thrift/parser"
)

type paymentService struct {
	_rep       payment.IPaymentRepo
	_orderRepo order.IOrderRepo
}

func NewPaymentService(rep payment.IPaymentRepo, orderRepo order.IOrderRepo) *paymentService {
	return &paymentService{
		_rep:       rep,
		_orderRepo: orderRepo,
	}
}

// 根据编号获取支付单
func (p *paymentService) GetPaymentOrderById(id int32) (*define.PaymentOrder, error) {
	po := p._rep.GetPaymentOrderById(id)
	if po != nil {
		v := po.GetValue()
		return parser.PaymentOrderDto(&v), nil
	}
	return nil, nil
}

// 根据支付单号获取支付单
func (p *paymentService) GetPaymentOrder(paymentNo string) (*define.PaymentOrder, error) {
	if po := p._rep.GetPaymentOrder(paymentNo); po != nil {
		v := po.GetValue()
		return parser.PaymentOrderDto(&v), nil
	}
	return nil, nil
}

// 创建支付单
func (p *paymentService) CreatePaymentOrder(s *define.PaymentOrder) (*define.Result_, error) {
	v := parser.PaymentOrder(s)
	o := p._rep.CreatePaymentOrder(v)
	id, err := o.Commit()
	r := &define.Result_{
		Result_: err == nil,
		ID:      id,
	}
	if err != nil {
		r.Message = err.Error()
	}
	return r, nil
}

// 调整支付单金额
func (p *paymentService) AdjustOrder(paymentNo string, amount float64) (*define.Result_, error) {
	var err error
	o := p._rep.GetPaymentOrder(paymentNo)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.Adjust(float32(amount))
	}
	return parser.Result(0, err), nil
}

func (p *paymentService) SetPrefixOfTradeNo(id int32, prefix string) error {
	o := p._rep.GetPaymentOrderById(id)
	if o == nil {
		return payment.ErrNoSuchPaymentOrder
	}
	return o.TradeNoPrefix(prefix)
}

// 积分抵扣支付单
func (p *paymentService) DiscountByIntegral(orderId int32,
	integral int32, ignoreOut bool) (r *define.DResult_, err error) {
	var amount float32
	o := p._rep.GetPaymentOrderById(orderId)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		amount, err = o.IntegralDiscount(int(integral), ignoreOut)
	}
	return parser.DResult(float64(amount), err), nil
}

// 余额抵扣
func (p *paymentService) DiscountByBalance(orderId int32, remark string) (*define.Result_, error) {
	var err error
	o := p._rep.GetPaymentOrderById(orderId)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.BalanceDiscount(remark)
	}
	return parser.Result(0, err), nil
}

// 赠送账户支付
func (p *paymentService) PaymentByPresent(orderId int32, remark string) (r *define.Result_, err error) {
	o := p._rep.GetPaymentOrderById(orderId)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.PaymentByPresent(remark)
	}
	return parser.Result(0, err), nil
}

// 完成支付单支付，并传入支付方式及外部订单号
func (p *paymentService) FinishPayment(tradeNo string, spName string,
	outerNo string) (r *define.Result_, err error) {
	o := p._rep.GetPaymentOrder(tradeNo)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.PaymentFinish(spName, outerNo)
	}
	return parser.Result(0, err), nil
}
