const express = require('express');
const app = express();
const port = 8081;

app.get('/ping', (req, res) => {
    res.json({ status: 'success' });
});

app.listen(port, () => {
    console.log(`Auth Service running on port ${port}`);
});
