package lambdaHandson

import (
	"testing"

	"github.com/k-nishijima/lambda-handson"
)

func checkErrorMsg(t *testing.T, err error, msg string) {
	if err != nil {
		if msg != err.Error() {
			t.Errorf("No match : expected[%v] / actual[%v]", msg, err.Error())
		}
	} else {
		if msg != "" {
			t.Error("Why? this request is valid. expected[" + msg + "]")
		}
	}
}

func TestValidateAddValueRequest(t *testing.T) {
	var dao lambdaHandson.DynamoDBClient
	var err error

	// 値なし：エラー
	err = dao.ValidateRequest(lambdaHandson.AddValueRequest{})
	checkErrorMsg(t, err, "Stage: non zero value required;Email: non zero value required;Message: non zero value required;")

	// 正しいメールアドレスが必要
	err = dao.ValidateRequest(lambdaHandson.AddValueRequest{
		Stage:   "dev",
		Email:   "invalid email",
		Message: "hello world",
	})
	checkErrorMsg(t, err, "Email: invalid email does not validate as email;")

	// all OK.
	err = dao.ValidateRequest(lambdaHandson.AddValueRequest{
		Stage:   "dev",
		Email:   "foo@bar.com",
		Message: "hello world",
	})
	checkErrorMsg(t, err, "")
}

func TestPut(t *testing.T) {
	var dao lambdaHandson.DynamoDBClient
	request := lambdaHandson.AddValueRequest{
		Stage:   "dev",
		Email:   "foo@example.com",
		Message: "Hello golang lambda world.",
	}
	err := dao.Put(request)
	if err != nil {
		t.Error(err)
	}

	items, err := dao.GetItems(request)
	if err != nil {
		t.Error(err)
	}
	for _, i := range items {
		t.Logf("%v : %v", i.Timestamp, i.Email)
	}
}
