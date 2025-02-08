const btns = [document.getElementById('Post-Created'), document.getElementById('Likes'), document.getElementById('Create-Post'), document.getElementById('commentSubmit'), document.getElementById('commentContent')]


btns.forEach(btn => {
    if (btn) {
        btn.classList.add('Permission-Denied')
        btn.setAttribute('readonly',true)
        btn.addEventListener('click', (e)=>{
            e.preventDefault()
        })
    }
});