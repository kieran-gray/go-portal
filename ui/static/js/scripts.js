function getAllCookies() {
    var allCookies = document.cookie;
    cookieArray = allCookies.split(';');
    var cookies = new Map();
    for(let cookie of cookieArray) {
        var splitCookie = cookie.split("=");
        cookies.set(splitCookie[0].trim(), splitCookie[1]);
    }
    return cookies
}

function getCookie(name) {
    cookies = getAllCookies()
    return cookies.get(name) || undefined
}

function setCookie(name, value, expiry) {
    const d = new Date();
    d.setTime(d.getTime() + (expiry*24*60*60*1000));
    let expires="expires="+d.toUTCString();
    document.cookie = name + "=" + value + ";" + expires + ";path=/;samesite=lax";
}

function search(cardType) {
    var cardType = cardType ? cardType : "portal-card"
    const input = document.getElementById("searchBar");
    const filter = input.value.toUpperCase();
    const cards = document.getElementsByClassName(cardType);
    for (let i = 0; i < cards.length; i++) {
        const text = cards[i].getAttribute("data-search");
        if (text.toUpperCase().indexOf(filter) > -1) {
            cards[i].classList.remove("hidden-search");
        } else {
            cards[i].classList.add("hidden-search");
        }
    }
}

function getPortalCardColor(theme, devMode) {
    if (theme === "dark") {
        return devMode ? "#242424" : "#181818";
    } else {
        return devMode ? "#ececec" : "#fff";
    }
}

function setTheme(theme, devMode) {
    document.documentElement.style.setProperty("--dev-mode", devMode ? "block" : "none");
    document.body.setAttribute('data-theme', theme);
    document.documentElement.style.setProperty("--portal-background-color", getPortalCardColor(theme, devMode));
}

function toggleDarkMode() {
    const devToggle = document.getElementById("dev-mode-toggle");
    const devMode = devToggle ? devToggle.checked : false
    const darkMode = document.getElementById("dark-mode-toggle").checked;
    setCookie("darkMode", darkMode, 10);
    setTheme(darkMode ? 'dark': 'light', devMode);
}

function setDarkModeToggle() {
    const darkToggle = document.getElementById("dark-mode-toggle");
    const prefersDarkScheme = window.matchMedia("(prefers-color-scheme: dark)");
    const darkMode = getCookie("darkMode") || `${prefersDarkScheme.matches}`;
    darkToggle.checked = darkMode == "true" ? true : false;
}

function setDevModeToggle() {
    const devToggle = document.getElementById("dev-mode-toggle");
    const devMode = getCookie("devMode") || false;
    devToggle.checked = devMode == "true" ? true : false;
}