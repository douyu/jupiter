// Copyright 2022 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gorm

import (
	"context"
	"testing"
	"time"

	"fmt"

	"github.com/douyu/jupiter/pkg/core/imeta"
	"github.com/stretchr/testify/assert"
)

var (
	dsn = "root:123456@tcp(localhost:3306)/mysql?timeout=5s&readTimeout=5s&parseTime=true"
)

type LzTest struct {
	Model
	// ID            uint
	Name          string
	Sex           uint
	Age           uint
	Role          string
	LzTestRealted LzTestRealted `gorm:"foreignkey:ID"`
}

type LzTestRealted struct {
	Model
	Job      string
	Realname string
	Role     string
}

type Lz2_Test struct {
	Model
	Name  string
	Sex   uint
	Age   uint
	Role  string
	Email string
}

func (l LzTest) TableName() string {
	if l.Role == "admin" {
		return "admin_lztest"
	}
	if l.Role == "shadow" {
		return "lztest__stress_test__"
	}
	return "lztest"
}

func (l Lz2_Test) TableName() string {
	if l.Role == "admin" {
		return "admin_lz2test"
	}
	if l.Role == "shadow" {
		return "lz2test__stress_test__"
	}
	return "lz2test"
}

func (l LzTestRealted) TableName() string {
	if l.Role == "admin" {
		return "admin_lztestrealted"
	}
	if l.Role == "shadow" {
		return "lztestrelated__stress_test__"
	}
	return "lztestrelated"
}

func openTestingDB(t *testing.T, options *Config) (*DB, error) {
	got, err := open(options)
	if err != nil {
		return nil, err
	}
	got.Migrator().DropTable(&LzTest{})
	got.AutoMigrate(&LzTest{})

	got.Migrator().DropTable(&Lz2_Test{})
	got.AutoMigrate(&Lz2_Test{})

	got.Migrator().DropTable(&LzTestRealted{})
	got.AutoMigrate(&LzTestRealted{})

	return got, nil
}

func TestOpen(t *testing.T) {
	type args struct {
		dialect string
		options *Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				dialect: "mysql",
				options: &Config{
					DSN:             dsn,
					Debug:           false,
					MaxIdleConns:    10,
					MaxOpenConns:    100,
					ConnMaxLifetime: time.Second * 300,
					OnDialError:     "panic",
					SlowThreshold:   time.Millisecond * 300,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := openTestingDB(t, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.NotNil(t, got)
		})
	}
}

/*
create
*/
func TestCreate(t *testing.T) {
	//test data
	lztest := LzTest{
		Name: "lz1",
		Sex:  1,
		Age:  1,
	}
	lztestshadow := LzTest{
		Name: "lz1",
		Sex:  1,
		Age:  1,
	}

	//release
	got, err := openTestingDB(t, &Config{
		DSN:             dsn,
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Second * 300,
		OnDialError:     "panic",
		SlowThreshold:   time.Millisecond * 1,
		DetailSQL:       true,
	})
	assert.Nil(t, err)
	assert.Nil(t, got.Create(&lztest).Error)

	//stress
	gotshadow, err := openTestingDB(t, &Config{
		DSN:             dsn,
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Second * 300,
		OnDialError:     "panic",
		SlowThreshold:   time.Millisecond * 300,
		AutoShadowTable: true,
	})
	md := imeta.New(nil)
	md["x-dyctx-label"] = []string{"1"}
	ctx := imeta.WithContext(context.Background(), md)
	assert.Nil(t, err)
	assert.Nil(t, gotshadow.WithContext(ctx).Create(&lztestshadow).Error)
	assert.Equal(t, lztest.Name, lztestshadow.Name)

}

