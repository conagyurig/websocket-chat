package store

import (
	"context"
	"database/sql"
	"fmt"
	"websocket-chat/internal/models"

	"github.com/google/uuid"
)

type SQLStore struct {
	DB *sql.DB
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		DB: db,
	}
}

func (s *SQLStore) CreateRoom(name string) (*models.Room, error) {
	roomID := uuid.New().String()

	room := &models.Room{
		RoomID: roomID,
		Name:   name,
	}
	query := `INSERT INTO Rooms (RoomID, Name) VALUES (?, ?);`

	_, err := s.DB.ExecContext(context.Background(), query, room.RoomID, room.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}

	return room, nil
}

func (s *SQLStore) GetRoomByID(roomID string) (*models.Room, error) {
	query := `SELECT RoomID, Name FROM Rooms WHERE RoomID = ?;`

	room := &models.Room{}

	err := s.DB.QueryRowContext(context.Background(), query, roomID).Scan(&room.RoomID, &room.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("room not found")
		}
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	return room, nil
}

func (s *SQLStore) CreateUser(roomID, displayName string) (*models.User, error) {
	userID := uuid.New().String()

	user := &models.User{
		UserID:      userID,
		RoomID:      roomID,
		DisplayName: displayName,
	}
	query := `INSERT INTO Users (UserID, RoomID, DisplayName) VALUES (?, ?, ?);`

	_, err := s.DB.ExecContext(context.Background(), query, user.UserID, user.RoomID, user.DisplayName)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *SQLStore) GetUserByID(userID string) (*models.User, error) {
	query := `SELECT UserID, RoomID, DisplayName FROM Users WHERE UserID = ?;`

	user := &models.User{}

	err := s.DB.QueryRowContext(context.Background(), query, userID).Scan(&user.UserID, &user.RoomID, &user.DisplayName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *SQLStore) CreateOption(roomID, userID, content string) (*models.Option, error) {
	optionID := uuid.New().String()

	option := &models.Option{
		OptionID: optionID,
		RoomID:   roomID,
		UserID:   userID,
		Content:  content,
	}
	query := `INSERT INTO Options (OptionID, RoomID, UserID, Content) VALUES (?, ?, ?, ?);`

	_, err := s.DB.ExecContext(context.Background(), query, option.OptionID, option.RoomID, option.UserID, option.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to create option: %w", err)
	}

	return option, nil
}

func (s *SQLStore) GetOption(optionID string) (*models.Option, error) {
	query := `SELECT OptionID, RoomID, UserID, content FROM Options WHERE OptionID = ?;`

	option := &models.Option{}

	err := s.DB.QueryRowContext(context.Background(), query, optionID).Scan(&option.OptionID, &option.RoomID, &option.UserID, &option.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("option not found")
		}
		return nil, fmt.Errorf("failed to get option: %w", err)
	}

	return option, nil
}

func (s *SQLStore) CreateVote(optionID, userID string) (*models.Vote, error) {
	voteID := uuid.New().String()

	vote := &models.Vote{
		VoteID:   voteID,
		OptionID: optionID,
		UserID:   userID,
	}
	query := `INSERT INTO Votes ( VoteID, OptionID, UserID ) VALUES (?, ?, ?);`

	_, err := s.DB.ExecContext(context.Background(), query, vote.VoteID, vote.OptionID, vote.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create vote: %w", err)
	}

	return vote, nil
}

func (s *SQLStore) GetVote(voteID string) (*models.Vote, error) {
	query := `SELECT VoteID, OptionID, UserID FROM Votes WHERE VoteID = ?;`

	vote := &models.Vote{}

	err := s.DB.QueryRowContext(context.Background(), query, voteID).Scan(&vote.VoteID, &vote.OptionID, &vote.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vote not found")
		}
		return nil, fmt.Errorf("failed to get vote: %w", err)
	}

	return vote, nil
}

func (s *SQLStore) GetFullRoomState(roomID string) (*FullRoomStateMessage, error) {
	room, err := s.GetRoomByID(roomID)
	if err != nil {
		return nil, err
	}

	users, err := s.GetUsersByRoomID(roomID)
	if err != nil {
		return nil, err
	}

	options, err := s.GetOptionsByRoomID(roomID)
	if err != nil {
		return nil, err
	}

	votes, err := s.GetVotesByRoomID(roomID)
	if err != nil {
		return nil, err
	}

	fullState := &FullRoomStateMessage{
		RoomName:    room.Name,
		Users:       users,
		Options:     options,
		Votes:       votes,
		RevealVotes: false,
	}

	return fullState, nil
}

func (s *SQLStore) GetUsersByRoomID(roomID string) ([]models.User, error) {
	query := `SELECT UserID, DisplayName FROM Users WHERE RoomID = ?;`

	rows, err := s.DB.Query(query, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID, &user.DisplayName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		user.RoomID = roomID
		users = append(users, user)
	}
	return users, nil
}

func (s *SQLStore) GetOptionsByRoomID(roomID string) ([]models.Option, error) {
	query := `SELECT OptionID, UserID, Content FROM Options WHERE RoomID = ?;`

	rows, err := s.DB.Query(query, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get options: %w", err)
	}
	defer rows.Close()

	var options []models.Option
	for rows.Next() {
		var option models.Option
		err := rows.Scan(&option.OptionID, &option.UserID, &option.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to scan option: %w", err)
		}
		option.RoomID = roomID
		options = append(options, option)
	}
	return options, nil
}

