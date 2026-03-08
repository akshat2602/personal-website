var fuse;
var resList = document.getElementById('searchResults');
var sInput = document.getElementById('searchInput');
var first, last, current_elem = null;
var resultsAvailable = false;

window.searchReady = false;
window.onload = function () {
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                var data = JSON.parse(xhr.responseText);
                if (data) {
                    var options = {
                        distance: 100,
                        threshold: 0.4,
                        ignoreLocation: true,
                        keys: ['title', 'permalink', 'summary', 'content']
                    };
                    fuse = new Fuse(data, options);
                    window.searchReady = true;
                }
            }
        }
    };
    xhr.open('GET', "/index.json");
    xhr.send();
};

function activeToggle(ae) {
    document.querySelectorAll('.focus').forEach(function (element) {
        element.classList.remove("focus");
    });
    if (ae) {
        ae.focus();
        document.activeElement = current_elem = ae;
        ae.parentElement.classList.add("focus");
    } else {
        document.activeElement.parentElement.classList.add("focus");
    }
}

function reset() {
    resultsAvailable = false;
    resList.innerHTML = sInput.value = '';
    sInput.focus();
}

if (sInput) {
    sInput.onkeyup = function (e) {
        if (fuse) {
            const results = fuse.search(this.value.trim());
            if (results.length !== 0) {
                let resultSet = '';
                for (let item in results) {
                    resultSet += `<li class="post-entry"><header class="entry-header">${results[item].item.title}&nbsp;»</header>` +
                        `<a href="${results[item].item.permalink}" aria-label="${results[item].item.title}"></a></li>`;
                }
                resList.innerHTML = resultSet;
                resultsAvailable = true;
                first = resList.firstChild;
                last = resList.lastChild;
            } else {
                resultsAvailable = false;
                resList.innerHTML = '';
            }
        }
    };

    sInput.addEventListener('search', function (e) {
        if (!this.value) reset();
    });
}

document.onkeydown = function (e) {
    let key = e.key;
    var ae = document.activeElement;
    var searchbox = document.getElementById("searchbox");
    if (!searchbox) return;

    let inbox = searchbox.contains(ae);

    if (ae === sInput) {
        var elements = document.getElementsByClassName('focus');
        while (elements.length > 0) {
            elements[0].classList.remove('focus');
        }
    } else if (current_elem) {
        ae = current_elem;
    }

    if (key === "Escape") {
        reset();
    } else if (!resultsAvailable || !inbox) {
        return;
    } else if (key === "ArrowDown") {
        e.preventDefault();
        if (ae == sInput) {
            activeToggle(resList.firstChild.lastChild);
        } else if (ae.parentElement != last) {
            activeToggle(ae.parentElement.nextSibling.lastChild);
        }
    } else if (key === "ArrowUp") {
        e.preventDefault();
        if (ae.parentElement == first) {
            activeToggle(sInput);
        } else if (ae != sInput) {
            activeToggle(ae.parentElement.previousSibling.lastChild);
        }
    } else if (key === "ArrowRight") {
        ae.click();
    }
};
