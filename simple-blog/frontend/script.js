let page = 1;
let loading = false;

async function loadBlogs() {
  if (loading) return;
  loading = true;

  const res = await fetch(`/api/blogs?page=${page}`);
  const data = await res.json();
  const container = document.getElementById('blog-container');

  data.forEach(blog => {
    const div = document.createElement('div');
    div.className = 'blog';
    div.innerHTML = `
      <p><strong>时间:</strong> ${new Date(blog.created_at).toLocaleString()}</p>
      <p>${blog.text}</p>
      ${blog.image_url ? `<img src="${blog.image_url}">` : ''}
      ${blog.video_url ? `<video controls src="${blog.video_url}"></video>` : ''}
    `;
    container.appendChild(div);
  });

  if (data.length < 5) {
    document.getElementById('loading').innerText = "没有更多了";
    window.removeEventListener('scroll', onScroll);
  }

  page++;
  loading = false;
}

function onScroll() {
  if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight - 50) {
    loadBlogs();
  }
}

document.getElementById("blogForm").addEventListener("submit", async function(e) {
  e.preventDefault();
  const formData = new FormData(this);
  await fetch("/api/blog", {
    method: "POST",
    body: formData
  });
  this.reset();
  document.getElementById('blog-container').innerHTML = '';
  page = 1;
  window.addEventListener('scroll', onScroll);
  loadBlogs();
});

window.addEventListener('scroll', onScroll);
loadBlogs();