/*
FirstOrInit
*/
func TestFirstOrInit(t *testing.T) {
	//release
	got, err := openTestingDB(t, &Config{
		DSN:             dsn,
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Second * 300,
		OnDialError:     "panic",
		SlowThreshold:   time.Millisecond * 300,
	})
	assert.Nil(t, err)

	var lztest []LzTest
	var lztest1 LzTest
	assert.Nil(t, got.WithContext(context.Background()).Where(LzTest{Name: "LzFirstOrInit1"}).Attrs(LzTest{Age: 15}).FirstOrInit(&lztest1).Error)
	assert.Nil(t, got.WithContext(context.Background()).Create(&lztest1).Error)
	assert.Nil(t, got.WithContext(context.Background()).Where("name = ?", "LzFirstOrInit1").Find(&lztest).Error)
	assert.Equal(t, uint(15), lztest[0].Age)

	assert.Nil(t, got.WithContext(context.Background()).Where(LzTest{Name: "LzFirstOrInit1"}).Attrs("age", 30).FirstOrInit(&lztest1).Error)
	assert.Equal(t, uint(15), lztest1.Age)

	assert.Nil(t, got.WithContext(context.Background()).Where(LzTest{Name: "LzFirstOrInit1"}).Assign(map[string]interface{}{"age": 16}).FirstOrInit(&lztest1).Error)
	assert.Equal(t, uint(16), lztest1.Age)

	//stress
	gotshadow, err := openTestingDB(t, &Config{
		DSN:             dsn,
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Second * 300,
		OnDialError:     "panic",
		SlowThreshold:   time.Millisecond * 300,
		AutoShadowTable: true,
	})
	md := imeta.New(nil)
	md["x-dyctx-label"] = []string{"1"}
	ctx := imeta.WithContext(context.Background(), md)
	assert.Nil(t, err)

	var lztestshadow []LzTest
	var lztestshadow1 LzTest

	// where(struct) unsupport
	//assert.Nil(t,WithContext(ctx,gotshadow).Where(LzTest{Name:"LzFirstOrInit1"}).Attrs(LzTest{Age:15}).FirstOrInit(&lztestshadow1).Error)
	assert.Nil(t, gotshadow.WithContext(ctx).Where(map[string]interface{}{"name": "LzFirstOrInit1"}).Attrs(LzTest{Age: 15}).FirstOrInit(&lztestshadow1).Error)
	assert.Nil(t, gotshadow.WithContext(ctx).Create(&lztestshadow1).Error)
	assert.Nil(t, gotshadow.WithContext(ctx).Where("name = ?", "LzFirstOrInit1").Find(&lztestshadow).Error)
	assert.Equal(t, uint(15), lztestshadow[0].Age)

	//assert.Nil(t,WithContext(ctx,gotshadow).Where(LzTest{Name:"LzFirstOrInit1"}).Attrs("age",30).FirstOrInit(&lztestshadow1).Error)
	assert.Nil(t, gotshadow.WithContext(ctx).Where(map[string]interface{}{"name": "LzFirstOrInit1"}).Attrs("age", 30).FirstOrInit(&lztestshadow1).Error)
	assert.Equal(t, uint(15), lztestshadow1.Age)

	assert.Nil(t, gotshadow.WithContext(ctx).Where(map[string]interface{}{"name": "LzFirstOrInit1"}).Assign(map[string]interface{}{"age": 16}).FirstOrInit(&lztestshadow1).Error)
	assert.Equal(t, uint(16), lztestshadow1.Age)

}

/*
FirstOrCreate
*/
func TestFirstOrCreate(t *testing.T) {
	//release
	got, err := openTestingDB(t, &Config{
		DSN:             dsn,
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Second * 300,
		OnDialError:     "panic",
		SlowThreshold:   time.Millisecond * 300,
	})
	assert.Nil(t, err)

	var Lztest []LzTest
	var lztest1 LzTest
	assert.Nil(t, got.WithContext(context.Background()).Where(LzTest{Name: "LzFirstOrInit2"}).Attrs(LzTest{Age: 15}).FirstOrCreate(&lztest1).Error)
	assert.Nil(t, got.WithContext(context.Background()).Where("name = ?", "LzFirstOrInit2").Find(&Lztest).Error)
	assert.Equal(t, uint(15), Lztest[0].Age)

	assert.Nil(t, got.WithContext(context.Background()).Where(LzTest{Name: "LzFirstOrInit2"}).Attrs("age", 30).FirstOrCreate(&lztest1).Error)
	assert.Nil(t, got.WithContext(context.Background()).Where("name = ?", "LzFirstOrInit2").Find(&Lztest).Error)
	assert.Equal(t, uint(15), Lztest[0].Age)

	assert.Nil(t, got.WithContext(context.Background()).Where(LzTest{Name: "LzFirstOrInit2"}).Assign(map[string]interface{}{"age": 16}).FirstOrCreate(&lztest1).Error)
	assert.Nil(t, got.WithContext(context.Background()).Where("name = ?", "LzFirstOrInit2").Find(&Lztest).Error)
	assert.Equal(t, uint(16), Lztest[0].Age)

	//stress
	gotshadow, err := openTestingDB(t, &Config{
		DSN:             dsn,
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Second * 300,
		OnDialError:     "panic",
		SlowThreshold:   time.Millisecond * 300,
		AutoShadowTable: true,
	})
	md := imeta.New(nil)
	md["x-dyctx-label"] = []string{"1"}
	ctx := imeta.WithContext(context.Background(), md)
	assert.Nil(t, err)

	var Lztestshadow []LzTest
	var lztestshadow1 LzTest
	//assert.Nil(t,WithContext(ctx,gotshadow).Where(LzTest{Name:"LzFirstOrInit2"}).Attrs(LzTest{Age:15}).FirstOrCreate(&lztestshadow1).Error)
	assert.Nil(t, gotshadow.WithContext(ctx).Where(map[string]interface{}{"name": "LzFirstOrInit2"}).Attrs(LzTest{Age: 15}).FirstOrCreate(&lztestshadow1).Error)
	assert.Nil(t, gotshadow.WithContext(ctx).Where("name = ?", "LzFirstOrInit2").Find(&Lztestshadow).Error)
	assert.Equal(t, uint(15), Lztestshadow[0].Age)

	// where(struct) unsuport
	//assert.Nil(t,WithContext(ctx,gotshadow).Where(LzTest{Name:"LzFirstOrInit2"}).Attrs("age",30).FirstOrCreate(&lztestshadow1).Error)
	assert.Nil(t, gotshadow.WithContext(ctx).Where(map[string]interface{}{"name": "LzFirstOrInit2"}).Attrs("age", 30).FirstOrCreate(&lztestshadow1).Error)
	assert.Nil(t, gotshadow.WithContext(ctx).Where("name = ?", "LzFirstOrInit2").Find(&Lztestshadow).Error)
	assert.Equal(t, uint(15), Lztestshadow[0].Age)

	//assert.Nil(t,WithContext(ctx,gotshadow).Where(LzTest{Name:"LzFirstOrInit2"}).Assign(map[string]interface{}{"age": 16}).FirstOrCreate(&lztestshadow1).Error)
	assert.Nil(t, gotshadow.WithContext(ctx).Where("name = ?", "LzFirstOrInit2").Assign(map[string]interface{}{"age": 16}).FirstOrCreate(&lztestshadow1).Error)
	assert.Nil(t, gotshadow.WithContext(ctx).Where("name = ?", "LzFirstOrInit2").Find(&Lztestshadow).Error)
	assert.Equal(t, uint(16), Lztestshadow[0].Age)
}