func (s *SQLStore) GetVotesByRoomID(roomID string) ([]models.Vote, error) {
	query := `
        SELECT v.VoteID, v.OptionID, v.UserID
        FROM Votes v
        JOIN Options o ON v.OptionID = o.OptionID
        WHERE o.RoomID = ?;	
    `

	rows, err := s.DB.Query(query, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get votes: %w", err)
	}
	defer rows.Close()

	var votes []models.Vote
	for rows.Next() {
		var vote models.Vote
		err := rows.Scan(&vote.VoteID, &vote.OptionID, &vote.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vote: %w", err)
		}
		votes = append(votes, vote)
	}
	return votes, nil
}

func (s *SQLStore) GetOptionByUserID(userID string) ([]models.Option, error) {
	query := `SELECT OptionID, RoomID, Content FROM Options WHERE UserID = ?;`

	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get options: %w", err)
	}
	defer rows.Close()

	var options []models.Option
	for rows.Next() {
		var option models.Option
		err := rows.Scan(&option.OptionID, &option.RoomID, &option.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to scan option: %w", err)
		}
		option.UserID = userID
		options = append(options, option)
	}
	return options, nil
}

func (s *SQLStore) ChangeVote(userID, newOptionID string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var existingVoteID string
	err = tx.QueryRow(`
        SELECT VoteID FROM Votes WHERE UserID = ?
    `, userID).Scan(&existingVoteID)
	if err != nil {
		if err == sql.ErrNoRows {
			voteID := uuid.New().String()
			_, err = tx.Exec(`
                INSERT INTO Votes (VoteID, UserID, OptionID) VALUES (?, ?, ?)
            `, voteID, userID, newOptionID)
			if err != nil {
				return fmt.Errorf("failed to insert new vote: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check existing vote: %w", err)
		}
	} else {
		_, err = tx.Exec(`
            UPDATE Votes SET OptionID = ? WHERE VoteID = ?
        `, newOptionID, existingVoteID)
		if err != nil {
			return fmt.Errorf("failed to update vote: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *SQLStore) ChangeOption(userID, roomID, newContent string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var existingOptionID string
	err = tx.QueryRow(`
        SELECT OptionID FROM Options WHERE UserID = ?
    `, userID).Scan(&existingOptionID)
	if err != nil {
		if err == sql.ErrNoRows {
			optionID := uuid.New().String()
			_, err = tx.Exec(`
                INSERT INTO Options (OptionID, UserID, OptionID) VALUES (?, ?, ?)
            `, optionID, roomID, userID, newContent)
			if err != nil {
				return fmt.Errorf("failed to insert new option: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check existing option: %w", err)
		}
	} else {
		_, err = tx.Exec(`
            UPDATE Options SET Content = ? WHERE OptionID = ?
        `, newContent, existingOptionID)
		if err != nil {
			return fmt.Errorf("failed to update option: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *SQLStore) ChangeUserName(userID, roomID, newName string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var existingUserID string
	err = tx.QueryRow(`
        SELECT UserID FROM Users WHERE UserID = ?
    `, userID).Scan(&existingUserID)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("failed to find user: %w", err)
		} else {
			return fmt.Errorf("failed to check existing user: %w", err)
		}
	} else {
		_, err = tx.Exec(`
            UPDATE Users SET DisplayName = ? WHERE UserID = ?
        `, newName, existingUserID)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *SQLStore) CreateDate(roomID, userID, dateContent string) (*models.Date, error) {
	dateID := uuid.New().String()

	date := &models.Date{
		DateID: dateID,
		RoomID: roomID,
		UserID: userID,
		Date:   dateContent,
	}
	query := `INSERT INTO Dates (DateID, RoomID, UserID, Date) VALUES (?, ?, ?, ?);`

	_, err := s.DB.ExecContext(context.Background(), query, date.DateID, date.RoomID, date.UserID, date.Date)
	if err != nil {
		return nil, fmt.Errorf("failed to create date: %w", err)
	}

	return date, nil
}

func (s *SQLStore) DeleteUserDates(roomID, userID string) error {

	query := `DELETE FROM Dates WHERE UserID = ? AND RoomID = ?;`

	_, err := s.DB.ExecContext(context.Background(), query, userID, roomID)
	if err != nil {
		return fmt.Errorf("failed to create date: %w", err)
	}

	return nil
}

func (s *SQLStore) GetDateByUserID(userID string) ([]models.Date, error) {
	query := `SELECT DateID, RoomID, Date FROM Dates WHERE UserID = ?;`

	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dates: %w", err)
	}
	defer rows.Close()

	var dates []models.Date
	for rows.Next() {
		var date models.Date
		err := rows.Scan(&date.DateID, &date.RoomID, &date.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to scan date: %w", err)
		}
		date.UserID = userID
		dates = append(dates, date)
	}
	return dates, nil
}

func (s *SQLStore) GetDatesByRoomID(roomID string) (*models.RoomDatesResponse, error) {
	users, err := s.GetUsersByRoomID(roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users for room %s: %w", roomID, err)
	}
	dateWithUsersList := []models.DateWithUsers{}

	for _, user := range users {
		dates, err := s.GetDateByUserID(user.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get dates for user %s: %w", user.UserID, err)
		}
		for _, date := range dates {
			found := false
			for i, dateWithUsers := range dateWithUsersList {
				if dateWithUsers.Date == date.Date {
					dateWithUsersList[i].Users = append(dateWithUsersList[i].Users, user)
					found = true
					break
				}
			}
			if !found {
				dateWithUsersList = append(dateWithUsersList, models.DateWithUsers{
					Date:  date.Date,
					Users: []models.User{user},
				})
			}
		}
	}

	return &models.RoomDatesResponse{
		RoomID: roomID,
		Dates:  dateWithUsersList,
	}, nil
}
