package adapter

import (
	"api/internal/entity"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type postgresAdapter struct {
	db *sql.DB
}

func NewPostgreAdapter(db *sql.DB) *postgresAdapter {
	return &postgresAdapter{
		db: db,
	}
}

func (a *postgresAdapter) Upsert(ctx context.Context, user entity.User) error {
	query := `
	INSERT INTO challenge.users (id, first_name, last_name, email_address, created_at, deleted_at, merged_at, parent_user_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (email_address)
	DO UPDATE SET 
		first_name = CASE WHEN EXCLUDED.created_at > users.created_at THEN EXCLUDED.first_name ELSE users.first_name END,
		last_name = CASE WHEN EXCLUDED.created_at > users.created_at THEN EXCLUDED.last_name ELSE users.last_name END,
		created_at = CASE WHEN EXCLUDED.created_at > users.created_at THEN EXCLUDED.created_at ELSE users.created_at END,
		deleted_at = CASE WHEN EXCLUDED.created_at > users.created_at THEN EXCLUDED.deleted_at ELSE users.deleted_at END,
		merged_at = CASE WHEN EXCLUDED.created_at > users.created_at THEN EXCLUDED.merged_at ELSE users.merged_at END,
		parent_user_id = CASE WHEN EXCLUDED.created_at > users.created_at THEN EXCLUDED.parent_user_id ELSE users.parent_user_id END
	RETURNING id;
    `
	_, err := a.db.Exec(
		query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.CreatedAt,
		user.DeletedAt,
		user.MergedAt,
		user.ParentUserID,
	)

	return err
}

func (a *postgresAdapter) GetByID(ctx context.Context, id string) (*entity.User, error) {
	query := `
    SELECT id, first_name, last_name, email_address, created_at, deleted_at, merged_at, parent_user_id
    FROM challenge.users
    WHERE id = $1;
    `

	var user entity.User
	err := a.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.DeletedAt,
		&user.MergedAt,
		&user.ParentUserID,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

type QueryOpts struct {
	FirstName string
	LastName  string
	Email     string

	Fields string
}

func (a *postgresAdapter) Search(ctx context.Context, queryOpts QueryOpts) ([]entity.User, error) {
	baseQuery := `
	   SELECT id, first_name, last_name, email_address, created_at, deleted_at, merged_at, parent_user_id
	   FROM challenge.users
	   `
	// baseQuery := fmt.Sprintf("SELECT %s FROM users\n", queryOpts.Fields)

	conditions := []string{}
	args := []interface{}{}
	argID := 1

	if queryOpts.FirstName != "" {
		conditions = append(conditions, fmt.Sprintf("first_name ILIKE $%d", argID))
		args = append(args, "%"+queryOpts.FirstName+"%")
		argID++
	}
	if queryOpts.LastName != "" {
		conditions = append(conditions, fmt.Sprintf("last_name ILIKE $%d", argID))
		args = append(args, "%"+queryOpts.LastName+"%")
		argID++
	}
	if queryOpts.Email != "" {
		conditions = append(conditions, fmt.Sprintf("email_address ILIKE $%d", argID))
		args = append(args, "%"+queryOpts.Email+"%")
		argID++
	}

	if len(conditions) > 0 {
		baseQuery += "WHERE " + strings.Join(conditions, " AND ")
	}

	fmt.Println(baseQuery)

	rows, err := a.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fmt.Println(rows)

	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(
			// scanFields(&user, strings.Split(queryOpts.Fields, ",")...),
			&user.ID, &user.FirstName, &user.LastName,
			&user.Email, &user.CreatedAt,
			&user.DeletedAt, &user.MergedAt, &user.ParentUserID,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// func scanFields(user *entity.User, fields ...string) []interface{} {
// 	scanArgs := make([]interface{}, len(fields))
//
// 	fmt.Println(scanArgs...)
//
// 	for i, field := range fields {
// 		switch strings.TrimSpace(field) {
// 		case "id":
// 			scanArgs[i] = &user.ID
// 		case "first_name":
// 			scanArgs[i] = &user.FirstName
// 		case "last_name":
// 			scanArgs[i] = &user.LastName
// 		case "email_address":
// 			scanArgs[i] = &user.Email
// 		}
// 		// field validation here... if default, invalid
// 	}
// 	return scanArgs
// }