/*
query
1.first
2.last
3.find
4.where //where(struct unsupport)
5.or
6.not
7.where or not
8.inner find
9.select order
10.multi-table query
11.related-table query
12.limit offset
13.group&having
14.pluck
15.preload
*/
func TestQuery(t *testing.T) {
	//release
	got, err := openTestingDB(t, &Config{
		DSN:             dsn,
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Second * 300,
		OnDialError:     "panic",
		SlowThreshold:   time.Millisecond * 300,
	})
	assert.Nil(t, err)
	assert.Nil(t, got.WithContext(context.Background()).Create(&LzTest{
		Name: "lzfirst",
		Sex:  1,
		Age:  100,
	}).Error)

	assert.Nil(t, got.WithContext(context.Background()).Create(&LzTest{
		Name: "lzfirst",
		Sex:  1,
		Age:  101,
	}).Error)

	assert.Nil(t, got.WithContext(context.Background()).Create(&LzTest{
		Name: "lzwhere",
		Sex:  1,
		Age:  100,
	}).Error)

	assert.Nil(t, got.WithContext(context.Background()).Create(&Lz2_Test{
		Name:  "lzfirst",
		Sex:   1,
		Age:   25,
		Email: "xxx",
	}).Error)

	assert.Nil(t, got.WithContext(context.Background()).Create(&LzTestRealted{
		Job:      "test",
		Realname: "lizhou",
	}).Error)

	assert.Nil(t, got.WithContext(context.Background()).Create(&LzTestRealted{
		Job:      "test",
		Realname: "lizhou2",
	}).Error)

	var lztestlen []LzTest
	//assert.Nil(t,got.WithContext(ctx).First(&lztestshadow,2).Error)
	//assert.Equal(t,uint(101),lztestshadow.Age)

	assert.Nil(t, got.WithContext(context.Background()).Where("name = ?", "lzfirst").First(&lztestlen).Error)
	assert.Equal(t, uint(100), lztestlen[0].Age)

	assert.Nil(t, got.WithContext(context.Background()).Where("name = ?", "lzfirst").Last(&lztestlen).Error)
	assert.Equal(t, uint(101), lztestlen[0].Age)

	assert.Nil(t, got.WithContext(context.Background()).Where("name = ?", "lzfirst").Find(&lztestlen).Error)
	assert.Equal(t, 2, len(lztestlen))

	assert.Nil(t, got.WithContext(context.Background()).Where("name <> ?", "lzfirst").Find(&lztestlen).Error)
	assert.Equal(t, 1, len(lztestlen))

	assert.Nil(t, got.WithContext(context.Background()).Where("name in (?)", []string{"lzwhere", "lzfirst"}).Find(&lztestlen).Error)
	assert.Equal(t, 3, len(lztestlen))

	assert.Nil(t, got.WithContext(context.Background()).Where("name LIKE ?", "%lzwhere%").Find(&lztestlen).Error)
	assert.Equal(t, 1, len(lztestlen))

	assert.Nil(t, got.WithContext(context.Background()).Where("name = ? AND age > ?", "lzfirst", 100).Find(&lztestlen).Error)
	assert.Equal(t, 1, len(lztestlen))

	//where 连用
	assert.Nil(t, got.WithContext(context.Background()).Where("name = ? ", "lzfirst").Where("age > ?", 100).Find(&lztestlen).Error)
	assert.Equal(t, 1, len(lztestlen))

	//where(struct)
	assert.Nil(t, got.WithContext(context.Background()).Where(&LzTest{Name: "lzfirst", Age: 100}).Find(&lztestlen).Error)
	assert.Equal(t, 1, len(lztestlen))

	//where(map)
	assert.Nil(t, got.WithContext(context.Background()).Where(map[string]interface{}{"name": "lzfirst", "age": 100}).Find(&lztestlen).Error)
	assert.Equal(t, 1, len(lztestlen))

	//or(struct)
	assert.Nil(t, got.WithContext(context.Background()).Where("name = ?", "lzfirst").Or(LzTest{Name: "lzwhere"}).Find(&lztestlen).Error)
	assert.Equal(t, 3, len(lztestlen))

	//or(map)
	assert.Nil(t, got.WithContext(context.Background()).Where("name = ?", "lzfirst").Or(map[string]interface{}{"name": "lzwhere"}).Find(&lztestlen).Error)
	assert.Equal(t, 3, len(lztestlen))

	//not
	assert.Nil(t, got.WithContext(context.Background()).Not("name", "lzfirst").Find(&lztestlen).Error)
	assert.Equal(t, 1, len(lztestlen))

	//not(slice)
	assert.Nil(t, got.WithContext(context.Background()).Not("name", []string{"lzfirst", "lzwhere"}).Find(&lztestlen).Error)
	assert.Equal(t, 0, len(lztestlen))

	//not(struct)
	assert.Nil(t, got.WithContext(context.Background()).Not(&LzTest{Name: "lzfirst"}).Find(&lztestlen).Error)
	assert.Equal(t, 1, len(lztestlen))

	//where or not(gorm bug)
	//assert.Nil(t,WithContext(context.Background(),got).Where("name = ?","lzfirst").Or("age >= ?","100").Not("name = ?","lzwhere").Find(&lztestlen).Error)
	//assert.Equal(t,2,len(lztestlen))

	//inner find
	assert.Nil(t, got.WithContext(context.Background()).Find(&lztestlen, "name <> ? AND age > ?", "lzwhere", 99).Error)
	assert.Equal(t, 2, len(lztestlen))

	//find(struct)
	assert.Nil(t, got.WithContext(context.Background()).Find(&lztestlen, LzTest{Age: 100}).Error)
	assert.Equal(t, 2, len(lztestlen))

	//find(map[])
	assert.Nil(t, got.WithContext(context.Background()).Find(&lztestlen, map[string]interface{}{"age": 101}).Error)
	assert.Equal(t, 1, len(lztestlen))

	//select order
	assert.Nil(t, got.WithContext(context.Background()).Select("sex,age").Where("name = ?", "lzfirst").Order("age desc").Find(&lztestlen).Error)
	assert.Equal(t, uint(100), lztestlen[1].Age)
	assert.Equal(t, "", lztestlen[1].Name)

	assert.Nil(t, got.WithContext(context.Background()).Select([]string{"sex,age"}).Where("name = ?", "lzfirst").Order("age desc").Find(&lztestlen).Error)
	assert.Equal(t, uint(100), lztestlen[1].Age)
	assert.Equal(t, "", lztestlen[1].Name)

	//multi-table query
	got.WithContext(context.Background()).Table("lztest").Select("lztest.name,lz2test.name").Joins("left join lz2test on lz2test.name = lztest.name").Find(&lztestlen)
	assert.Equal(t, "", lztestlen[2].Name)

	//related query
	var lztest LzTest
	assert.Nil(t, got.WithContext(context.Background()).Where("id = ?", 2).Preload("LzTestRealted").First(&lztest).Error)
	// assert.Nil(t, got.WithContext(context.Background()).Model(&lztest).Related(&lztest.LzTestRealted, "ID").Error)
	assert.Equal(t, "lizhou2", lztest.LzTestRealted.Realname)

	//limit offset count
	var count int64
	assert.Nil(t, got.WithContext(context.Background()).Limit(1).Offset(1).Order("age desc").Find(&lztestlen).Error)
	assert.Nil(t, got.WithContext(context.Background()).Table("lztest").Limit(1).Count(&count).Error)
	assert.Nil(t, got.WithContext(context.Background()).Limit(2).Offset(1).Find(&lztestlen).Error)
	assert.Equal(t, "lzwhere", lztestlen[1].Name)
	assert.Equal(t, int64(3), count)

	//group&having
	rows, err := got.WithContext(context.Background()).Table("lztest").Select("name,sum(age) as totalage").Group("name").Having("sum(age) > ?", 200).Rows()
	assert.Nil(t, err)
	for rows.Next() {
		fmt.Print(rows.Columns())
	}

	//pluck
	var names []string
	assert.Nil(t, got.WithContext(context.Background()).Model(&LzTest{}).Pluck("name", &names).Error)
	assert.Equal(t, "lzfirst", names[1])

	//preload
	assert.Nil(t, got.WithContext(context.Background()).Where("name = ?", "lzfirst").Preload("LzTestRealted", "realname NOT IN (?)", "lizhou2").Find(&lztestlen).Error)
	assert.Equal(t, "lizhou", lztestlen[0].LzTestRealted.Realname)

	//stress
	gotshadow, err := openTestingDB(t, &Config{
		DSN:             dsn,
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Second * 300,
		OnDialError:     "panic",
		SlowThreshold:   time.Millisecond * 300,
		AutoShadowTable: true,
	})
	md := imeta.New(nil)
	md["x-dyctx-label"] = []string{"1"}
	ctx := imeta.WithContext(context.Background(), md)
	assert.Nil(t, err)

	assert.Nil(t, gotshadow.WithContext(ctx).Create(&LzTest{
		Name: "lzfirst",
		Sex:  1,
		Age:  100,
	}).Error)

	assert.Nil(t, gotshadow.WithContext(ctx).Create(&LzTest{
		Name: "lzfirst",
		Sex:  1,
		Age:  101,
	}).Error)

	assert.Nil(t, gotshadow.WithContext(ctx).Create(&LzTest{
		Name: "lzwhere",
		Sex:  1,
		Age:  100,
	}).Error)

	assert.Nil(t, gotshadow.WithContext(ctx).Create(&Lz2_Test{
		Name:  "lzfirst",
		Sex:   1,
		Age:   25,
		Email: "xxx",
	}).Error)

	assert.Nil(t, gotshadow.WithContext(ctx).Create(&LzTestRealted{
		Job:      "test",
		Realname: "lizhou",
	}).Error)

	assert.Nil(t, gotshadow.WithContext(ctx).Create(&LzTestRealted{
		Job:      "test",
		Realname: "lizhou2",
	}).Error)

	var lztestshadowlen []LzTest
	//assert.Nil(t,got.WithContext(ctx).First(&lztestshadow,2).Error)
	//assert.Equal(t,uint(101),lztestshadow.Age)

	assert.Nil(t, gotshadow.WithContext(ctx).Where("name = ?", "lzfirst").First(&lztestshadowlen).Error)
	assert.Equal(t, uint(100), lztestshadowlen[0].Age)

	assert.Nil(t, gotshadow.WithContext(ctx).Where("name = ?", "lzfirst").Last(&lztestshadowlen).Error)
	assert.Equal(t, uint(101), lztestshadowlen[0].Age)

	assert.Nil(t, gotshadow.WithContext(ctx).Where("name = ?", "lzfirst").Find(&lztestshadowlen).Error)
	assert.Equal(t, 2, len(lztestshadowlen))

	assert.Nil(t, gotshadow.WithContext(ctx).Where("name <> ?", "lzfirst").Find(&lztestshadowlen).Error)
	assert.Equal(t, 1, len(lztestshadowlen))

	assert.Nil(t, gotshadow.WithContext(ctx).Where("name in (?)", []string{"lzwhere", "lzfirst"}).Find(&lztestshadowlen).Error)
	assert.Equal(t, 3, len(lztestshadowlen))

	assert.Nil(t, gotshadow.WithContext(ctx).Where("name LIKE ?", "%lzwhere%").Find(&lztestshadowlen).Error)
	assert.Equal(t, 1, len(lztestshadowlen))

	assert.Nil(t, gotshadow.WithContext(ctx).Where("name = ? AND age > ?", "lzfirst", 100).Find(&lztestshadowlen).Error)
	assert.Equal(t, 1, len(lztestshadowlen))

	//where 连用
	assert.Nil(t, gotshadow.WithContext(ctx).Where("name = ? ", "lzfirst").Where("age > ?", 100).Find(&lztestshadowlen).Error)
	assert.Equal(t, 1, len(lztestshadowlen))

	//where(struct)--unsupport
	//assert.Nil(t,WithContext(ctx,gotshadow).Where(&LzTest{Name:"lzfirst",Age:100}).Find(&lztestshadowlen).Error)
	//assert.Equal(t,2,len(lztestshadowlen))

	//where(map)
	assert.Nil(t, gotshadow.WithContext(ctx).Where(map[string]interface{}{"name": "lzfirst", "age": 100}).Find(&lztestshadowlen).Error)
	assert.Equal(t, 1, len(lztestshadowlen))

	//or(struct)-unsupport
	//assert.Nil(t,WithContext(ctx,gotshadow).Where("name = ?","lzfirst").Or(LzTest{Name:"lzwhere"}).Find(&lztestshadowlen).Error)
	//assert.Equal(t,3,len(lztestshadowlen))

	//or(map) --unsupport
	assert.Nil(t, gotshadow.WithContext(ctx).Where("name = ?", "lzfirst").Or(map[string]interface{}{"name": "lzwhere"}).Find(&lztestshadowlen).Error)
	assert.Equal(t, 3, len(lztestshadowlen))

	//not
	assert.Nil(t, gotshadow.WithContext(ctx).Not("name", "lzfirst").Find(&lztestshadowlen).Error)
	assert.Equal(t, 1, len(lztestshadowlen))

	//not(slice)
	assert.Nil(t, gotshadow.WithContext(ctx).Not("name", []string{"lzfirst", "lzwhere"}).Find(&lztestshadowlen).Error)
	assert.Equal(t, 0, len(lztestshadowlen))

	////not(struct)--unsupport
	//assert.Nil(t,WithContext(ctx,gotshadow).Not(&LzTest{Name:"lzfirst"}).Find(&lztestshadowlen).Error)
	//assert.Equal(t,1,len(lztestshadowlen))
	//
	////where or not--gorm bug
	//assert.Nil(t,WithContext(ctx,gotshadow).Where("name = ?","lzfirst").Or("age >= ?","100").Not("name = ?","lzwhere").Find(&lztestshadowlen).Error)
	//assert.Equal(t,2,len(lztestshadowlen))

	//inner find
	assert.Nil(t, gotshadow.WithContext(ctx).Find(&lztestshadowlen, "name <> ? AND age > ?", "lzwhere", 99).Error)
	assert.Equal(t, 2, len(lztestshadowlen))

	//find(struct)-unsupport
	//assert.Nil(t,WithContext(ctx,gotshadow).Find(&lztestshadowlen,LzTest{Age:100}).Error)
	//assert.Equal(t,1,len(lztestshadowlen))

	//find(map[])
	assert.Nil(t, gotshadow.WithContext(ctx).Find(&lztestshadowlen, map[string]interface{}{"age": 101}).Error)
	assert.Equal(t, 1, len(lztestshadowlen))

	//select order
	assert.Nil(t, gotshadow.WithContext(ctx).Select("sex,age").Where("name = ?", "lzfirst").Order("age desc").Find(&lztestshadowlen).Error)
	assert.Equal(t, uint(100), lztestshadowlen[1].Age)
	assert.Equal(t, "", lztestshadowlen[1].Name)

	assert.Nil(t, gotshadow.WithContext(ctx).Select([]string{"sex,age"}).Where("name = ?", "lzfirst").Order("age desc").Find(&lztestshadowlen).Error)
	assert.Equal(t, uint(100), lztestshadowlen[1].Age)
	assert.Equal(t, "", lztestshadowlen[1].Name)

	//multi-table query--unsupport
	//WithContext(ctx,gotshadow).Table("lztest").Select("lztest.name,lz2_test.name").Joins("left join lz2_test on lz2_test.name = lztest.name").Find(&lztestshadowlen)
	//assert.Equal(t,2,len(lztestshadowlen))

	//related query
	var lztestshadow LzTest
	assert.Nil(t, gotshadow.WithContext(ctx).Where("id = ?", 2).Preload("LzTestRealted").First(&lztestshadow).Error)
	// assert.Nil(t, gotshadow.WithContext(ctx).Model(&lztestshadow).Related(&lztestshadow.LzTestRealted, "ID").Error)
	assert.Equal(t, "lizhou2", lztestshadow.LzTestRealted.Realname)

	//limit offset count
	var countshadow int64
	assert.Nil(t, gotshadow.WithContext(ctx).Limit(1).Offset(1).Order("age desc").Find(&lztestshadowlen).Error)
	assert.Nil(t, gotshadow.WithContext(ctx).Table("lztest").Limit(1).Count(&countshadow).Error)
	assert.Nil(t, gotshadow.WithContext(ctx).Limit(2).Offset(1).Find(&lztestshadowlen).Error)
	assert.Equal(t, "lzwhere", lztestshadowlen[1].Name)
	assert.Equal(t, int64(3), countshadow)

	//group&having
	rowsshadow, err := gotshadow.WithContext(ctx).Table("lztest").Select("name,sum(age) as totalage").Group("name").Having("sum(age) > ?", 200).Rows()
	assert.Nil(t, err)
	for rowsshadow.Next() {
		fmt.Print(rowsshadow.Columns())
	}

	//pluck
	var namesshadow []string
	assert.Nil(t, gotshadow.WithContext(ctx).Model(&LzTest{}).Pluck("name", &namesshadow).Error)
	assert.Equal(t, "lzfirst", namesshadow[1])

	//preload
	assert.Nil(t, gotshadow.WithContext(ctx).Where("name = ?", "lzfirst").Preload("LzTestRealted", "realname NOT IN (?)", "lizhou2").Find(&lztestshadowlen).Error)
	assert.Equal(t, "lizhou", lztestshadowlen[0].LzTestRealted.Realname)
}

