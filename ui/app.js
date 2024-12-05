document.addEventListener("DOMContentLoaded", function() {
    fetch("/api/issues")
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok ' + response.statusText);
            }
            return response.json();
        })
        .then(data => {
            console.log('Fetched Data:', data); // Log the fetched data for inspection
            const issuesDiv = document.getElementById("issues");
            issuesDiv.innerHTML = ''; // Clear previous content

            data.forEach(issue => {
                const issueDiv = document.createElement("div");
                issueDiv.className = "issue";
                issueDiv.innerHTML = `
                    <h2><a href="${issue.html_url}" target="_blank">${issue.title}</a></h2>
                    <p>${issue.body}</p>
                `;
                issuesDiv.appendChild(issueDiv);
            });
        })
        .catch(error => console.error('Error fetching issues:', error));
});