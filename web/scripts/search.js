/**
 * @param {string} searchWord
 */
const navigateToSearch = searchWord => {
    let url = new URL(window.location)
    const queryString = url.search;
    const urlParams = new URLSearchParams(queryString);

    if (searchWord !== '') {
        urlParams.set('search', searchWord)
    } else {
        urlParams.delete('search')
    }

    urlParams.set("offset", "0")

    window.location = `${window.location.protocol}//${window.location.host}${window.location.pathname}?${urlParams.toString()}`.toString()
};

const showLoadingSpinner = () => {
    let spinner = document.getElementById('loading-spinner')
    spinner.classList.toggle('is-hidden');
}

document.addEventListener("DOMContentLoaded", () => {
    let searchInput = document.getElementById('search-input')
    let searchButtonSubmit = document.getElementById('search-input-submit')

    // For every new page fill the input with the value of the search
    const urlParams = new URLSearchParams(window.location.search);
    searchInput.value = urlParams.get('search')

    searchInput.addEventListener("keydown", e => {
        if (e.key !== "Enter") {
            return;
        }
        navigateToSearch(searchInput.value)
    });

    searchButtonSubmit.addEventListener('click', () => {
        navigateToSearch(searchInput.value)
    })
})