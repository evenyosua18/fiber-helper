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

	//to fiber
	f, ok := ctx.(*fiber.Ctx)

	if !ok {
		return nil, ""
	}

	//check transaction name, use it if exist
	transactionName, ok := f.Locals("transaction_name").(string)

	if !ok {
		transactionName = f.Route().Name
	}

	return f.UserContext(), transactionName
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

func (h *FiberResponseImpl) ResponseSuccess(ctx, response interface{}, code ...interface{}) error {
	f := ctx.(*fiber.Ctx)

	//set http status code
	msg := ErrorResponse{
		ResponseCode: 200,
		CustomCode:   200,
	}

	if len(code) == 1 {
		if err := mapstructure.Decode(code[0], &msg); err != nil {
			return err
		}
	}

	f.Status(msg.ResponseCode)

	//set response
	return f.JSON(HttpResponse{
		Code:         msg.CustomCode,
		Message:      msg.ResponseMessage,
		ErrorMessage: "",
		Data:         response,
	})
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
