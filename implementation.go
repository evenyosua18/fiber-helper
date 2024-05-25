package fiber_helper

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
)

type FiberImpl struct{}

func (h *FiberImpl) GetContextName(ctx interface{}) (context.Context, string) {
	//check nil
	if ctx == nil {
		return nil, ""
	}

	//if fiber context
	if f, ok := ctx.(*fiber.Ctx); ok {
		//check transaction name, use it if exist
		transactionName, ok := f.Locals("transaction_name").(string)

		if !ok {
			transactionName = f.Route().Name
		}

		return f.UserContext(), transactionName
	}

	if c, ok := ctx.(context.Context); ok {
		//check transaction name, use it if exist
		transactionName, ok := c.Value("transaction_name").(string)

		if !ok {
			transactionName = "Unknown Parent Name"
		}

		return c, transactionName
	}

	return nil, ""
}

func (h *FiberImpl) GetInfo(ctx interface{}) (info map[string]interface{}) {
	//check nil
	if ctx == nil {
		return
	}

	//to fiber
	f, ok := ctx.(*fiber.Ctx)

	if !ok {
		return
	}

	//set info
	info = make(map[string]interface{})
	info["header"] = f.GetReqHeaders()
	info["endpoint"] = f.OriginalURL()
	info["context"] = f.Context().String()

	return
}

type FiberResponseImpl struct{}

// ResponseSuccess
// 1 code, 2 paginate, 3 filter
func (h *FiberResponseImpl) ResponseSuccess(ctx, response interface{}, additions ...interface{}) error {
	f := ctx.(*fiber.Ctx)

	// set response
	res := HttpResponse{
		Data: response,
	}

	//set http status code
	msg := ErrorResponse{
		ResponseCode: 200,
		CustomCode:   200,
	}

	// code
	if len(additions) >= 1 {
		if err := mapstructure.Decode(additions[0], &msg); err != nil {
			return err
		}
	}

	// set message
	res.Code = msg.CustomCode
	res.Message = msg.ResponseMessage

	// pagination
	if len(additions) >= 2 {
		res.Pagination = additions[1]
	}

	// filters
	if len(additions) >= 3 {
		res.Filters = additions[2]
	}

	//set response
	return f.Status(msg.ResponseCode).JSON(res)
}

func (h *FiberResponseImpl) ResponseFailed(ctx, code interface{}) error {
	f := ctx.(*fiber.Ctx)

	//set code response
	msg := ErrorResponse{}

	if err := mapstructure.Decode(code, &msg); err != nil {
		return err
	}

	//set http code
	f.Status(msg.ResponseCode)

	//set response
	return f.JSON(HttpResponse{
		Code:         msg.CustomCode,
		Message:      msg.ResponseMessage,
		ErrorMessage: msg.ErrorMessage,
	})
}

func (h *FiberResponseImpl) ResponseErrors(ctx, code interface{}, errs interface{}) error {
	f := ctx.(*fiber.Ctx)

	//set code response
	msg := ErrorResponse{}

	if err := mapstructure.Decode(code, &msg); err != nil {
		return err
	}

	//set http code
	f.Status(msg.ResponseCode)

	//set response
	return f.JSON(HttpResponse{
		Code:         msg.CustomCode,
		Message:      msg.ResponseMessage,
		ErrorMessage: msg.ErrorMessage,
		Errors:       errs,
	})
}
