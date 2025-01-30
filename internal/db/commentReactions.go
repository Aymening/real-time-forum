package db

import "fmt"

func (d *Database) InsertCommentReactionInDB(userName string, userReaction UserReaction) error {
	insertReaction := fmt.Sprintf(`INSERT INTO likes(username, post_id, comment_id, reaction) VALUES(?, null, ?,  "%s")`, userReaction.Reaction)
	stmnt, err := d.db.Prepare(insertReaction)
	if err != nil {
		return err
	}

	if _, err := stmnt.Exec(userName, userReaction.CommentId); err != nil {
		return err
	}

	return nil
}

func (d *Database) GetCommentReactionFromDB(userName string, userReaction UserReaction) (string, error) {
	var postReaction string

	selectPost := `SELECT reaction FROM likes WHERE post_id IS NULL AND comment_id = ? AND username = ? ;`

	row := d.db.QueryRow(selectPost, userReaction.CommentId, userName)
	err := row.Scan(&postReaction)
	if err != nil {
		return "", err
	}

	return postReaction, nil
}

func (d *Database) UpdateCommentReactionInDB(userName string, userReaction UserReaction) error {
	expression := `UPDATE likes SET reaction = ? WHERE post_id IS NULL AND comment_id = ? AND username = ?;`
	_, err := d.db.Exec(expression, userReaction.Reaction, userReaction.CommentId, userName)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) DeleteCommentReactionInDB(userName string, userReaction UserReaction) error {
	expression := `DELETE FROM likes  WHERE post_id IS NULL AND comment_id = ? AND username = ?;`
	_, err := d.db.Exec(expression, userReaction.CommentId, userName)
	if err != nil {
		return err
	}

	return nil
}