/*
update
*/

func TestUpdate(t *testing.T) {
	//test data
	lztest := LzTest{
		Name: "lz1",
		Sex:  1,
		Age:  1,
	}
	lztestshadow := LzTest{
		Name: "lz1",
		Sex:  1,
		Age:  1,
	}

	//release
	got, err := openTestingDB(t, &Config{
		DSN:             dsn,
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Second * 300,
		OnDialError:     "panic",
		SlowThreshold:   time.Millisecond * 300,
	})
	assert.Nil(t, err)
	assert.Nil(t, got.WithContext(context.Background()).Create(&lztest).Error)

	var lztest1 LzTest
	assert.Nil(t, got.WithContext(context.Background()).Table("lztest").Where("name = ?", "lz1").Update("age", 2).Scan(&lztest1).Error)

	assert.Nil(t, got.WithContext(context.Background()).Table("lztest").Where("name= ? ", "lz1").Scan(&lztest1).Error)
	assert.Equal(t, uint(2), lztest1.Age)

	//stress
	gotshadow, err := openTestingDB(t, &Config{
		DSN:             dsn,
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Second * 300,
		OnDialError:     "panic",
		SlowThreshold:   time.Millisecond * 300,
		AutoShadowTable: true,
	})
	md := imeta.New(nil)
	md["x-dyctx-label"] = []string{"1"}
	ctx := imeta.WithContext(context.Background(), md)
	assert.Nil(t, err)

	assert.Nil(t, gotshadow.WithContext(ctx).Create(&lztestshadow).Error)

	var lztestshadow1 LzTest
	assert.Nil(t, gotshadow.WithContext(ctx).Table("lztest").Where("name = ?", "lz1").Update("age", 2).Scan(&lztestshadow1).Error)

	assert.Nil(t, gotshadow.WithContext(ctx).Table("lztest").Where("name= ? ", "lz1").Scan(&lztestshadow1).Error)
	assert.Equal(t, uint(2), lztestshadow1.Age)

}

