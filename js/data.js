const time = new EventSource('/data');
content = document.getElementById("content")
time.addEventListener('data', (e) => {
    content.innerText = content.innerText +e.data + "\n"

}, false);