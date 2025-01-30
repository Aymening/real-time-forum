package db

func (db *Database) CreateAllTablesInDatabase() []string {

	TableUsers := `
			-- DROP TABLE IF EXISTS users;
			-- DROP TABLE IF EXISTS posts;
			-- DROP TABLE IF EXISTS comments;
			-- DROP TABLE IF EXISTS categories;
			-- DROP TABLE IF EXISTS likes; 
			-- DROP TABLE IF EXISTS post_category;
			-- DROP TABLE IF EXISTS categories;

			CREATE TABLE IF NOT EXISTS users (
    				id INTEGER PRIMARY KEY AUTOINCREMENT,
    				username TEXT NOT NULL,
    				email TEXT NOT NULL,
    				password TEXT NOT NULL,
					token TEXT 
			);
	`

	TableCategories := `
					CREATE TABLE IF NOT EXISTS categories (
    					id INTEGER PRIMARY KEY AUTOINCREMENT,
    					category TEXT NOT NULL UNIQUE
					);
	`

	TablePosts := `
	
				CREATE TABLE IF NOT EXISTS posts (
    				id INTEGER PRIMARY KEY AUTOINCREMENT,
    				title TEXT NOT NULL,
    				content TEXT NOT NULL,
    				user_id INTEGER NOT NULL,
    				creation_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    				FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
				);
	`

	TablePostsCategories := `
					CREATE TABLE IF NOT EXISTS post_categories (
    				post_id INTEGER NOT NULL,
    				category_id INTEGER NOT NULL,
    				PRIMARY KEY (post_id, category_id),
    				FOREIGN KEY (post_id ) REFERENCES posts (id) ON DELETE CASCADE ON UPDATE CASCADE,
    				FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE ON UPDATE CASCADE
				);
	
	`

	TableComments := `
				CREATE TABLE IF NOT EXISTS comments (
    				id INTEGER PRIMARY KEY AUTOINCREMENT,
    				content TEXT NOT NULL,
    				user_id INTEGER NOT NULL,
    				post_id INTEGER NOT NULL,
    				create_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    				FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    				FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE ON UPDATE CASCADE
				);
	`

	TableLikes := `
			CREATE TABLE IF NOT EXISTS likes (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					username INTEGER NOT NULL,
					post_id INTEGER,
					comment_id INTEGER,
					reaction TEXT CHECK(reaction IN ("like", "dislike")),
					FOREIGN KEY (username) REFERENCES users (username) ON DELETE CASCADE ON UPDATE CASCADE,
					FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE ON UPDATE CASCADE,
					FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE CASCADE ON UPDATE CASCADE,
					UNIQUE(username, post_id),
					UNIQUE(username, comment_id)
			);
	`

	TableInsertCategories := `
			INSERT OR IGNORE INTO
    			categories (category)
				VALUES
					('sport'),
					('games'),
					('news'),
					('lifestyle'),
					('food')
			`

	return []string{TableUsers, TableCategories, TablePosts, TablePostsCategories, TableLikes, TableComments, TableInsertCategories}
}

func (db *Database) ExecuteAllTableInDataBase() error {
	
	for _, Table := range db.CreateAllTablesInDatabase() {
		_, err := db.db.Exec(Table)
		if err != nil {
			return err
		}
	}
	return nil
}
