package gotools

import (
	"fmt"
	"go/token"
	"go/types"
)

func SafeSignature(sig *types.Signature) *types.Signature {
	var newParamVars []*types.Var

	params := sig.Params()
	for i := 0; i < params.Len(); i++ {
		param := params.At(i)
		newVar := types.NewVar(token.NoPos, nil, fmt.Sprintf("z%04d", i+1), param.Type())
		newParamVars = append(newParamVars, newVar)
	}

	var newResultVars []*types.Var

	results := sig.Results()
	for i := 0; i < results.Len(); i++ {
		result := results.At(i)
		newVar := types.NewVar(token.NoPos, nil, "", result.Type())
		newResultVars = append(newResultVars, newVar)
	}

	return types.NewSignatureType(
		nil, nil, nil,
		types.NewTuple(newParamVars...),
		types.NewTuple(newResultVars...),
		sig.Variadic(),
	)

}

// Removes all parameter names, only types left.
func AnonymSignature(sig *types.Signature) *types.Signature {
	var newParamVars []*types.Var

	params := sig.Params()
	for i := 0; i < params.Len(); i++ {
		param := params.At(i)
		newVar := types.NewVar(token.NoPos, nil, "", param.Type())
		newParamVars = append(newParamVars, newVar)
	}

	var newResultVars []*types.Var

	results := sig.Results()
	for i := 0; i < results.Len(); i++ {
		result := results.At(i)
		newVar := types.NewVar(token.NoPos, nil, "", result.Type())
		newResultVars = append(newResultVars, newVar)
	}

	recvTypeParams := make([]*types.TypeParam, 0, sig.RecvTypeParams().Len())
	for i := 0; i < sig.RecvTypeParams().Len(); i++ {
		recvTypeParams = append(recvTypeParams, sig.RecvTypeParams().At(i))
	}
	typeParams := make([]*types.TypeParam, 0, sig.TypeParams().Len())
	for i := 0; i < sig.TypeParams().Len(); i++ {
		typeParams = append(typeParams, sig.TypeParams().At(i))
	}

	return types.NewSignatureType(
		sig.Recv(), recvTypeParams, typeParams,
		types.NewTuple(newParamVars...),
		types.NewTuple(newResultVars...),
		sig.Variadic(),
	)
}
