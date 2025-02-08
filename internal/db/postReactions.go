package db

import (
	"fmt"
	"strconv"
)

type UserReaction struct {
	// Id        int64  `json:"id"`
	// UserName  string `json:"userId"`
	PostId    string `json:"postId"`
	CommentId string `json:"commentId"`
	Reaction  string `json:"reaction"`
}

func (d *Database) InsertPostReactionInDB(userName string, userReaction UserReaction) error {
	insertReaction := fmt.Sprintf(`INSERT INTO likes(username, post_id, comment_id, reaction) VALUES(?, ?, null,  "%s")`, userReaction.Reaction)
	stmnt, err := d.db.Prepare(insertReaction)
	if err != nil {
		return err
	}
	postId, err := strconv.Atoi(userReaction.PostId)
	if err != nil {
		return err
	}
	if _, err := stmnt.Exec(userName, postId); err != nil {
		return err
	}

	return nil
}

func (d *Database) GetPostReactionFromDB(userName string, userReaction UserReaction) (string, error) {
	var postReaction string

	selectPost := `SELECT reaction FROM likes WHERE post_id = ? AND username = ? AND comment_id IS NULL;`

	row := d.db.QueryRow(selectPost, userReaction.PostId, userName)
	err := row.Scan(&postReaction)
	if err != nil {
		return "", err
	}

	return postReaction, nil
}

func (d *Database) UpdatePostReactionInDB(userName string, userReaction UserReaction) error {
	expression := `UPDATE likes SET reaction = ? WHERE post_id = ? AND username = ? AND comment_id IS NULL;`
	_, err := d.db.Exec(expression, userReaction.Reaction, userReaction.PostId, userName)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) DeletePostReactionInDB(userName string, userReaction UserReaction) error {
	expression := `DELETE FROM likes  WHERE post_id = ? AND username = ? AND comment_id IS NULL;`
	_, err := d.db.Exec(expression, userReaction.PostId, userName)
	if err != nil {
		return err
	}

	return nil
}
