package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"Users/model"
)

type Repository interface {
	GetUser() ([]model.User, error)
	GetAlltask(id int) ([]model.TaskRes, error)
	GetTaskById(id int, taskId int) (model.TaskRes, error)
	CreateTask(arg model.TaskReq, parsedDate time.Time, imageURL string, id int) (model.TaskRes, error)
	UpdateTasks(arg model.TaskReq, parseDate time.Time, ImageURL string, id int) (model.TaskRes, error)
	DeleteTask(Id int) error
	BulkDeleteTask(taskIds []int, Id int) error
	Regis(email string, HasPassword string) (model.UserRegisRespon, error)
	Login(email string) (model.LoginRes, error)
	Logout(reqToken string) error
	SaveToken(token string, userId int) error
	CountTasks(id int, keywoard string, parsedDate time.Time) (int, error)
	SearchTasks(id int, keywoard string, parsedDate time.Time, limit, offset int) ([]model.TaskRes, error)
	CountTask(Id int) (model.Count, error)
	getUserIDByEmail(email string) (int, error)
	GetUserByEmail(email string) (user model.User, err error)
	StoreToken(token string, expirationTime time.Time, id int) (err error)
	CekToken(token string) (data model.ForgotPassword, err error)
	ResetPassword(Password string, Id int) error
	DeleteToken(token string) error
	CreateTableContoh() error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db}
}

func (r *repository) GetUser() ([]model.User, error) {
	var (
		db    = r.db
		users = []model.User{}
	)
	query := `SELECT id,name,email,umur,created_at,updated_at FROM users`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User

		err = rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Umur,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *repository) GetAlltask(id int) ([]model.TaskRes, error) {
	var (
		db    = r.db
		tasks = []model.TaskRes{}
	)

	query := `SELECT tasks.id, tasks.title, tasks.description, tasks.status, tasks.date, tasks.image, tasks.created_at, tasks.updated_at, tasks.id_user, tasks.category_id, category.name_category, tasks.important
		FROM tasks
		LEFT JOIN category
		ON tasks.category_id = category.id
		WHERE tasks.id_user = $1
		ORDER BY tasks.important ASC`

	row, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	for row.Next() {
		var task model.TaskRes
		err = row.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Date,
			&task.Image,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.IdUser,
			&task.CategoryId,
			&task.CategoryName,
			&task.Important,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)

	}
	return tasks, nil
}

func (r *repository) CreateTask(arg model.TaskReq, parsedDate time.Time, imageURL string, id int) (model.TaskRes, error) {
	var (
		db   = r.db
		task = model.TaskRes{}
	)

	query := `
		INSERT INTO tasks (title, description, status, date, image, created_at, id_user, category_id, important)
		VALUES ($1, $2, $3, $4, $5, now(), $6, $7, $8)  
		RETURNING id, title, description, status, date, image, created_at, updated_at, id_user, category_id, important
		`

	rows := db.QueryRowx(query, arg.Title, arg.Description, arg.Status, parsedDate, imageURL, id, arg.CategoryId, arg.Important)
	err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.Date, &task.Image, &task.CreatedAt, &task.UpdatedAt, &task.IdUser, &task.CategoryId, &task.Important)
	if err != nil {
		return model.TaskRes{}, err
	}

	return task, nil
}

func (r *repository) UpdateTasks(arg model.TaskReq, parseDate time.Time, imageURL string, id int) (model.TaskRes, error) {
	var (
		db   = r.db
		task = model.TaskRes{}
	)

	query := `UPDATE tasks SET title = $1, description = $2, status = $3, date = $4, image = $5, updated_at = now(), id_user = $6, category_id = $7, important = $8
		WHERE id = $9 AND id_user = $10
		RETURNING id, title, description, status, date, image, created_at, updated_at, id_user, category_id, important`

	rows := db.QueryRowx(query, arg.Title, arg.Description, arg.Status, parseDate, imageURL, id, arg.CategoryId, arg.Important, id, id)
	err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.Date, &task.Image, &task.CreatedAt, &task.UpdatedAt, &task.IdUser, &task.CategoryId, &task.Important)
	if err != nil {
		return model.TaskRes{}, err
	}

	return task, nil
}

