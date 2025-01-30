package db

import (
	"database/sql"
	"slices"
	"strconv"
	"strings"
)

type Post struct {
	Id           int               `json:"id"`
	User         string            `json:"user"`
	Title        string            `json:"title"`
	Content      string            `json:"content"`
	CreationDate string            `json:"creationDate"`
	Categories   []string          `json:"categories,omitempty"`
	Errors       map[string]string `json:"-,omitempty"`
	Reactions    map[string]string `json:"reactions,omitempty"`
	Likes        int               `json:"likes,omitempty"`
	Dislikes     int               `json:"dislikes,omitempty"`
}

func (d *Database) InsertPost(userId int64, title, content, creationDate string) (int64, error) {
	insertPost := `INSERT INTO posts(user_id, title, content, creation_date) VALUES(?, ?, ?, DATETIME('now'));`
	stmnt, err := d.db.Prepare(insertPost)
	if err != nil {
		return 0, err
	}
	row, err := stmnt.Exec(userId, title, content)
	if err != nil {
		return 0, err
	}

	lastInsertId, err := row.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertId, nil
}

func (d *Database) InsertCategories(postId int, categories []string) error {
	for _, category := range categories {
		insertCategories := `INSERT INTO post_categories(post_id, category_id) VALUES(?,?)`
		stmnt, err := d.db.Prepare(insertCategories)
		if err != nil {
			return err
		}
		_, err = stmnt.Exec(postId, category)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Database) SelectCategories(categories []string) ([]string, error) {
	var categoriesIDs []string
	selectCategories := `SELECT * FROM categories`

	rows, err := d.db.Query(selectCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var temp string
		err := rows.Scan(&id, &temp)
		if err != nil {
			return nil, err
		}
		for _, category := range categories {
			if strings.ToLower(category) == temp {
				categoriesIDs = append(categoriesIDs, id)
			}
		}
	}

	if rows.Err() != nil {
		return nil, err
	}

	return categoriesIDs, nil
}

func (d *Database) InsertPostWithCategories(users map[string]*User, userName string, title, content, creationDate string, categories []string) (*Post, error) {
	postId, err := d.InsertPost(users[userName].Id, title, content, creationDate)
	if err != nil {
		return &Post{}, err
	}

	selectedCategories, err := d.SelectCategories(categories)
	if err != nil {
		return &Post{}, err
	}

	if err := d.InsertCategories(int(postId), selectedCategories); err != nil {
		return &Post{}, err
	}

	return &Post{
		Id:           int(postId),
		User:         userName,
		Title:        title,
		Content:      content,
		CreationDate: creationDate,
		Categories:   categories,
	}, nil
}

func (d *Database) GetPostsFromDB(users map[string]*User, posts *[]*Post) error {
	selectUsers := `SELECT DISTINCT
    users.id,
    users.username,
    posts.id,
    posts.title,
    posts.content,
    posts.creation_date,
    post_categories.category_id,
    categories.category,
    likes.username AS like_username,
    likes.post_id,
    likes.reaction
FROM
    users
    LEFT JOIN posts ON users.id = posts.user_id
    LEFT JOIN post_categories ON posts.id = post_categories.post_id
    LEFT JOIN categories ON post_categories.category_id = categories.id
    LEFT JOIN likes ON posts.id = likes.post_id
ORDER BY
    posts.creation_date DESC`

	rows, err := d.db.Query(selectUsers)
	if err != nil {
		return err
	}
	defer rows.Close()

	postsMap := make(map[int]*Post)

	for rows.Next() {
		var (
			u User

			likedPost    sql.NullString
			reaction     sql.NullString
			likedUser    sql.NullString
			postId       sql.NullInt64
			title        sql.NullString
			content      sql.NullString
			creationDate sql.NullTime
			categoryId   sql.NullInt64
			category     sql.NullString
		)

		if err := rows.Scan(&u.Id, &u.UserName, &postId, &title, &content, &creationDate, &categoryId, &category, &likedUser, &likedPost, &reaction); err != nil {
			return err
		}

		u.Reactions = make(map[string][]int)
		if _, userExists := users[u.UserName]; !userExists {
			users[u.UserName] = &u
		}

		post, exists := postsMap[int(postId.Int64)]
		if postId.Valid && title.Valid && content.Valid && creationDate.Valid {
			if !exists {
				post = &Post{
					Id:           int(postId.Int64),
					User:         u.UserName,
					Title:        title.String,
					Content:      content.String,
					CreationDate: creationDate.Time.Format("Mon Jan 02 15:04:05 2006"),
					Reactions:    make(map[string]string),
				}

				postsMap[int(postId.Int64)] = post
				*posts = append(*posts, post)
				users[u.UserName].Posts = append(users[u.UserName].Posts, post)

			}
		}
		if categoryId.Valid && category.Valid {
			post.Categories = append(post.Categories, category.String)
		}

		if likedUser.Valid && likedPost.Valid && reaction.Valid {
			for _, post := range *posts {
				if strconv.Itoa(post.Id) == likedPost.String {
					if _, exists := post.Reactions[likedUser.String]; !exists {
						post.Reactions[likedUser.String] = reaction.String
						if reaction.String == "like" {
							post.Likes += 1
						} else if reaction.String == "dislike" {
							post.Dislikes += 1
						}
					}
				}
			}

			user := users[likedUser.String]
			if _, exists := user.Reactions[likedPost.String]; !exists {
				postId, err := strconv.Atoi(likedPost.String)
				if err != nil {
					return err
				}
				if !slices.Contains(user.Reactions[reaction.String], postId) {
					user.Reactions[reaction.String] = append(user.Reactions[reaction.String], postId)
				}

			}
		}
	}
	if rows.Err() != nil {
		return err
	}
	return nil
}
