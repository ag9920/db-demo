package biz

import (
	"context"
	"fmt"

	"gorm.io/hints"

	"github.com/ag9920/db-demo/gendemo/dal"
	"github.com/ag9920/db-demo/gendemo/dal/model"
	"github.com/ag9920/db-demo/gendemo/dal/query"
)

var q = query.Use(dal.DB.Debug())

func Create(ctx context.Context) {
	var err error
	ud := q.User.WithContext(ctx)

	userData := &model.User{ID: 1, Name: "modi"}
	// INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`name`,`age`,`role`,`id`) VALUES ('2021-09-13 20:05:51.389','2021-09-13 20:05:51.389',NULL,'modi',0,'',1)
	err = ud.Create(userData)

	userDataArray := []*model.User{{ID: 2, Name: "A"}, {ID: 3, Name: "B"}}
	err = ud.CreateInBatches(userDataArray, 2)
	// INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`name`,`age`,`role`,`id`) VALUES ('2021-09-13 20:05:51.403','2021-09-13 20:05:51.403',NULL,'A',0,'',2),('2021-09-13 20:05:51.403','2021-09-13 20:05:51.403',NULL,'B',0,'',3)

	userData.Name = "new name"
	err = ud.Save(userData)
	// INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`name`,`age`,`role`,`id`) VALUES ('2021-09-13 20:05:51.389','2021-09-13 20:05:51.409',NULL,'new name',0,'',1) ON DUPLICATE KEY UPDATE `updated_at`=VALUES(`updated_at`),`deleted_at`=VALUES(`deleted_at`),`name`=VALUES(`name`),`age`=VALUES(`age`),`role`=VALUES(`role`)

	fmt.Println(err)
}

func Delete(ctx context.Context) {
	var err error
	u, ud := q.User, q.User.WithContext(ctx)

	_, err = ud.Where(u.ID.Eq(1)).Delete()
	// UPDATE `users` SET `deleted_at`='2021-09-13 20:05:51.418' WHERE `users`.`id` = 1 AND `users`.`deleted_at` IS NULL

	_, err = ud.Where(u.ID.In(2, 3)).Delete()
	// UPDATE `users` SET `deleted_at`='2021-09-13 20:05:51.428' WHERE `users`.`id` IN (2,3) AND `users`.`deleted_at` IS NULL

	_, err = ud.Where(u.ID.Gt(100)).Unscoped().Delete()
	// DELETE FROM `users` WHERE `users`.`id` > 100

	fmt.Println(err)
}

