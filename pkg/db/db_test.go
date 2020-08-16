package db

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"

	"go.etcd.io/bbolt"
)

func TestNew(t *testing.T) {

	db := "/tmp/db_new"
	service, err := New(db)
	if err != nil {
		t.Errorf("create db failed %s", err.Error())
		return
	}
	defer os.Remove(db)

	_, err = os.Stat(db)
	if os.IsNotExist(err) {
		t.Errorf("db file not exist %s", err.Error())
		return
	}

	if service.path != db {
		t.Errorf("invalid service")
		return
	}

	if service.DB == nil {
		t.Errorf("invalid database")
	}

	db = "./tmp/./.db_new"
	_, err = New(db)
	if err == nil {
		t.Error("create db should fail with invalid path")
		return
	}

}

func TestService_Close(t *testing.T) {

	db := "/tmp/db_close"
	service, err := New(db)
	if err != nil {
		t.Errorf("create db failed %s", err.Error())
		return
	}
	defer os.Remove(db)

	err = service.Close()
	if err != nil {
		t.Errorf("close db failed %s", err.Error())
	}
}

func TestService_CreateBucket(t *testing.T) {

	db := "/tmp/db_create_bucket"
	service, err := New(db)
	if err != nil {
		t.Errorf("create db failed %s", err.Error())
		return
	}
	defer os.Remove(db)

	bucket := "test"
	b, err := service.CreateBucket(bucket)
	if err != nil {
		t.Errorf("create bucket failed %s", err.Error())
		return
	}

	if b == nil {
		t.Error("invalid bucket")
		return
	}

	_ = service.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			t.Errorf("bucket %s not exist", bucket)
		}
		return nil
	})

	bucket = ""
	_, err = service.CreateBucket(bucket)
	if err == nil {
		t.Errorf("create bucket should fail %s", err.Error())
		return
	}
}

func TestService_DeleteBucket(t *testing.T) {

	db := "/tmp/db_delete_bucket"
	service, err := New(db)
	if err != nil {
		t.Errorf("create db failed %s", err.Error())
		return
	}
	defer os.Remove(db)

	bucket := "test_delete"
	b, err := service.CreateBucket(bucket)
	if err != nil {
		t.Errorf("create bucket failed %s", err.Error())
		return
	}

	if b == nil {
		t.Error("invalid bucket")
		return
	}

	err = service.DeleteBucket(bucket)
	if err != nil {
		t.Errorf("delete bucket failed %s", err.Error())
		return
	}

	_ = service.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			t.Errorf("bucket %s still exist", bucket)
		}
		return nil
	})

	err = service.DeleteBucket(bucket)
	if err == nil {
		t.Error("delete bucket should failed when not exist")
		return
	}
}

func createDBandBucket(db, bucket string) (*Service, error) {

	service, err := New(db)
	if err != nil {
		return nil, err
	}
	defer os.Remove(db)

	b, err := service.CreateBucket(bucket)
	if err != nil {
		return nil, err
	}

	if b == nil {
		return nil, fmt.Errorf("create bucket %s failed", bucket)
	}
	return service, nil
}

func TestService_Set(t *testing.T) {

	db := "/tmp/db_set"
	bucket := "test_set"
	service, err := createDBandBucket(db, bucket)
	if err != nil {
		t.Error("error in prepare db and bucket")
		return
	}
	defer os.Remove(db)

	key := "test_key_set"
	val := []byte("this is a test")

	err = service.Set(key, bucket, val)
	if err != nil {
		t.Error("create key failed")
	}

	_ = service.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b.Get([]byte(key)) == nil {
			t.Error("key is not found in db")
		}
		return nil
	})

	bucket = "not-exist"
	// test when bucket in not exist
	err = service.Set(key, bucket, val)
	if err == nil {
		t.Error("create key should fail when bucket not exist")
	}
}

func TestService_GetOne(t *testing.T) {

	db := "/tmp/db_get_one"
	bucket := "test_get_one"
	service, err := createDBandBucket(db, bucket)
	if err != nil {
		t.Error("error in prepare db and bucket")
		return
	}
	defer os.Remove(db)

	key := "test_key_get_one"
	val := []byte("this is a test")

	err = service.Set(key, bucket, val)
	if err != nil {
		t.Error("create key failed")
		return
	}

	getVal, err := service.GetOne(key, bucket)
	if err != nil {
		t.Errorf("GetOne failed with error %s", err.Error())
		return
	}

	if bytes.Compare(val, getVal) != 0 {
		t.Error("read value is not same as wrote")
	}

	bucket = "not-exist"
	// test when bucket in not exist
	_, err = service.GetOne(key, bucket)
	if err == nil {
		t.Error("get key should fail when bucket not exist")
	}
}

func TestService_GetAll(t *testing.T) {
	db := "/tmp/db_get_all"
	bucket := "test_get_all"
	service, err := createDBandBucket(db, bucket)
	if err != nil {
		t.Error("error in prepare db and bucket")
		return
	}
	defer os.Remove(db)

	kv := []KeyVal{{Key: "test1", Val: []byte("test1")}, {Key: "test2", Val: []byte("test2")}}

	for _, kVal := range kv {
		err = service.Set(kVal.Key, bucket, kVal.Val)
		if err != nil {
			t.Error("create key failed")
			return
		}
	}

	result, err := service.GetAll(bucket)
	if !reflect.DeepEqual(result, kv) {
		t.Error("get all returned different values")
	}

	bucket = "not-exist"
	// test when bucket in not exist
	_, err = service.GetAll(bucket)
	if err == nil {
		t.Error("get all key should fail when bucket not exist")
	}
}

