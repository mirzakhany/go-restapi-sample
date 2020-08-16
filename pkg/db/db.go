package db

import (
	"encoding/json"
	"fmt"
	"reflect"

	"go.etcd.io/bbolt"
)

type KeyVal struct {
	Key string
	Val []byte
}

const ContextKey = "db"

// Service will hold bbolt db and its settings
type Service struct {
	path string
	DB   *bbolt.DB
}

func New(path string) (*Service, error) {
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}
	return &Service{
		path: path,
		DB:   db,
	}, nil
}

func (s *Service) Close() error {
	return s.DB.Close()
}

func (s *Service) CreateBucket(bucketName string) (*bbolt.Bucket, error) {
	var b *bbolt.Bucket
	var err error
	err = s.DB.Update(func(tx *bbolt.Tx) error {
		b, err = tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: `%s`", err)
		}
		return nil
	})
	return b, err
}

func (s *Service) IsExist(key, bucketName string) (bool, error) {
	var exist bool
	err := s.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket `%s` not exist", bucketName)
		}
		exist = b.Get([]byte(key)) != nil
		return nil
	})
	return exist, err
}

func (s *Service) DeleteBucket(bucketName string) error {
	err := s.DB.Update(func(tx *bbolt.Tx) error {
		err := tx.DeleteBucket([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("delete bucket: `%s`", err)
		}
		return nil
	})
	return err
}

func (s *Service) Set(key, bucketName string, value []byte) error {
	err := s.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket `%s` not exist", bucketName)
		}
		return b.Put([]byte(key), value)
	})
	return err
}

func (s *Service) GetOne(key, bucketName string) ([]byte, error) {
	var v []byte
	err := s.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket `%s` not exist", bucketName)
		}
		v = b.Get([]byte(key))
		if v == nil {
			return fmt.Errorf("key `%s` not exist", key)
		}
		return nil
	})
	return v, err
}

func (s *Service) GetAll(bucketName string) ([]KeyVal, error) {
	var result []KeyVal
	err := s.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket `%s` not exist", bucketName)
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			result = append(result, KeyVal{
				Key: string(k),
				Val: v,
			})
		}
		return nil
	})
	return result, err
}

func (s *Service) Delete(key, bucketName string) error {
	err := s.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket `%s` not exist", bucketName)
		}
		return b.Delete([]byte(key))
	})
	return err
}

func (s *Service) BatchDelete(keys []string, bucketName string) error {
	err := s.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket `%s` not exist", bucketName)
		}
		for _, key := range keys {
			err := b.Delete([]byte(key))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (s *Service) SetJson(key, bucketName string, value interface{}) error {
	// Marshal and save the encoded data.
	buf, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.Set(key, bucketName, buf)
}

func (s *Service) GetJson(key, bucketName string, ret interface{}) error {
	data, err := s.GetOne(key, bucketName)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &ret); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetJsonList(bucketName string, ret interface{}) error {

	v := reflect.ValueOf(ret)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("non-pointer %v", v.Type())
	}
	// get the value that the pointer v points to.
	v = v.Elem()
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("can't fill non-slice value")
	}

	data, err := s.GetAll(bucketName)
	if err != nil {
		return err
	}
	dataLen := len(data)

	v.Set(reflect.MakeSlice(v.Type(), dataLen, dataLen))
	realVal := reflect.New(v.Type().Elem()).Interface()

	for i, kVal := range data {
		if err := json.Unmarshal(kVal.Val, realVal); err != nil {
			return err
		}
		v.Index(i).Set(reflect.ValueOf(realVal).Elem())
	}
	return nil
}
