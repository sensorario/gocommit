package web

const MainHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>gocommit!!</title>
  <style>
    body { font-family: sans-serif; padding: 2rem; }
    button { display: block; margin: 0.25rem 0; padding: 0.4rem 1rem; cursor: pointer; }
    #message { margin-top: 1rem; font-weight: bold; }
  </style>
</head>
<body>
  <h1>Hello World</h1>
  <h2>Branches</h2>
  <div id="branches"></div>
  <div id="message"></div>
  <script>
    function checkout(branch) {
      fetch('/checkout', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ branch })
      })
      .then(r => {
        if (!r.ok) {
          return r.text().then(t => { throw new Error(t); });
        }
        return r.json();
      })
      .then(data => {
        document.getElementById('message').textContent = 'Switched to: ' + data.branch;
      })
      .catch(err => {
        document.getElementById('message').textContent = 'Error: ' + err.message;
      });
    }

    fetch('/branches')
      .then(r => r.json())
      .then(branches => {
        const container = document.getElementById('branches');
        branches.forEach(b => {
          const btn = document.createElement('button');
          btn.textContent = b;
          btn.onclick = () => checkout(b);
          container.appendChild(btn);
        });
      });
  </script>
</body>
</html>`