func (r *repository) DeleteTask(Id int) error {
	var db = r.db

	query := `DELETE FROM tasks WHERE id = $1`
	_, err := db.Exec(query, Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) BulkDeleteTask(taskIds []int, Id int) error {
	var db = r.db

	for _, taskId := range taskIds {
		query := `DELETE FROM tasks WHERE id = $1 AND id_user = $2`
		_, err := db.Exec(query, taskId, Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) Regis(email string, HasPassword string) (model.UserRegisRespon, error) {
	var db = r.db
	var regis = model.UserRegisRespon{}

	fmt.Println(email, HasPassword)
	query := `
		INSERT INTO users (email, password, created_at)
		VALUES ( $1, $2, now())
		RETURNING id, email, created_at`

	row := db.QueryRowx(query, email, HasPassword)
	err := row.Scan(&regis.ID, &regis.Email, &regis.CreatedAt)
	if err != nil {
		return model.UserRegisRespon{}, err
	}

	return regis, nil
}

func (r *repository) Login(email string) (model.LoginRes, error) {
	var db = r.db
	var login = model.LoginRes{}

	query := `SELECT id, email, created_at, updated_at, password FROM users WHERE email = $1`
	row := db.QueryRowx(query, email)
	err := row.Scan(&login.Id, &login.Email, &login.CreatedAt, &login.UpdatedAt, &login.Password)
	if err != nil {
		return model.LoginRes{}, err
	}

	return login, nil
}

func (r *repository) SaveToken(token string, userId int) error {
	var db = r.db

	const query2 = `INSERT INTO user_token (user_id, token) VALUES ($1, $2)`
	_ = db.QueryRowx(query2, userId, token)

	return nil
}

func (r *repository) GetTaskById(id int, taskId int) (model.TaskRes, error) {
	var (
		db    = r.db
		tasks = model.TaskRes{}
	)

	query := `SELECT id, title, description, status, date, image, created_at, updated_at, id_user FROM tasks WHERE id_user = $1 AND id = $2`
	rows, err := db.Query(query, id, taskId)
	if err != nil {
		return model.TaskRes{}, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(
			&tasks.Id,
			&tasks.Title,
			&tasks.Description,
			&tasks.Status,
			&tasks.Date,
			&tasks.Image,
			&tasks.CreatedAt,
			&tasks.UpdatedAt,
			&tasks.IdUser,
		)
		if err != nil {
			return model.TaskRes{}, err
		}
	}

	return tasks, err
}

func (r *repository) GetUserByEmail(email string) (user model.User, err error) {
	var db = r.db

	const statement = `SELECT id FROM users WHERE email = $1`
	row := db.QueryRowx(statement, email)
	if err != nil {
		return
	}

	err = row.Scan(&user.Id)
	if err != nil {
		return
	}

	return

}

func (r *repository) Logout(reqToken string) error {
	var (
		db = r.db
	)

	query := `DELETE FROM user_token WHERE token = $1`

	_, err := db.Exec(query, reqToken)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) CountTasks(id int, keywoard string, parsedDate time.Time) (int, error) {
	var (
		db = r.db
	)
	totalQuery := `SELECT COUNT(*) FROM tasks WHERE id_user = $1`

	if !parsedDate.IsZero() {
		totalQuery += fmt.Sprintf(" AND date::date = '%s'", parsedDate.Format("2006-01-02"))
	}

	if keywoard != "" {
		totalQuery += fmt.Sprintf(" AND (tittle ILIKE '%s' OR description ILIKE '%s')", keywoard, keywoard)
	}

	var count int
	err := db.Get(&count, totalQuery, id)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repository) SearchTasks(id int, keywoard string, parsedDate time.Time, limit, offset int) ([]model.TaskRes, error) {
	var (
		db = r.db
	)

	query := `SELECT id, tittle, description, status, date, image, created_at, updated_at, id_user FROM tasks WHERE id_user = $1`

	keywoard = "%" + keywoard + "%"

	if !parsedDate.IsZero() {
		query += fmt.Sprintf(" AND date::date = '%s'", parsedDate.Format("2006-01-02"))
	}

	if keywoard != "" {
		query += fmt.Sprintf(" AND (tittle ILIKE '%s' OR description ILIKE '%s')", keywoard, keywoard)
	}

	query += " LIMIT $2 OFFSET $3"

	rows, err := db.Query(query, id, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.TaskRes

	for rows.Next() {
		var user model.TaskRes
		err = rows.Scan(
			&user.Id,
			&user.Title,
			&user.Description,
			&user.Status,
			&user.Date,
			&user.Image,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.IdUser,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *repository) CountTask(Id int) (model.Count, error) {
	var (
		db    = r.db
		count = model.Count{}
	)

	query := `SELECT 
				SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) AS pending,
				SUM(CASE WHEN status = 'progress' THEN 1 ELSE 0 END) AS progress,
				SUM(CASE WHEN status = 'done' THEN 1 ELSE 0 END) AS done
			FROM tasks 
			WHERE id_user = $1`

	err := db.Get(&count, query, Id)
	if err != nil {
		return model.Count{}, err
	}

	return count, err
}

func (r *repository) getUserIDByEmail(email string) (int, error) {
	var (
		userID int
		db     = r.db
	)
	// Lakukan pencarian pengguna di database Anda
	// Misalnya, jika menggunakan SQL:
	query := "SELECT id FROM users WHERE email = $1"
	err := db.QueryRow(query, email).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *repository) StoreToken(token string, expirationTime time.Time, id int) (err error) {
	var db = r.db
	_, err = db.Exec("INSERT INTO token_forgot (user_id, token, expired_at) VALUES ($1, $2, $3)", id, token, expirationTime)
	if err != nil {
		return
	}

	return
}

func (r *repository) CekToken(token string) (data model.ForgotPassword, err error) {
	var db = r.db

	query := `SELECT user_id, expired_at FROM token_forgot WHERE token = $1`
	err = db.QueryRow(query, token).Scan(&data.UserId, &data.ExpiredAt)
	if err != nil {
		return
	}

	return
}

func (r *repository) CreateTableContoh() error {
	var db = r.db

	query := `CREATE TABLE CONTOH (name varchar(255) NOT NULL)`

	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *repository) ResetPassword(Password string, Id int) error {
	var db = r.db

	query := `UPDATE users SET password = $1 WHERE id = $2`

	_, err := db.Exec(query, Password, Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) DeleteToken(token string) error {
	var db = r.db

	query := `DELETE FROM token_forgot WHERE  token = $1`
	_, err := db.Exec(query, token)
	if err != nil {
		return err
	}

	return nil
}
