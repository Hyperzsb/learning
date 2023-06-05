package model

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"time"
)

// Token defines the authentication token used for user authentication and
// tracking purposes. It uses base64 for encoding and SHA256 for hashing.
type Token struct {
	UserID int       `json:"-"`
	Scope  string    `json:"-"`
	Text   string    `json:"text"`
	Hash   [32]byte  `json:"-"`
	Expiry time.Time `json:"expiry"`
}

// NewToken generates a new token for a specific user within a specific scope,
// which will last for a ttl period.
func NewToken(userID int, scope string, ttl time.Duration) (Token, error) {
	token := Token{
		UserID: userID,
		Scope:  scope,
		Expiry: time.Now().Add(ttl),
	}

	nonce := make([]byte, 16)
	_, err := rand.Read(nonce)
	if err != nil {
		return token, err
	}

	text := base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(nonce)
	hash := sha256.Sum256([]byte(text))

	token.Text = text
	token.Hash = hash

	return token, nil
}

// CreateToken creates a new token record and persists it into the database.
// If the persistence succeeds, the ID of the new record will be returned.
// If not, an error will be returned.
func (m *Model) CreateToken(t Token) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	statement := `
		insert into tokens (user_id, scope, hash, expiry)
		values (?, ?, ?, ?)
	`
	result, err := m.db.ExecContext(ctx, statement,
		t.UserID,
		t.Scope,
		t.Hash[:],
		t.Expiry,
	)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetToken gets the token record with a specific ID from the database.
// If there is no such record, an EmptyQueryError will be returned.
// In other failure cases, an arbitrary error will be returned.
func (m *Model) GetToken(id int) (Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	statement := `
		select user_id, scope, hash, expiry
		from tokens
		where id = ?
	`
	row := m.db.QueryRowContext(ctx, statement, id)

	token := Token{}
	err := row.Scan(
		&token.UserID,
		&token.Scope,
		&token.Hash,
		&token.Expiry,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return token, &EmptyQueryError{err.Error()}
		}

		return token, err
	}

	return token, nil
}

// GetTokenByUserID gets the very last token created for a specific user.
func (m *Model) GetTokenByUserID(userID int) (Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	statement := `
		select user_id, scope, hash, expiry
		from tokens
		where user_id = ?
		order by id desc
		limit 1
	`
	row := m.db.QueryRowContext(ctx, statement, userID)

	token := Token{}
	err := row.Scan(
		&token.UserID,
		&token.Scope,
		&token.Hash,
		&token.Expiry,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return token, &EmptyQueryError{err.Error()}
		}

		return token, err
	}

	return token, nil
}

// DeleteTokensByUserID deletes all tokens created for a specific user. As in
// general cases, persisting all pre-created tokens for a specific is meaningless,
// storing the very last token is enough.
func (m *Model) DeleteTokensByUserID(userID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	statement := `
		delete from tokens
		where user_id = ?
	`
	result, err := m.db.ExecContext(ctx, statement, userID)
	if err != nil {
		return -1, err
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(cnt), nil
}

// GetUserByToken gets the corresponding user bonded with a specific.
// The parameter text is the plaintext of this token. Note that if
// the token has expired, no user will be returned.
func (m *Model) GetUserByToken(text string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	hash := sha256.Sum256([]byte(text))

	statement := `
		select id, first_name, last_name, email, password
		from users
		inner join (
		    select user_id
		    from tokens
		    where hash = ? and expiry > current_timestamp
		) as user_token
		on users.id = user_token.user_id
	`
	row := m.db.QueryRowContext(ctx, statement, hash[:])

	user := User{}
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return user, &EmptyQueryError{err.Error()}
		}

		return user, err
	}

	return user, nil
}

func (m *Model) ExpireToken(text string) (int, error) {
	ctxR, cancelR := context.WithTimeout(context.Background(), time.Second)
	defer cancelR()

	hash := sha256.Sum256([]byte(text))

	statementR := `
		select user_id, scope, expiry
		from tokens
		where hash = ? and expiry > ?
	`
	row := m.db.QueryRowContext(ctxR, statementR, hash[:], time.Now())

	token := Token{}
	err := row.Scan(
		&token.UserID,
		&token.Scope,
		&token.Expiry,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return -1, &EmptyQueryError{err.Error()}
		}

		return -1, err
	}

	ctxU, cancelU := context.WithTimeout(context.Background(), time.Second)
	defer cancelU()

	statementU := `
		update tokens
		set expiry = ?
		where hash = ?
	`
	result, err := m.db.ExecContext(ctxU, statementU, time.Now(), hash[:])
	if err != nil {
		return -1, err
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(cnt), nil
}
