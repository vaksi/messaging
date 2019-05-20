package repositories

import (
	"encoding/json"
	"fmt"

	"github.com/vaksi/messaging/configs"
	"github.com/vaksi/messaging/internal/models"
	"github.com/vaksi/messaging/pkg/kafka"
	"github.com/vaksi/messaging/pkg/mysql"
)

type MessageRepository struct {
	db mysql.MySqlFactory
	kafka kafka.Kafka
	config *configs.Config
}

type IMessageRepository interface {
	Store(*models.Message) (uint64, error)
	Publish(request PublishRequest) (err error)
	Update(*models.Message) (error)
	GetMessageByUserIDStatusSent(uint64,uint8) ([]models.Message, error)
	GetMessageByUserIDInbox(uint64) ([]models.Message, error)
}

func NewMessageRepository(mysqlDB mysql.MySqlFactory, k kafka.Kafka, cfg *configs.Config) *MessageRepository {
	return &MessageRepository{db: mysqlDB, kafka: k, config: cfg}
}

// CreateMessage store create message
func (mr *MessageRepository) Store(message *models.Message) (lastID uint64, err error) {
	db, err := mr.db.GetDB()
	if err != nil {
		return
	}

	tx, err := db.Begin()
	if err != nil {
		return
	}

	stmt, err := tx.Prepare("INSERT INTO messages (user_id, subject, message_text, to_user_id, status) VALUES (?,?,?,?,?)")
	if err != nil {
		return
	}
	res, err := stmt.Exec(&message.UserID, &message.Subject, &message.Text, &message.ToUserID, &message.Status)
	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return 0, errTx
		}
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return 0, errTx
		}
		return
	}

	lastID = uint64(id)
	err = tx.Commit()

	return
}

type PublishRequest struct {
	MessageID uint64 `json:"message_id"`
}

// Publish this service for publish message to kafka
func (mr *MessageRepository) Publish(request PublishRequest) (err error) {
	msg, _ := json.Marshal(request)
	err = mr.kafka.SendMessage(mr.config.Kafka.MessagingConsumer.Topic, []byte{}, msg)
	return
}

// Update this for update data message
func (mr *MessageRepository) Update(message *models.Message) (err error) {
	db, err := mr.db.GetDB()
	if err != nil {
		return
	}

	tx, err := db.Begin()
	if err != nil {
		return
	}

	stmt, err := tx.Prepare("UPDATE messages SET status=? WHERE id=?")
	if err != nil {
		return
	}
	_, err = stmt.Exec(&message.Status, &message.ID)
	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return errTx
		}
		return
	}

	err = tx.Commit()

	return
}

func (mr *MessageRepository) GetMessageByUserIDStatusSent(userID uint64, status uint8) (messages []models.Message, err error) {
	db, err := mr.db.GetDB()
	if err != nil {
		return
	}

	var strQueryStatus string
	if status != 0 {
		strQueryStatus = fmt.Sprintf("AND status = %d", status)
	}
	queryStr := fmt.Sprintf(`
		SELECT
			id,
			to_user_id,
			message_text,
			subject,
			status,	
			created_at,
			updated_at
		FROM messages 
		WHERE user_id=? %s
		ORDER BY created_at DESC
	`, strQueryStatus)

	rows, err := db.Query(queryStr, userID)

	if err != nil {
		return
	}

	for rows.Next() {
		message := models.Message{}
		err = rows.Scan(
			&message.ID,
			&message.ToUserID,
			&message.Text,
			&message.Subject,
			&message.Status,
			&message.CreatedAt,
			&message.UpdatedAt,
		)
		// skip when scan error
		if err != nil {
			continue
		}
		messages = append(messages, message)
	}

	// Good Habit to Close
	err = rows.Close()
	return
}

func (mr *MessageRepository) GetMessageByUserIDInbox(userID uint64) (messages []models.Message, err error) {
	db, err := mr.db.GetDB()
	if err != nil {
		return
	}

	rows, err := db.Query(fmt.Sprintf(`
		SELECT
			id,
			user_id,
			message_text,
			subject,
			status,	
			created_at,
			updated_at
		FROM messages 
		WHERE to_user_id=? AND status=4
		ORDER BY created_at DESC
	`), userID)

	if err != nil {
		return
	}

	for rows.Next() {
		message := models.Message{}
		err = rows.Scan(
			&message.ID,
			&message.UserID,
			&message.Text,
			&message.Subject,
			&message.Status,
			&message.CreatedAt,
			&message.UpdatedAt,
		)
		// skip when scan error
		if err != nil {
			continue
		}
		messages = append(messages, message)
	}

	// Good Habit to Close
	err = rows.Close()
	return
}
