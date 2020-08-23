function handleReadMore() {
    let y = windowHeight - headerHeight;
    window.scroll({
        top: y,
        behavior: 'smooth'
    });
}

let windowHeight = window.innerHeight;

let header = document.getElementById('header');
let headerHeight = header.offsetHeight;

let readMoreElement = document.getElementById('atf-more');
readMoreElement.addEventListener('click', handleReadMore);

let atfTextElement = document.getElementById('atf-centre-text');

window.addEventListener('scroll', function(e) {

    let atfTextPosition = atfTextElement.getBoundingClientRect();

    if (atfTextPosition.y < (headerHeight + 20)) {
        header.classList.add('whiteHeader');
    } else {
        header.classList.remove('whiteHeader');
    }
});