document.addEventListener('DOMContentLoaded', function() {
    const headings = document.querySelectorAll('.post-content h1[id], .post-content h2[id], .post-content h3[id], .post-content h4[id], .post-content h5[id], .post-content h6[id]');

    headings.forEach(function(heading) {
        const anchor = document.createElement('a');
        anchor.className = 'anchor-link';
        anchor.href = '#' + heading.id;
        anchor.textContent = '#';
        heading.appendChild(anchor);
    });
});