func Query(ctx context.Context) {
	var err error
	var user *model.User
	var users []*model.User

	u, ud := q.User, q.User.WithContext(ctx)

	/*--------------Basic query-------------*/
	user, err = ud.Take()
	// SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL LIMIT 1
	fmt.Printf("query 1 item: %+v", user)

	user, err = ud.Where(u.ID.Gt(100), u.Name.Like("%T%")).Take()
	// SELECT * FROM `users` WHERE `users`.`id` > 100 AND `users`.`name` LIKE '%T%' AND `users`.`deleted_at` IS NULL LIMIT 1
	fmt.Printf("query conditions got: %+v", user)

	user, err = ud.Where(ud.Columns(u.ID).In(ud.Select(u.ID.Min()))).First()
	// SELECT * FROM `users` WHERE `users`.`id` IN (SELECT MIN(`users`.`id`) FROM `users` WHERE `users`.`deleted_at` IS NULL) AND `users`.`deleted_at` IS NULL
	// ORDER BY `users`.`id` LIMIT 1
	fmt.Printf("subquery 1 got item: %+v", user)

	user, err = ud.Where(ud.Columns(u.ID).Eq(ud.Select(u.ID.Max()))).First()
	// SELECT * FROM `users` WHERE `users`.`id` = (SELECT MAX(`users`.`id`) FROM `users` WHERE `users`.`deleted_at` IS NULL) AND `users`.`deleted_at` IS NULL
	// ORDER BY `users`.`id` LIMIT 1
	fmt.Printf("subquery 2 got item: %+v", user)

	users, err = ud.Distinct(u.Name).Find()
	// SELECT DISTINCT `users`.`name` FROM `users` WHERE `users`.`deleted_at` IS NULL
	fmt.Printf("select distinct got: %d", len(users))

	/*--------------Diy query-------------*/
	user, err = ud.FindByNameAndAge("tom", 29)
	// SELECT * FROM `users` WHERE name='tom' and age=29 AND `users`.`deleted_at` IS NULL
	fmt.Printf("FindByNameAndAge: %+v", user)

	users, err = ud.FindBySimpleName()
	// select id,name,age from users where age>18
	fmt.Printf("FindBySimpleName: (%d)%+v", len(users), users)

	user, err = ud.FindByIDOrName(false, 0, "tom", "user")
	// select id,name,age from users where age>18
	fmt.Printf("FindByIDOrName: %+v", user)

	/*--------------Advanced query-------------*/
	users, err = ud.Clauses(hints.New("MAX_EXECUTION_TIME(10000)")).Find()
	// SELECT /*+ MAX_EXECUTION_TIME(10000) */ * FROM `users` WHERE `users`.`deleted_at` IS NULL
	fmt.Printf("find with hints 2: (%d)%+v", len(users), users)

	users, err = ud.Clauses(hints.New("idx_user_name")).Find()
	// SELECT /*+ idx_user_name */ * FROM `users` WHERE `users`.`deleted_at` IS NULL
	fmt.Printf("find with hints 2: (%d)%+v", len(users), users)

	users, err = ud.Clauses(hints.New("hint")).Select(u.ID, u.Name).Where(u.ID.IsNotNull(), u.Age.Gt(18)).Find()
	// SELECT `users`.`id`,`users`.`name` FROM `users` WHERE `users`.`id` IS NOT NULL AND `users`.`age` > 18 AND `users`.`deleted_at` IS NULL
	fmt.Printf("find with hints 3: (%d)%+v", len(users), users)

	user, err = ud.Select(u.ID, u.Name).Where(u.ID.Eq(1)).FirstOrInit()
	fmt.Printf("FirstOrInit got: %+v", user)

	user, err = ud.Select(u.ID, u.Name).Where(u.ID.Eq(1)).Attrs(u.Name.Value("modi")).FirstOrInit()
	fmt.Printf("FirstOrInit got: %+v", user)

	user, err = ud.Select(u.ID, u.Name).Where(u.ID.Eq(1)).Attrs(u.Name.Value("modi")).Assign(u.Age.Value(17)).FirstOrInit()
	fmt.Printf("FirstOrInit got: %+v", user)

	user, err = ud.Select(u.ID, u.Name).Where(u.ID.Eq(1)).FirstOrCreate()
	fmt.Printf("FirstOrCreate got: %+v", user)

	user, err = ud.Select(u.ID, u.Name).Where(u.ID.Eq(1)).Attrs(u.Name.Value("modi")).FirstOrCreate()
	fmt.Printf("FirstOrCreate got: %+v", user)

	user, err = ud.Select(u.ID, u.Name).Where(u.ID.Eq(1)).Attrs(u.Name.Value("modi")).Assign(u.Age.Value(17)).FirstOrCreate()
	// UPDATE `users` SET `age`=17 WHERE `users`.`id` = 1 AND `users`.`deleted_at` IS NULL
	fmt.Printf("FirstOrCreate got: %+v", user)

	fmt.Println(err)
}

func Update(ctx context.Context) {
	var err error

	u, ud := q.User, q.User.WithContext(ctx)

	user, err := ud.First()
	// SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1

	user.Name = "save test"
	err = ud.Save(user)
	// INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`name`,`age`,`role`,`id`) VALUES ('2021-09-13 20:12:18.655','2021-09-13 20:12:18.655',NULL,'save test',190,'',4) ON DUPLICATE KEY UPDATE `updated_at`=VALUES(`updated_at`),`deleted_at`=VALUES(`deleted_at`),`name`=VALUES(`name`),`age`=VALUES(`age`),`role`=VALUES(`role`)

	_, err = ud.Where(u.ID.Eq(user.ID)).Update(u.Name, "update test")
	// UPDATE `users` SET `name`='update test',`updated_at`='2021-09-13 20:12:18.664' WHERE `users`.`id` = 4 AND `users`.`deleted_at` IS NULL

	_, err = ud.Where(u.ID.Eq(user.ID)).Updates(model.User{Name: "updates test"})
	// UPDATE `users` SET `updated_at`='2021-09-28 20:14:41.139',`name`='updates test' WHERE `users`.`id` = 4 AND `users`.`deleted_at` IS NULL

	_, err = ud.Where(u.ID.Eq(user.ID)).UpdateSimple(u.Age.Add(1), u.CreatedAt.Null(), u.Name.Value("modi"), u.UpdatedAt.Zero())
	// UPDATE `users` SET `age`=`users`.`age`+1,`created_at`=NULL,`name`='modi',`updated_at`='0000-00-00 00:00:00'
	// WHERE `users`.`id` = 4 AND `users`.`deleted_at` IS NULL

	fmt.Println(err)
}