func TestService_Delete(t *testing.T) {
	db := "/tmp/db_get_delete"
	bucket := "test_get_delete"
	service, err := createDBandBucket(db, bucket)
	if err != nil {
		t.Error("error in prepare db and bucket")
		return
	}
	defer os.Remove(db)

	key := "test_key_delete"
	val := []byte("this is a test")

	err = service.Set(key, bucket, val)
	if err != nil {
		t.Error("create key failed")
		return
	}

	err = service.Delete(key, bucket)
	if err != nil {
		t.Errorf("delete key failed %s", err.Error())
		return
	}

	err = service.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v := b.Get([]byte(key))
		if v != nil {
			t.Error("deleted key still exist")
		}
		return nil
	})

	// delete bucket and test again
	_ = service.DeleteBucket(bucket)

	err = service.Delete(key, bucket)
	if err == nil {
		t.Error("delete key should fail when bucket not exist")
	}
}

func TestService_BatchDelete(t *testing.T) {
	db := "/tmp/db_batch_delete"
	bucket := "test_batch_delete"
	service, err := createDBandBucket(db, bucket)
	if err != nil {
		t.Error("error in prepare db and bucket")
		return
	}
	defer os.Remove(db)

	kv := []KeyVal{{Key: "test1", Val: []byte("test1")}, {Key: "test2", Val: []byte("test2")}}

	for _, kVal := range kv {
		err = service.Set(kVal.Key, bucket, kVal.Val)
		if err != nil {
			t.Error("create key failed")
			return
		}
	}
	keys := []string{"tes1", "test2"}
	err = service.BatchDelete(keys, bucket)
	if err != nil {
		t.Error("batch delete failed")
	}

	err = service.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		for _, k := range keys {
			v := b.Get([]byte(k))
			if v != nil {
				t.Error("deleted key still exist")
			}
		}
		return nil
	})

	// delete bucket and test again
	_ = service.DeleteBucket(bucket)

	err = service.BatchDelete(keys, bucket)
	if err == nil {
		t.Error("batch delete key should fail when bucket not exist")
	}
}

func TestService_SetJsons(t *testing.T) {
	db := "/tmp/db_set_json"
	bucket := "test_set_json"
	service, err := createDBandBucket(db, bucket)
	if err != nil {
		t.Error("error in prepare db and bucket")
		return
	}
	defer os.Remove(db)

	type User struct {
		Username string
		Email    string
	}

	var user = User{
		Username: "test1",
		Email:    "test1@test.com",
	}

	err = service.SetJson(user.Username, bucket, user)
	if err != nil {
		t.Error("create key failed")
		return
	}

	data, err := service.GetOne(user.Username, bucket)
	if err != nil || data == nil {
		t.Error("json object not exist in db")
	}
}

func TestService_GetJson(t *testing.T) {
	db := "/tmp/db_get_json"
	bucket := "test_get"
	service, err := createDBandBucket(db, bucket)
	if err != nil {
		t.Error("error in prepare db and bucket")
		return
	}
	defer os.Remove(db)

	type User struct {
		Username string
		Email    string
	}

	var users = []User{
		{
			Username: "test1",
			Email:    "test1@test.com",
		},
		{
			Username: "test2",
			Email:    "test2@test.com",
		},
	}

	for _, user := range users {
		err = service.SetJson(user.Username, bucket, user)
		if err != nil {
			t.Error("create key failed")
			return
		}
	}

	var retUser User
	err = service.GetJson("test1", bucket, &retUser)
	if err != nil {
		t.Error("get json object failed")
	}

	if retUser.Username != users[0].Username || retUser.Email != users[0].Email {
		t.Errorf("get json object returned wrong data, %v", retUser)
	}
}

func TestService_GetJsonList(t *testing.T) {
	db := "/tmp/db_get_json_list"
	bucket := "test_get_json_list"
	service, err := createDBandBucket(db, bucket)
	if err != nil {
		t.Error("error in prepare db and bucket")
		return
	}
	defer os.Remove(db)

	type User struct {
		Username string
		Email    string
	}

	var users = []User{
		{
			Username: "test1",
			Email:    "test1@test.com",
		},
		{
			Username: "test2",
			Email:    "test2@test.com",
		},
	}

	for _, user := range users {
		err = service.SetJson(user.Username, bucket, user)
		if err != nil {
			t.Error("create user failed")
			return
		}
	}

	var retUsers []User

	err = service.GetJsonList(bucket, &retUsers)
	if err != nil {
		t.Errorf("get json list failed, %s", err)
		return
	}

	for i := 0; i < len(users); i++ {
		if users[i].Username != retUsers[i].Username || users[i].Email != retUsers[i].Email {
			t.Error("returned json list is not same as input")
		}
	}
}
