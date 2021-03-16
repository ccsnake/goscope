/**
 * @param {string} searchWord
 * @param {number} searchMode
 */
const navigateToSearch = (searchWord, searchMode) => {
    let url = new URL(window.location)
    const queryString = url.search;
    const urlParams = new URLSearchParams(queryString);

    if (searchWord !== '') {
        urlParams.set('search', searchWord)
        urlParams.set('search-mode', searchMode.toString())
    } else {
        urlParams.delete('search')
        urlParams.delete('search-mode')
    }

    urlParams.set("offset", "0")

    window.location = `${window.location.protocol}//${window.location.host}${window.location.pathname}?${urlParams.toString()}`.toString()
};

const toggleNavigationShow = () => {
    document.getElementById('navbarBasicExample').classList.toggle('is-active')
    document.getElementById('navbarBurger').classList.toggle('is-active')
}

const showLoadingSpinner = () => {
    let spinner = document.getElementById('loading-spinner')
    spinner.classList.toggle('is-hidden');
}

document.addEventListener("DOMContentLoaded", () => {
    feather.replace()
    let searchInput = document.getElementById('search-input')
    let searchModeSelect = document.getElementById('search-mode-select')
    let searchButtonSubmit = document.getElementById('search-input-submit')

    // For every new page fill the input with the value of the search
    const urlParams = new URLSearchParams(window.location.search);

    if (searchModeSelect) {
        searchModeSelect.value = urlParams.get('search-mode') ?? 8
    }

    if (searchInput) {
        searchInput.value = urlParams.get('search')

        searchInput.addEventListener("keydown", e => {
            if (e.key !== "Enter") {
                return;
            }
            showLoadingSpinner()
            navigateToSearch(searchInput.value, searchModeSelect?.value ?? 8)
        });

        searchButtonSubmit.addEventListener('click', () => {
            showLoadingSpinner()
            navigateToSearch(searchInput.value, searchModeSelect?.value ?? 8)
        })
    }
})