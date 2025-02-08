import { showAlert } from "/static/js/alert.js";

export async function createPosts(userName, displayPostsWithReactions) {

    if (!userName) {
        return
    }

    const create_post = [document.getElementById('Create-Post'), document.getElementById('create-post-btn')]

    create_post.forEach(btn => {

        btn.addEventListener('click', function (e) {
            document.getElementById('overlay').style.display = 'flex';
            e.preventDefault()
        });
    })

    document.getElementById('Create-Post').addEventListener('click', function () {
        document.getElementById('overlay').style.display = 'flex';
    });

    document.getElementById('closeFormBtn').addEventListener('click', function () {
        document.getElementById('overlay').style.display = 'none';
    });

    document.getElementById('postForm').addEventListener('submit', async function (event) {
        const title = document.getElementById('title').value.trim();
        const content = document.getElementById('content').value.trim();
        event.preventDefault();
        
        if (title.length > 250) {
            return showAlert("Post title is too long")
        }
        if (content.length > 3000) {
            return showAlert("Post content is too long")
        }

        await addPost();
        document.getElementById('overlay').style.display = 'none';
        displayPostsWithReactions()
    });

    const addPost = async () => {
        const selectedCategories = Array.from(
            document.querySelectorAll('input[name="categories"]:checked')
        ).map(checkbox => checkbox.value);
        const title = document.getElementById('title').value.trim();
        const content = document.getElementById('content').value.trim();
        if (!title || !content) {
            return showAlert("Please fill in all fields correctly!");
        }
        if (selectedCategories.length === 0) {
            return showAlert("Please select at least one category!");
        }

        const requestBody = {
            title,
            content,
            selectedCategories
        };

        try {
            const response = await fetch('/create/post', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestBody),
            });

            if (!response.ok) {
                const errorMessage = await response.json();
                return showAlert(`${errorMessage.message || 'Failed to create post'}`);
            }

            const responseData = await response.json();
            showAlert(responseData.message);
            document.getElementById('title').value = '';
            document.getElementById('content').value = '';
            document.querySelectorAll('input[name="categories"]').forEach(checkbox => {
                checkbox.checked = false;
            });

        } catch (error) {
            console.error('Unexpected error:', error);
            showAlert('An unexpected error occurred. Please try again later.');
        }
    }

};