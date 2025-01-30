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

// const createScrollPagination = (posts, displayCallback) => {
//   let startIndex = 0;
//   let endIndex = 5;
//   let isLoading = false;

//   displayCallback(posts.slice(0, endIndex));

//   const handleScroll = () => {
//     if (isLoading) return;

//     const scrollable = document.documentElement.scrollHeight - window.innerHeight;
//     const scrolled = window.scrollY;

//     if (Math.ceil(scrolled) >= scrollable && posts.length > endIndex) {
//       isLoading = true;

//       startIndex = endIndex;
//       endIndex = Math.min(endIndex + 5, posts.length);

//       displayCallback(posts.slice(startIndex, endIndex), true);

//       isLoading = false;
//     }
//   };

//   // Clean up previous event listeners before adding new one
//   window.removeEventListener('scroll', handleScroll);
//   window.addEventListener('scroll', handleScroll);

//   return () => window.removeEventListener('scroll', handleScroll);
// };

// const scrollingPosts = (posts) => {
//   if (!posts || !posts.length) return;
//   return createScrollPagination(posts, DisplayAllPosts);
// };

export const getPosts = async (UserName, displayCallback) => {
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
      DisplayAllPosts(currentPosts)
      displayCallback()
    }
  };

  let currentPosts = await loadPosts()

  if (currentPosts) {
    DisplayAllPosts(currentPosts)
    displayCallback()
  }

  document.removeEventListener('click', handleClick);
  document.addEventListener('click', handleClick);
};

// export const GoToTop = () => {
//   if ('scrollRestoration' in history) {
//     history.scrollRestoration = 'manual';
//   }
//   window.scrollTo(0, 0);
// }

export const loadPosts = async (input, flag) => {
  let apiData;
  let posts = [];
  let postsId = [];
  const categories = ["sport", "games", "news", "lifestyle", "food"]


  const isGategory = categories.includes(input)
  const isUser = input && !isGategory;


  if (isUser) {
    apiData = await fetchApi(`/api/users/${input}`);
  } else {
    apiData = await fetchApi(`/api/posts`)
  }
  if (isUser) {
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
  return allPosts.filter(obj => obj.categories.includes(category.toLowerCase()))
}

// Function to create a post element
export const RenderPosts = function (post) {
  const postElement = document.createElement("div");
  postElement.className = "post";
  postElement.dataset.postId = post.id;
  const categoryLinks = post.categories
    .map(cat => `<a>${cat}</a>`)
    .join(" | ");
  postElement.innerHTML = `
      <div>
          <div class="headers">
              <span class="username">${post.user}</span>
              <span class="date">${post.creationDate}</span>
          </div>
          <div class="category-label">
              Categories: ${categoryLinks}
          </div>
          <div class="title">${post.title}</div>
          <div class="content">${post.content}</div>
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

// Display posts dynamically :
const DisplayAllPosts = function (posts) {

  const postContainer = document.getElementById("post-container");
  if (!postContainer) {
    console.error("post-container element not found!");
    return;
  }
  // if (!isLoadPosts) {
  postContainer.innerHTML = ""
  posts.forEach(post => {
    const postElement = RenderPosts(post);
    postContainer.appendChild(postElement);
  });

};