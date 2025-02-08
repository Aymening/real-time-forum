import { escapeHTML } from "/static/js/createComment.js"
import { showAlert } from "/static/js/alert.js";

export const fetchApi = async (url) => {
  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        action: "fetchComments"
      })
    });
    return await response.json();
  } catch (error) {
    console.error("Error fetching data:", error);
    return null;
  }
};

export const getPosts = async (UserName, userReactions) => {

  const handleClick = async (e) => {
    const category = e.target.dataset.category;
    const isClickOnAllPosts = e.target.id === 'All-Posts';
    const isCLickOnMyLikes = e.target.id === 'Likes';
    const isClickOnMyPosts = e.target.id === 'Post-Created'
    let currentPosts

    if (category || isClickOnAllPosts || isClickOnMyPosts || isCLickOnMyLikes) {
      let input = '';
      let flag = '';

      if (category) {
        input = category;
      } else if (isClickOnMyPosts || isCLickOnMyLikes) {
        if (!UserName) {
          return
        }
        input = UserName
        if (isClickOnMyPosts) {
          flag = "myPosts"
        } else {
          flag = "myLikes"
        }
      }

      currentPosts = await loadPosts(input, flag)
    }
    if (currentPosts) {
      DisplayAllPosts(currentPosts, userReactions)
    }
  };

  let currentPosts = await loadPosts()

  if (currentPosts) {
    DisplayAllPosts(currentPosts, userReactions)
  }

  document.removeEventListener('click', handleClick);
  document.addEventListener('click', handleClick);
};

const loadPosts = async (input, flag) => {
  let apiData;
  let posts = [];
  let postsId = [];
  const categories = ["sport", "games", "news", "lifestyle", "food"]


  const isGategory = categories.includes(input)
  const isUser = input && !isGategory

  
  if (isUser) {
    apiData = await fetchApi(`/api/users/${input}`);
  } else {
    apiData = await fetchApi(`/api/posts`)
  }
  
  if (input && !isGategory && isUser && !flag) {
    return showAlert("Please enter the right category")
  }
  
  if (isUser && flag) {
    if (flag === 'myPosts') {
      posts = apiData?.posts || [];
    } else if (flag === 'myLikes') {
      postsId = apiData?.reactions || [];
      postsId = postsId.like || []

      posts = await GetPostFromIds(postsId)

    }
  } else if (isGategory) {
    posts = FilterByCategory(apiData, input) || [];
  } else {
    posts = apiData || [];
  }

  return posts
};

const GetPostFromIds = async (postsId) => {
  const posts = await Promise.all(
    postsId.map(async (id) => {
      const apiData = await fetchApi(`/api/posts/${id}`)
      return apiData
    })
  )
  return posts.reverse()
}

const FilterByCategory = (allPosts, category) => {
  if (allPosts) {
    return allPosts.filter(obj => obj.categories.includes(category.toLowerCase()))
  }
 
}

const RenderPosts = function (post) {
  const postElement = document.createElement("div");
  postElement.className = "post";
  postElement.dataset.postId = post.id;
  const categoryLinks = post.categories || []
    .map(cat => `<a>${cat}</a>`)
    .join(" | ");
  postElement.innerHTML = `
      <div>
          <div class="headers">
              <span class="username">${escapeHTML(post.user)}</span>
              <span class="date">${post.creationDate}</span>
          </div>
          <div class="category-label">
              Categories: <a>${categoryLinks}</a>
          </div>
          <div class="title">${escapeHTML(post.title)}</div>
          <div class="content">${escapeHTML(post.content)}</div>
          <a href="/post?id=${post.id}" class="comment-link">See All Comments</a>
          <div class="reactions">
          <div class="like-div">
            <button class="btn">
              <span class="material-icons">thumb_up</span>
            </button>
            <span class="count">${post.likes || 0}</span>
          </div>
          <div class="dislike-div">
            <button class="btn">
              <span class="material-icons">thumb_down</span>
            </button>
            <span class="count">${post.dislikes || 0}</span>
          </div>
        </div>
      </div>
      </div>
  `;
  return postElement;
}

const DisplayAllPosts = function (posts, callBackReactions) {

  const postContainer = document.getElementById("post-container");
  if (!postContainer) {
    console.error("post-container element not found!");
    return;
  }
  
  postContainer.innerHTML = ""

  posts.forEach(post => {
    const postElement = RenderPosts(post);
    postContainer.appendChild(postElement);
  });

  callBackReactions(posts)

};