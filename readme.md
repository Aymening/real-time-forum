## <span style="color:magenta; size : 20px">*REAL-TIME-FORUM*</span>

### <span style="color:pink">Description

Forum consists in creating a web forum that allows :

- communication between users.
- associating categories to posts.
- liking and disliking posts and comments.
- filtering posts.


### <span style="color:pink">Authors

- Ayman Hiti
- Hasnae Lamrani


### <span style="color:pink">Usage

1. Clone the repository:
``` 
git clone https://learn.zone01oujda.ma/git/ekerzazi/forum
```
2. Navigate to the project directory:
```
cd forum
```
3. Run the program :
```
go run cmd/*
```
4. Open your web browser and go to http://localhost:8080

### <span style="color:pink">Implementation Details

- #### <span style="color:#40E0D0">SQLite
The SQLite database library will be utilized to save data in the forum (including users, posts, comments, etc.).

It is necessary to utilize at least one SELECT, one CREATE, and one INSERT query. To obtain additional details regarding SQLite, one can refer to the SQLite page.

- #### <span style="color:#40E0D0">Authentication
Ask for an email address.

Return an error if the email is already in use.

Ask for a username.

Ask for a password.

The password needs to be encrypted during storage (bonus task).

The forum should confirm whether the supplied email is present in the database and that all credentials are accurate. If the password is incorrect, an error message should be issued.

- #### <span style="color:#40E0D0">Communication
All users, both registered and non-registered, will be able to see the posts and comments, but non-registered users will only have the option to view them.

Only users who are registered can post and comment

- #### <span style="color:#40E0D0">Likes and Dislikes
Only registered users will be able to like or dislike posts and comments.
The number of likes and dislikes should be visible by all users (registered or not).

- #### <span style="color:#40E0D0">Filter

You must set up a filtering system that enables users to sort the shown posts by:

      - categories
      - made posts
      - posts that were liked

- #### <span style="color:#40E0D0">Docker
The forum project requires Docker.

### <span style="color:pink">Requirements

   - The backend must be written in Go.
   - SQLite must be used.
   - The site and server cannot crash at any time.
   - All pages must function properly, and any mistakes must be fixed.
   - Handle HTTP status and website issues.
   - The code must respect the good practices.
   - It is recommended to have test files for unit testing.
   - The website must be appealing, interactive, and intuitive, using CSS.
