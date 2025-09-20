package utils

import (
	"errors"

	"buf.build/go/protovalidate"
	"github.com/irvandMR/go-grpc-ecommerce-BE/pb/common"
	"google.golang.org/protobuf/proto"
)

func CheckValidtion(req proto.Message) ([]*common.ValidationError, error) {
	if err := protovalidate.Validate(req); err != nil {
		var validateError *protovalidate.ValidationError
		if errors.As(err, &validateError) {
			var validateErrors []*common.ValidationError = make([]*common.ValidationError, 0)
			for _, violation := range validateError.Violations {
				validateErrors = append(validateErrors, &common.ValidationError{
					Field:       *violation.Proto.Field.Elements[0].FieldName,
					Description: *violation.Proto.Message,
				})
			}
			return validateErrors, nil
		}
		// return  nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
		return  nil, err
	}
	return nil, nil
}