/*
delete
*/

func TestDelete(t *testing.T) {
	//test data
	lztest := LzTest{
		Name: "lz1",
		Sex:  1,
		Age:  1,
	}
	lztestshadow := LzTest{
		Name: "lz1",
		Sex:  1,
		Age:  1,
	}

	//release
	got, err := openTestingDB(t, &Config{
		DSN:             dsn,
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Second * 300,
		OnDialError:     "panic",
		SlowThreshold:   time.Millisecond * 300,
	})
	assert.Nil(t, err)
	assert.Nil(t, got.WithContext(context.Background()).Create(&lztest).Error)

	var lztestlen []LzTest
	assert.Nil(t, got.WithContext(context.Background()).Find(&lztestlen).Error)
	assert.Equal(t, 1, len(lztestlen))

	assert.Nil(t, got.WithContext(context.Background()).Delete(lztestlen, "name = ?", "lz1").Error)

	assert.Nil(t, got.WithContext(context.Background()).Find(&lztestlen).Error)
	assert.Equal(t, 0, len(lztestlen))

	//stress
	gotshadow, err := openTestingDB(t, &Config{
		DSN:             dsn,
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Second * 300,
		OnDialError:     "panic",
		SlowThreshold:   time.Millisecond * 300,
		AutoShadowTable: true,
	})
	md := imeta.New(nil)
	md["x-dyctx-label"] = []string{"1"}
	ctx := imeta.WithContext(context.Background(), md)
	assert.Nil(t, err)

	assert.Nil(t, gotshadow.WithContext(ctx).Create(&lztestshadow).Error)

	var lztestshadowlen []LzTest
	assert.Nil(t, gotshadow.WithContext(ctx).Find(&lztestshadowlen).Error)
	assert.Equal(t, 1, len(lztestshadowlen))

	assert.Nil(t, gotshadow.WithContext(ctx).Delete(lztestshadowlen, "name = ?", "lz1").Error)

	assert.Nil(t, gotshadow.WithContext(ctx).Find(&lztestshadowlen).Error)
	assert.Equal(t, 0, len(lztestshadowlen))

}

