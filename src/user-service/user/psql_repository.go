package user

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/user/model"

	_ "github.com/lib/pq"
)

type PsqlRepository struct {
	db *sql.DB
}

func NewPsqlRepository(config database.Config) (*PsqlRepository, error) {
	dsn := config.Dsn()
	fmt.Println(dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &PsqlRepository{db}, nil
}

const createUsersTable = `
create table if not exists users (
    id				serial primary key, 
	email			varchar(100) not null unique,
	username    	varchar(16) not null unique,
	password 		bytea not null,
	profile_name 	varchar(100) not null,
	balance 		int not null default 0
                            
)
`

func (repo *PsqlRepository) Migrate() error {
	_, err := repo.db.Exec(createUsersTable)
	return err
}

const createUsersBatchQuery = `
insert into users (email, username, password, profile_name) values %s
`

func (repo *PsqlRepository) Create(users []*model.DbUser) error {
	placeholders := make([]string, len(users))
	values := make([]interface{}, len(users)*4)

	for i := 0; i < len(users); i++ {
		placeholders[i] = fmt.Sprintf("($%d,$%d,$%d,$%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		values[i*4+0] = users[i].Email
		values[i*4+1] = users[i].Username
		values[i*4+2] = users[i].Password
		values[i*4+3] = users[i].ProfileName
	}

	query := fmt.Sprintf(createUsersBatchQuery, strings.Join(placeholders, ","))
	_, err := repo.db.Exec(query, values...)
	return err
}

const updateUserBatchQuery = `
update users set profile_name = $1, balance = $2 where username = $3 returning id
`

func (repo *PsqlRepository) Update(username string, user *model.UpdateUser) error {
	_, err := repo.db.Exec(updateUserBatchQuery, user.ProfileName, user.Balance, username)
	return err
}

const findAllUsersQuery = `
select id, email, password, username, profile_name, balance from users
`

func (repo *PsqlRepository) FindAll() ([]*model.DbUser, error) {
	rows, err := repo.db.Query(findAllUsersQuery)
	if err != nil {
		return nil, err
	}

	var users []*model.DbUser
	for rows.Next() {
		user := model.DbUser{}
		if err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.ProfileName, &user.Balance); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

const findUsersByEmailQuery = `
select id, email, password, username, profile_name, balance from users where email = $1
`

func (repo *PsqlRepository) FindByEmail(email string) ([]*model.DbUser, error) {
	rows, err := repo.db.Query(findUsersByEmailQuery, email)
	if err != nil {
		return nil, err
	}

	var users []*model.DbUser
	for rows.Next() {
		user := model.DbUser{}
		if err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.ProfileName, &user.Balance); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

const findUsersByUsernameQuery = `
select id, email, password, username, profile_name, balance from users where username = $1
`

func (repo *PsqlRepository) FindByUsername(username string) ([]*model.DbUser, error) {
	rows, err := repo.db.Query(findUsersByUsernameQuery, username)
	if err != nil {
		return nil, err
	}

	var users []*model.DbUser
	for rows.Next() {
		user := model.DbUser{}
		if err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.ProfileName, &user.Balance); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

const deleteUsersBatchQuery = `
delete from users where username in (%s)
`

func (repo *PsqlRepository) Delete(users []*model.DbUser) error {
	placeholders := make([]string, len(users))
	usernames := make([]interface{}, len(users))

	for i := 0; i < len(users); i++ {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		usernames[i] = users[i].Username
	}

	query := fmt.Sprintf(deleteUsersBatchQuery, strings.Join(placeholders, ","))
	_, err := repo.db.Exec(query, usernames...)
	return err
}
