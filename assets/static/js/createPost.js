// import { getPosts } from "/static/js/displayPosts.js";

async function createPosts() {

    const selectedCategories = Array.from(
        document.querySelectorAll('input[name="categories"]:checked')
    ).map(checkbox => checkbox.value);
    // const postsContainer=document.getElementById('post-container')
    const title = document.getElementById('title').value.trim();
    const content = document.getElementById('content').value.trim();

    // Validate inputs
    if (!title || !content) {
        return alert("Please fill in all fields correctly!");
    }
    if (selectedCategories.length === 0) {
        return alert("Please select at least one category!");
    }

    const requestBody = {
        title,
        content,
        selectedCategories
    };

    console.log("Data to be sent:", requestBody);
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
            console.log('Error response:', errorMessage.message);
            return alert(`${errorMessage.message || 'Failed to create post'}`);
        }

        // Process JSON response
        const responseData = await response.json();
        console.log('Response:', responseData);
        alert(responseData.message);
        // postsContainer.innerHTML=''
        // getPosts()
        // Clear form fields and checkboxes
        document.getElementById('title').value = '';
        document.getElementById('content').value = '';
        document.querySelectorAll('input[name="categories"]').forEach(checkbox => {
            checkbox.checked = false;
        });

    } catch (error) {
        console.error('Unexpected error:', error);
        alert('An unexpected error occurred. Please try again later.');
    }

};

const btns = [document.getElementById('Post-Created'), document.getElementById('Likes'), document.getElementById('Create-Post')]
const create_post = [document.getElementById('Create-Post'),document.getElementById('create-post-btn')]

btns.forEach(btn => {
    if (btn.classList.contains('Permission-Denied')) {
        btn.classList.remove('Permission-Denied')
    }
});

create_post.forEach(btn=>{
    
    btn.addEventListener('click', function (e) {
        document.getElementById('overlay').style.display = 'flex';
        e.preventDefault()
    });
})

// Show the form when the button is clicked
document.getElementById('Create-Post').addEventListener('click', function () {
    document.getElementById('overlay').style.display = 'flex';
});

// Hide the form when the close button is clicked
document.getElementById('closeFormBtn').addEventListener('click', function () {
    document.getElementById('overlay').style.display = 'none';
});

// Handle form submission
document.getElementById('postForm').addEventListener('submit', async function (event) {
    event.preventDefault(); // Prevent form from reloading the page
    // console.log('Post Created:', postData);
    // alert('Post created successfully!');
    await createPosts();
    // Hide the form after submission
    document.getElementById('overlay').style.display = 'none';
    location.reload()

    // Optionally, reset the form
    // this.reset();
});