/*
row sql(unsupported---cannot recognize tablename)
*/

func TestRowSql(t *testing.T) {

	//release
	got, err := openTestingDB(t, &Config{
		DSN:             dsn,
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Second * 300,
		OnDialError:     "panic",
		SlowThreshold:   time.Millisecond * 300,
	})
	assert.Nil(t, err)

	var sql string
	var lztests []LzTest
	var name string
	var name_update string
	name = "lzsql01"
	name_update = "lzsql123"
	sql = "INSERT INTO lztest(`name`,`sex`,`age`,`role`) values ('lzsql01','1','1','lz');"
	assert.Nil(t, got.WithContext(context.Background()).Raw(sql).Scan(&lztests).Error)

	sql = fmt.Sprintf("select * from lztest where name = \"%s\"", name)
	assert.Nil(t, got.WithContext(context.Background()).Raw(sql).Scan(&lztests).Error)
	assert.Equal(t, name, lztests[0].Name)

	sql = fmt.Sprintf("update lztest set name = \"%s\" where name = \"%s\"", name_update, name)
	assert.Nil(t, got.WithContext(context.Background()).Raw(sql).Scan(&lztests).Error)
	sql = fmt.Sprintf("select * from lztest where name = \"%s\"", name_update)
	assert.Nil(t, got.WithContext(context.Background()).Raw(sql).Scan(&lztests).Error)
	assert.Equal(t, name_update, lztests[0].Name)

	sql = fmt.Sprintf("delete from lztest where name = \"%s\"", name_update)
	assert.Nil(t, got.WithContext(context.Background()).Raw(sql).Debug().Scan(&lztests).Error)
	sql = fmt.Sprintf("select * from lztest where name = \"%s\"", name_update)
	// TODO：确认gorm scan的逻辑是否符合预期
	lztests = []LzTest{}
	assert.Nil(t, got.WithContext(context.Background()).Raw(sql).Debug().Scan(&lztests).Error)
	assert.Equal(t, 0, len(lztests))
	fmt.Println("=======delete:", len(lztests))

	//stress
	gotshadow, err := openTestingDB(t, &Config{
		DSN:             dsn,
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Second * 300,
		OnDialError:     "panic",
		SlowThreshold:   time.Millisecond * 300,
		AutoShadowTable: true,
	})
	md := imeta.New(nil)
	md["x-dyctx-label"] = []string{"1"}
	ctx := imeta.WithContext(context.Background(), md)
	assert.Nil(t, err)

	var sqlshadow string
	var lztestshadows []LzTest
	var nameshadow string
	var name_update_shadow string
	nameshadow = "lzsql01"
	name_update_shadow = "lzsql123"
	sqlshadow = "INSERT INTO lztest(`name`,`sex`,`age`,`role`) values ('lzsql01','1','1','lz');"
	assert.Nil(t, gotshadow.WithContext(ctx).Raw(sqlshadow).Scan(&lztestshadows).Error)

	sqlshadow = fmt.Sprintf("select * from lztest where name = \"%s\"", nameshadow)
	assert.Nil(t, gotshadow.WithContext(ctx).Raw(sqlshadow).Scan(&lztestshadows).Error)
	assert.Equal(t, nameshadow, lztestshadows[0].Name)

	sqlshadow = fmt.Sprintf("update lztest set name = \"%s\" where name = \"%s\"", name_update_shadow, nameshadow)
	assert.Nil(t, gotshadow.WithContext(ctx).Raw(sqlshadow).Scan(&lztestshadows).Error)
	sqlshadow = fmt.Sprintf("select * from lztest where name = \"%s\"", name_update_shadow)
	assert.Nil(t, gotshadow.WithContext(ctx).Raw(sqlshadow).Scan(&lztestshadows).Error)
	assert.Equal(t, name_update_shadow, lztestshadows[0].Name)

	sqlshadow = fmt.Sprintf("delete from lztest where name = \"%s\"", name_update_shadow)
	assert.Nil(t, gotshadow.WithContext(ctx).Raw(sqlshadow).Scan(&lztestshadows).Error)
	sqlshadow = fmt.Sprintf("select * from lztest where name = \"%s\"", name_update_shadow)
	// TODO：确认gorm scan的逻辑是否符合预期
	lztestshadows = []LzTest{}
	assert.Nil(t, gotshadow.WithContext(ctx).Raw(sqlshadow).Scan(&lztestshadows).Error)
	assert.Equal(t, 0, len(lztestshadows))
	fmt.Println("=======delete:", len(lztestshadows))

}
